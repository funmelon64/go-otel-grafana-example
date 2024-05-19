package config

import (
	"github.com/caarlos0/env/v11"
	"log"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
)

type Config struct {
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8080"`
	BookingAddr string `env:"BOOKING_ADDR,required"`
	LoggingCfg  logging.Config
	TracingCfg  tracing.Config
}

func MustLoadConfig() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panicf("failed to load env config: %v", err)
	}
	log.Println("LoggingCfg: ", cfg.LoggingCfg)
	return cfg
}
