package main

import (
	"fmt"
	"log"
	"net/http"

	"auth-service/config"
	"auth-service/internal/db"
	"auth-service/internal/handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment vars")
	}

	// Load config
	cfg := config.LoadConfig()

	// Init MariaDB
	err := db.InitMariaDB(
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MariaDB: %v", err)
	}

	// Init router
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/auth/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/verify", handler.VerifyHandler).Methods("POST")
	r.HandleFunc("/auth/resend", handler.ResendVerificationHandler).Methods("POST")
	r.HandleFunc("/admin/users", handler.AdminCreateUserHandler).Methods("POST")

	r.HandleFunc("/admin/users/{id}", handler.UpdateUserHandler).Methods("PUT")

	r.HandleFunc("/users", handler.GetAllUsersHandler).Methods("GET")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"}, // or "*" for dev
		AllowedMethods:   []string{"GET", "PUT", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(r)

	// Start server
	fmt.Printf("✅ Server running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handlerWithCORS))
}
