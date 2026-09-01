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

// ListTransactions retrieves transactions from the database based on the provided filters.
func ListTransactions(filters TransactionFilters) ([]models.Transaction, error) {
	query := `
		SELECT "Id", "Date", "Amount", "CategoryId", "AreaId", "Type"
		FROM "Transactions"
		WHERE ($1::text IS NULL OR "Type" = $1)
		AND ($2::integer IS NULL OR "CategoryId" = $2)
		AND ($3::integer IS NULL OR "AreaId" = $3)
		AND ($4::date IS NULL OR "Date" >= $4)
		AND ($5::date IS NULL OR "Date" <= $5)
	`

	rows, err := Pool.QueryContext(context.Background(), query, filters.Type, filters.CategoryID, filters.AreaID, filters.From, filters.To)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, errors.New("Falha ao consultar o banco de dados, em caso de persistencia contatar o suporte.")
	} else if rows.Err() != nil {
		log.Println("Error with rows:", rows.Err())
		return nil, errors.New("Falha ao processar os resultados da consulta, em caso de persistencia contatar o suporte.")
	}
	defer rows.Close()

	var result []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.CategoryID, &t.AreaID, &t.Type); err != nil {
			log.Println("Error scanning row:", err)
			return nil, errors.New("Falha indexar os resultados, em caso de persistencia contatar o suporte.")
		}
		result = append(result, t)
	}

	return result, nil
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
		return models.Transaction{}, nil
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, errors.New("Falha ao consultar o banco de dados, em caso de persistencia contatar o suporte.")
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
		return models.Transaction{}, errors.New("Falha ao registrar a transação no banco de dados, em caso de persistencia contatar o suporte.")
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
		return models.Transaction{}, errors.New("Falha ao atualizar a transação: transação não encontrada")
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Transaction{}, errors.New("Falha ao atualizar a transação no banco de dados, em caso de persistencia contatar o suporte.")
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
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error executing query:", err)
		return errors.New("Falha ao deletar a transação no banco de dados, em caso de persistencia contatar o suporte.")
	}

	if rowsAffected == 0 {
		return errors.New("Transação não encontrada")
	}

	return nil
}
