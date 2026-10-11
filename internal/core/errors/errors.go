package core_errors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid agument")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrUnauthtorize    = errors.New("unauthtorize")
	ErrForbidden       = errors.New("forbidden")
)
