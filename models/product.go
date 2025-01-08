package models

type Product struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Price         float64   `json:"price" db:"price"`
	Category      string    `json:"category" db:"category"`
	Description   string    `json:"description" db:"description"`
	Image         string    `json:"image" db:"image"`
}
