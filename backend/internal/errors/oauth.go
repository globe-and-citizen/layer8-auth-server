package errors

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
)

type OAuthError struct {
	Code        consts.OAuthErrorCode
	Description string
	StatusCode  int
	Err         error // internal wrapped error
}

func (e *OAuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return string(e.Code)
}

func (e *OAuthError) Unwrap() error {
	return e.Err
}
