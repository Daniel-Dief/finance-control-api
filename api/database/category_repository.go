package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"
)

// ListCategories retrieves a paginated list of categories from the database,
// optionally filtered by name. If name is empty, all categories are returned.
func ListCategories(name *string, pag Pagination) (*PaginatedResult[models.Category], error) {
	where := `
		WHERE ($1::text IS NULL OR "Name" ILIKE '%' || $1 || '%')
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "Categories" `+where, name).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Name"
		FROM "Categories"
		` + where + `
		ORDER BY "Id"
		LIMIT $2 OFFSET $3
	`

	rows, err := Pool.QueryContext(context.Background(), query, name, pag.Limit, pag.Offset())
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, ErrGenericDatabase
	}

	defer rows.Close()
	if rows.Err() != nil {
		log.Println("Error with rows:", rows.Err())
		return nil, ErrProcessQuery
	}

	result := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, c)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
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

	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, ErrGenericDatabase
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
		log.Println("Error executing query:", err)
		return models.Category{}, ErrRegisterObject("categoria")
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
		return models.Category{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, ErrUpdateObject("categoria")
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
		log.Println("Error executing query:", err)
		return ErrDeleteObject("categoria")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("categoria")
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
