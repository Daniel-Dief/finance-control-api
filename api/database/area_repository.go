package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
)

// ListAreas retrieves a list of areas from the database, optionally filtered by name.
// If name is nil, all areas are returned.
func ListAreas(name *string) ([]models.Area, error) {
	query := `
		SELECT "Id", "Name"
		FROM "Areas"
		WHERE ($1::text IS NULL OR "Name" ILIKE '%' || $1 || '%')
	`

	rows, err := Pool.QueryContext(context.Background(), query, name)
	if err != nil {
		return nil, err
	} else if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var result []models.Area
	for rows.Next() {
		var a models.Area
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

// GetAreaByID retrieves an area by its ID from the database.
func GetAreaByID(id int) (models.Area, error) {
	query := `
		SELECT "Id", "Name"
		FROM "Areas"
		WHERE "Id" = $1
	`

	var a models.Area
	err := Pool.QueryRowContext(context.Background(), query, id).Scan(&a.ID, &a.Name)
	if err != nil {
		return models.Area{}, err
	}

	return a, nil
}

// CreateArea inserts a new area into the database and returns the created area.
func CreateArea(name string) (models.Area, error) {
	query := `
		INSERT INTO "Areas" ("Name")
		VALUES ($1)
		RETURNING "Id", "Name"
	`

	var a models.Area
	err := Pool.QueryRowContext(context.Background(), query, name).Scan(&a.ID, &a.Name)
	if err != nil {
		return models.Area{}, err
	}

	return a, nil
}

// UpdateArea updates the name of an existing area in the database and returns the updated area.
func UpdateArea(id int, name string) (models.Area, error) {
	query := `
		UPDATE "Areas"
		SET "Name" = $2
		WHERE "Id" = $1
		RETURNING "Id", "Name"
	`

	var a models.Area
	err := Pool.QueryRowContext(context.Background(), query, id, name).Scan(&a.ID, &a.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Area{}, errors.New("Area not found")
	} else if err != nil {
		return models.Area{}, err
	}

	return a, nil
}

// DeleteArea deletes an area from the database by its ID.
func DeleteArea(id int) error {
	query := `
		DELETE FROM "Areas"
		WHERE "Id" = $1
	`

	res, err := Pool.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Area not found")
	}

	return nil
}
