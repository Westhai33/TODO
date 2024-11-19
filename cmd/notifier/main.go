package main

import (
	"TODO/internal/config"
	"TODO/internal/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	consumerGroup, err := kafka.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		log.Fatalf("Ошибка создания consumer группы: %v", err)
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Println("Получен сигнал завершения. Завершаем работу...")
		os.Exit(0)
	}()

	handler := kafka.NotifierHandler{}

	kafka.StartConsumer(consumerGroup, []string{cfg.KafkaTopic}, handler)

	select {}
}
