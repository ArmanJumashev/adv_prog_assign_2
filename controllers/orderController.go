package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
)

type Order struct {
    ID        int       `json:"id"`
    UserEmail string    `json:"user_email"`
    ProductID int       `json:"product_id"`
    Quantity  int       `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
}

func GetOrders(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT id, user_email, product_id, quantity, created_at FROM orders")
        if err != nil {
            http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var orders []Order
        for rows.Next() {
            var order Order
            if err := rows.Scan(&order.ID, &order.UserEmail, &order.ProductID, &order.Quantity, &order.CreatedAt); err != nil {
                http.Error(w, "Failed to parse order", http.StatusInternalServerError)
                return
            }
            orders = append(orders, order)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(orders)
    }
}
