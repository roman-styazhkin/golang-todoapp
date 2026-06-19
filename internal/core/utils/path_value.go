package core_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	path := r.PathValue(key)

	if path == "" {
		return 0, fmt.Errorf(
			"failed to get path by key=%s, err: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	pathInt, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get path by key=%s, err: %v, %w",
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return pathInt, nil
}
