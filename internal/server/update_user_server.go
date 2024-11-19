package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser обновляет данные пользователя
func (s *APIServiceServer) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	if err := controller.UpdateUser(ctx, s.userService, req.UserId, req.Username); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обновления пользователя: %v", err)
	}

	return &v1.UpdateUserResponse{Message: "Пользователь успешно обновлен"}, nil
}
