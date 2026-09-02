package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"
)

type BudgetProps struct {
	Year   *int
	Month  *int
	AreaId *int
	Amount *int
}

// ListBudgets returns a paginated list of budgets. You can filter by year,
// month and/or areaId.
func ListBudgets(props BudgetProps, pag Pagination) (*PaginatedResult[models.Budget], error) {
	where := `
		WHERE ($1::integer IS NULL OR "Year" = $1)
		AND ($2::integer IS NULL OR "Month" = $2)
		AND ($3::integer IS NULL OR "AreaId" = $3)
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "MonthlyBudget" `+where, props.Year, props.Month, props.AreaId).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Year", "Month", "AreaId", "Amount"
		FROM "MonthlyBudget"
		` + where + `
		ORDER BY "Id"
		LIMIT $4 OFFSET $5
	`

	rows, err := Pool.QueryContext(context.Background(), query, props.Year, props.Month, props.AreaId, pag.Limit, pag.Offset())
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, ErrGenericDatabase
	}

	defer rows.Close()
	if rows.Err() != nil {
		log.Println("Error with rows:", rows.Err())
		return nil, ErrProcessQuery
	}

	result := make([]models.Budget, 0)
	for rows.Next() {
		var b models.Budget
		if err := rows.Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, b)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
}

// GetBudgetByID retrieves a budget by its ID from the database.
func GetBudgetByID(id int) (models.Budget, error) {
	query := `
		SELECT "Id", "Year", "Month", "AreaId", "Amount"
		FROM "MonthlyBudget"
		WHERE "Id" = $1
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, id).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Budget{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Budget{}, ErrGenericDatabase
	}

	return b, nil
}

// CreateBudget inserts a new budget into the database and returns the created budget.
func CreateBudget(budget BudgetProps) (models.Budget, error) {
	query := `
		INSERT INTO "MonthlyBudget" ("Year", "Month", "AreaId", "Amount")
		VALUES ($1, $2, $3, $4)
		RETURNING "Id", "Year", "Month", "AreaId", "Amount"
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, budget.Year, budget.Month, budget.AreaId, budget.Amount).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount)
	if err != nil {
		log.Println("Error executing query:", err)
		return models.Budget{}, ErrRegisterObject("orçamento")
	}

	return b, nil
}

// UpdateBudget updates an existing budget in the database and returns the updated budget.
func UpdateBudget(id int, budget BudgetProps) (models.Budget, error) {
	query := `
		UPDATE "MonthlyBudget"
		SET "Year" = COALESCE($1, "Year"),
			"Month" = COALESCE($2, "Month"),
			"AreaId" = COALESCE($3, "AreaId"),
			"Amount" = COALESCE($4, "Amount")
		WHERE "Id" = $5
		RETURNING "Id", "Year", "Month", "AreaId", "Amount"
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, budget.Year, budget.Month, budget.AreaId, budget.Amount, id).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Budget{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Budget{}, ErrUpdateObject("orçamento")
	}

	return b, nil
}

// DeleteBudget deletes a budget from the database by its ID.
func DeleteBudget(id int) error {
	query := `
		DELETE FROM "MonthlyBudget"
		WHERE "Id" = $1
	`

	res, err := Pool.ExecContext(context.Background(), query, id)
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("orçamento")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error executing query:", err)
		return ErrDeleteObject("orçamento")
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
