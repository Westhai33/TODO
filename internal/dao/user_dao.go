package dao

import (
	"TODO/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// CreateUser создает нового пользователя.
func CreateUser(ctx context.Context, user model.User, pool *pgxpool.Pool) (int64, error) {
	tm := NewTransactionManager(pool)
	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
		}
	}()

	query := `INSERT INTO users (username, created_at) VALUES ($1, $2) RETURNING id`
	var userID int64
	err = tx.QueryRow(ctx, query, user.Username, user.CreatedAt).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return userID, nil
}

// GetUserByID возвращает пользователя по его ID.
func GetUserByID(ctx context.Context, userID int64, pool *pgxpool.Pool) (*model.User, error) {
	tm := NewTransactionManager(pool)
	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = tx.QueryRow(ctx, `SELECT id, username, created_at FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &user, nil
}

// UpdateUser обновляет данные пользователя.
func UpdateUser(ctx context.Context, user model.User, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)
	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE users SET username = $2 WHERE id = $1`, user.ID, user.Username)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			return fmt.Errorf("ошибка обновления пользователя с ID %d: %w", user.ID, err)
		}
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeleteUser удаляет пользователя по его ID.
func DeleteUser(ctx context.Context, userID int64, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)
	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// GetAllUsers возвращает список всех пользователей.
func GetAllUsers(ctx context.Context, pool *pgxpool.Pool) ([]model.User, error) {
	tm := NewTransactionManager(pool)
	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `SELECT id, username, created_at FROM users`)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения пользователей: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования пользователя: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка итерации по строкам пользователей: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return users, nil
}

// GetUserNameByID returns the username of a user by their ID.
func GetUserNameByID(ctx context.Context, userID int64, pool *pgxpool.Pool) (string, error) {

	tx, conn, err := NewTransactionManager(pool).BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return "", err
	}

	var username string
	err = tx.QueryRow(ctx, `SELECT username FROM users WHERE id = $1`, userID).Scan(&username)
	if err != nil {

		if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return "", fmt.Errorf("ошибка получения имени пользователя с ID %d: %w", userID, err)
	}

	if err = NewTransactionManager(pool).CommitTransaction(ctx, tx, conn); err != nil {
		return "", fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return username, nil
}
