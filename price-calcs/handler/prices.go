package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"math/rand"
	"net/http"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
	"price-calcs/storage/pricespg"
)

type PricesHnd struct {
	db *pricespg.Storage
}

func NewPricesHnd(db *pricespg.Storage) *PricesHnd {
	return &PricesHnd{db: db}
}

func (b *PricesHnd) GetBookingPrice(c *gin.Context) {
	ctx := c.Request.Context()
	logging.Debug("GetBookingPrice()")

	// Создаём новый span на основе текущего контекста
	spanCtx, span := tracing.NewSpan(ctx, "Booking Price Calculation")
	defer span.End() // Обязательно, иначе будет висеть в памяти

	// Рандомный водятел
	driverId := fmt.Sprint(rand.Intn(pricespg.DRIVERS_COUNT))

	// Получем цену водителя из базы данных
	price, err := b.db.GetDriverPrice(spanCtx, driverId)
	if err != nil {
		// Добавляем информацию об ошибке в span
		span.AddError("db.GetDriverPrice returns error", err)
		logging.ErrorErr("db.GetDriverPrice returns error", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "db.GetDriverPrice returns error"})
	}

	// Получаем скидки водителя из базы данных
	discounts, err := b.db.GetDriverDiscounts(spanCtx, driverId)
	if err != nil {
		// Добавляем информацию об ошибке в span
		span.AddError("db.GetDriverDiscounts returns error", err)
		logging.ErrorErr("db.GetDriverDiscounts returns error", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "db.GetDriverDiscounts returns error"})
	}

	// Вычисляем общую цену
	totalPrice := price
	for _, discount := range discounts {
		totalPrice -= float64(discount)
	}

	// Добавляем информацию что цена посчитана в span и добавляем цену в атрибуты
	span.AddEvent("Price Calculated", slog.Float64("price", price))

	c.JSON(http.StatusOK, gin.H{"price": totalPrice})
}
