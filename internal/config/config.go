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
	envTest = "testing"
)

type Config struct {
	Environment     *Environment  `yaml:"environment" env-default:"development"`
	LogLevel        slog.Level    `yaml:"log_level" env-default:"DEBUG"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
	HTTPServer      HTTPServer    `yaml:"http_server" env-prefix:"HTTP_SERVER_"`
	Postgres        Postgres      `yaml:"postgres" env-prefix:"POSTGRES_"`
	Kafka           Kafka         `yaml:"kafka" env-prefix:"KAFKA_"`
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
	ConnectionURL string `env:"CONNECTION_URL" env-required:"true" env-description:"required"`
	Migrate       bool   `yaml:"migrate" env:"MIGRATE" env-required:"true" env-description:"required"`
}

type Kafka struct {
	ClientID string   `yaml:"client_id" env:"CLIENT_ID" env-required:"true" env-description:"required"`
	Brokers  []string `yaml:"brokers" env:"BROKERS" env-required:"true" env-description:"required"`

	Producers struct {
		Messages KafkaProducer `yaml:"messages" env-prefix:"MESSAGES_"`
	} `yaml:"producers" env-prefix:"PRODUCER_"`
	Consumers struct {
		ProcessedMessages KafkaConsumer `yaml:"processed_messages" env-prefix:"PROCESSED_MESSAGES"`
	} `yaml:"consumers" env-prefix:"CONSUMER_"`
}

// KafkaProducer : some from sarama.NewConfig and sarama.Config
type KafkaProducer struct {
	Topic string `yaml:"topic" env:"TOPIC" env-required:"true" env-description:"required"`

	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"TIMEOUT"`
	Retry   struct {
		// The total number of times to retry sending a message (default 3).
		// Similar to the `message.send.max.retries` setting of the JVM producer.
		Max int `yaml:"max" env-default:"3" env:"MAX"`
		// How long to wait for the cluster to settle between retries
		// (default 100ms). Similar to the `retry.backoff.ms` setting of the
		// JVM producer.
		Backoff time.Duration `yaml:"backoff" env-default:"100ms" env:"BACKOFF"`
	} `yaml:"retry" env-prefix:"RETRY_"`

	Flush struct {
		// The best-effort number of bytes needed to trigger a flush. Use the
		// global sarama.MaxRequestSize to set a hard upper limit.
		Bytes int `yaml:"bytes" env:"BYTES"`
		// The best-effort number of messages needed to trigger a flush. Use
		// `MaxMessages` to set a hard upper limit.
		Messages int `yaml:"messages" env:"MESSAGES"`
		// The best-effort frequency of flushes. Equivalent to
		// `queue.buffering.max.ms` setting of JVM producer.
		Frequency time.Duration `yaml:"frequency" env:"FREQUENCY"`
		// The maximum number of messages the producer will send in a single
		// broker request. Defaults to 0 for unlimited. Similar to
		// `queue.buffering.max.messages` in the JVM producer.
		MaxMessages int `yaml:"max_messages" env:"MAX_MESSAGES"`
	} `yaml:"flush" env-prefix:"FLUSH_"`
}

type KafkaConsumer struct {
	Group  string   `yaml:"group" env:"GROUP" env-required:"true" env-description:"required"`
	Topics []string `yaml:"topics" env:"TOPICS" env-required:"true" env-description:"required"`

	Retry struct {
		Backoff time.Duration `yaml:"backoff" env:"BACKOFF" env-default:"2s"`
	} `yaml:"retry" env-prefix:"RETRY_"`

	MaxWaitTime time.Duration `yaml:"max_wait_time" env:"MAX_WAIT_TIME" env-default:"500ms"`
	Fetch       struct {
		Min     int32 `yaml:"min" env:"MIN" env-default:"1"`
		Default int32 `yaml:"default" env:"DEFAULT" env-default:"1048576"`
		Max     int32 `yaml:"max" env:"MAX"`
	} `yaml:"fetch" env-prefix:"FETCH_"`
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
	case envDev, envProd, envTest:
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
	return e.String() == envDev || e.String() == envTest
}

func (e *Environment) IsProd() bool {
	return e.String() == envProd
}
