package dto

import (
	"encoding/json"
	"messagio_assignment/internal/domain/message"
)

type MessageValue struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Processed bool   `json:"processed"`

	bytes []byte
	err   error
}

func NewMessageValue(msg *message.Message) *MessageValue {
	v := &MessageValue{}

	v.ID = msg.ID
	v.Content = msg.Content
	v.Processed = msg.Processed

	var err error
	v.bytes, err = json.Marshal(v)
	if err != nil {
		v.err = err
	}

	return v
}

func (v *MessageValue) Encode() ([]byte, error) {
	if v.err != nil {
		return v.bytes, v.err
	}
	return v.bytes, nil
}

func (v *MessageValue) Length() int {
	return len(v.bytes)
}
