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

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(dbPool *pgxpool.Pool) *UserService {
	return &UserService{
		pool: dbPool,
	}
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, username string) (int64, error) {
	var userID int64

	newUser := model.User{
		Username:  username,
		CreatedAt: time.Now().UTC(),
	}

	userID, err := dao.CreateUser(ctx, newUser, s.pool)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	log.Printf("Создан новый пользователь с ID %d", userID)
	return userID, nil
}

// GetUserByID возвращает пользователя по его ID
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	user, err := dao.GetUserByID(ctx, userID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	return user, nil
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := dao.GetAllUsers(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех пользователей: %w", err)
	}

	log.Printf("Получены все пользователи: %+v", users)
	return users, nil
}

// UpdateUser обновляет данные пользователя
func (s *UserService) UpdateUser(ctx context.Context, userID int64, username string) error {
	user, err := dao.GetUserByID(ctx, userID, s.pool)
	if err != nil {
		return fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	user.Username = username

	if err := dao.UpdateUser(ctx, *user, s.pool); err != nil {
		return fmt.Errorf("ошибка обновления данных пользователя с ID %d: %w", userID, err)
	}

	log.Printf("Пользователь с ID %d был обновлен", userID)
	return nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	if err := dao.DeleteUser(ctx, userID, s.pool); err != nil {
		return fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
	}

	log.Printf("Пользователь с ID %d удален", userID)
	return nil
}

// GetUserNameByID возвращает имя пользователя по его ID
func (s *UserService) GetUserNameByID(ctx context.Context, userID int64) (string, error) {
	username, err := dao.GetUserNameByID(ctx, userID, s.pool)
	if err != nil {
		return "", fmt.Errorf("ошибка получения имени пользователя с ID %d: %w", userID, err)
	}

	return username, nil
}
