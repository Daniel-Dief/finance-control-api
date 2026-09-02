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

// ListTransactions retrieves a paginated list of transactions from the database
// based on the provided filters.
func ListTransactions(filters TransactionFilters, pag Pagination) (*PaginatedResult[models.Transaction], error) {
	where := `
		WHERE ($1::text IS NULL OR "Type" = $1)
		AND ($2::integer IS NULL OR "CategoryId" = $2)
		AND ($3::integer IS NULL OR "AreaId" = $3)
		AND ($4::date IS NULL OR "Date" >= $4)
		AND ($5::date IS NULL OR "Date" <= $5)
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "Transactions" `+where,
		filters.Type, filters.CategoryID, filters.AreaID, filters.From, filters.To).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
		FROM "Transactions"
		` + where + `
		ORDER BY "Id"
		LIMIT $6 OFFSET $7
	`

	rows, err := Pool.QueryContext(context.Background(), query,
		filters.Type, filters.CategoryID, filters.AreaID, filters.From, filters.To, pag.Limit, pag.Offset())
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, ErrGenericDatabase
	}

	defer rows.Close()
	if rows.Err() != nil {
		log.Println("Error with rows:", rows.Err())
		return nil, ErrProcessQuery
	}

	result := make([]models.Transaction, 0)
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, t)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
}

// GetTransactionByID retrieves a single transaction from the database based on its ID.
func GetTransactionByID(id int) (models.Transaction, error) {
	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
		FROM "Transactions"
		WHERE "Id" = $1
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(), query, id).Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Transaction{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, ErrGenericDatabase
	}

	return t, nil
}

// CreateTransaction inserts a new transaction into the database and returns the created transaction.
func CreateTransaction(props TransactionProps) (models.Transaction, error) {
	query := `
		INSERT INTO "Transactions" ("Date", "Amount", "CategoryId", "AreaId", "Type")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(),
		query, props.Date, props.Amount, props.CategoryID, props.AreaID, props.Type).
		Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type)
	if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, ErrRegisterObject("transação")
	}

	return t, nil
}

// UpdateTransaction updates an existing transaction in the database based on its ID and the provided properties.
func UpdateTransaction(id int, props TransactionProps) (models.Transaction, error) {
	query := `
		UPDATE "Transactions"
		SET "Date" = COALESCE($1, "Date"),
			"Amount" = COALESCE($2, "Amount"),
			"CategoryId" = COALESCE($3, "CategoryId"),
			"AreaId" = COALESCE($4, "AreaId"),
			"Type" = COALESCE($5, "Type")
		WHERE "Id" = $6
		RETURNING "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(), query, props.Date, props.Amount, props.CategoryID, props.AreaID, props.Type, id).
		Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Transaction{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, ErrUpdateObject("transação")
	}

	return t, nil
}

// DeleteTransaction removes a transaction from the database based on its ID.
func DeleteTransaction(id int) error {
	query := `
		DELETE FROM "Transactions"
		WHERE "Id" = $1
	`

	res, err := Pool.ExecContext(context.Background(), query, id)
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
