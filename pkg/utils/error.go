package utils

import "github.com/pkg/errors"

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// fixme this is abstraction, accepted for now, should be updated later
func StackError(err error) error {
	if err == nil {
		return nil
	}

	_, ok := err.(stackTracer)
	if ok {
		return err
	}

	return errors.WithStack(err)
}
