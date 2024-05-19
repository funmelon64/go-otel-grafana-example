package main

import (
	"booking/config"
	"booking/handler"
	"booking/storage/bookingpg"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
)

func main() {
	cfg := config.LoadConfig()

	logging.InitLogging(cfg.LoggingCfg)

	// Initialize Jaeger tracer
	if err := tracing.InitTracer(cfg.TracingCfg, "bookings"); err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	// Set up Gin router
	router := gin.Default()
	// Будет принимать из запроса или создавать новый трейс при каждом запросек
	tracing.AddOtelMiddleware(router, "bookings")

	bookingStorage, err := bookingpg.NewStorage(cfg.PgAddr, cfg.PgDb, cfg.PgUser, cfg.PgPass)
	if err != nil {
		log.Panicf("fail to create storage: %v", err)
	}

	client := tracing.NewOtelHttpClient()

	bookingHandler := handler.NewBookingHnd(client, bookingStorage, cfg)

	// Routes
	router.POST("/add-booking", func(c *gin.Context) { bookingHandler.AddBooking(c) })
	router.GET("/get-booking", func(c *gin.Context) { bookingHandler.GetBooking(c) })

	// Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.HttpPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
