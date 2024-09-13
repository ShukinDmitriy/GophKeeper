package repositories

import "errors"

var errNotFound = errors.New("not found")

type NotFoundError struct {
	err error
}

func (e *NotFoundError) Error() string {
	return e.err.Error()
}
