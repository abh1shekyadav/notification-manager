package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/abh1shekyadav/notification-manager/internal/auth"
	"github.com/abh1shekyadav/notification-manager/internal/db"
	"github.com/abh1shekyadav/notification-manager/internal/notification"
	"github.com/abh1shekyadav/notification-manager/internal/user"
	"github.com/abh1shekyadav/notification-manager/middleware"
)

func main() {
	mux := http.NewServeMux()
	db, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	if db == nil {
		log.Fatal("Database connection is nil. Check DB_CONN environment variable")
	}
	defer db.Close()
	secret := getEnvOrFatal("JWT_SECRET")
	validator := auth.NewHMACValidator(secret)
	exempt := map[string]bool{
		"/users/register": true,
		"/auth/login":     true,
		"/users":          false,
		"/notify":         false,
		"/notification":   false,
	}
	userRepo := user.NewPostgresRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	mux.HandleFunc("/users/register", middleware.Chain(userHandler.RegisterUser,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))
	mux.HandleFunc("/users", middleware.Chain(userHandler.FindUserByEmail,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))

	authService := auth.NewAuthService(userRepo, secret)
	authHandler := auth.NewAuthHandler(authService)
	mux.HandleFunc("/auth/login", middleware.Chain(authHandler.Login,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))

	notificationRepo := notification.NewNotifcationRepo(db)
	brokers := strings.Split(getEnvOrFatal("KAFKA_BROKERS"), ",")
	topic := getEnvOrFatal("KAFKA_TOPIC")
	producer := notification.NewKafkaProducer(brokers, topic)
	notificationService := notification.NewNotificationService(notificationRepo, producer)
	notificationHandler := notification.NewNotificationHandler(notificationService)
	notification.StartConsumer(brokers, topic, notificationRepo)
	mux.HandleFunc("/notify", middleware.Chain(notificationHandler.Notify,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))
	mux.HandleFunc("/notification", middleware.Chain(notificationHandler.FindNotificationById,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return val
}
