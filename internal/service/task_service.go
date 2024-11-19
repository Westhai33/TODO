package service

import (
	"TODO/internal/dao"
	"TODO/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskService struct {
	pool *pgxpool.Pool
}

func NewTaskService(dbPool *pgxpool.Pool) *TaskService {
	return &TaskService{
		pool: dbPool,
	}
}

// CreateTask создаёт новую задачу
func (s *TaskService) CreateTask(ctx context.Context, userID int64, title, note string) (int64, error) {
	var taskID int64

	newTask := model.Task{
		UserID:    userID,
		Title:     title,
		Note:      note,
		Done:      false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Вставка задачи в базу данных
	taskID, err := dao.CreateTask(ctx, newTask, s.pool)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания задачи: %w", err)
	}

	log.Printf("Создана задача с ID %d", taskID)
	return taskID, nil
}

// UpdateTask обновляет задачу
func (s *TaskService) UpdateTask(ctx context.Context, taskID int64, title, note string, done bool) error {
	task, err := dao.GetTaskByID(ctx, taskID, s.pool)
	if err != nil {
		return fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
	}

	task.Title = title
	task.Note = note
	task.Done = done
	task.UpdatedAt = time.Now()

	if err := dao.UpdateTask(ctx, *task, s.pool); err != nil {
		return fmt.Errorf("ошибка обновления данных задачи с ID %d: %w", taskID, err)
	}

	log.Printf("Обновлена задача с ID %d", taskID)
	return nil
}

// DeleteTask удаляет задачу
func (s *TaskService) DeleteTask(ctx context.Context, taskID int64) error {
	if err := dao.DeleteTask(ctx, taskID, s.pool); err != nil {
		return fmt.Errorf("ошибка удаления задачи с ID %d: %w", taskID, err)
	}

	log.Printf("Удалена задача с ID %d", taskID)
	return nil
}

// GetTask получает задачу по ID
func (s *TaskService) GetTask(ctx context.Context, taskID int64) (*model.Task, error) {
	task, err := dao.GetTaskByID(ctx, taskID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
	}

	return task, nil
}

// GetAllTasks получает все задачи
func (s *TaskService) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	tasks, err := dao.GetAllTasks(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех задач: %w", err)
	}

	return tasks, nil
}
