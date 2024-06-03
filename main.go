package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"product/controller/auth"
	"product/controller/categoryAdmin"
	"product/controller/user"
	"product/controller/category"


	"product/database"

	"github.com/rs/cors"
	
)

func main() {
	database.InitDB()

	router := mux.NewRouter()

	router.HandleFunc("/users", user.GetUsers).Methods("GET")

	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/regis", auth.Register).Methods("POST")

	// Route for category
	router.HandleFunc("/category", category.GetCategory).Methods("GET")
	router.HandleFunc("/category", auth.JWTAuth(category.AddCategory)).Methods("POST")

	// Route for category admin with authentication middleware
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
	
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}


