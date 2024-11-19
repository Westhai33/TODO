package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

// GetAllTasks возвращает все задачи
func (s *APIServiceServer) GetAllTasks(ctx context.Context, _ *emptypb.Empty) (*v1.GetAllTasksResponse, error) {

	tasks, err := controller.GetAllTasks(ctx, s.taskService)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения задач: %v", err)
	}

	response := &v1.GetAllTasksResponse{}
	for _, task := range tasks {
		response.Tasks = append(response.Tasks, &v1.Task{
			TaskId:    task.ID,
			UserId:    task.UserID,
			Title:     task.Title,
			Note:      task.Note,
			CreatedAt: task.CreatedAt.Format(time.RFC3339),
			UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
		})
	}
	return response, nil
}
