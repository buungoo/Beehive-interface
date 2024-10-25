package authentication

import(
	"beehive_api/utils"
	"net/http"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"strings"
)

// This needs to be made properly in the near future
var secretKey = []byte("secret-key")

// Creates, signs and encodes the token with given claims and username
func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

return tokenString, nil

}

// Verifies the token
func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	
return nil

}

// This acts as "middleware" by wrapping the handlers of the secure endpoints.
// Will only let verified user pass
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

		// Verify the token
		err := verifyToken(tokenString)
		if err != nil {
			utils.SendErrorResponse(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// continue to next handler
		next(w,r)
	}
}