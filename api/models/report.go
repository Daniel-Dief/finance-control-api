package models

type ReportKPIs struct {
	Year    int `json:"year"`
	Month   int `json:"month"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
	Balance int `json:"balance"`
}

type LineChartPoint struct {
	Day    int    `json:"day"`
	Date   string `json:"date"`
	Value  int    `json:"value"`
	Median int    `json:"median"`
}

type LineChart struct {
	Year          int               `json:"year"`
	Month         int               `json:"month"`
	Type          string            `json:"type"`
	HistoryMonths int               `json:"history_months"`
	Points        []LineChartPoint  `json:"points"`
}

type RadarArea struct {
	AreaID   int    `json:"area_id"`
	AreaName string `json:"area_name"`
	Value    int    `json:"value"`
	Median   int    `json:"median"`
}

type RadarChart struct {
	Year          int          `json:"year"`
	Month         int          `json:"month"`
	Type          string       `json:"type"`
	HistoryMonths int          `json:"history_months"`
	Areas         []RadarArea  `json:"areas"`
}

type BudgetProgress struct {
	ID        int     `json:"id"`
	AreaID    int     `json:"area_id"`
	AreaName  string  `json:"area_name"`
	Amount    int     `json:"amount"`
	Spent     int     `json:"spent"`
	Remaining int     `json:"remaining"`
	Percentage float64 `json:"percentage"`
	Exceeded  bool    `json:"exceeded"`
}

type BudgetReport struct {
	Year    int               `json:"year"`
	Month   int               `json:"month"`
	Budgets []BudgetProgress  `json:"budgets"`
}
