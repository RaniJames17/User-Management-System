package routes

import (
	"net/http"
	"user-management-system/controllers"

	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health Check
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is up and running"))
	}).Methods("GET")

	// User Sign-Up
	router.HandleFunc("/api/signup", controllers.SignUp).Methods("POST")

	return router
}
