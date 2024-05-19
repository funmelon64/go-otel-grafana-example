package tracing

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"net/http"
)

func AddOtelMiddleware(r *gin.Engine, serviceName string) {
	// Middleware который будет создавать новый или брать из заголовков трейс при каждом запросе
	r.Use(otelgin.Middleware(serviceName))
}

func NewOtelHttpClient() *http.Client {
	// Клиент который будет передавать трейс в запросе (нужно глобально установить propagator)
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func InitTracer(cfg Config, serviceName string) error {
	// Создаём экспортер в Tempo
	exporter, err := newTempoExporter(cfg.TempoAddr)
	if err != nil {
		return err
	}

	// Создаём ресурс для трейсера (имя сервиса и прочие данные отображаемые в Jaeger)
	res, err := resource.New(context.TODO(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return fmt.Errorf("could not set up resource: %v", err)
	}

	// Создание трейсер провайдера, с Jaeger, который принимает все span'ы (AlwaysSample)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// Установка Propagator'а для корректного распространения трейса через запросы в другие сервисы
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func newTempoExporter(addr string) (sdktrace.SpanExporter, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(addr), // otlp http port
		otlptracehttp.WithInsecure(),     // HTTP instead of HTTPS
	)
	return otlptrace.New(context.TODO(), client)
}
