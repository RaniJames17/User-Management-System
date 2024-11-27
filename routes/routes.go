package routes

import (
	"net/http"
	"user-management-system/controllers"
	"user-management-system/middlewares"

	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health Check
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is up and running"))
	}).Methods("GET")

	// Protected routes (use the Authenticate middleware)
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middlewares.Authenticate)
	protected.HandleFunc("/protected-resource", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access to protected resource granted"))
	}).Methods("GET")

	// User Sign-Up
	router.HandleFunc("/api/signup", controllers.SignUp).Methods("POST")

	//User Sign-In
	router.HandleFunc("/api/signin", controllers.SignIn).Methods("POST")

	return router
}
