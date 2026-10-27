
package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type Repository struct {
	db *sql.DB
}

//This creates a repository and gives it the database connection.
func NewRepository(database *sql.DB) *Repository {
	return &Repository{db: database}
}
//create function takes in a context, name, email, and password hash, and returns a User struct and an error. It inserts the user into the database and retrieves the created user.
func (repository *Repository) Create(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (User, error) {
	var createdUser User

	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES (?, ?, ?)
	`
//Executing the query
	result, err := repository.db.ExecContext(
		ctx,
		query,
		name,
		email,
		passwordHash,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return User{}, ErrEmailAlreadyExists
		}

		return User{}, fmt.Errorf("insert user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("get new user ID: %w", err)
	}

	err = repository.db.QueryRowContext(
		ctx,
		`SELECT id, name, email, created_at FROM users WHERE id = ?`,
		userID,
	).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.Email,
		&createdUser.CreatedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("read created user: %w", err)
	}

	return createdUser, nil
}
