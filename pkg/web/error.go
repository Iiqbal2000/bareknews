package web

import "errors"

type RequestError struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) *RequestError {
	return &RequestError{err, status}
}

func (rr RequestError) Error() string {
	return rr.Err.Error()
}

func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
