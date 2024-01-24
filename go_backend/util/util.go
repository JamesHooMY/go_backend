package util

import "fmt"

func ErrorWrapper(existErr, newErr error) error {
	if newErr != nil && existErr == nil {
		return newErr
	}

	if newErr == nil && existErr != nil {
		return existErr
	}

	return fmt.Errorf("%w, %w", existErr, newErr)
}
