package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"
)

// ListCategories retrieves a paginated list of categories belonging to the
// given user, optionally filtered by name. If name is empty, all categories
// are returned.
func ListCategories(userID int, name *string, pag Pagination) (*PaginatedResult[models.Category], error) {
	where := `
		WHERE "StatusId" = 1 AND "UserId" = $1
		AND ($2::text IS NULL OR "Name" ILIKE '%' || $2 || '%')
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "Categories" `+where, userID, name).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Name", "UserId"
		FROM "Categories"
		` + where + `
		ORDER BY "Id"
		LIMIT $3 OFFSET $4
	`

	rows, err := Pool.QueryContext(context.Background(), query, userID, name, pag.Limit, pag.Offset())
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
		if err := rows.Scan(&c.ID, &c.Name, &c.UserID); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, c)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
}

// GetCategoryByID retrieves a category by its ID for the given user.
func GetCategoryByID(userID, id int) (models.Category, error) {
	query := `
		SELECT "Id", "Name", "UserId"
		FROM "Categories"
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = 1
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, id, userID).Scan(&c.ID, &c.Name, &c.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, ErrGenericDatabase
	}

	return c, nil
}

// CreateCategory inserts a new category for the given user and returns the created category.
func CreateCategory(userID int, name string) (models.Category, error) {
	query := `
		INSERT INTO "Categories" ("Name", "UserId")
		VALUES ($1, $2)
		RETURNING "Id", "Name", "UserId"
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, name, userID).Scan(&c.ID, &c.Name, &c.UserID)
	if isUniqueViolation(err) {
		return models.Category{}, ErrDuplicateCategory
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, ErrRegisterObject("categoria")
	}

	return c, nil
}

// UpdateCategory updates the name of an existing category belonging to the
// given user and returns the updated category.
func UpdateCategory(userID, id int, name string) (models.Category, error) {
	query := `
		UPDATE "Categories"
		SET "Name" = $3
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = 1
		RETURNING "Id", "Name", "UserId"
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, id, userID, name).Scan(&c.ID, &c.Name, &c.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, ErrNotFound
	} else if isUniqueViolation(err) {
		return models.Category{}, ErrDuplicateCategory
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, ErrUpdateObject("categoria")
	}

	return c, nil
}

// DeleteCategory soft-deletes a category belonging to the given user by
// setting its StatusId to 3 (deleted).
func DeleteCategory(userID, id int) error {
	query := `
		UPDATE "Categories"
		SET "StatusId" = 3
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = 1
	`

	res, err := Pool.ExecContext(context.Background(), query, id, userID)
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
