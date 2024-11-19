package view

import (
	v1 "TODO/internal/api/v1"
	"TODO/internal/client"
	"TODO/internal/pool"
	"context"
	"fmt"
	"strconv"
	"strings"
)

// handleTaskCommands обрабатывает команды, связанные с задачами
func handleTaskCommands(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	if len(args) < 1 {
		printTaskUsage()
		return
	}

	switch args[0] {
	case "create-task":
		handleCreateTaskCommand(ctx, args, grpcWrapper, workerPool)
	case "get-task":
		handleGetTaskCommand(ctx, args, grpcWrapper, workerPool)
	case "get-tasks":
		handleGetAllTasksCommand(ctx, args, grpcWrapper, workerPool)
	case "update-task":
		handleUpdateTaskCommand(ctx, args, grpcWrapper, workerPool)
	case "delete-task":
		handleDeleteTaskCommand(ctx, args, grpcWrapper, workerPool)
	default:
		fmt.Printf("Неизвестная команда для задач: %s\n", args[0])
	}
}

func printTaskUsage() {
	fmt.Println("Недостаточно аргументов. Доступные команды: create-task, get-task, get-tasks, update-task, delete-task")
}

func handleCreateTaskCommand(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	if len(args) != 4 {
		fmt.Println("Использование: create-task [userID] [title] [note]")
		return
	}
	workerPool.SubmitTask(func() {
		handleCreateTask(ctx, args[1], args[2], args[3], grpcWrapper)
	})
}

func handleGetTaskCommand(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	if len(args) != 2 {
		fmt.Println("Использование: get-task [taskID]")
		return
	}
	workerPool.SubmitTask(func() {
		handleGetTask(ctx, args[1], grpcWrapper)
	})
}

func handleGetAllTasksCommand(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	if len(args) != 1 {
		fmt.Println("Использование: get-tasks")
		return
	}
	workerPool.SubmitTask(func() {
		handleGetAllTasks(ctx, grpcWrapper)
	})
}

func handleUpdateTaskCommand(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	if len(args) != 5 {
		fmt.Println("Использование: update-task [taskID] [title] [note] [done]")
		return
	}
	workerPool.SubmitTask(func() {
		handleUpdateTask(ctx, args[1], args[2], args[3], args[4], grpcWrapper)
	})
}

func handleDeleteTaskCommand(ctx context.Context, args []string, grpcWrapper *client.APIServiceClientWrapper, workerPool *pool.WorkerPool) {
	if len(args) != 2 {
		fmt.Println("Использование: delete-task [taskID]")
		return
	}
	workerPool.SubmitTask(func() {
		handleDeleteTask(ctx, args[1], grpcWrapper)
	})
}

// handleCreateTask создает новую задачу
func handleCreateTask(ctx context.Context, userIDStr, title, note string, grpcWrapper *client.APIServiceClientWrapper) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID пользователя: %v\n", err)
		return
	}

	req := &v1.CreateTaskRequest{
		UserId: userID,
		Title:  title,
		Note:   note,
	}
	resp, err := grpcWrapper.CreateTask(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка создания задачи: %v\n", err)
		return
	}
	fmt.Printf("Задача успешно создана с ID: %d\n", resp.TaskId)
}

// handleGetTask получает задачу по ID
func handleGetTask(ctx context.Context, taskIDStr string, grpcWrapper *client.APIServiceClientWrapper) {
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID задачи: %v\n", err)
		return
	}

	req := &v1.GetTaskRequest{TaskId: taskID}
	resp, err := grpcWrapper.GetTask(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка получения задачи: %v\n", err)
		return
	}
	fmt.Printf("Задача: ID=%d, UserID=%d, Title=%s, Note=%s, Done=%t, CreatedAt=%s\n",
		resp.TaskId, resp.UserId, resp.Title, resp.Note, resp.Done, resp.CreatedAt)
}

// handleGetAllTasks получает список всех задач
func handleGetAllTasks(ctx context.Context, grpcWrapper *client.APIServiceClientWrapper) {
	resp, err := grpcWrapper.GetAllTasks(ctx)
	if err != nil {
		fmt.Printf("Ошибка получения списка задач: %v\n", err)
		return
	}

	fmt.Printf("%-10s %-10s %-20s %-30s %-5s %-25s\n", "TaskID", "UserID", "Title", "Note", "Done", "CreatedAt")
	fmt.Println(strings.Repeat("-", 85))
	for _, task := range resp.Tasks {
		fmt.Printf("%-10d %-10d %-20s %-30s %-5t %-25s\n",
			task.TaskId, task.UserId, task.Title, task.Note, task.Done, task.CreatedAt)
	}
}

// handleUpdateTask обновляет задачу
func handleUpdateTask(ctx context.Context, taskIDStr, title, note, doneStr string, grpcWrapper *client.APIServiceClientWrapper) {
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID задачи: %v\n", err)
		return
	}

	done, err := strconv.ParseBool(doneStr)
	if err != nil {
		fmt.Printf("Ошибка преобразования статуса задачи (done): %v\n", err)
		return
	}

	req := &v1.UpdateTaskRequest{
		TaskId: taskID,
		Title:  title,
		Note:   note,
		Done:   done,
	}

	resp, err := grpcWrapper.UpdateTask(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка обновления задачи: %v\n", err)
		return
	}

	fmt.Println(resp.Message)
}

// handleDeleteTask удаляет задачу
func handleDeleteTask(ctx context.Context, taskIDStr string, grpcWrapper *client.APIServiceClientWrapper) {
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования ID задачи: %v\n", err)
		return
	}

	req := &v1.DeleteTaskRequest{TaskId: taskID}
	err = grpcWrapper.DeleteTask(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка удаления задачи: %v\n", err)
		return
	}
	fmt.Println("Задача успешно удалена.")
}
