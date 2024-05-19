package tracing

type Config struct {
	TempoAddr string `env:"TEMPO_ADDR,required"`
}
