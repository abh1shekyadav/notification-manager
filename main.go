package main

import (
	"log"
	"net/http"

	"github.com/abh1shekyadav/notification-manager/internal/db"
	"github.com/abh1shekyadav/notification-manager/internal/user"
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
	userRepo := user.NewPostgresRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", userHandler.RegisterUser)
	mux.HandleFunc("/users", userHandler.FindUserByEmail)
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
