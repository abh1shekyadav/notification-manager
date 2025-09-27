package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abh1shekyadav/notification-manager/internal/auth"
	"github.com/abh1shekyadav/notification-manager/internal/db"
	"github.com/abh1shekyadav/notification-manager/internal/user"
	"github.com/abh1shekyadav/notification-manager/middleware"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	if db == nil {
		log.Fatal("Database connection is nil. Check DB_CONN environment variable")
	}
	defer db.Close()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Secret is required")
	}
	validator := auth.NewHMACValidator(secret)
	exempt := map[string]bool{
		"/users/register": true,
		"/auth/login":     true,
		"/users":          false,
	}
	userRepo := user.NewPostgresRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	authService := auth.NewAuthService(userRepo, secret)
	authHandler := auth.NewAuthHandler(authService)
	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", middleware.Chain(userHandler.RegisterUser,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))
	mux.HandleFunc("/users", middleware.Chain(userHandler.FindUserByEmail,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))
	mux.HandleFunc("/auth/login", middleware.Chain(authHandler.Login,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
