package middleware

import (
	"net/http"

	"github.com/dhanushd-27/blog_go/helper/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		tokenErr := auth.VerifyToken(cookie.Value)

		if tokenErr != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	})
}