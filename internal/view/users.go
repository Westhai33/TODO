package view

import (
	"TODO/internal/controller"
	"TODO/internal/pool"
	"TODO/internal/service"
	"context"
	"fmt"
	"strconv"
)

// handleUserCommands обрабатывает команды, связанные с пользователями
func handleUserCommands(ctx context.Context, args []string, userService *service.UserService, workerPool *pool.WorkerPool) {
	switch args[0] {
	case "create-user":
		if len(args) != 2 {
			fmt.Println("Использование: create-user [username]")
			return
		}
		workerPool.SubmitTask(func() {
			handleCreateUser(ctx, args[1], userService)
		})
	case "get-user":
		if len(args) != 2 {
			fmt.Println("Использование: get-user [userID]")
			return
		}
		workerPool.SubmitTask(func() {
			handleGetUser(ctx, args[1], userService)
		})
	case "get-users":
		workerPool.SubmitTask(func() {
			handleGetAllUsers(ctx, userService)
		})
	case "update-user":
		if len(args) != 3 {
			fmt.Println("Использование: update-user [userID] [username]")
			return
		}
		workerPool.SubmitTask(func() {
			handleUpdateUser(ctx, args[1], args[2], userService)
		})
	case "delete-user":
		if len(args) != 2 {
			fmt.Println("Использование: delete-user [userID]")
			return
		}
		workerPool.SubmitTask(func() {
			handleDeleteUser(ctx, args[1], userService)
		})
	default:
		fmt.Println("Неизвестная команда для пользователя.")
	}
}

// handleCreateUser создает нового пользователя
func handleCreateUser(ctx context.Context, username string, userService *service.UserService) {
	userID, err := controller.CreateUser(ctx, userService, username)
	if err != nil {
		fmt.Printf("Ошибка создания пользователя: %v\n", err)
		return
	}
	fmt.Printf("Пользователь успешно создан с ID: %d\n", userID)
}

// handleGetUser получает пользователя по ID
func handleGetUser(ctx context.Context, userIDStr string, userService *service.UserService) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}

	user, err := controller.GetUserByID(ctx, userService, userID)
	if err != nil {
		fmt.Printf("Ошибка получения пользователя: %v\n", err)
		return
	}
	fmt.Printf("Пользователь: ID=%d, Имя=%s\n", user.ID, user.Username)
}

// handleGetAllUsers получает список всех пользователей
func handleGetAllUsers(ctx context.Context, userService *service.UserService) {
	users, err := controller.GetAllUsers(ctx, userService)
	if err != nil {
		fmt.Printf("Ошибка получения списка пользователей: %v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("ID: %d, Имя: %s\n", user.ID, user.Username)
	}
}

// handleUpdateUser обновляет пользователя
func handleUpdateUser(ctx context.Context, userIDStr, username string, userService *service.UserService) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}

	err = controller.UpdateUser(ctx, userService, userID, username)
	if err != nil {
		fmt.Printf("Ошибка обновления пользователя: %v\n", err)
	} else {
		fmt.Println("Пользователь успешно обновлен.")
	}
}

// handleDeleteUser удаляет пользователя
func handleDeleteUser(ctx context.Context, userIDStr string, userService *service.UserService) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}

	err = controller.DeleteUser(ctx, userService, userID)
	if err != nil {
		fmt.Printf("Ошибка удаления пользователя: %v\n", err)
	} else {
		fmt.Println("Пользователь успешно удален.")
	}
}
