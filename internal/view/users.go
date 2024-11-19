package view

import (
	v1 "TODO/internal/api/v1"
	"TODO/internal/client"
	"TODO/internal/pool"
	"context"
	"fmt"
	"strconv"
)

// handleUserCommands обрабатывает команды, связанные с пользователями
func handleUserCommands(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	switch args[0] {
	case "create-user":
		if len(args) != 2 {
			fmt.Println("Использование: create-user [username]")
			return
		}
		workerPool.SubmitTask(func() {
			handleCreateUser(ctx, args[1], grpcWrapper)
		})
	case "get-user":
		if len(args) != 2 {
			fmt.Println("Использование: get-user [userID]")
			return
		}
		workerPool.SubmitTask(func() {
			handleGetUser(ctx, args[1], grpcWrapper)
		})
	case "get-users":
		workerPool.SubmitTask(func() {
			handleGetAllUsers(ctx, grpcWrapper)
		})
	case "update-user":
		if len(args) != 3 {
			fmt.Println("Использование: update-user [userID] [username]")
			return
		}
		workerPool.SubmitTask(func() {
			handleUpdateUser(ctx, args[1], args[2], grpcWrapper)
		})
	case "delete-user":
		if len(args) != 2 {
			fmt.Println("Использование: delete-user [userID]")
			return
		}
		workerPool.SubmitTask(func() {
			handleDeleteUser(ctx, args[1], grpcWrapper)
		})
	default:
		fmt.Println("Неизвестная команда для пользователя.")
	}
}

// handleCreateUser создает нового пользователя
func handleCreateUser(ctx context.Context, username string, grpcWrapper *client.APIServiceClientWrapper) {
	req := &v1.CreateUserRequest{Username: username}
	resp, err := grpcWrapper.CreateUser(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка создания пользователя: %v\n", err)
		return
	}
	fmt.Printf("Пользователь успешно создан с ID: %d\n", resp.UserId)
}

// handleGetUser получает пользователя по ID
func handleGetUser(ctx context.Context, userIDStr string, grpcWrapper *client.APIServiceClientWrapper) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}
	req := &v1.GetUserRequest{UserId: int64(userID)}
	resp, err := grpcWrapper.GetUser(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка получения пользователя: %v\n", err)
		return
	}
	fmt.Printf("Пользователь: ID=%d, Имя=%s\n", resp.UserId, resp.Username)
}

// handleGetAllUsers получает список всех пользователей
func handleGetAllUsers(ctx context.Context, grpcWrapper *client.APIServiceClientWrapper) {
	resp, err := grpcWrapper.GetAllUsers(ctx)
	if err != nil {
		fmt.Printf("Ошибка получения списка пользователей: %v\n", err)
		return
	}
	for _, user := range resp.Users {
		fmt.Printf("ID: %d, Имя: %s\n", user.UserId, user.Username)
	}
}

// handleUpdateUser обновляет пользователя
func handleUpdateUser(ctx context.Context, userIDStr, username string, grpcWrapper *client.APIServiceClientWrapper) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}
	req := &v1.UpdateUserRequest{UserId: int64(userID), Username: username}
	err = grpcWrapper.UpdateUser(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка обновления пользователя: %v\n", err)
	} else {
		fmt.Println("Пользователь успешно обновлен.")
	}
}

// handleDeleteUser удаляет пользователя
func handleDeleteUser(ctx context.Context, userIDStr string, grpcWrapper *client.APIServiceClientWrapper) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}
	req := &v1.DeleteUserRequest{UserId: int64(userID)}
	err = grpcWrapper.DeleteUser(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка удаления пользователя: %v\n", err)
	} else {
		fmt.Println("Пользователь успешно удален.")
	}
}
