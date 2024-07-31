package kafkaprod

import (
	"errors"
	"github.com/IBM/sarama"
	"log/slog"
	"messagioassignment/internal/config"
	"messagioassignment/internal/logger"
)

type KafkaProducers struct {
	log *slog.Logger

	messagesProducer *MessageProducer
}

func New(log *slog.Logger, saramaCfg *sarama.Config, kafkaConf config.Kafka) (*KafkaProducers, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	log = log.With(slog.String("component", "adapters/kafkaprod"))

	messagesProducer, err := NewMessageProducer(log, kafkaConf.Brokers,
		saramaCfg, kafkaConf.Producers.Messages)
	if err != nil {
		return nil, err
	}

	mq := &KafkaProducers{log: log, messagesProducer: messagesProducer}

	return mq, nil
}

func (p *KafkaProducers) Close() error {
	if p.messagesProducer != nil {
		return p.messagesProducer.Close()
	}
	return errors.New("KafkaProducers.Close: messagesProducer is nil")
}

func (p *KafkaProducers) Messages() *MessageProducer {
	if p.messagesProducer == nil {
		p.log.Error("messages producer is nil")
	}

	return p.messagesProducer
}
