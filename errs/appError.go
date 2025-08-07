package errs

import "strings"

type HTTPError struct {
	Message    string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(message string, code int) *HTTPError {
	return &HTTPError{
		Message:    strings.ToLower(message),
		StatusCode: code,
	}
}
