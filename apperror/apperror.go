package apperror

import (
	"strings"

	"github.com/pkg/errors"
)

// AppError .
type AppError struct {
	Raw       error
	ErrorCode int64
	HTTPCode  int
	Info      string
	Message   string
	IsSentry  bool
}

func (e AppError) Error() string {
	if e.Raw != nil {
		return errors.Wrap(e.Raw, e.Info).Error()
	}

	return e.Message
}

func (e AppError) Is(target error) bool {
	if e.Raw != nil {
		return errors.Is(e.Raw, target)
	}

	return strings.Contains(e.Error(), target.Error())
}

func (e AppError) As(target error) bool {
	return errors.As(e.Raw, target)
}

// NewError ErrorCode = {1+ digit of service}{3 digits of model}{4 digits of error} .
func NewError(err error, httpCode int, errCode int64, message string, info string, isSentry bool) AppError {
	if info == "" {
		info = message
	}

	return AppError{
		Raw:       err,
		ErrorCode: errCode,
		HTTPCode:  httpCode,
		Message:   message,
		Info:      info,
		IsSentry:  isSentry,
	}
}
