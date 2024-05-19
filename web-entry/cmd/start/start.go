package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"otel-jaeger-learn/pkg/tracing"
	"web-entry/config"
	"web-entry/handler"
)

const ServiceName = "web-entry"

func main() {
	cfg := config.MustLoadConfig()

	err := tracing.InitTracer(cfg.TracingCfg, ServiceName)
	if err != nil {
		log.Fatalf("failed to initialize logging: %v", err)
	}

	client := tracing.NewOtelHttpClient()
	bookingHandler := handler.NewBookingHnd(client, cfg)

	router := gin.Default()

	// Будет принимать из запроса или создавать новый трейс при каждом запросе
	tracing.AddOtelMiddleware(router, ServiceName)

	router.POST("/bookings", func(c *gin.Context) { bookingHandler.AddBooking(c) })
	router.GET("/bookings/:id", func(c *gin.Context) { bookingHandler.GetBookingByID(c) })

	err = router.Run(":" + cfg.HTTPPort)
	if err != nil {
		log.Panicf("fail to run router: %v", err)
	}
}
