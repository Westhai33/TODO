package dao

import (
	"TODO/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// CreateTask создает новую задачу.
func CreateTask(ctx context.Context, task model.Task, pool *pgxpool.Pool) (int64, error) {
	tx, conn, err := NewTransactionManager(pool).BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
		}
	}()

	query := `INSERT INTO tasks (user_id, title, note, done, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var taskID int64
	err = tx.QueryRow(ctx, query, task.UserID, task.Title, task.Note, task.Done, task.CreatedAt, task.UpdatedAt).Scan(&taskID)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания задачи: %w", err)
	}

	if err = NewTransactionManager(pool).CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return taskID, nil
}

// GetTaskByID извлекает задачу по ее идентификатору .
func GetTaskByID(ctx context.Context, taskID int64, pool *pgxpool.Pool) (*model.Task, error) {
	tx, conn, err := NewTransactionManager(pool).BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}

	var task model.Task
	err = tx.QueryRow(ctx, `SELECT id, user_id, title, note, done, created_at, updated_at FROM tasks WHERE id = $1`, taskID).
		Scan(&task.ID, &task.UserID, &task.Title, &task.Note, &task.Done, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения задачи с ID %d: %w", taskID, err)
	}

	if err = NewTransactionManager(pool).CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &task, nil
}

// UpdateTask обновление данных задачи.
func UpdateTask(ctx context.Context, task model.Task, pool *pgxpool.Pool) error {
	tx, conn, err := NewTransactionManager(pool).BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE tasks SET title = $2, note = $3, done = $4, updated_at = $5 WHERE id = $1`,
		task.ID, task.Title, task.Note, task.Done, task.UpdatedAt)
	if err != nil {
		if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return fmt.Errorf("ошибка обновления задачи с ID %d: %w", task.ID, err)
	}

	if err = NewTransactionManager(pool).CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeleteTask удаляет задачу.
func DeleteTask(ctx context.Context, taskID int64, pool *pgxpool.Pool) error {
	tx, conn, err := NewTransactionManager(pool).BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, taskID)
	if err != nil {
		if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return fmt.Errorf("ошибка удаления задачи с ID %d: %w", taskID, err)
	}

	if err = NewTransactionManager(pool).CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// GetAllTasks извлекает все задания.
func GetAllTasks(ctx context.Context, pool *pgxpool.Pool) ([]model.Task, error) {
	tx, conn, err := NewTransactionManager(pool).BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `SELECT id, user_id, title, note, done, created_at, updated_at FROM tasks`)
	if err != nil {
		if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения всех задач: %w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Note, &task.Done, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования задачи: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		if rollbackErr := NewTransactionManager(pool).RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка итерации по строкам задач: %w", err)
	}

	if err = NewTransactionManager(pool).CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return tasks, nil
}
