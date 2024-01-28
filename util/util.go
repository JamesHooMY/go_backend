package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseValidateError(err error) error {
	var errMsg error

	errValidation, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	for _, fieldErr := range errValidation {
		errMsg = ErrorWrapper(errMsg, fmt.Errorf("%s: %s", fieldErr.Field(), fieldErr.Tag()))
	}

	return errMsg
}

func ErrorWrapper(existErr, newErr error) error {
	if newErr != nil && existErr == nil {
		return newErr
	}

	if newErr == nil && existErr != nil {
		return existErr
	}

	return fmt.Errorf("%w, %w", existErr, newErr)
}
