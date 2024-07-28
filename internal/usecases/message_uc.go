package usecases

import (
	"context"
	"messagio_assignment/internal/domain/message"
)

type MessageUC struct {
	MessageRepo message.Repository
}

func NewMessageUC(messageRepo message.Repository) *MessageUC {
	return &MessageUC{MessageRepo: messageRepo}
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
