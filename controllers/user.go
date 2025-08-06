package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dhanushd-27/blog_go/helper/auth"
	"github.com/dhanushd-27/blog_go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserSignup(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode json data fron r body
		var Body models.User
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

		// Check if email exists
		var existingUser models.User
		if err := db.Where("email = ?", Body.Email).First(&existingUser).Error; err == nil {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		// Hash password
		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 14)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		Body.Password = string(HashedPassword)

		// Create user
		if err := db.Create(&Body).Error; err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Send success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully",
		})
	}
}

func UserLogin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Body from request body, parse body
		var user models.User
		bodyBytes, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &user)
		if err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
		}

		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))

		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := auth.CreateToken(existingUser.Username, existingUser.Email, existingUser.ID)

		if err != nil {
			w.Write([]byte("Error in creating token"))
			return
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}
