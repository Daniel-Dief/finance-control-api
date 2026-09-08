package models

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}
