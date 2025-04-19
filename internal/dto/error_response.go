package dto

import "time"

// ErrorResponse represents a standardized error response.
//
// Fields:
// - Message (string): A human-readable error message describing the error.
// - ErrorDetails (string): The technical error details (optional). This field is omitted in the JSON response if empty.
// - Timestamp (time.Time): The time when the error occurred, useful for debugging and logging.
type ErrorResponse struct {
	Message      string    `json:"message"`
	ErrorDetails string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// Error implements the error interface for ErrorResponse.
//
// Behavior:
// - If ErrorDetails is not empty, it concatenates the Message and ErrorDetails fields.
// - If ErrorDetails is empty, it returns only the Message.
//
// Returns:
// - string: A string representation of the error.
func (e ErrorResponse) Error() string {
	if e.ErrorDetails != "" {
		return e.Message + ": " + e.ErrorDetails
	}
	return e.Message
}

// NewErrorResponse creates a new instance of ErrorResponse.
//
// Parameters:
// - message (string): A human-readable error message.
// - err (error): The technical error (optional). If provided, its message is included in the ErrorDetails field.
//
// Behavior:
// - If the `err` parameter is nil, the ErrorDetails field is left empty.
// - The Timestamp field is set to the current time.
//
// Returns:
// - ErrorResponse: A new error response object with the provided message, error details, and timestamp.
func NewErrorResponse(message string, err error) ErrorResponse {
	errorDetails := ""
	if err != nil {
		errorDetails = err.Error()
	}

	return ErrorResponse{
		Message:      message,
		ErrorDetails: errorDetails,
		Timestamp:    time.Now(),
	}
}
