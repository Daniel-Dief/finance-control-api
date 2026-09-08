package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"
)

// ListAreas retrieves a paginated list of areas belonging to the given user,
// optionally filtered by name. If name is empty, all areas are returned.
func ListAreas(userID int, name string, pag Pagination) (*PaginatedResult[models.Area], error) {
	where := `
		WHERE "StatusId" = 1 AND "UserId" = $1
		AND ($2::text IS NULL OR "Name" ILIKE '%' || $2 || '%')
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "Areas" `+where, userID, name).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Name", "UserId"
		FROM "Areas"
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

	result := make([]models.Area, 0)
	for rows.Next() {
		var a models.Area
		if err := rows.Scan(&a.ID, &a.Name, &a.UserID); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, a)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
}

// GetAreaByID retrieves an area by its ID for the given user.
func GetAreaByID(userID, id int) (models.Area, error) {
	query := `
		SELECT "Id", "Name", "UserId"
		FROM "Areas"
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = 1
	`

	var a models.Area
	err := Pool.QueryRowContext(context.Background(), query, id, userID).Scan(&a.ID, &a.Name, &a.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Area{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Area{}, ErrGenericDatabase
	}

	return a, nil
}

// CreateArea inserts a new area for the given user and returns the created area.
func CreateArea(userID int, name string) (models.Area, error) {
	query := `
		INSERT INTO "Areas" ("Name", "UserId")
		VALUES ($1, $2)
		RETURNING "Id", "Name", "UserId"
	`

	var a models.Area
	err := Pool.QueryRowContext(context.Background(), query, name, userID).Scan(&a.ID, &a.Name, &a.UserID)
	if isUniqueViolation(err) {
		return models.Area{}, ErrDuplicateArea
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Area{}, ErrRegisterObject("área")
	}

	return a, nil
}

// UpdateArea updates the name of an existing area belonging to the given user.
func UpdateArea(userID, id int, name string) (models.Area, error) {
	query := `
		UPDATE "Areas"
		SET "Name" = $3
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = 1
		RETURNING "Id", "Name", "UserId"
	`

	var a models.Area
	err := Pool.QueryRowContext(context.Background(), query, id, userID, name).Scan(&a.ID, &a.Name, &a.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Area{}, ErrNotFound
	} else if isUniqueViolation(err) {
		return models.Area{}, ErrDuplicateArea
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Area{}, ErrUpdateObject("área")
	}

	return a, nil
}

// DeleteArea soft-deletes an area belonging to the given user by setting its
// StatusId to 3 (deleted).
func DeleteArea(userID, id int) error {
	query := `
		UPDATE "Areas"
		SET "StatusId" = 3
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = 1
	`

	res, err := Pool.ExecContext(context.Background(), query, id, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("área")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("área")
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
