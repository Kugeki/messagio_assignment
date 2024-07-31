package kafkacons

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log/slog"
	"messagioassignment/internal/config"
	"messagioassignment/internal/domain/message"
	"messagioassignment/internal/logger"
	"messagioassignment/internal/ports/kafkacons/dto"
)

type MessagesUsecase interface {
	UpdateProcessedMessage(ctx context.Context, msg *message.Message) error
}

type ProcessedMsgConsumer struct {
	log *slog.Logger

	msgUC  MessagesUsecase
	cg     sarama.ConsumerGroup
	topics []string
}

func NewProcessedMsgConsumer(log *slog.Logger, msgUC MessagesUsecase, brokerList []string,
	saramaCfg *sarama.Config, consumerCfg config.KafkaConsumer) (*ProcessedMsgConsumer, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	log = log.With(slog.String("component", "ports/kafkacons/processed_message_consumer"))

	if saramaCfg == nil {
		saramaCfg = sarama.NewConfig()
	}

	conf := *saramaCfg

	conf.Consumer.Retry.Backoff = consumerCfg.Retry.Backoff

	conf.Consumer.MaxWaitTime = consumerCfg.MaxWaitTime
	conf.Consumer.Fetch.Min = consumerCfg.Fetch.Min
	conf.Consumer.Fetch.Default = consumerCfg.Fetch.Default
	conf.Consumer.Fetch.Max = consumerCfg.Fetch.Max

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, consumerCfg.Group, &conf)
	if err != nil {
		return nil, err
	}

	topics := consumerCfg.Topics

	return &ProcessedMsgConsumer{log: log, msgUC: msgUC, cg: consumerGroup, topics: topics}, nil
}

func (c *ProcessedMsgConsumer) Close() error {
	if c.cg != nil {
		return c.cg.Close()
	}
	return errors.New("ProcessedMsgConsumer.Close: cg is nil")
}

// StartConsume is blocking
func (c *ProcessedMsgConsumer) StartConsume(ctx context.Context) {
	for {
		if err := c.cg.Consume(ctx, c.topics, c); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			c.log.Error("error from consume",
				slog.String("function", "StartConsume"),
				logger.Err(err),
			)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (c *ProcessedMsgConsumer) ConsumeClaim(ses sarama.ConsumerGroupSession, cm sarama.ConsumerGroupClaim) error {
	log := c.log.With(
		slog.String("handler", "processed message"),
		slog.String("member_id", ses.MemberID()),
		slog.Int64("generation_id", int64(ses.GenerationID())),
	)
	for {
		select {
		case claimMsg, ok := <-cm.Messages():
			if !ok {
				log.Info("claim message channel was closed")
				return nil
			}

			log = c.log.With(
				slog.String("handler", "processed message"),
				slog.String("topic", claimMsg.Topic),
				slog.Int64("partition", int64(claimMsg.Partition)),
				slog.Int64("offset", claimMsg.Offset),
				slog.Time("timestamp", claimMsg.Timestamp),
				slog.Time("block_timestamp", claimMsg.BlockTimestamp),
			)
			// go func??
			err := c.HandleMessage(ses.Context(), log, claimMsg)
			if err != nil {
				log.Error("handle claim message", logger.Err(err))
			}

			ses.MarkMessage(claimMsg, "")
		case <-ses.Context().Done():
			ses.Commit()
			return nil
		}
	}
}

func (c *ProcessedMsgConsumer) HandleMessage(ctx context.Context, log *slog.Logger, claimMsg *sarama.ConsumerMessage) error {
	mv, err := dto.MessageValueFromBytes(claimMsg.Value)
	if err != nil {
		log.Warn("message value from bytes", logger.Err(err))
		return err
	}

	msg := &message.Message{
		ID:        mv.ID,
		Processed: true,
	}

	err = c.msgUC.UpdateProcessedMessage(ctx, msg)
	if err != nil {
		log.Error("update processed message", logger.Err(err))
		return err
	}

	return nil
}

func (c *ProcessedMsgConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ProcessedMsgConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
