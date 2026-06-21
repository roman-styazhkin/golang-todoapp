package core_utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	intQueryParam, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get int query param by key=%s, err: %v, %w",
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &intQueryParam, nil
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil, nil
	}

	format := "2006-01-02"

	date, err := time.Parse(format, param)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse date, param=%s by key=%s, err: %v, %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}
