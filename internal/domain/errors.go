package domain

import "errors"

var (
	// create errors
	ErrAlreadyExists = errors.New("already exists")
	ErrNotCreated    = errors.New("not created")

	// get errors
	ErrNotFound = errors.New("not found")
)
