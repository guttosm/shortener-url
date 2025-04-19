package dto

// LoginRequest represents the request body for the login endpoint.
//
// Fields:
// - Username (string): The username for authentication (required).
// - Password (string): The password for authentication (required).
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
