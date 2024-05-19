package logging

type Config struct {
	Level   string `env:"LOG_LEVEL" envDefault:"debug"` // debug, warn, info, error
	FileOut string `env:"LOG_FILE,unset"`               // if omitted - log to stdout
}
