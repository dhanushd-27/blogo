package helper

import (
	"fmt"
	"net/http"
	"github.com/dhanushd-27/blog_go/routes"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ApiServer struct {
	addr string
	db   *gorm.DB
}

func NewApiServer(addr string, db *gorm.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Register user routes
	routes.UserRoutes(subrouter, s.db)

	fmt.Println("Server started on port 8080")
	return http.ListenAndServe(s.addr, router)
}
