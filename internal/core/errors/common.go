package core_errors

import "errors"

var (
	ErrConflict        = errors.New("err conflict")
	ErrInvalidArgument = errors.New("err invalid")
	ErrNotFound        = errors.New("err not found")
)
