package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "TODO/internal/api/v1"
	"TODO/internal/controller"
)

// UpdateTask обновляет данные задачи
func (s *APIServiceServer) UpdateTask(ctx context.Context, req *v1.UpdateTaskRequest) (*v1.UpdateTaskResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	err := controller.UpdateTask(ctx, s.taskService, req.TaskId, req.Title, req.Note, req.Done)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обновления задачи: %v", err)
	}

	return &v1.UpdateTaskResponse{
		Message: "Задача успешно обновлена",
	}, nil
}
