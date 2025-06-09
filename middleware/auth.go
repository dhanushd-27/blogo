package middleware

import (
	"net/http"

	"github.com/dhanushd-27/blog_go/helper/auth"
	"github.com/gorilla/context"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			http.Error(w, "Token Doesn't Exist", http.StatusForbidden)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		id, tokenErr := auth.VerifyToken(cookie.Value)

		if tokenErr != nil {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		context.Set(r, "id", id)

		next.ServeHTTP(w, r)
	})
}

