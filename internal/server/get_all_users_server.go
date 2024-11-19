package server

import (
	"TODO/internal/api/v1"
	"TODO/internal/controller"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetAllUsers возвращает всех пользователей
func (s *APIServiceServer) GetAllUsers(ctx context.Context, _ *emptypb.Empty) (*v1.GetAllUsersResponse, error) {
	users, err := controller.GetAllUsers(ctx, s.userService)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения пользователей: %v", err)
	}

	response := &v1.GetAllUsersResponse{}
	for _, user := range users {
		response.Users = append(response.Users, &v1.User{
			UserId:    user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.String(),
		})
	}
	return response, nil
}
