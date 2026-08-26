package models

type MonthlyBudget struct {
	ID     int   `json:"id"`
	Year   int   `json:"year"`
	Month  int   `json:"month"`
	AreaID int   `json:"area_id"`
	Amount int64 `json:"amount"`
}
