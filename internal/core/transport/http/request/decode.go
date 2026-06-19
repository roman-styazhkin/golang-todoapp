package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"failed to decode request, %v, err: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	var (
		err error
	)

	val, ok := dest.(validatable)
	if ok {
		err = val.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"failed to validate request, %v, err: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
