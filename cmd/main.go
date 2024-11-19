package main

import (
	"TODO/internal/config"
	"TODO/internal/dao"
	"TODO/internal/service"
	"TODO/internal/view"
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefulShutdown(cancel)

	initDatabase(cfg)
	defer dao.Closedb()

	dbPool := dao.GetPool()

	userService, taskService := initServices(dbPool)

	view.RunInteractiveMode(ctx, taskService, userService)
}

func initDatabase(cfg *config.Config) {
	port, err := strconv.Atoi(cfg.DBPort)
	if err != nil {
		log.Fatalf("Некорректный формат порта базы данных: %v", err)
	}

	dao.Initdb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, port)
	log.Println("Подключение к базе данных успешно выполнено.")
}

// Инициализация сервисов
func initServices(dbPool *pgxpool.Pool) (*service.UserService, *service.TaskService) {
	userService := service.NewUserService(dbPool)
	taskService := service.NewTaskService(dbPool)
	return userService, taskService
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
