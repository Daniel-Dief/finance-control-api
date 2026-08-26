package models

type Transaction struct {
	ID         int    `json:"id"`
	Date       string `json:"date"`
	Amount     int64  `json:"amount"`
	CategoryID int    `json:"category_id"`
	AreaID     int    `json:"area_id"`
	Type       string `json:"type"` // 'income' or 'expense'
}
