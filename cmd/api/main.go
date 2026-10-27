
package main

import (
	"log"
	"net/http"
	"time"

	"auth-api/internal/database"
	"auth-api/internal/handler"
	"auth-api/internal/user"
)

func main() {
	databaseConnection, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer databaseConnection.Close()

	userRepository := user.NewRepository(databaseConnection)
	authHandler := handler.NewAuthHandler(userRepository)

	router := http.NewServeMux()
	authHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Authentication API running at http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("API server failed: ", err)
	}
}
