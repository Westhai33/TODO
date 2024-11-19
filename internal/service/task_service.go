package service

import (
	"TODO/internal/cache"
	"TODO/internal/dao"
	"TODO/internal/kafka"
	"TODO/internal/model"
	"TODO/internal/pool"
	"context"
	"fmt"
	"log"
	"time"

	"TODO/internal/tracing"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

// TaskService управляет задачами
type TaskService struct {
	pool          *pgxpool.Pool
	wp            *pool.WorkerPool
	taskCache     *cache.RedisCache[string, model.Task]
	tracer        trace.Tracer
	kafkaProducer *kafka.Producer
}

// NewTaskService создаёт новый TaskService с необходимыми зависимостями
func NewTaskService(dbPool *pgxpool.Pool, wp *pool.WorkerPool, taskCache *cache.RedisCache[string, model.Task], kafkaProducer *kafka.Producer) *TaskService {
	return &TaskService{
		pool:          dbPool,
		wp:            wp,
		taskCache:     taskCache,
		tracer:        tracing.GetTracer(),
		kafkaProducer: kafkaProducer,
	}
}

// sendKafkaMessage отправляет сообщение о задаче в Kafka
func (s *TaskService) sendKafkaMessage(operation string, taskID int64, userID int64, title, note string, done bool) error {
	orderMessage := kafka.TaskMessage{
		TimeStamp: time.Now(),
		Operation: operation,
		TaskID:    taskID,
		UserID:    userID,
		Title:     title,
		Note:      note,
		Done:      done,
	}

	if err := s.kafkaProducer.SendTaskMessage(orderMessage); err != nil {
		return fmt.Errorf("ошибка отправки сообщения о задаче в Kafka: %v", err)
	}

	return nil
}

// CreateTask создаёт новую задачу через общий worker pool и отправляет сообщение в Kafka
func (s *TaskService) CreateTask(ctx context.Context, userID int64, title, note string) (int64, error) {
	ctx, span := s.tracer.Start(ctx, "CreateTask")
	defer span.End()

	var taskID int64
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		newTask := model.Task{
			UserID:    userID,
			Title:     title,
			Note:      note,
			Done:      false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		var err error
		taskID, err = dao.CreateTask(ctx, newTask, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка создания задачи: %w", err)
			return
		}

		if err := s.sendKafkaMessage("create-task", taskID, userID, title, note, false); err != nil {
			log.Printf("Ошибка отправки сообщения о задаче в Kafka: %v", err)
		}

		cacheKey := fmt.Sprintf("task_%d", taskID)
		if err := s.taskCache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша задачи с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	err := <-errCh
	return taskID, err
}

// UpdateTask обновляет задачу и сбрасывает кэш, отправляя обновление в Kafka
func (s *TaskService) UpdateTask(ctx context.Context, taskID int64, title, note string, done bool) error {
	ctx, span := s.tracer.Start(ctx, "UpdateTask")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		task, err := dao.GetTaskByID(ctx, taskID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
			return
		}

		task.Title = title
		task.Note = note
		task.Done = done
		task.UpdatedAt = time.Now()

		if err := dao.UpdateTask(ctx, *task, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка обновления данных задачи с ID %d: %w", taskID, err)
			return
		}

		if err := s.sendKafkaMessage("update-task", taskID, task.UserID, title, note, done); err != nil {
			log.Printf("Ошибка отправки сообщения о задаче в Kafka: %v", err)
		}

		cacheKey := fmt.Sprintf("task_%d", taskID)
		if err := s.taskCache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша задачи с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	return <-errCh
}

// DeleteTask удаляет задачу и сбрасывает кэш, отправляя удаление в Kafka
func (s *TaskService) DeleteTask(ctx context.Context, taskID int64) error {
	ctx, span := s.tracer.Start(ctx, "DeleteTask")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		if err := dao.DeleteTask(ctx, taskID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления задачи с ID %d: %w", taskID, err)
			return
		}

		if err := s.sendKafkaMessage("delete-task", taskID, 0, "", "", false); err != nil {
			log.Printf("Ошибка отправки сообщения о задаче в Kafka: %v", err)
		}

		cacheKey := fmt.Sprintf("task_%d", taskID)
		if err := s.taskCache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша задачи с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	return <-errCh
}

// GetTask получает задачу по ID
func (s *TaskService) GetTask(ctx context.Context, taskID int64) (*model.Task, error) {
	ctx, span := s.tracer.Start(ctx, "GetTask")
	defer span.End()

	task, err := dao.GetTaskByID(ctx, taskID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
	}

	return task, nil
}

// GetAllTasks получает все задачи
func (s *TaskService) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllTasks")
	defer span.End()

	tasks, err := dao.GetAllTasks(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех задач: %w", err)
	}

	return tasks, nil
}
