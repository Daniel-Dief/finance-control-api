package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"

	"github.com/lib/pq"
)

// CreateUser inserts a new active user and returns the created record.
// It returns ErrDuplicateLogin when the login is already taken.
func CreateUser(name, login, passwordHash string) (models.User, error) {
	query := `
		INSERT INTO "Users" ("StatusId", "Name", "Login", "PasswordHash")
		VALUES ($1, $2, $3, $4)
		RETURNING "Id", "Name", "Login", "PasswordHash"
	`

	var u models.User
	err := Pool.QueryRowContext(context.Background(), query, StatusActive, name, login, passwordHash).Scan(
		&u.ID, &u.Name, &u.Login, &u.PasswordHash,
	)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return models.User{}, ErrDuplicateLogin
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.User{}, ErrGenericDatabase
	}

	return u, nil
}

// GetUserByLogin retrieves a user by their login from the database.
func GetUserByLogin(login string) (models.User, error) {
	query := `
		SELECT "Id", "Name", "Login", "PasswordHash"
		FROM "Users"
		WHERE "Login" = $1 AND "StatusId" = $2
	`

	var u models.User
	err := Pool.QueryRowContext(context.Background(), query, login, StatusActive).Scan(
		&u.ID, &u.Name, &u.Login, &u.PasswordHash,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.User{}, ErrGenericDatabase
	}

	return u, nil
}

// GetUserByID retrieves a user by their ID from the database.
func GetUserByID(id int) (models.User, error) {
	query := `
		SELECT "Id", "Name", "Login", "PasswordHash"
		FROM "Users"
		WHERE "Id" = $1 AND "StatusId" = $2
	`

	var u models.User
	err := Pool.QueryRowContext(context.Background(), query, id, StatusActive).Scan(
		&u.ID, &u.Name, &u.Login, &u.PasswordHash,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.User{}, ErrGenericDatabase
	}

	return u, nil
}
