package database

import (
	"context"
	"errors"
	"finance-control-api/models"
)

type BudgetProps struct {
	year   *int
	month  *int
	areaId *int
	amount *int64
}

// ListBudgets returns all budgets. You can filter by year, month and/or areaId.
func ListBudgets(props BudgetProps) ([]models.Budget, error) {
	query := `
		SELECT "Id", "Year", "Month", "AreaId", "Amount"'
		FROM "MonthlyBudget"
		WHERE ($1::integer IS NULL OR "Year" = $1)
		AND ($2::integer IS NULL OR "Month" = $2)
		AND ($3::integer IS NULL OR "AreaId" = $3)
	`

	rows, err := Pool.QueryContext(context.Background(), query, props.year, props.month, props.areaId)
	if err != nil {
		return nil, err
	} else if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var result []models.Budget
	for rows.Next() {
		var b models.Budget
		if err := rows.Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount); err != nil {
			return nil, err
		}
		result = append(result, b)
	}

	return result, nil
}

// GetBudgetByID retrieves a budget by its ID from the database.
func GetBudgetByID(id int) (models.Budget, error) {
	query := `
		SELECT "Id", "Year", "Month", "AreaId", "Amount"
		FROM "MonthlyBudget"
		WHERE "Id" = $1
	`

	rows, err := Pool.QueryContext(context.Background(), query, id)
	if err != nil {
		return models.Budget{}, err
	} else if rows.Err() != nil {
		return models.Budget{}, rows.Err()
	}
	defer rows.Close()

	if !rows.Next() {
		return models.Budget{}, nil
	}

	var b models.Budget
	if err := rows.Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount); err != nil {
		return models.Budget{}, err
	}

	return b, nil
}

// CreateBudget inserts a new budget into the database and returns the created budget.
func CreateBudget(budget BudgetProps) (models.Budget, error) {
	query := `
		INSERT INTO "MonthlyBudget" ("Year", "Month", "AreaId", "Amount")
		VALUES ($1, $2, $3, $4)
		RETURNING "Id"
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, budget.year, budget.month, budget.areaId, budget.amount).Scan(&b.ID)
	if err != nil {
		return models.Budget{}, err
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
	err := Pool.QueryRowContext(context.Background(), query, budget.year, budget.month, budget.areaId, budget.amount, id).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount)
	if err != nil {
		return models.Budget{}, err
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
