package consts

import (
	"errors"
	"net/http"
)

var (
	ErrUnknownRepositoryFactory = errors.New("unknown repository factory")
	ErrRequiredFieldMissing     = errors.New("required field missing")
	ErrInvalidField             = errors.New("invalid field")
	ErrInvalidEmail             = errors.New("invalid email")
	ErrAlreadyExists            = errors.New("already exists")
	ErrNotFound                 = errors.New("not found")
	ErrInvalidPassword          = errors.New("password is invalid")
	ErrMissingFields            = errors.New("missing fields")
	ErrInternalServer           = errors.New("internal server error")
	ErrBadRequest               = errors.New("bad request")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrDuplicateKey             = errors.New("duplicate key")
	ErrRequestCanceled          = errors.New("request canceled")
	ErrRequestTimeout           = errors.New("request timeout")
)

func MapErrorToStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrInternalServer), err == nil:
		return http.StatusInternalServerError
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}
