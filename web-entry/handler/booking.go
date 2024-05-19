package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"math/rand"
	"net/http"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
	"time"
	"web-entry/config"
)

type bookingSchema struct {
	ID   string `json:"id"`
	Time string `json:"time"`
}

type BookingHnd struct {
	client *http.Client
	cfg    config.Config
}

func NewBookingHnd(client *http.Client, cfg config.Config) *BookingHnd {
	return &BookingHnd{client: client, cfg: cfg}
}

func (b *BookingHnd) sendRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	// Создаем запрос с контекстом в котором уже есть трейс, он будет обработан otelhttp transport'ом
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return b.client.Do(req)
}

func (b *BookingHnd) AddBooking(c *gin.Context) {
	logging.Info("AddBooking()",
		slog.String("method", c.Request.Method),
		slog.String("url", c.Request.URL.String()))

	// Создаем контекст с трейсом
	ctx, span := tracing.NewSpan(c.Request.Context(), "Handler.AddBooking")
	defer span.End()

	// Пример лога с traceID
	tracing.TraceLogger(ctx).Debug("AddBooking() with traceID")

	// Новое событие в трейс
	span.AddEvent("Starting new booking")

	var newBooking bookingSchema
	if err := c.BindJSON(&newBooking); err != nil {
		// Записываем ошибку в трейс
		span.AddError("error in BindJSON", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestBody, err := json.Marshal(newBooking)
	if err != nil {
		span.AddError("marshaling booking JSON", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Отправляем запрос на сервис booking
	resp, err := b.sendRequest(ctx, "POST", b.cfg.BookingAddr+"/add-booking", requestBody)
	if err != nil {
		// Записываем ошибку в трейс
		span.AddError("sending request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	// Проверяем код ответа от сервиса booking
	if resp.StatusCode != http.StatusOK {
		// Записываем ошибку в трейс
		span.AddError("unexpected status code from booking service", nil)
		// Пример лога с traceID
		tracing.TraceLogger(ctx).
			Warn("unexpected status code from booking service", slog.Int("status", resp.StatusCode))

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Возвращаем ответ от сервиса booking
	c.JSON(http.StatusOK, gin.H{"message": "booking added successfully"})
}

func (b *BookingHnd) GetBookingByID(c *gin.Context) {
	// Создание контекста с трассировочным спаном
	_, span := tracing.NewSpan(c.Request.Context(), "Handler.GetBookingByID")
	defer span.End()

	booking := bookingSchema{ID: "123", Time: time.Now().Format(time.DateTime)}
	found := rand.Intn(5) == 0
	if !found {
		span.AddEvent("booking not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}
	c.JSON(http.StatusOK, booking)
}
