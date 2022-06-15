package apperror

import (
	"net/http"
)

// ErrUnauthorized .
func ErrUnauthorized(err error) AppError {
	return AppError{
		Raw:       err,
		HTTPCode:  http.StatusUnauthorized,
		ErrorCode: 100000001,
		Info:      "unauthorized",
		Message:   "Unauthorized.",
	}
}

// ErrUnauthorizedExpiredToken .
func ErrUnauthorizedExpiredToken(err error) AppError {
	return AppError{
		Raw:       err,
		HTTPCode:  http.StatusUnauthorized,
		ErrorCode: 100000002,
		Info:      "expired_token",
		Message:   "Token expired.",
	}
}

// ErrCommitTransaction .
func ErrCommitTransaction(err error) AppError {
	return AppError{
		Raw:       err,
		HTTPCode:  http.StatusInternalServerError,
		ErrorCode: 100000003,
		Info:      "error commit sql transaction",
		Message:   "Fail to commit SQL transaction",
		IsSentry:  true,
	}
}

// ErrInvalidInput .
func ErrInvalidInput(err error) AppError {
	return AppError{
		Raw:       err,
		HTTPCode:  http.StatusNotAcceptable,
		ErrorCode: 100000004,
		Info:      "invalid input",
		Message:   "Invalid input",
	}
}
