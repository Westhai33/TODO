package controller

import (
	"TODO/internal/model"
	"TODO/internal/service"
	"context"
	"fmt"
)

// CreateTask создает новую задачу
func CreateTask(ctx context.Context, taskService *service.TaskService, userID int64, title, note string) (int64, error) {

	taskID, err := taskService.CreateTask(ctx, userID, title, note)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

// GetTask получает задачу по ID
func GetTask(ctx context.Context, taskService *service.TaskService, taskID int64) (*model.Task, error) {

	task, err := taskService.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
	}

	return task, nil
}

// GetAllTasks получает все задачи
func GetAllTasks(ctx context.Context, taskService *service.TaskService) ([]model.Task, error) {

	tasks, err := taskService.GetAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех задач: %w", err)
	}

	return tasks, nil
}

// UpdateTask обновляет задачу
func UpdateTask(ctx context.Context, taskService *service.TaskService, taskID int64, title, note string, done bool) error {
	err := taskService.UpdateTask(ctx, taskID, title, note, done)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTask удаляет задачу
func DeleteTask(ctx context.Context, taskService *service.TaskService, taskID int64) error {
	err := taskService.DeleteTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	return nil
}
