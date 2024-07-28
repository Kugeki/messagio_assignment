package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"time"
)

var (
	envDev  = "development"
	envProd = "production"
)

type Config struct {
	Environment     *Environment  `yaml:"environment" env-default:"development"`
	LogLevel        slog.Level    `yaml:"log_level" env-default:"DEBUG"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
	HTTPServer      HTTPServer    `yaml:"http_server" env-prefix:"HTTP_SERVER_"`
	Postgres        Postgres      `yaml:"postgres" env-prefix:"POSTGRES_"`
}

type HTTPServer struct {
	Addr           string `yaml:"addr" env:"ADDR" env-required:"true"`
	MaxHeaderBytes int    `yaml:"max_header_bytes" env:"MAX_HEADER_BYTES"`
	Timeouts       struct {
		Read       time.Duration `yaml:"read" env:"READ_TIMEOUT"`
		ReadHeader time.Duration `yaml:"read_header" env:"READ_HEADER_TIMEOUT"`
		Write      time.Duration `yaml:"write" env:"WRITE_TIMEOUT"`
		Idle       time.Duration `yaml:"idle" env:"IDLE_TIMEOUT"`
	} `yaml:"timeouts"`
}

type Postgres struct {
	ConnectionURL string `env:"CONNECTION_URL" env-required:"true"`
	Migrate       bool   `yaml:"migrate" env:"MIGRATE" env-required:"true"`
}

func ReadConfig(path string) (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

type Environment string

func (e *Environment) UnmarshalText(data []byte) error {
	str := string(data)
	switch str {
	case envDev, envProd:
		*e = Environment(str)
		return nil
	default:
		return errors.New("unknown environment")
	}
}

func (e *Environment) String() string {
	return string(*e)
}

func (e *Environment) IsDev() bool {
	return e.String() == envDev
}

func (e *Environment) IsProd() bool {
	return e.String() == envProd
}
