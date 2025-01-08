package models

type Order struct {
    ID         int     `json:"id"`
    UserID     int     `json:"user_id"`
    ProductIDs []int   `json:"product_ids"`
    Status     string  `json:"status"`
    TotalPrice float64 `json:"total_price"`
}