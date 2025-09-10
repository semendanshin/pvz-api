package domain

import "errors"

var (
	// ErrAlreadyExists is an error for already existing entity
	ErrAlreadyExists = errors.New("entity already exists")
	// ErrNotFound is an error for not found entity
	ErrNotFound = errors.New("entity not found")
	// ErrInvalidArgument is an error for invalid argument
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrInternal is an error for internal server error
	ErrInternal = errors.New("internal server error")
)
