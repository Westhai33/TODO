package service

import (
	"TODO/internal/dao"
	"TODO/internal/model"
	"TODO/internal/pool"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserService struct {
	pool *pgxpool.Pool
	wp   *pool.WorkerPool
}

func NewUserService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool) *UserService {
	return &UserService{
		pool: dbPool,
		wp:   workerPool,
	}
}

// CreateUser создает нового пользователя через общий worker pool
func (s *UserService) CreateUser(ctx context.Context, username string) (int64, error) {
	var userID int64
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		newUser := model.User{
			Username:  username,
			CreatedAt: time.Now().UTC(),
		}

		var err error
		userID, err = dao.CreateUser(ctx, newUser, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка создания пользователя: %w", err)
			return
		}

		errCh <- nil
	})

	err := <-errCh
	return userID, err
}

// GetUserByID возвращает пользователя по его ID
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	user, err := dao.GetUserByID(ctx, userID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	return user, nil
}

// GetAllUsers возвращает всех пользователей через worker pool
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
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		user, err := dao.GetUserByID(ctx, userID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
			return
		}

		user.Username = username

		if err := dao.UpdateUser(ctx, *user, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка обновления данных пользователя с ID %d: %w", userID, err)
			return
		}

		errCh <- nil
	})

	return <-errCh
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		if err := dao.DeleteUser(ctx, userID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
			return
		}

		errCh <- nil
	})

	return <-errCh
}

// GetUserNameByID возвращает имя пользователя по его ID через worker pool
func (s *UserService) GetUserNameByID(ctx context.Context, userID int64) (string, error) {
	username, err := dao.GetUserNameByID(ctx, userID, s.pool)
	if err != nil {
		return "", fmt.Errorf("ошибка получения имени пользователя с ID %d: %w", userID, err)
	}

	return username, nil
}
