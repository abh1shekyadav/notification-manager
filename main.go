package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/abh1shekyadav/notification-manager/internal/auth"
	"github.com/abh1shekyadav/notification-manager/internal/db"
	"github.com/abh1shekyadav/notification-manager/internal/kafka"
	"github.com/abh1shekyadav/notification-manager/internal/notification"
	"github.com/abh1shekyadav/notification-manager/internal/notifier"
	"github.com/abh1shekyadav/notification-manager/internal/user"
	"github.com/abh1shekyadav/notification-manager/middleware"
)

func main() {
	mux := http.NewServeMux()

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	if dbConn == nil {
		log.Fatal("Database connection is nil. Check DB_CONN environment variable")
	}
	defer dbConn.Close()

	// JWT validator
	secret := getEnvOrFatal("JWT_SECRET")
	validator := auth.NewHMACValidator(secret)

	// Routes exempted from auth
	exempt := map[string]bool{
		"/users/register": true,
		"/auth/login":     true,
		"/users":          false,
		"/notify":         false,
		"/notification":   false,
	}

	// User management
	userRepo := user.NewPostgresRepo(dbConn)
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

	// Auth
	authService := auth.NewAuthService(userRepo, secret)
	authHandler := auth.NewAuthHandler(authService)
	mux.HandleFunc("/auth/login", middleware.Chain(authHandler.Login,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))

	// Notifications
	notificationRepo := notification.NewNotificationRepo(dbConn)
	brokers := strings.Split(getEnvOrFatal("KAFKA_BROKERS"), ",")
	topic := getEnvOrFatal("KAFKA_TOPIC")
	producer := kafka.NewKafkaProducer(brokers, topic)
	notificationService := notification.NewNotificationService(notificationRepo, producer)
	notificationHandler := notification.NewNotificationHandler(notificationService)

	// --- Notifiers (Twilio SMS + SendGrid Email) ---
	twilioSID := getEnvOrFatal("TWILIO_ACCOUNT_SID")
	twilioToken := getEnvOrFatal("TWILIO_AUTH_TOKEN")
	twilioFrom := getEnvOrFatal("TWILIO_FROM_NUMBER")
	smsNotifier := notifier.NewTwilioSMSNotifier(twilioSID, twilioToken, twilioFrom)

	sendgridKey := getEnvOrFatal("SENDGRID_API_KEY")
	sendgridFrom := getEnvOrFatal("SENDGRID_FROM_EMAIL")
	emailNotifier := notifier.NewSendGridEmailNotifier(sendgridKey, sendgridFrom)

	// Start Kafka consumer
	consumer := kafka.NewConsumer(brokers, topic, "notification-consumer-group")

	go func() {
		log.Println("Starting Kafka consumer...")
		handler := notification.NewConsumerHandler(notificationRepo, smsNotifier, emailNotifier)
		err := consumer.Start(context.Background(), handler)
		if err != nil {
			log.Fatalf("Kafka consumer stopped: %v", err)
		}
	}()
	// Notification routes
	mux.HandleFunc("/notify", middleware.Chain(notificationHandler.Notify,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware(exempt, validator),
	))
	mux.HandleFunc("/notification", middleware.Chain(notificationHandler.FindNotificationByID,
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
