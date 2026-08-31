package database

import (
	"context"
	"errors"
	"finance-control-api/models"
)

type TransactionProps struct {
	Date       *string
	Amount     *int64
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

func ListTransactions(filters TransactionFilters) ([]models.Transaction, error) {
	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
		FROM "Transactions"
		WHERE ($1::text IS NULL OR "Type" = $1)
		AND ($2::integer IS NULL OR "CategoryId" = $2)
		AND ($3::integer IS NULL OR "AreaId" = $3)
		AND ($4::text IS NULL OR "Date" >= $4)
		AND ($5::text IS NULL OR "Date" <= $5)
	`

	rows, err := Pool.QueryContext(context.Background(), query, filters.Type, filters.CategoryID, filters.AreaID, filters.From, filters.To)
	if err != nil {
		return nil, err
	} else if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var result []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func GetTransactionByID(id int) (models.Transaction, error) {
	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
		FROM "Transactions"
		WHERE "Id" = $1
	`

	var t models.Transaction
	err := Pool.QueryRowContext(context.Background(), query, id).Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type)
	if err != nil {
		return models.Transaction{}, err
	}

	return t, nil
}

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
		return models.Transaction{}, err
	}

	return t, nil
}

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
	if err != nil {
		return models.Transaction{}, err
	}

	return t, nil
}

func DeleteTransaction(id int) error {
	query := `
		DELETE FROM "Transactions"
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
		return errors.New("Transaction not found")
	}

	return nil
}
