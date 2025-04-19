package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/dto"
)


// ErrorHandler is a middleware that handles errors and formats the response.
//
// Behavior:
// - Processes the request and checks for any errors registered in the Gin context.
// - If errors are present, it creates a standardized error response using the `dto.NewErrorResponse` function.
// - Responds with an HTTP 500 status code and the error details in JSON format.
//
// Parameters:
// - c (*gin.Context): The Gin context containing the HTTP request and response.
func ErrorHandler(c *gin.Context) {
    c.Next()
    if len(c.Errors) > 0 {
        errResponse := dto.NewErrorResponse("An error occurred", c.Errors[0].Err)

        c.JSON(http.StatusInternalServerError, errResponse)
    }
}