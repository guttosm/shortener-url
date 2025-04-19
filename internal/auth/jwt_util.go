package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTSecret is the secret key used to sign the JWT tokens.
var JWTSecret = []byte("your-secret-key") // TODO: replace with secure env-based key in production

// GenerateToken generates a JWT token for a given user ID.
//
// Parameters:
// - userID (string): The ID of the user.
//
// Returns:
// - string: The signed JWT token.
// - error: An error if token generation fails.
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// TokenValidator defines the interface for validating JWT tokens.
type TokenValidator interface {
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

// JWTValidator is the concrete implementation of TokenValidator using HMAC signing.
type JWTValidator struct{}

// ValidateToken validates a JWT token and returns its claims.
//
// Parameters:
// - tokenString (string): The JWT token string.
//
// Returns:
// - map[string]interface{}: The claims extracted from the token.
// - error: An error if validation fails.
func (j *JWTValidator) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token uses HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
