package errs

type HTTPError struct {
	Message    string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(message string, code int) *HTTPError {
	return &HTTPError{
		Message:    message,
		StatusCode: code,
	}
}
