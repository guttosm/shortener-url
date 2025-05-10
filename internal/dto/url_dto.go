package dto

// ShortenRequest represents the request body for shortening a URL.
//
// Fields:
//   - URL (string): The original long URL provided by the user.
//     It must be a valid URL and is required in the request body.
type ShortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// ShortenResponse represents the response body for a shortened URL.
//
// Fields:
// - ShortID (string): The unique identifier for the shortened URL.
// - ShortURL (string): The full shortened URL.
type ShortenResponse struct {
	ShortID  string `json:"short_id"`
	ShortURL string `json:"short_url"`
}
