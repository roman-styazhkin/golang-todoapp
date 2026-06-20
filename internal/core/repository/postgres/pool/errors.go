package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("err no rows")
	ErrViolatesForeignKey = errors.New("err violates foreign key custom")
	ErrUnknown            = errors.New("unknown error")
)
