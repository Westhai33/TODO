package view

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"TODO/internal/pool"
	"TODO/internal/service"
)

// RunInteractiveMode запускает интерактивный режим управления задачами и пользователями.
func RunInteractiveMode(ctx context.Context, taskService *service.TaskService, userService *service.UserService, workerPool *pool.WorkerPool) {
	fmt.Println("Добро пожаловать в систему управления задачами и пользователями!")
	fmt.Println("Введите 'help' для получения списка команд.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		if args[0] == "exit" {
			fmt.Println("Завершение работы...")
			break
		}

		handleCommand(ctx, args, taskService, userService, workerPool)
	}
}

// handleCommand обрабатывает введенные команды
func handleCommand(ctx context.Context, args []string, taskService *service.TaskService, userService *service.UserService, workerPool *pool.WorkerPool) {
	if len(args) == 0 {
		fmt.Println("Команда не указана. Введите 'help' для получения справки.")
		return
	}

	switch args[0] {
	case "help":
		printHelp()
	case "set-workers":
		handleSetWorkersCommand(ctx, args, workerPool)
	default:
		handleOtherCommands(ctx, args, taskService, userService, workerPool)
	}
}

// handleSetWorkersCommand обрабатывает команду изменения количества воркеров
func handleSetWorkersCommand(ctx context.Context, args []string, workerPool *pool.WorkerPool) {
	if len(args) != 2 {
		fmt.Println("Ошибка: Неверное количество аргументов. Используйте: set-workers [количество]")
		return
	}

	newWorkerCount, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Ошибка: Неверное значение для количества воркеров.")
		return
	}

	workerPool.SetWorkerCount(newWorkerCount)
	fmt.Printf("Количество воркеров обновлено до: %d\n", newWorkerCount)
}

// handleOtherCommands обрабатывает команды для пользователей и задач
func handleOtherCommands(ctx context.Context, args []string, taskService *service.TaskService, userService *service.UserService, workerPool *pool.WorkerPool) {
	switch args[0] {
	case "create-user", "get-user", "get-users", "update-user", "delete-user":
		workerPool.SubmitTask(func() {
			handleUserCommands(ctx, args, userService, workerPool)
		})
	case "create-task", "get-task", "get-tasks", "update-task", "delete-task":
		workerPool.SubmitTask(func() {
			handleTaskCommands(ctx, args, taskService, userService, workerPool)
		})
	default:
		fmt.Println("Неизвестная команда. Введите 'help' для получения справки.")
	}
}

// printHelp выводит доступные команды
func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("Пользователи:")
	fmt.Println("  create-user [username] - Создать нового пользователя")
	fmt.Println("  get-user [userID] - Получить пользователя по ID")
	fmt.Println("  get-users - Получить список всех пользователей")
	fmt.Println("  update-user [userID] [username] - Обновить данные пользователя")
	fmt.Println("  delete-user [userID] - Удалить пользователя")
	fmt.Println("Задачи:")
	fmt.Println("  create-task [userID] [title] [note] - Создать новую задачу")
	fmt.Println("  get-task [taskID] - Получить задачу по ID")
	fmt.Println("  get-tasks - Получить список всех задач")
	fmt.Println("  update-task [taskID] [title] [note] [done] - Обновить задачу")
	fmt.Println("  delete-task [taskID] - Удалить задачу")
	fmt.Println("Системные команды:")
	fmt.Println("  set-workers [количество] - Изменить количество воркеров")
	fmt.Println("  exit - Выйти из программы")
}
