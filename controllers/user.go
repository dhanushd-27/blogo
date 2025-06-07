package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dhanushd-27/blog_go/helper/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSignupType struct {
	Username string
	Email    string
	Password string
}

func UserSignup(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode json data fron r body
		var Body UserSignupType
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &Body)
		if err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 14)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		Body.Password = string(HashedPassword)

		// result := db.Create(&Body)

		// check if email exists
		// hash password
		// create user
		// send json response as response

	}
}

func UserLogin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.CreateToken("dhanush", "dhanush@gmail.com")

		if err != nil {
			w.Write([]byte("Error in creating token"))
			return
		}

		w.Write([]byte(token))
	}
}
