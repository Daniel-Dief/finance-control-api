package database

import (
	"context"
	"finance-control-api/models"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getCurrentYearMonth() (int, int) {
	now := time.Now()
	return now.Year(), int(now.Month())
}

func getDaysInMonth(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func computeMedianInt(values []int64) int {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]int64, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	n := len(sorted)
	if n%2 == 1 {
		return int(math.Round(float64(sorted[n/2])))
	}
	mid1 := sorted[n/2-1]
	mid2 := sorted[n/2]
	avg := float64(mid1+mid2) / 2.0
	return int(math.Round(avg))
}

func GetKPIs(userID, year, month int) (models.ReportKPIs, error) {
	if year == 0 || month == 0 {
		y, m := getCurrentYearMonth()
		year, month = y, m
	}

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN "Type" = 'income' THEN "Amount" ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN "Type" = 'expense' THEN "Amount" ELSE 0 END), 0) as expense
		FROM "Transactions"
		WHERE "StatusId" = $1 AND "UserId" = $2
		AND EXTRACT(YEAR FROM "Date") = $3
		AND EXTRACT(MONTH FROM "Date") = $4
	`

	var income, expense int64
	err := Pool.QueryRowContext(context.Background(), query, StatusActive, userID, year, month).Scan(&income, &expense)
	if err != nil {
		log.Println("Error executing kpi query:", err)
		return models.ReportKPIs{}, ErrGenericDatabase
	}

	return models.ReportKPIs{
		Year:    year,
		Month:   month,
		Income:  int(income),
		Expense: int(expense),
		Balance: int(income - expense),
	}, nil
}

func GetLineChart(userID, year, month int, txType string) (models.LineChart, error) {
	if year == 0 || month == 0 {
		y, m := getCurrentYearMonth()
		year, month = y, m
	}
	if txType == "" {
		txType = "expense"
	}

	daysInMonth := getDaysInMonth(year, month)
	points := make([]models.LineChartPoint, 0, daysInMonth)

	for day := 1; day <= daysInMonth; day++ {
		dateStr := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

		// Get current month value for this day
		var currentValue int64
		queryCurrent := `
			SELECT COALESCE(SUM("Amount"), 0)
			FROM "Transactions"
			WHERE "StatusId" = $1 AND "UserId" = $2 AND "Type" = $3
			AND EXTRACT(YEAR FROM "Date") = $4
			AND EXTRACT(MONTH FROM "Date") = $5
			AND EXTRACT(DAY FROM "Date") = $6
		`
		err := Pool.QueryRowContext(context.Background(), queryCurrent, StatusActive, userID, txType, year, month, day).Scan(&currentValue)
		if err != nil {
			log.Println("Error executing line chart current query:", err)
			return models.LineChart{}, ErrGenericDatabase
		}

		// Get median from previous 6 months
		md := getDayMedian(userID, txType, year, month, day)

		points = append(points, models.LineChartPoint{
			Day:    day,
			Date:   dateStr,
			Value:  int(currentValue),
			Median: md,
		})
	}

	return models.LineChart{
		Year:          year,
		Month:         month,
		Type:          txType,
		HistoryMonths: 6,
		Points:        points,
	}, nil
}

func getDayMedian(userID int, txType string, year, month, day int) int {
	// Look at previous 6 months
	values := make([]int64, 0, 6)
	for i := 1; i <= 6; i++ {
		prevTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		pYear := prevTime.Year()
		pMonth := int(prevTime.Month())

		// Check if month had any movement of this type
		var monthTotal int64
		queryMonthTotal := `
			SELECT COALESCE(SUM("Amount"), 0)
			FROM "Transactions"
			WHERE "StatusId" = $1 AND "UserId" = $2 AND "Type" = $3
			AND EXTRACT(YEAR FROM "Date") = $4
			AND EXTRACT(MONTH FROM "Date") = $5
		`
		err := Pool.QueryRowContext(context.Background(), queryMonthTotal, StatusActive, userID, txType, pYear, pMonth).Scan(&monthTotal)
		if err != nil {
			continue
		}
		if monthTotal == 0 {
			continue // ignore months with no movement
		}

		// Get total for same day in this month
		var dayTotal int64
		queryDayTotal := `
			SELECT COALESCE(SUM("Amount"), 0)
			FROM "Transactions"
			WHERE "StatusId" = $1 AND "UserId" = $2 AND "Type" = $3
			AND EXTRACT(YEAR FROM "Date") = $4
			AND EXTRACT(MONTH FROM "Date") = $5
			AND EXTRACT(DAY FROM "Date") = $6
		`
		err = Pool.QueryRowContext(context.Background(), queryDayTotal, StatusActive, userID, txType, pYear, pMonth, day).Scan(&dayTotal)
		if err != nil {
			continue
		}
		// If month had movement but this day had none, count as 0
		values = append(values, dayTotal)
	}

	if len(values) == 0 {
		return 0
	}
	return computeMedianInt(values)
}

func GetRadarChart(userID, year, month int, txType string, areaIDs []int) (models.RadarChart, error) {
	if year == 0 || month == 0 {
		y, m := getCurrentYearMonth()
		year, month = y, m
	}
	if txType == "" {
		txType = "expense"
	}

	// Get all active areas for user
	areas := make([]models.RadarArea, 0)
	
	whereArea := `WHERE "StatusId" = $1 AND "UserId" = $2`
	args := []interface{}{StatusActive, userID}
	argPos := 3
	
	if len(areaIDs) > 0 {
		placeholders := make([]string, len(areaIDs))
		for i, aid := range areaIDs {
			placeholders[i] = "$" + strconv.Itoa(argPos)
			args = append(args, aid)
			argPos++
		}
		whereArea += ` AND "Id" IN (` + strings.Join(placeholders, ",") + `)`
	}
	
	queryAreas := `SELECT "Id", "Name" FROM "Areas" ` + whereArea + ` ORDER BY "Name"`
	rows, err := Pool.QueryContext(context.Background(), queryAreas, args...)
	if err != nil {
		log.Println("Error executing radar areas query:", err)
		return models.RadarChart{}, ErrGenericDatabase
	}
	defer rows.Close()
	
	type areaInfo struct {
		id   int
		name string
	}
	areaList := make([]areaInfo, 0)
	for rows.Next() {
		var a areaInfo
		if err := rows.Scan(&a.id, &a.name); err != nil {
			log.Println("Error scanning area:", err)
			return models.RadarChart{}, ErrBindQuery
		}
		areaList = append(areaList, a)
	}
	if rows.Err() != nil {
		log.Println("Error with area rows:", rows.Err())
		return models.RadarChart{}, ErrProcessQuery
	}

	for _, a := range areaList {
		// Get current month value for this area
		var currentValue int64
		queryCurrent := `
			SELECT COALESCE(SUM("Amount"), 0)
			FROM "Transactions"
			WHERE "StatusId" = $1 AND "UserId" = $2 AND "Type" = $3 AND "AreaId" = $4
			AND EXTRACT(YEAR FROM "Date") = $5
			AND EXTRACT(MONTH FROM "Date") = $6
		`
		err := Pool.QueryRowContext(context.Background(), queryCurrent, StatusActive, userID, txType, a.id, year, month).Scan(&currentValue)
		if err != nil {
			log.Println("Error executing radar current query:", err)
			return models.RadarChart{}, ErrGenericDatabase
		}

		md := getAreaMedian(userID, txType, a.id, year, month)

		areas = append(areas, models.RadarArea{
			AreaID:   a.id,
			AreaName: a.name,
			Value:    int(currentValue),
			Median: md,
		})
	}

	return models.RadarChart{
		Year:          year,
		Month:         month,
		Type:          txType,
		HistoryMonths: 6,
		Areas:         areas,
	}, nil
}

func getAreaMedian(userID int, txType string, areaID, year, month int) int {
	values := make([]int64, 0, 6)
	for i := 1; i <= 6; i++ {
		prevTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		pYear := prevTime.Year()
		pMonth := int(prevTime.Month())

		// Check if area had movement this month
		var monthTotal int64
		queryMonthTotal := `
			SELECT COALESCE(SUM("Amount"), 0)
			FROM "Transactions"
			WHERE "StatusId" = $1 AND "UserId" = $2 AND "Type" = $3 AND "AreaId" = $4
			AND EXTRACT(YEAR FROM "Date") = $5
			AND EXTRACT(MONTH FROM "Date") = $6
		`
		err := Pool.QueryRowContext(context.Background(), queryMonthTotal, StatusActive, userID, txType, areaID, pYear, pMonth).Scan(&monthTotal)
		if err != nil {
			continue
		}
		if monthTotal == 0 {
			continue
		}
		values = append(values, monthTotal)
	}

	if len(values) == 0 {
		return 0
	}
	return computeMedianInt(values)
}

func GetBudgetReport(userID, year, month int) (models.BudgetReport, error) {
	if year == 0 || month == 0 {
		y, m := getCurrentYearMonth()
		year, month = y, m
	}

	budgets := make([]models.BudgetProgress, 0)

	// Get budgets for this month where area is active
	query := `
		SELECT b."Id", b."AreaId", a."Name", b."Amount"
		FROM "MonthlyBudget" b
		JOIN "Areas" a ON a."Id" = b."AreaId"
		WHERE b."StatusId" = $1 AND b."UserId" = $2
		AND b."Year" = $3 AND b."Month" = $4
		AND a."StatusId" = $1 AND a."UserId" = $2
		ORDER BY a."Name"
	`
	rows, err := Pool.QueryContext(context.Background(), query, StatusActive, userID, year, month)
	if err != nil {
		log.Println("Error executing budget report query:", err)
		return models.BudgetReport{}, ErrGenericDatabase
	}
	defer rows.Close()

	for rows.Next() {
		var b models.BudgetProgress
		var amount int
		if err := rows.Scan(&b.ID, &b.AreaID, &b.AreaName, &amount); err != nil {
			log.Println("Error scanning budget:", err)
			return models.BudgetReport{}, ErrBindQuery
		}
		b.Amount = amount

		// Get spent (only expenses) for this area in this month
		var spent int64
		querySpent := `
			SELECT COALESCE(SUM("Amount"), 0)
			FROM "Transactions"
			WHERE "StatusId" = $1 AND "UserId" = $2 AND "Type" = 'expense' AND "AreaId" = $3
			AND EXTRACT(YEAR FROM "Date") = $4
			AND EXTRACT(MONTH FROM "Date") = $5
		`
		err := Pool.QueryRowContext(context.Background(), querySpent, StatusActive, userID, b.AreaID, year, month).Scan(&spent)
		if err != nil {
			log.Println("Error executing budget spent query:", err)
			return models.BudgetReport{}, ErrGenericDatabase
		}
		b.Spent = int(spent)

		if amount == 0 {
			if b.Spent == 0 {
				b.Remaining = 0
				b.Percentage = 0
				b.Exceeded = false
			} else {
				b.Remaining = -b.Spent
				b.Percentage = 100
				b.Exceeded = true
			}
		} else {
			b.Remaining = amount - b.Spent
			if amount > 0 {
				b.Percentage = math.Round((float64(b.Spent)/float64(amount))*100*100) / 100
			} else {
				b.Percentage = 0
			}
			b.Exceeded = b.Spent >= amount
		}

		budgets = append(budgets, b)
	}
	if rows.Err() != nil {
		log.Println("Error with budget rows:", rows.Err())
		return models.BudgetReport{}, ErrProcessQuery
	}

	return models.BudgetReport{
		Year:    year,
		Month:   month,
		Budgets: budgets,
	}, nil
}
