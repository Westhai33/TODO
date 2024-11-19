package view

import (
	"TODO/internal/controller"
	"TODO/internal/service"
	"context"
	"fmt"
	"strconv"
	"strings"
)

// handleTaskCommands обрабатывает команды, связанные с задачами
func handleTaskCommands(ctx context.Context, args []string, taskService *service.TaskService) {
	if len(args) < 1 {
		printTaskUsage()
		return
	}

	switch args[0] {
	case "create-task":
		handleCreateTaskCommand(ctx, args, taskService)
	case "get-task":
		handleGetTaskCommand(ctx, args, taskService)
	case "get-tasks":
		handleGetAllTasksCommand(ctx, taskService)
	case "update-task":
		handleUpdateTaskCommand(ctx, args, taskService)
	case "delete-task":
		handleDeleteTaskCommand(ctx, args, taskService)
	default:
		fmt.Printf("Неизвестная команда для задач: %s\n", args[0])
	}
}

func printTaskUsage() {
	fmt.Println("Недостаточно аргументов. Доступные команды: create-task, get-task, get-tasks, update-task, delete-task")
}

func handleCreateTaskCommand(ctx context.Context, args []string, taskService *service.TaskService) {
	if len(args) != 4 {
		fmt.Println("Использование: create-task [userID] [title] [note]")
		return
	}

	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}

	taskID, err := controller.CreateTask(ctx, taskService, userID, args[2], args[3])
	if err != nil {
		fmt.Printf("Ошибка создания задачи: %v\n", err)
		return
	}
	fmt.Printf("Задача успешно создана с ID: %d\n", taskID)
}

func handleGetTaskCommand(ctx context.Context, args []string, taskService *service.TaskService) {
	if len(args) != 2 {
		fmt.Println("Использование: get-task [taskID]")
		return
	}

	taskID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID задачи: %v\n", err)
		return
	}

	task, err := controller.GetTask(ctx, taskService, taskID)
	if err != nil {
		fmt.Printf("Ошибка получения задачи: %v\n", err)
		return
	}
	fmt.Printf("Задача: %+v\n", task)
}

func handleGetAllTasksCommand(ctx context.Context, taskService *service.TaskService) {
	tasks, err := controller.GetAllTasks(ctx, taskService)
	if err != nil {
		fmt.Printf("Ошибка получения списка задач: %v\n", err)
		return
	}

	fmt.Printf("%-10s %-10s %-20s %-30s %-5s\n", "TaskID", "UserID", "Title", "Note", "Done")
	fmt.Println(strings.Repeat("-", 85))
	for _, task := range tasks {
		fmt.Printf("%-10d %-10d %-20s %-30s %-5t\n", task.ID, task.UserID, task.Title, task.Note, task.Done)
	}
}

func handleUpdateTaskCommand(ctx context.Context, args []string, taskService *service.TaskService) {
	if len(args) != 5 {
		fmt.Println("Использование: update-task [taskID] [title] [note] [done]")
		return
	}

	taskID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID задачи: %v\n", err)
		return
	}

	done, err := strconv.ParseBool(args[4])
	if err != nil {
		fmt.Printf("Ошибка преобразования статуса задачи (done): %v\n", err)
		return
	}

	err = controller.UpdateTask(ctx, taskService, taskID, args[2], args[3], done)
	if err != nil {
		fmt.Printf("Ошибка обновления задачи: %v\n", err)
		return
	}
	fmt.Println("Задача успешно обновлена.")
}

func handleDeleteTaskCommand(ctx context.Context, args []string, taskService *service.TaskService) {
	if len(args) != 2 {
		fmt.Println("Использование: delete-task [taskID]")
		return
	}

	taskID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID задачи: %v\n", err)
		return
	}

	err = controller.DeleteTask(ctx, taskService, taskID)
	if err != nil {
		fmt.Printf("Ошибка удаления задачи: %v\n", err)
		return
	}
	fmt.Println("Задача успешно удалена.")
}
