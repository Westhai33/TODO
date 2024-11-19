package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "TODO/internal/api/v1"
	"TODO/internal/controller"
)

// DeleteUser удаляет пользователя
func (s *APIServiceServer) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {

	if err := req.Validate(); err != nil {
		log.Printf("Валидация DeleteUserRequest не прошла: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	userID := int64(req.UserId)

	if err := controller.DeleteUser(ctx, s.userService, userID); err != nil {
		log.Printf("Ошибка удаления пользователя: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка удаления пользователя: %v", err)
	}

	return &emptypb.Empty{}, nil
}
