package controller

import (
	"TODO/internal/metrics"
	"TODO/internal/model"
	"TODO/internal/service"
	"TODO/internal/tracing"
	"context"
	"fmt"
)

// CreateTask создает новую задачу с проверкой существования пользователя и трассировкой.
func CreateTask(ctx context.Context, taskService *service.TaskService, userService *service.UserService, userID int64, title, note string) (int64, error) {
	ctx, span := tracing.GetTracer().Start(ctx, "CreateTask")
	defer span.End()

	span.AddEvent("Начинаем создание задачи")

	_, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return 0, fmt.Errorf("пользователь с ID %d не найден: %w", userID, err)
	}

	taskID, err := taskService.CreateTask(ctx, userID, title, note)
	if err != nil {
		span.RecordError(err)
		return 0, fmt.Errorf("ошибка при создании задачи: %w", err)
	}

	metrics.IncrementTaskCreated("created")

	span.AddEvent("Задача успешно создана")

	return taskID, nil
}

// GetTask получает задачу по ID с трассировкой.
func GetTask(ctx context.Context, taskService *service.TaskService, taskID int64) (*model.Task, error) {
	ctx, span := tracing.GetTracer().Start(ctx, "GetTask")
	defer span.End()

	span.AddEvent("Начинаем получение задачи")

	task, err := taskService.GetTask(ctx, taskID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
	}

	span.AddEvent("Задача успешно получена")
	return task, nil
}

// GetAllTasks получает все задачи с трассировкой.
func GetAllTasks(ctx context.Context, taskService *service.TaskService) ([]model.Task, error) {
	ctx, span := tracing.GetTracer().Start(ctx, "GetAllTasks")
	defer span.End()

	span.AddEvent("Начинаем получение всех задач")

	tasks, err := taskService.GetAllTasks(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("ошибка получения всех задач: %w", err)
	}

	span.AddEvent("Все задачи успешно получены")
	return tasks, nil
}

// UpdateTask обновляет задачу с проверкой существования и трассировкой.
func UpdateTask(ctx context.Context, taskService *service.TaskService, taskID int64, title, note string, done bool) error {
	ctx, span := tracing.GetTracer().Start(ctx, "UpdateTask")
	defer span.End()

	span.AddEvent("Начинаем обновление задачи")

	_, err := taskService.GetTask(ctx, taskID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("задача с ID %d не найдена: %w", taskID, err)
	}

	err = taskService.UpdateTask(ctx, taskID, title, note, done)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("ошибка обновления задачи с ID %d: %w", taskID, err)
	}

	span.AddEvent("Задача успешно обновлена")
	return nil
}

// DeleteTask удаляет задачу с проверкой существования и трассировкой.
func DeleteTask(ctx context.Context, taskService *service.TaskService, taskID int64) error {
	ctx, span := tracing.GetTracer().Start(ctx, "DeleteTask")
	defer span.End()

	span.AddEvent("Начинаем удаление задачи")

	_, err := taskService.GetTask(ctx, taskID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("задача с ID %d не найдена: %w", taskID, err)
	}

	err = taskService.DeleteTask(ctx, taskID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("ошибка удаления задачи с ID %d: %w", taskID, err)
	}

	span.AddEvent("Задача успешно удалена")
	return nil
}
