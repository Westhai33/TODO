package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitTracer инициализирует трейсинг и возвращает функцию завершения работы
func InitTracer(serviceName, endpoint string) func(context.Context) error {
	exp, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при создании OTLP экспортера: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(serviceName))),
	)

	otel.SetTracerProvider(tp)
	tracer = otel.Tracer(serviceName)

	return tp.Shutdown
}

// GetTracer возвращает Tracer для использования в других частях приложения
func GetTracer() trace.Tracer {
	return tracer
}
