package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"
)

type TransactionProps struct {
	Date       *string
	Amount     *int
	CategoryID *int
	AreaID     *int
	Type       *string
}

type TransactionFilters struct {
	Type       *string
	CategoryID *int
	AreaID     *int
	From       *string
	To         *string
}

// ListTransactions retrieves a paginated list of transactions for the given
// user based on the provided filters.
func ListTransactions(userID int, filters TransactionFilters, pag Pagination) (*PaginatedResult[models.Transaction], error) {
	where := `
		WHERE "StatusId" = $1 AND "UserId" = $2
		AND ($3::text IS NULL OR "Type" = $3)
		AND ($4::integer IS NULL OR "CategoryId" = $4)
		AND ($5::integer IS NULL OR "AreaId" = $5)
		AND ($6::date IS NULL OR "Date" >= $6)
		AND ($7::date IS NULL OR "Date" <= $7)
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "Transactions" `+where,
		StatusActive, userID, filters.Type, filters.CategoryID, filters.AreaID, filters.From, filters.To).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type", "UserId"
		FROM "Transactions"
		` + where + `
		ORDER BY "Id"
		LIMIT $8 OFFSET $9
	`

	rows, err := Pool.QueryContext(context.Background(), query,
		StatusActive, userID, filters.Type, filters.CategoryID, filters.AreaID, filters.From, filters.To, pag.Limit, pag.Offset())
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, ErrGenericDatabase
	}

	if rows.Err() != nil {
		log.Println("Error with rows:", rows.Err())
		return nil, ErrProcessQuery
	}
	defer rows.Close()

	result := make([]models.Transaction, 0)
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type, &t.UserID); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, t)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
}

// GetTransactionByID retrieves a single transaction for the given user based on its ID.
func GetTransactionByID(userID, id int) (models.Transaction, error) {
	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type", "UserId"
		FROM "Transactions"
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = $3
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(), query, id, userID, StatusActive).Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type, &t.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Transaction{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, ErrGenericDatabase
	}

	return t, nil
}

// CreateTransaction inserts a new transaction for the given user and returns the created transaction.
func CreateTransaction(userID int, props TransactionProps) (models.Transaction, error) {
	query := `
		INSERT INTO "Transactions" ("Date", "Amount", "CategoryId", "AreaId", "Type", "UserId")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "Id", "Date", "Amount", "CategoryId", "AreaId", "Type", "UserId"
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(),
		query, props.Date, props.Amount, props.CategoryID, props.AreaID, props.Type, userID).
		Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type, &t.UserID)
	if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, ErrRegisterObject("transação")
	}

	return t, nil
}

// UpdateTransaction updates an existing transaction belonging to the given user
// based on its ID and the provided properties.
func UpdateTransaction(userID, id int, props TransactionProps) (models.Transaction, error) {
	query := `
		UPDATE "Transactions"
		SET "Date" = COALESCE($1, "Date"),
			"Amount" = COALESCE($2, "Amount"),
			"CategoryId" = COALESCE($3, "CategoryId"),
			"AreaId" = COALESCE($4, "AreaId"),
			"Type" = COALESCE($5, "Type")
		WHERE "Id" = $6 AND "UserId" = $7 AND "StatusId" = $8
		RETURNING "Id", "Date", "Amount", "CategoryId", "AreaId", "Type", "UserId"
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(), query, props.Date, props.Amount, props.CategoryID, props.AreaID, props.Type, id, userID, StatusActive).
		Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type, &t.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Transaction{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, ErrUpdateObject("transação")
	}

	return t, nil
}

// DeleteTransaction soft-deletes a transaction belonging to the given user by
// setting its StatusId to 3 (deleted).
func DeleteTransaction(userID, id int) error {
	query := `
		UPDATE "Transactions"
		SET "StatusId" = $3
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = $4
	`

	res, err := Pool.ExecContext(context.Background(), query, id, userID, StatusDeleted, StatusActive)
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("transação")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("transação")
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
