package main

import (
	"log"
	"net/http"
	"os"

	"product/controller/auth"
	"product/controller/category"
	"product/controller/categoryAdmin"
	"product/controller/user"
	"product/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the router
	router := mux.NewRouter()

	// Define your routes
	router.HandleFunc("/users", user.GetUsers).Methods("GET")
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/regis", auth.Register).Methods("POST")

	router.HandleFunc("/category", category.GetCategory).Methods("GET")
	router.HandleFunc("/category", auth.JWTAuth(category.AddCategory)).Methods("POST")

	router.HandleFunc("/admin-category", categoryAdmin.GetCategoryAdmins).Methods("GET")
	router.HandleFunc("/admin-category", auth.JWTAuth(categoryAdmin.AddCategoryAdmin)).Methods("POST")
	router.HandleFunc("/admin-category/{id}", auth.JWTAuth(categoryAdmin.UpdateCategoryAdmin)).Methods("PUT")
	router.HandleFunc("/admin-category/{id}", auth.JWTAuth(categoryAdmin.DeleteCategoryAdmin)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug: true,
	})

	handler := c.Handler(router)

	handler.ServeHTTP(w, r)
}

func main() {
	// Load environment variables from .env file if running locally
	if os.Getenv("VERCEL_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	database.InitDB()

	http.HandleFunc("/", Handler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
