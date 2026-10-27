package main

import (
	"log"

	"auth-api/internal/database"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	log.Println("Successfully connected to XAMPP MySQL!")
}
