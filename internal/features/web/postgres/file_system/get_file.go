package web_repository

import (
	"fmt"
	"os"

	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"file not found, %v, %w",
				err,
				core_errors.ErrNotFound,
			)
		}

		return nil, fmt.Errorf(
			"failed to get file, %s, %w",
			filePath,
			err,
		)
	}

	return file, nil
}
