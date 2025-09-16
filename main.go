package main

import (
	"log"
	"net/http"

	"github.com/abh1shekyadav/notification-manager/internal/user"
)

func main() {
	userRepo := user.NewInMemoryRepo()
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
