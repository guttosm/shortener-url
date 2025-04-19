package dto

import "time"

// ErrorResponse represents a standardized error response.
//
// Fields:
// - Message (string): A human-readable error message.
// - Error (string): The technical error details (optional).
// - Timestamp (time.Time): The time when the error occurred.
type ErrorResponse struct {
    Message   string    `json:"message"`
    Error     string    `json:"error,omitempty"`
    Timestamp time.Time `json:"timestamp"`
}

// NewErrorResponse creates a new instance of ErrorResponse.
//
// Parameters:
// - message (string): A human-readable error message.
// - err (error): The technical error (optional).
//
// Returns:
// - ErrorResponse: A new error response object.
func NewErrorResponse(message string, err error) ErrorResponse {
    errorDetails := ""
    if err != nil {
        errorDetails = err.Error()
    }

    return ErrorResponse{
        Message:   message,
        Error:     errorDetails,
        Timestamp: time.Now(),
    }
}