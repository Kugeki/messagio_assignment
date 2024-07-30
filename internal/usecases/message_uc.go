package usecases

import (
	"context"
	"messagio_assignment/internal/domain/message"
)

type MessageUC struct {
	MessageRepo      message.Repository
	MessagesProducer message.Producer
}

func NewMessageUC(messageRepo message.Repository, messagesProd message.Producer) *MessageUC {
	return &MessageUC{MessageRepo: messageRepo, MessagesProducer: messagesProd}
}

func (uc *MessageUC) CreateMessage(ctx context.Context, msg *message.Message) error {
	return uc.MessageRepo.Create(ctx, msg)
}

func (uc *MessageUC) GetStats(ctx context.Context) (*message.Stats, error) {
	return uc.MessageRepo.GetStats(ctx)
}

func (uc *MessageUC) UpdateProcessedMessage(ctx context.Context, msg *message.Message) error {
	return uc.MessageRepo.UpdateProcessed(ctx, msg)
}
