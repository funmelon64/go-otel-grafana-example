package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
	"price-calcs/config"
	"price-calcs/handler"
	"price-calcs/storage/pricespg"
)

func main() {
	// Load environment variables
	cfg := config.LoadConfig()

	// Initialize logging
	logging.InitLogging(cfg.LoggingCfg)

	// Initialize Jaeger tracer
	if err := tracing.InitTracer(cfg.TracingCfg, "price-calcs"); err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	// Set up Gin router
	router := gin.Default()
	// Будет принимать из запроса или создавать новый трейс при каждом запросе
	tracing.AddOtelMiddleware(router, "price-calcs")

	bookingStorage, err := pricespg.NewStorage(cfg.PgAddr, cfg.PgDb, cfg.PgUser, cfg.PgPass)
	if err != nil {
		log.Panicf("fail to create storage: %v", err)
	}

	priceHandler := handler.NewPricesHnd(bookingStorage)

	// Routes
	router.GET("/booking-price", func(c *gin.Context) { priceHandler.GetBookingPrice(c) })

	// Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.HttpPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
