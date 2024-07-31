package kafkaprod

import (
	"errors"
	"github.com/IBM/sarama"
	"log/slog"
	"messagio_assignment/internal/adapters/kafkaprod/dto"
	"messagio_assignment/internal/config"
	"messagio_assignment/internal/domain/message"
	"messagio_assignment/internal/logger"
)

type MessageProducer struct {
	p     sarama.AsyncProducer
	log   *slog.Logger
	topic string
}

func NewMessagesProducer(log *slog.Logger, brokerList []string,
	saramaCfg *sarama.Config, producerCfg config.KafkaProducer) (*MessageProducer, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	log = log.With(slog.String("component", "adapters/kafkaprod/message_producer"))

	if saramaCfg == nil {
		saramaCfg = sarama.NewConfig()
	}

	conf := *saramaCfg

	conf.Producer.Timeout = producerCfg.Timeout

	conf.Producer.Retry.Max = producerCfg.Retry.Max
	conf.Producer.Retry.Backoff = producerCfg.Retry.Backoff

	conf.Producer.Flush.Bytes = producerCfg.Flush.Bytes
	conf.Producer.Flush.Messages = producerCfg.Flush.Messages
	conf.Producer.Flush.Frequency = producerCfg.Flush.Frequency
	conf.Producer.Flush.MaxMessages = producerCfg.Flush.MaxMessages

	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewAsyncProducer(brokerList, &conf)
	if err != nil {
		return nil, err
	}

	mp := &MessageProducer{p: producer, log: log, topic: producerCfg.Topic}

	go func() {
		for err = range producer.Errors() {
			mp.log.Error("messages producer error", logger.Err(err))
		}
	}()

	return mp, nil
}

func (p *MessageProducer) Close() error {
	if p.p != nil {
		return p.p.Close()
	}
	return errors.New("MessageProducer.Close: async producer is nil")
}

func (p *MessageProducer) Produce(msg *message.Message) {
	p.p.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   nil, // sarama.StringEncoder(strconv.Itoa(msg.ID)),
		Value: dto.NewMessageValue(msg),
	}
}
