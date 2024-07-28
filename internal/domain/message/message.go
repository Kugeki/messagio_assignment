package message

import (
	"context"
	"fmt"
)

type Message struct {
	ID        int
	Content   string
	Processed bool
}

type Repository interface {
	Create(ctx context.Context, msg *Message) error
	GetByID(ctx context.Context, id int) (*Message, error)
	GetStats(ctx context.Context) (*Stats, error)
	UpdateProcessed(ctx context.Context, msg *Message) error
}

type Error struct {
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("message: %v", e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

type ErrorWithID struct {
	ID  int
	Err error
}

func (e *ErrorWithID) Error() string {
	return fmt.Sprintf("message with %q id: %v", e.ID, e.Err)
}

func (e *ErrorWithID) Unwrap() error {
	return e.Err
}
