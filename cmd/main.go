package main

import (
	"TODO/internal/config"
	"TODO/internal/dao"
	"TODO/internal/pool"
	"TODO/internal/service"
	"TODO/internal/view"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefulShutdown(cancel)

	initDatabase(cfg)
	defer dao.Closedb()

	dbPool := dao.GetPool()

	wp := pool.NewWorkerPool(2)

	userService, taskService := initServices(dbPool, wp)

	view.RunInteractiveMode(ctx, taskService, userService)

	log.Println("Система завершила работу")
}

func initDatabase(cfg *config.Config) {
	port, err := strconv.Atoi(cfg.DBPort)
	if err != nil {
		log.Fatalf("Некорректный формат порта базы данных: %v", err)
	}

	dao.Initdb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, port)
	log.Println("Подключение к базе данных успешно выполнено.")
}

func initServices(dbPool *pgxpool.Pool, wp *pool.WorkerPool) (*service.UserService, *service.TaskService) {
	userService := service.NewUserService(dbPool, wp)
	taskService := service.NewTaskService(dbPool, wp)
	return userService, taskService
}

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
