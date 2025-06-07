package helper

import (
	"fmt"
	"net/http"

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

	subrouter.HandleFunc("/", getUsers).Methods("GET")

	// subrouter.HandleFunc("/users", s.handleGetUsers).Methods("GET")

	fmt.Println("Server started on port 8080")
	return http.ListenAndServe(s.addr, router)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
