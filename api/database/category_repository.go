package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
)

// ListCategories retrieves a list of categories from the database, optionally filtered by name.
func ListCategories(name *string) ([]models.Category, error) {
	query := `
		SELECT "Id", "Name"
		FROM "Categories"
		WHERE ($1::text IS NULL OR "Name" ILIKE '%' || $1 || '%')
	`

	rows, err := Pool.QueryContext(context.Background(), query, name)
	if err != nil {
		return nil, err
	} else if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var result []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

// GetCategoryByID retrieves a category from the database by its ID.
func GetCategoryByID(id int) (models.Category, error) {
	query := `
		SELECT "Id", "Name"
		FROM "Categories"
		WHERE "Id" = $1
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, id).Scan(&c.ID, &c.Name)
	if err != nil {
		return models.Category{}, err
	}

	return c, nil
}

// CreateCategory inserts a new category into the database and returns the created category.
func CreateCategory(name string) (models.Category, error) {
	query := `
		INSERT INTO "Categories" ("Name")
		VALUES ($1)
		RETURNING "Id", "Name"
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, name).Scan(&c.ID, &c.Name)
	if err != nil {
		return models.Category{}, err
	}

	return c, nil
}

// UpdateCategory updates the name of an existing category in the database and returns the updated category.
func UpdateCategory(id int, name string) (models.Category, error) {
	query := `
		UPDATE "Categories"
		SET "Name" = $2
		WHERE "Id" = $1
		RETURNING "Id", "Name"
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, id, name).Scan(&c.ID, &c.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, errors.New("Category not found")
	} else if err != nil {
		return models.Category{}, err
	}

	return c, nil
}

// DeleteCategory removes a category from the database by its ID.
func DeleteCategory(id int) error {
	query := `
		DELETE FROM "Categories"
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
		return errors.New("Category not found")
	}

	return nil
}
