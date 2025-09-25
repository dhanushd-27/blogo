package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// helper to create a signed JWT string with desired claims for tests
func MakeTestJWT(secret string, userID int32, username, email string, expiresIn time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = time.Now().Add(expiresIn).Unix()
	return token.SignedString([]byte(secret))
}
