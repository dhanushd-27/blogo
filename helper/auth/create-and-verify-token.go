package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, email string) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"email":    email,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if _, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		claims := token.Claims.(*jwt.MapClaims)

		email, ok := (*claims)["email"].(string)

		if ok {
			fmt.Println("Email", email)
			return email, nil
		}
	} else {
		return "", fmt.Errorf("invalid token")
	}

	return token.Raw, nil
}
