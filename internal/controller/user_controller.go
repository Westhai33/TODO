package controller

import (
	"TODO/internal/model"
	"TODO/internal/service"
	"TODO/internal/tracing"
	"context"
	"fmt"
)

// CreateUser создает нового пользователя с трассировкой.
func CreateUser(ctx context.Context, userService *service.UserService, username string) (int64, error) {
	ctx, span := tracing.GetTracer().Start(ctx, "CreateUser")
	defer span.End()

	span.AddEvent("Начинаем создание пользователя")

	userID, err := userService.CreateUser(ctx, username)
	if err != nil {
		span.RecordError(err)
		return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	span.AddEvent("Пользователь успешно создан")

	return userID, nil
}

// GetUserByID возвращает пользователя по его ID с трассировкой.
func GetUserByID(ctx context.Context, userService *service.UserService, userID int64) (*model.User, error) {
	ctx, span := tracing.GetTracer().Start(ctx, "GetUserByID")
	defer span.End()

	span.AddEvent("Начинаем получение пользователя по ID")

	user, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	span.AddEvent("Пользователь успешно получен")

	return user, nil
}

// GetAllUsers возвращает всех пользователей с трассировкой.
func GetAllUsers(ctx context.Context, userService *service.UserService) ([]model.User, error) {
	ctx, span := tracing.GetTracer().Start(ctx, "GetAllUsers")
	defer span.End()

	span.AddEvent("Начинаем получение всех пользователей")

	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("ошибка получения списка пользователей: %w", err)
	}

	span.AddEvent("Все пользователи успешно получены")

	return users, nil
}

// UpdateUser обновляет существующего пользователя с трассировкой.
func UpdateUser(ctx context.Context, userService *service.UserService, userID int64, username string) error {
	ctx, span := tracing.GetTracer().Start(ctx, "UpdateUser")
	defer span.End()

	span.AddEvent("Начинаем обновление пользователя")

	_, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("пользователь с ID %d не найден: %w", userID, err)
	}

	err = userService.UpdateUser(ctx, userID, username)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("ошибка обновления пользователя с ID %d: %w", userID, err)
	}

	span.AddEvent("Пользователь успешно обновлен")

	return nil
}

// DeleteUser удаляет пользователя по его ID с трассировкой.
func DeleteUser(ctx context.Context, userService *service.UserService, userID int64) error {
	ctx, span := tracing.GetTracer().Start(ctx, "DeleteUser")
	defer span.End()

	span.AddEvent("Начинаем удаление пользователя")

	_, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("пользователь с ID %d не найден: %w", userID, err)
	}

	err = userService.DeleteUser(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
	}

	span.AddEvent("Пользователь успешно удален")

	return nil
}
