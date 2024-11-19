package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser создает нового пользователя
func (s *APIServiceServer) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	userID, err := controller.CreateUser(ctx, s.userService, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания пользователя: %v", err)
	}

	return &v1.CreateUserResponse{
		UserId:  userID,
		Message: "Пользователь успешно создан",
	}, nil
}
