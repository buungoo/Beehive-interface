// Package authentication contains all the functions to both create a JSON Web Token and also verify it.
//
// This package is used when a user is logging in to the Api and needs to be verified.
// The func CreateToken is used to create a JWT for that user.
// JWTAuth is a wrapping function that acts as a middleware to protect secure endpoints.
// It is used when a user is trying to access an endpoint with this function wrapping it.
package authentication

import (
	"beehive_api/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// This needs to be made properly in the near future
var secretKey = []byte("secret-key")

// CreateToken signs and encodes the token with given claims and username.
func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// VerifyToken verifies a JSON Web Token and returns the claims of the token and a error. It take a token in string format as input.
func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil

}

// JWTAuth acts as middleware by wrapping the handlers of the secure endpoints. It takes a HandlerFunc as input and returns a Handlerfunc.
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get token from Authorization-header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.SendErrorResponse(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// In OAuth 2.0-specification
		if !strings.HasPrefix(tokenString, "Bearer ") {
			utils.SendErrorResponse(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Slice to remove "Bearer" to retreive token
		tokenString = tokenString[len("Bearer "):]

		// Verify the token and get the claims
		claims, err := verifyToken(tokenString)
		if err != nil {
			utils.SendErrorResponse(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract the username
		username, ok := claims["username"].(string)
		if !ok {
			utils.SendErrorResponse(w, "Error extrcting username", http.StatusInternalServerError)
			return
		}

		// Add username to request context
		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)

		// continue to next handler
		next(w, r)
	}
}
