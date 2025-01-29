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

func OrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

//
// func CreateOrder(db *sql.DB) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         var order models.Order
//         err := json.NewDecoder(r.Body).Decode(&order)
//         if err != nil {
//             http.Error(w, "Invalid order data", http.StatusBadRequest)
//             return
//         }
//
//         query := `INSERT INTO orders (user_id, total_price, shipping_address, status) VALUES ($1, $2, $3, $4) RETURNING id`
//         var orderId int
//         err = db.QueryRow(query, order.User.ID, order.TotalPrice, order.ShippingAddress, "pending").Scan(&orderId)
//         if err != nil {
//             http.Error(w, "Failed to create order", http.StatusInternalServerError)
//             return
//         }
//
//         for _, item := range order.Items {
//             query := `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`
//             _, err := db.Exec(query, orderId, item.ProductID, item.Quantity)
//             if err != nil {
//                 http.Error(w, "Failed to add order items", http.StatusInternalServerError)
//                 return
//             }
//         }
//
//         w.WriteHeader(http.StatusOK)
//         w.Write([]byte(fmt.Sprintf("Order %d has been successfully created.", orderId)))
//     }
// }
