package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"online-shop/models"
	"time"
)

func Register(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
		return
	}

	// Парсим дату рождения в тип time.Time
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

// Авторизация пользователя
func Login(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
            return
        }

        // Проверяем, существует ли пользователь с таким email и правильным паролем
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

        // Проверяем пароль
        if user.Password != dbUser.Password {
            http.Error(w, "Неверный пароль", http.StatusUnauthorized)
            return
        }

        // Возвращаем данные пользователя
        response := struct {
            FullName    string `json:"full_name"`
            Email       string `json:"email"`
            DateOfBirth string `json:"date_of_birth"`
        }{
            FullName:    dbUser.FullName,
            Email:       dbUser.Email,
            DateOfBirth: dbUser.DateOfBirth,
        }

        // Устанавливаем статус ответа и кодируем данные в JSON
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(response)
    }
}

