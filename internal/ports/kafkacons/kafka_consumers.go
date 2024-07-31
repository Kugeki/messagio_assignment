package kafkacons

import (
	"errors"
	"github.com/IBM/sarama"
	"log/slog"
	"messagioassignment/internal/config"
	"messagioassignment/internal/logger"
)

type KafkaConsumers struct {
	log *slog.Logger

	procMsgsConsumer *ProcessedMsgConsumer
}

func New(log *slog.Logger, msgUC MessagesUsecase,
	saramaCfg *sarama.Config, kafkaConf config.Kafka) (*KafkaConsumers, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	log = log.With(slog.String("component", "ports/kafkacons"))

	procMsgsConsumer, err := NewProcessedMsgConsumer(log, msgUC, kafkaConf.Brokers,
		saramaCfg, kafkaConf.Consumers.ProcessedMessages)
	if err != nil {
		return nil, err
	}

	kc := &KafkaConsumers{log: log, procMsgsConsumer: procMsgsConsumer}

	return kc, nil
}

func (c *KafkaConsumers) Close() error {
	if c.procMsgsConsumer != nil {
		return c.procMsgsConsumer.Close()
	}
	return errors.New("KafkaConsumers.Close: procMsgsConsumer is nil")
}

func (c *KafkaConsumers) ProcessedMsgs() *ProcessedMsgConsumer {
	if c.procMsgsConsumer == nil {
		c.log.Error("procMsgsConsumer is nil")
	}
	return c.procMsgsConsumer
}
