package config

import (
	"github.com/caarlos0/env/v11"
	"log"
	"otel-jaeger-learn/pkg/logging"
	"otel-jaeger-learn/pkg/tracing"
)

type Config struct {
	HttpPort       string `env:"HTTP_PORT" envDefault:"8080"`
	PgUser         string `env:"PG_USER" envDefault:"postgres"`
	PgPass         string `env:"PG_PASS" envDefault:"postgres"`
	PgAddr         string `env:"PG_ADDR" envDefault:"localhost:5432"`
	PgDb           string `env:"PG_DB" envDefault:"postgres"`
	CalcPricesAddr string `env:"CALC_PRICES_ADDR,required"`
	LoggingCfg     logging.Config
	TracingCfg     tracing.Config
}

func LoadConfig() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panicf("failed to load env config: %v", err)
	}
	return cfg
}
