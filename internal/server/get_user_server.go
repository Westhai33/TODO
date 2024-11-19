package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser возвращает информацию о пользователе по ID
func (s *APIServiceServer) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	user, err := controller.GetUserByID(ctx, s.userService, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "пользователь не найден: %v", err)
	}

	return &v1.GetUserResponse{
		UserId:    user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}
