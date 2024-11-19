package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateTask создает новую задачу
func (s *APIServiceServer) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	taskID, err := controller.CreateTask(ctx, s.taskService, s.userService, req.UserId, req.Title, req.Note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания задачи: %v", err)
	}

	return &v1.CreateTaskResponse{
		TaskId:  taskID,
		Message: "Задача успешно создана",
	}, nil
}
