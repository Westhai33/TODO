package client

import (
	"context"
	"log"

	v1 "TODO/internal/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// APIServiceClientWrapper обертка для APIServiceClient
type APIServiceClientWrapper struct {
	client v1.APIServiceClient
	conn   *grpc.ClientConn
}

// NewAPIServiceClientWrapper создает новый экземпляр обертки APIServiceClientWrapper
func NewAPIServiceClientWrapper(conn *grpc.ClientConn) (*APIServiceClientWrapper, error) {
	if conn == nil {
		log.Println("gRPC подключение не может быть nil")
		return nil, status.Error(codes.InvalidArgument, "gRPC подключение отсутствует")
	}

	client := v1.NewAPIServiceClient(conn)
	return &APIServiceClientWrapper{client: client, conn: conn}, nil
}

// CreateUser проксирует запрос к CreateUser gRPC методу
func (w *APIServiceClientWrapper) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	resp, err := w.client.CreateUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateUser: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetUser проксирует запрос к GetUser gRPC методу
func (w *APIServiceClientWrapper) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	resp, err := w.client.GetUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetUser: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllUsers проксирует запрос к GetAllUsers gRPC методу
func (w *APIServiceClientWrapper) GetAllUsers(ctx context.Context) (*v1.GetAllUsersResponse, error) {
	req := &emptypb.Empty{}
	resp, err := w.client.GetAllUsers(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllUsers: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateUser проксирует запрос к UpdateUser gRPC методу
func (w *APIServiceClientWrapper) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) error {
	resp, err := w.client.UpdateUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова UpdateUser: %v", err)
		return err
	}
	log.Println("Сообщение от сервера:", resp.Message)
	return nil
}

// DeleteUser проксирует запрос к DeleteUser gRPC методу
func (w *APIServiceClientWrapper) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) error {
	_, err := w.client.DeleteUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteUser: %v", err)
		return err
	}
	return nil
}

// CreateTask проксирует запрос к CreateTask gRPC методу
func (w *APIServiceClientWrapper) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error) {
	resp, err := w.client.CreateTask(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateTask: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTask проксирует запрос к GetTask gRPC методу
func (w *APIServiceClientWrapper) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.GetTaskResponse, error) {
	resp, err := w.client.GetTask(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetTask: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllTasks проксирует запрос к GetAllTasks gRPC методу
func (w *APIServiceClientWrapper) GetAllTasks(ctx context.Context) (*v1.GetAllTasksResponse, error) {
	req := &emptypb.Empty{} // Пустой запрос для получения всех задач
	resp, err := w.client.GetAllTasks(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllTasks: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateTask проксирует запрос к UpdateTask gRPC методу
func (w *APIServiceClientWrapper) UpdateTask(ctx context.Context, req *v1.UpdateTaskRequest) (*v1.UpdateTaskResponse, error) {
	resp, err := w.client.UpdateTask(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова UpdateTask: %v", err)
		return nil, err
	}
	return resp, nil
}

// DeleteTask проксирует запрос к DeleteTask gRPC методу
func (w *APIServiceClientWrapper) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) error {
	_, err := w.client.DeleteTask(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteTask: %v", err)
		return err
	}
	return nil
}

// Close закрывает соединение gRPC
func (w *APIServiceClientWrapper) Close() error {
	if w.conn != nil {
		if err := w.conn.Close(); err != nil {
			log.Printf("Ошибка при закрытии gRPC соединения: %v", err)
			return err
		}
		log.Println("gRPC соединение закрыто успешно")
	} else {
		log.Println("gRPC соединение уже закрыто или не было установлено")
	}
	return nil
}
