package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// // ExtractJWTFromCookie middleware extracts JWT from cookies
// func ExtractJWTFromCookie(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		cookie, err := c.Cookie("token")
// 		if err == nil && cookie.Value != "" {
// 			c.Request().Header.Set("Authorization", "Bearer "+cookie.Value)
// 		}

// 		return next(c)
// 	}
// }

func JWTCookieMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil || cookie.Value == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
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
