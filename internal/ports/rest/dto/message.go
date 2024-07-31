package dto

import "messagioassignment/internal/domain/message"

type CreateMessageReq struct {
	Content   string `json:"content"`
	Processed bool   `json:"processed,omitempty"`
}

func (r *CreateMessageReq) ToDomain() *message.Message {
	return &message.Message{
		Content:   r.Content,
		Processed: r.Processed,
	}
}

type CreateMessageResp struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Processed bool   `json:"processed"`
}

func (r *CreateMessageResp) FromDomain(msg *message.Message) {
	r.ID = msg.ID
	r.Content = msg.Content
	r.Processed = msg.Processed
}
