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
	mux.HandleFunc("/user/register", userHandler.RegisterUser)
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
