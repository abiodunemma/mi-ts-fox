package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "root"
	}

	password := os.Getenv("DB_PASSWORD")

	databaseName := os.Getenv("DB_DATABASE")
	if databaseName == "" {
		databaseName = "auth_api"
	}

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		username,
		password,
		host,
		port,
		databaseName,
	)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("open MySQL connection: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	contextWithTimeout, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(contextWithTimeout); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to MySQL: %w", err)
	}

	return db, nil
}
