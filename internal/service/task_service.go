package service

import (
	"TODO/internal/dao"
	"TODO/internal/model"
	"TODO/internal/pool"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskService struct {
	pool *pgxpool.Pool
	wp   *pool.WorkerPool
}

func NewTaskService(dbPool *pgxpool.Pool, wp *pool.WorkerPool) *TaskService {
	return &TaskService{
		pool: dbPool,
		wp:   wp,
	}
}

// CreateTask создаёт новую задачу через общий worker pool
func (s *TaskService) CreateTask(ctx context.Context, userID int64, title, note string) (int64, error) {
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

		errCh <- nil
	})

	err := <-errCh
	return taskID, err
}

// UpdateTask обновляет задачу
func (s *TaskService) UpdateTask(ctx context.Context, taskID int64, title, note string, done bool) error {
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

		errCh <- nil
	})

	return <-errCh
}

// DeleteTask удаляет задачу
func (s *TaskService) DeleteTask(ctx context.Context, taskID int64) error {
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		if err := dao.DeleteTask(ctx, taskID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления задачи с ID %d: %w", taskID, err)
			return
		}

		errCh <- nil
	})

	return <-errCh
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
