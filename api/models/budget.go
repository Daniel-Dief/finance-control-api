package models

type Budget struct {
	ID     int `json:"id"`
	Year   int `json:"year"`
	Month  int `json:"month"`
	AreaID int `json:"area_id"`
	Amount int `json:"amount"`
}
