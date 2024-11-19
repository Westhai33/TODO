package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/service"
)

// APIServiceServer представляет реализацию интерфейса v1.APIServiceServer
type APIServiceServer struct {
	v1.UnimplementedAPIServiceServer // Встраиваем не реализованный сервер
	userService                      *service.UserService
	taskService                      *service.TaskService
}

// NewAPIServiceServer создает новый APIServiceServer
func NewAPIServiceServer(
	userService *service.UserService,
	taskService *service.TaskService,
) *APIServiceServer {
	return &APIServiceServer{
		userService: userService,
		taskService: taskService,
	}
}
