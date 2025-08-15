package errs

import (
	"log"
	"runtime"
	"strings"
)

type HTTPError struct {
	Message    string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(message string, code int) *HTTPError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[CALLER] %s:%d", file, line)
	}
	return &HTTPError{
		Message:    strings.ToLower(message),
		StatusCode: code,
	}
}
