package auth

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
)

// JWTSecret is the secret key used to sign the JWT tokens.
var JWTSecret = []byte("your-secret-key") // Replace with a secure key

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
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JWTSecret)
}

// ValidateToken validates a JWT token and extracts the claims.
//
// Parameters:
// - tokenString (string): The JWT token to validate.
//
// Returns:
// - jwt.MapClaims: The claims extracted from the token.
// - error: An error if validation fails.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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