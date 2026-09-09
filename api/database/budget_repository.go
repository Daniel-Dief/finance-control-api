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

// ListBudgets returns a paginated list of budgets for the given user. You can
// filter by year, month and/or areaId.
func ListBudgets(userID int, props BudgetProps, pag Pagination) (*PaginatedResult[models.Budget], error) {
	where := `
		WHERE "StatusId" = $1 AND "UserId" = $2
		AND ($3::integer IS NULL OR "Year" = $3)
		AND ($4::integer IS NULL OR "Month" = $4)
		AND ($5::integer IS NULL OR "AreaId" = $5)
	`

	var total int
	if err := Pool.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "MonthlyBudget" `+where, StatusActive, userID, props.Year, props.Month, props.AreaId).Scan(&total); err != nil {
		log.Println("Error executing count query:", err)
		return nil, ErrGenericDatabase
	}

	query := `
		SELECT "Id", "Year", "Month", "AreaId", "Amount", "UserId"
		FROM "MonthlyBudget"
		` + where + `
		ORDER BY "Id"
		LIMIT $6 OFFSET $7
	`

	rows, err := Pool.QueryContext(context.Background(), query, StatusActive, userID, props.Year, props.Month, props.AreaId, pag.Limit, pag.Offset())
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
		if err := rows.Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount, &b.UserID); err != nil {
			log.Println("Error scanning row:", err)
			return nil, ErrBindQuery
		}
		result = append(result, b)
	}

	paginated := newPaginatedResult(result, pag, total)
	return &paginated, nil
}

// GetBudgetByID retrieves a budget by its ID for the given user.
func GetBudgetByID(userID, id int) (models.Budget, error) {
	query := `
		SELECT "Id", "Year", "Month", "AreaId", "Amount", "UserId"
		FROM "MonthlyBudget"
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = $3
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, id, userID, StatusActive).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount, &b.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Budget{}, ErrNotFound
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Budget{}, ErrGenericDatabase
	}

	return b, nil
}

// CreateBudget inserts a new budget for the given user and returns the created budget.
func CreateBudget(userID int, budget BudgetProps) (models.Budget, error) {
	query := `
		INSERT INTO "MonthlyBudget" ("Year", "Month", "AreaId", "Amount", "UserId")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "Id", "Year", "Month", "AreaId", "Amount", "UserId"
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, budget.Year, budget.Month, budget.AreaId, budget.Amount, userID).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount, &b.UserID)
	if isUniqueViolation(err) {
		return models.Budget{}, ErrDuplicateBudget
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Budget{}, ErrRegisterObject("orçamento")
	}

	return b, nil
}

// UpdateBudget updates an existing budget belonging to the given user and
// returns the updated budget.
func UpdateBudget(userID, id int, budget BudgetProps) (models.Budget, error) {
	query := `
		UPDATE "MonthlyBudget"
		SET "Year" = COALESCE($1, "Year"),
			"Month" = COALESCE($2, "Month"),
			"AreaId" = COALESCE($3, "AreaId"),
			"Amount" = COALESCE($4, "Amount")
		WHERE "Id" = $5 AND "UserId" = $6 AND "StatusId" = $7
		RETURNING "Id", "Year", "Month", "AreaId", "Amount", "UserId"
	`

	var b models.Budget
	err := Pool.QueryRowContext(context.Background(), query, budget.Year, budget.Month, budget.AreaId, budget.Amount, id, userID, StatusActive).Scan(&b.ID, &b.Year, &b.Month, &b.AreaID, &b.Amount, &b.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Budget{}, ErrNotFound
	} else if isUniqueViolation(err) {
		return models.Budget{}, ErrDuplicateBudget
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Budget{}, ErrUpdateObject("orçamento")
	}

	return b, nil
}

// DeleteBudget soft-deletes a budget belonging to the given user by setting
// its StatusId to 3 (deleted).
func DeleteBudget(userID, id int) error {
	query := `
		UPDATE "MonthlyBudget"
		SET "StatusId" = $3
		WHERE "Id" = $1 AND "UserId" = $2 AND "StatusId" = $4
	`

	res, err := Pool.ExecContext(context.Background(), query, id, userID, StatusDeleted, StatusActive)
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
