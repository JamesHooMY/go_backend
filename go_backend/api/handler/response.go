package handler

import (
	"fmt"

	"go_backend/util"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status int    `json:"-"`
	Code   int    `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func ParseValidateError(err error, req any) error {
	var errMsg error

	errValidation, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	for _, fieldErr := range errValidation {
		errMsg = util.ErrorWrapper(errMsg, fmt.Errorf("%s: %s", fieldErr.Field(), fieldErr.Tag()))
	}

	return errMsg
}
