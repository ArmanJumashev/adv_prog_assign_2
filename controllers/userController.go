package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

type User struct {
    Email    string `json:"email"`
    Username string `json:"username"`
}

func GetUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT email, username FROM users")
        if err != nil {
            http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []User
        for rows.Next() {
            var user User
            if err := rows.Scan(&user.Email, &user.Username); err != nil {
                http.Error(w, "Failed to parse user", http.StatusInternalServerError)
                return
            }
            users = append(users, user)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}
