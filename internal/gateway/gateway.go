package gateway

import (
	v1 "TODO/internal/api/v1"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

// RunGateway запускает HTTP-gateway, который работает как прокси для gRPC сервера.
func RunGateway(ctx context.Context, grpcEndpoint, httpEndpoint string) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := registerServices(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatalf("Не удалось зарегистрировать сервисы HTTP-gateway: %v", err)
	}

	handler := corsMiddleware(mux)

	log.Printf("HTTP Gateway запущен на %s, проксирует к gRPC на %s", httpEndpoint, grpcEndpoint)
	return http.ListenAndServe(httpEndpoint, handler)
}

// registerServices регистрирует единый APIService для HTTP-Gateway.
func registerServices(ctx context.Context, mux *runtime.ServeMux, grpcEndpoint string, opts []grpc.DialOption) error {

	if err := v1.RegisterAPIServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		return err
	}
	return nil
}

// corsMiddleware добавляет поддержку CORS, чтобы позволить запросы с других доменов.
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
