package routes

import (
	"github.com/dhanushd-27/blog_go/controllers"
	"github.com/dhanushd-27/blog_go/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func BlogRoutes(router *mux.Router, db *gorm.DB) {
	blogRouter := router.PathPrefix("/blog").Subrouter()
	blogRouter.Use(middleware.AuthMiddleware)

	blogRouter.HandleFunc("", controllers.CreateBlog(db)).Methods("POST")
	blogRouter.HandleFunc("/all", controllers.ListBlog(db)).Methods("GET")
	blogRouter.HandleFunc("/{id}", controllers.FindBlog(db)).Methods("GET")
}
