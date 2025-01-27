package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"online-shop/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("advanced")


func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
			return
		}

		_, err = time.Parse("2006-01-02", user.DateOfBirth)
		if err != nil {
			http.Error(w, "Неверный формат даты рождения", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(`
			INSERT INTO users (full_name, email, password, date_of_birth)
			VALUES ($1, $2, $3, $4)`, user.FullName, user.Email, user.Password, user.DateOfBirth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при добавлении пользователя: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("Пользователь успешно зарегистрирован")
	}
}

// Login and generate JWT token
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
			return
		}

		var dbUser models.User
		err = db.QueryRow(`
			SELECT full_name, email, password, date_of_birth
			FROM users WHERE email = $1`, user.Email).Scan(&dbUser.FullName, &dbUser.Email, &dbUser.Password, &dbUser.DateOfBirth)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			} else {
				http.Error(w, fmt.Sprintf("Ошибка при поиске пользователя: %v", err), http.StatusInternalServerError)
			}
			return
		}

		if user.Password != dbUser.Password {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": dbUser.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
			return
		}

		// Response with token
		response := struct {
			Token       string `json:"token"`
			FullName    string `json:"full_name"`
			Email       string `json:"email"`
			DateOfBirth string `json:"date_of_birth"`
		}{
			Token:       tokenString,
			FullName:    dbUser.FullName,
			Email:       dbUser.Email,
			DateOfBirth: dbUser.DateOfBirth,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

