package handler

import (
	"booking/config"
	"booking/storage/bookingpg"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
	"time"
)

type Booking struct {
	Time  time.Time `json:"time"`
	Price float64   `json:"price"`
}

type BookingHnd struct {
	client *http.Client
	db     *bookingpg.Storage
	cfg    config.Config
}

func NewBookingHnd(client *http.Client, db *bookingpg.Storage, cfg config.Config) *BookingHnd {
	return &BookingHnd{client: client, db: db, cfg: cfg}
}

func (b *BookingHnd) AddBooking(c *gin.Context) {
	ctx := c.Request.Context()
	logging.Debug("BookingHnd.AddBooking()")

	spanCtx, span := tracing.NewSpan(ctx, "Handler.AddBooking")
	defer span.End() // Обязательно, иначе будет висеть в памяти

	// Запрашиваем цену для бронирования у сервиса расчёта цен
	resp, err := b.sendRequest(ctx, http.MethodGet, fmt.Sprintf("%s/booking-price", b.cfg.CalcPricesAddr), nil)
	if err != nil {
		// Добавляем в трейс ошибку отправки запроса
		span.AddError("error sending request to price calc service", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	type PriceResponse struct {
		Price float64 `json:"price"`
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Добавляем в трейс ошибку чтения
		span.AddError("error reading response body", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	// Десериализация JSON
	var priceResponse PriceResponse
	err = json.Unmarshal(body, &priceResponse)
	if err != nil {
		// Добавляем в трейс ошибку десериализации
		span.AddError("error unmarshalling response body", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	}

	// Добавляем бронирование в базу данных
	_, err = b.db.AddBooking(spanCtx, float64(priceResponse.Price), time.Now())
	if err != nil {
		// Добавляем в трейс ошибку добавления в базу данных
		span.AddError("db.AddBooking returns error", err)
		logging.ErrorErr("db.AddBooking returns error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db.AddBooking returns error"})
	} else {
		// Если успешно, добавляем новое событие "Booking added" в трейс
		span.AddEvent("Booking added") // Новое событие в этот span
		c.JSON(http.StatusOK, gin.H{"message": "Booking added successfully"})
	}
}

func (b *BookingHnd) GetBooking(c *gin.Context) {
	logging.Debug("GetBooking()")
	spanCtx, span := tracing.NewSpan(c.Request.Context(), "GetBooking")
	defer spanCtx.Done()

	booking := Booking{
		Time:  time.Now().Add(-24 * time.Hour),
		Price: 200.0,
	}

	span.AddEvent("Booking retrieved")
	c.JSON(http.StatusOK, booking)
}

func (b *BookingHnd) sendRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	// Обязательно создаем запрос с контекстом в котором уже есть трейс
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Используем созданный через otelhttp клиент для отправки запроса, он добавит трейс в заголовки
	return b.client.Do(req)
}
