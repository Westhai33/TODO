package controller

import (
	"TODO/internal/model"
	"TODO/internal/service"
	"context"
	"fmt"
)

// CreateTask создает новую задачу с проверкой существования пользователя
func CreateTask(ctx context.Context, taskService *service.TaskService, userService *service.UserService, userID int64, title, note string) (int64, error) {
	_, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("пользователь с ID %d не найден: %w", userID, err)
	}

	taskID, err := taskService.CreateTask(ctx, userID, title, note)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании задачи: %w", err)
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

// UpdateTask обновляет задачу с проверкой существования
func UpdateTask(ctx context.Context, taskService *service.TaskService, taskID int64, title, note string, done bool) error {
	_, err := taskService.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("задача с ID %d не найдена: %w", taskID, err)
	}

	err = taskService.UpdateTask(ctx, taskID, title, note, done)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи с ID %d: %w", taskID, err)
	}

	return nil
}

// DeleteTask удаляет задачу с проверкой существования
func DeleteTask(ctx context.Context, taskService *service.TaskService, taskID int64) error {
	_, err := taskService.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("задача с ID %d не найдена: %w", taskID, err)
	}

	err = taskService.DeleteTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи с ID %d: %w", taskID, err)
	}

	return nil
}
