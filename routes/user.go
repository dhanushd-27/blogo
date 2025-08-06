package routes

import (
	"github.com/dhanushd-27/blog_go/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func UserRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/signup", controllers.UserSignup(db)).Methods("POST")
	router.HandleFunc("/login", controllers.UserLogin(db)).Methods("POST")
}
