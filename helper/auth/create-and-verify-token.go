package auth

import (
	"fmt"
	"os"
	"time"


	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, email string, userid uint) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       userid,
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

func VerifyToken(tokenString string) (uint, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if _, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		claims := token.Claims.(*jwt.MapClaims)

		id, ok := (*claims)["id"]

		if ok {
			// Convert interface{} to float64 first (since JSON numbers are decoded as float64)
			idFloat, ok := id.(float64)
			if !ok {
				return 0, fmt.Errorf("invalid id type in token")
			}
			// Convert float64 to uint
			return uint(idFloat), nil
		}
	}
	return 0, fmt.Errorf("invalid token")
}
