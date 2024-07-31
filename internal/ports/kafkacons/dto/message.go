package dto

import (
	"encoding/json"
)

type MessageValue struct {
	ID int `json:"id"`
}

func MessageValueFromBytes(data []byte) (*MessageValue, error) {
	mv := MessageValue{}

	err := json.Unmarshal(data, &mv)
	if err != nil {
		return nil, err
	}

	return &mv, nil
}
