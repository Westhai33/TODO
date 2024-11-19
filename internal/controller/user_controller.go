package controller

import (
	"TODO/internal/model"
	"TODO/internal/service"
	"context"
	"fmt"
)

// CreateUser создает нового пользователя.
func CreateUser(ctx context.Context, userService *service.UserService, username string) (int64, error) {
	userID, err := userService.CreateUser(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return userID, nil
}

// GetUserByID возвращает пользователя по его ID.
func GetUserByID(ctx context.Context, userService *service.UserService, userID int64) (*model.User, error) {
	user, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}
	return user, nil
}

// GetAllUsers возвращает всех пользователей.
func GetAllUsers(ctx context.Context, userService *service.UserService) ([]model.User, error) {
	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка пользователей: %w", err)
	}

	return users, nil
}

// UpdateUser обновляет существующего пользователя с проверкой его наличия.
func UpdateUser(ctx context.Context, userService *service.UserService, userID int64, username string) error {
	_, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("пользователь с ID %d не найден: %w", userID, err)
	}

	err = userService.UpdateUser(ctx, userID, username)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя с ID %d: %w", userID, err)
	}

	return nil
}

// DeleteUser удаляет пользователя по его ID с проверкой его наличия.
func DeleteUser(ctx context.Context, userService *service.UserService, userID int64) error {
	_, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("пользователь с ID %d не найден: %w", userID, err)
	}

	err = userService.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
	}

	return nil
}
