package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTCookieMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			claims := token.Claims.(jwt.MapClaims)
			c.Set("user", token)
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("email", claims["email"])

			return next(c)
		}
	}
}
