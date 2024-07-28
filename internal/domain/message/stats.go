package message

import "fmt"

type Stats struct {
	All       int
	Processed int
}

type StatsError struct {
	Err error
}

func (e *StatsError) Error() string {
	return fmt.Sprintf("message stats: %v", e.Err)
}

func (e *StatsError) Unwrap() error {
	return e.Err
}
