package controllers

import (
	"net/http"
	"github.com/dhanushd-27/blog_go/helper/auth"
	"gorm.io/gorm"
)

func UserSignup(db *gorm.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		token, err := auth.CreateToken("dhanush", "dhanush@gmail.com")

		if err != nil {
			w.Write([]byte("Error in creating token"))
			return
		}

		w.Write([]byte(token))
	}
}

func UserLogin(db *gorm.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Login</h1>"))
	}
}