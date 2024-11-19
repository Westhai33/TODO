package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// TaskMessage представляет структуру сообщения о задаче
type TaskMessage struct {
	TimeStamp time.Time `json:"timestamp"` // Время отправки сообщения
	Operation string    `json:"operation"` // Тип операции (create-task, update-task, delete-task)
	TaskID    int64     `json:"task_id"`   // ID задачи
	UserID    int64     `json:"user_id"`   // ID пользователя
	Title     string    `json:"title"`     // Заголовок задачи
	Note      string    `json:"note"`      // Описание задачи
	Done      bool      `json:"done"`      // Статус выполнения задачи
}

// Producer представляет Kafka продюсера
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer создает нового продюсера Kafka
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 1

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	producer := &Producer{
		producer: syncProducer,
		topic:    topic,
	}

	return producer, nil
}

// SendTaskMessage отправляет сообщение о задаче в Kafka
func (p Producer) SendTaskMessage(message TaskMessage) error {
	// Сериализация сообщения
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("p.producer.SendMessage: %w", err)
	}

	log.Printf("Сообщение успешно отправлено в Kafka. Тема: %s, Раздел: %d, Смещение: %d, Сообщение: %+v\n", p.topic, partition, offset, message)

	return nil
}

// SendKafkaErrorMessage отправляет сообщение об ошибке в Kafka
func (p Producer) SendKafkaErrorMessage(operation string, taskID int64, userID int64, description string) error {
	errorMessage := TaskMessage{
		TimeStamp: time.Now(),
		Operation: operation,
		TaskID:    taskID,
		UserID:    userID,
		Title:     "",
		Note:      description,
		Done:      false,
	}

	msg, err := json.Marshal(errorMessage)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("p.producer.SendMessage: %w", err)
	}

	log.Printf("Сообщение об ошибке успешно отправлено в Kafka. Тема: %s, Раздел: %d, Смещение: %d, Сообщение: %+v\n", p.topic, partition, offset, errorMessage)

	return nil
}

// Close закрывает продюсера
func (p Producer) Close() error {
	err := p.producer.Close()
	if err != nil {
		return fmt.Errorf("p.producer.Close: %w", err)
	}

	return nil
}
