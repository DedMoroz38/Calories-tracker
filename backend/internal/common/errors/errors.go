// Package errors defines the application's transport-agnostic error type.
// Services return *APIError so controllers can translate a failure into an
// HTTP status without re-deriving it from sentinel errors.
package errors

// APIError carries an HTTP status code alongside a client-safe message.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

func New(statusCode int, message string) *APIError {
	return &APIError{StatusCode: statusCode, Message: message}
}

func BadRequest(message string) *APIError   { return New(400, message) }
func Unauthorized(message string) *APIError { return New(401, message) }
func NotFound(message string) *APIError     { return New(404, message) }
func Internal(message string) *APIError     { return New(500, message) }
