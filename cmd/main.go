package main

import (
	"TODO/internal/model"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"TODO/internal/api/v1"
	"TODO/internal/cache"
	"TODO/internal/client"
	"TODO/internal/config"
	"TODO/internal/dao"
	"TODO/internal/gateway"
	"TODO/internal/kafka"
	"TODO/internal/metrics"
	"TODO/internal/pool"
	"TODO/internal/server"
	"TODO/internal/service"
	"TODO/internal/tracing"
	"TODO/internal/view"
)

func main() {
	err := addServersSectionToOpenAPI()
	if err != nil {
		log.Fatalf("Ошибка при добавлении секции servers в OpenAPI: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefulShutdown(cancel)

	cfg := config.LoadConfig()

	shutdownTracer := tracing.InitTracer(cfg.ServiceName, cfg.TracingURL)
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Fatalf("Ошибка завершения трейсинга: %v", err)
		}
	}()
	log.Println("Трейсинг инициализирован:", cfg.TracingURL)

	initDatabase(cfg)
	defer dao.Closedb()

	dbPool := dao.GetPool()

	redisClient := initRedis(cfg)
	defer redisClient.Close()

	go metrics.StartMetricsServer(cfg.MetricsAddr)
	log.Println("Prometheus сервер запущен. Адрес метрик:", cfg.MetricsAddr)

	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Ошибка при инициализации Kafka Producer: %v", err)
	}
	defer kafkaProducer.Close()

	wp := pool.NewWorkerPool(2)

	// Инициализация сервисов
	userService, taskService := initServices(dbPool, wp, kafkaProducer, redisClient)

	// Запуск серверов
	startServers(ctx, cfg, userService, taskService)

	// Запуск интерактивного режима
	grpcClients := setupGRPCClients(cfg.GrpcPort)
	view.RunInteractiveMode(ctx, grpcClients, wp)
}

// Функция для инициализации Redis клиента
func initRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})
}

// Функция для инициализации сервисов с Redis-кэшем и Kafka
func initServices(dbPool *pgxpool.Pool, wp *pool.WorkerPool, kafkaProducer *kafka.Producer, redisClient *redis.Client) (
	*service.UserService, *service.TaskService) {

	cacheConfig := cache.CacheConfig{
		DefaultTTL: 10 * time.Minute,
	}

	userCache := cache.NewRedisCache[string, model.User](redisClient, cacheConfig)
	taskCache := cache.NewRedisCache[string, model.Task](redisClient, cacheConfig)

	userService := service.NewUserService(dbPool, wp, userCache)
	taskService := service.NewTaskService(dbPool, wp, taskCache, kafkaProducer)

	return userService, taskService
}

// Запуск gRPC и HTTP Gateway серверов
func startServers(
	ctx context.Context, cfg *config.Config,
	userService *service.UserService, taskService *service.TaskService) {

	go func() {
		if err := startGRPCServer(cfg.GrpcPort, userService, taskService); err != nil {
			log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
		}
		log.Println("gRPC сервер завершил работу")
	}()

	go func() {
		log.Printf("Запуск HTTP Gateway на порту %s", cfg.HttpPort)
		if err := gateway.RunGateway(ctx, "localhost:"+cfg.GrpcPort, "localhost:"+cfg.HttpPort); err != nil {
			log.Fatalf("Ошибка при запуске HTTP Gateway: %v", err)
		}
	}()
}

// Инициализация базы данных
func initDatabase(cfg *config.Config) {
	port, err := strconv.Atoi(cfg.DBPort)
	if err != nil {
		log.Fatalf("Некорректный формат порта: %v", err)
	}
	dao.Initdb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, port)
}

// Запуск gRPC сервера
func startGRPCServer(grpcPort string, userService *service.UserService, taskService *service.TaskService) error {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return fmt.Errorf("не удалось начать слушать порт %s: %w", grpcPort, err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(50*1024*1024),
		grpc.MaxSendMsgSize(50*1024*1024),
	)

	// Убираем WorkerPool из параметров
	v1.RegisterAPIServiceServer(grpcServer, server.NewAPIServiceServer(userService, taskService))

	log.Printf("gRPC сервер запущен на порту %s", grpcPort)
	return grpcServer.Serve(lis)
}

// Подключение к gRPC клиентам
func setupGRPCClients(grpcPort string) *client.APIServiceClientWrapper {
	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*1024*1024),
			grpc.MaxCallSendMsgSize(50*1024*1024),
		),
	}

	grpcConn, err := grpc.NewClient("localhost:"+grpcPort, clientOptions...)
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC серверу: %v", err)
	}

	grpcClientWrapper, err := client.NewAPIServiceClientWrapper(grpcConn)
	if err != nil {
		log.Fatalf("Ошибка при создании обертки gRPC клиента: %v", err)
	}

	return grpcClientWrapper
}

// Обработка сигналов завершения работы
func gracefulShutdown(cancelFunc context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Получен сигнал завершения. Завершаем работу...")
		cancelFunc()
		os.Exit(0)
	}()
}

func addServersSectionToOpenAPI() error {
	openAPIFile := "./internal/api/v1/openapi.yaml"

	data, err := os.ReadFile(openAPIFile)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	content := string(data)

	if strings.Contains(content, "info:") {
		content = strings.Replace(content, `title: ""`, `title: "Order Management API"`, 1)
		content = strings.Replace(content, `version: 0.0.1`, `version: 1.0.0`, 1)
		content = strings.Replace(content, `description: ""`, `description: "API для управления заказами, упаковками, возвратами, статусами и пользователями"`, 1)
	} else {
		infoSection := `
info:
  title: "Order Management API"
  description: "API для управления заказами, упаковками, возвратами, статусами и пользователями"
  version: "1.0.0"
`
		content = infoSection + content
	}

	if !strings.Contains(content, "servers:") {
		serversSection := `
servers:
  - url: http://localhost:8080
    description: "Local HTTP Gateway"
`
		content += serversSection
	}

	err = os.WriteFile(openAPIFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("ошибка при записи в файл: %w", err)
	}

	return nil
}
