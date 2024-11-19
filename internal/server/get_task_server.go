package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// GetTask возвращает задачу по ID
func (s *APIServiceServer) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.GetTaskResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	task, err := controller.GetTask(ctx, s.taskService, req.TaskId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "задача не найдена: %v", err)
	}

	return &v1.GetTaskResponse{
		TaskId:    task.ID,
		Title:     task.Title,
		Note:      task.Note,
		UserId:    task.UserID,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}, nil
}
