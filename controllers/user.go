package controllers

import (
	"net/http"
	"gorm.io/gorm"
)

func UserSignup(db *gorm.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Hello, World!</h1>"))
	}
}

func UserLogin(db *gorm.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Login</h1>"))
	}
}