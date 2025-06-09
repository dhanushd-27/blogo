package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/dhanushd-27/blog_go/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blog models.Blog
		blogBytes, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(blogBytes, &blog)
		if err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}
		email, ok := context.Get(r, "email").(string)
		if !ok {
			http.Error(w, "Invalid email in context", http.StatusBadRequest)
			return
		}

		blog.Email = email

		var user models.User
		if err := db.Where("email = ?", blog.Email).First(&user).Error; err != nil {
			http.Error(w, "Error finding user", http.StatusInternalServerError)
			return
		}

		blog.User = user

		if err := db.Create(&blog).Error; err != nil {
			http.Error(w, "Error creating blog", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(blog)
	}
}

func ListBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blogs []models.Blog
		if err := db.Find(&blogs).Error; err != nil {
			http.Error(w, "Error fetching blogs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(blogs)
	}
}

func FindBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid blog ID", http.StatusBadRequest)
			return
		}

		var blog models.Blog
		if err := db.First(&blog, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Blog not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error fetching blog", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(blog)
	}
}
