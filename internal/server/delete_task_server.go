package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "TODO/internal/api/v1"
	"TODO/internal/controller"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteTask удаляет задачу
func (s *APIServiceServer) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) (*emptypb.Empty, error) {

	if err := req.Validate(); err != nil {
		log.Printf("Валидация DeleteTaskRequest не прошла: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	err := controller.DeleteTask(ctx, s.taskService, req.TaskId) // Передаем только taskID
	if err != nil {
		log.Printf("Ошибка удаления задачи: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка удаления задачи: %v", err)
	}

	return &emptypb.Empty{}, nil
}
