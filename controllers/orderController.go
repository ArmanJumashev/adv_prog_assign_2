package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
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

type OrderRequest struct {
	User       User        `json:"user"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"totalPrice"`
}

type User struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
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
		// Only allow POST requests.
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body.
		var orderReq OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the request data.
		if orderReq.User.Email == "" || len(orderReq.Items) == 0 {
			http.Error(w, "Invalid order data", http.StatusBadRequest)
			return
		}

		// Save the order to the database.
		for _, item := range orderReq.Items {
			_, err := db.Exec(`
				INSERT INTO orders (user_email, product_id, quantity, created_at)
				VALUES ($1, $2, $3, $4)
			`, orderReq.User.Email, item.ProductID, item.Quantity, time.Now())
			if err != nil {
				log.Printf("Error saving order to database: %v", err)
				http.Error(w, "Failed to save order", http.StatusInternalServerError)
				return
			}
		}

		// Respond with a success message.
		response := map[string]string{"message": "Order placed successfully!"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
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
