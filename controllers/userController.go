package controllers

import (
//     "bytes"
	"database/sql"
	"encoding/json"
// 	"io"
	"log"
	"net/http"
    "strconv"
	"strings"

//     "net/smtp"
//     "mime/multipart"
//     "mime/quotedprintable"


    "online-shop/models"

// 	"gopkg.in/gomail.v2"
)

func GetUserProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		var user models.User
		err := db.QueryRow("SELECT full_name, email, date_of_birth FROM users WHERE id = $1", userID).
			Scan(&user.FullName, &user.Email, &user.DateOfBirth)
		if err != nil {
			http.Error(w, "Error fetching user details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}


func UpdateUserProfile(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user models.User

        log.Println("updating user ...")
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        query := "UPDATE users SET "
        params := []interface{}{}
        setClauses := []string{}




        // Добавляем только те поля, которые были переданы в запросе
        if user.FullName != "" {
            setClauses = append(setClauses, "full_name = $"+strconv.Itoa(len(params)+1))
            params = append(params, user.FullName)
        }
        if user.Email != "" {
            setClauses = append(setClauses, "email = $"+strconv.Itoa(len(params)+1))
            params = append(params, user.Email)
        }
        if user.Password != "" {
            setClauses = append(setClauses, "password = $"+strconv.Itoa(len(params)+1))
            params = append(params, user.Password)
        }
        if user.DateOfBirth != "" {
            setClauses = append(setClauses, "date_of_birth = $"+strconv.Itoa(len(params)+1))
            params = append(params, user.DateOfBirth)
        }

        // Если нет данных для обновления, то не выполняем запрос
        if len(setClauses) == 0 {
            http.Error(w, "No fields to update", http.StatusBadRequest)
            return
        }

        // Собираем запрос и добавляем условие для поиска пользователя
        query += strings.Join(setClauses, ", ") + " WHERE email = $"+strconv.Itoa(len(params)+1)
        params = append(params, user.Email)

        // Выполняем запрос
        _, err = db.Exec(query, params...)
        if err != nil {
            http.Error(w, "Error updating user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("User profile updated successfully"))
    }
}



func GetUserOrders(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, user_id, product_ids, status, total_price FROM orders WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Error fetching orders", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []models.Order
		for rows.Next() {
			var order models.Order
			var productIDs string
			err := rows.Scan(&order.ID, &order.UserID, &productIDs, &order.Status, &order.TotalPrice)
			if err != nil {
				http.Error(w, "Error reading orders", http.StatusInternalServerError)
				return
			}

			// Assuming `productIDs` is stored as a comma-separated string in the database
			order.ProductIDs = parseProductIDs(productIDs)
			orders = append(orders, order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

func parseProductIDs(ids string) []int {
	var productIDs []int
	for _, id := range strings.Split(ids, ",") {
		parsedID, err := strconv.Atoi(id)
		if err == nil {
			productIDs = append(productIDs, parsedID)
		}
	}
	return productIDs
}


// func SendSupportMessage() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		message := r.FormValue("message")
// 		if message == "" {
// 			http.Error(w, "Message is required", http.StatusBadRequest)
// 			return
// 		}
// 		file, _, err := r.FormFile("file")
// 		if err != nil && err != http.ErrMissingFile {
// 			http.Error(w, "Error processing file", http.StatusInternalServerError)
// 			return
// 		}
// 		err = sendEmailToSupport("support@example.com", message, file)
// 		if err != nil {
// 			http.Error(w, "Error sending support message", http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Message sent successfully"))
// 	}
// }
//
// func sendEmailToSupport(email, message string, file io.Reader) error {
// 	// Данные для SMTP-сервера
// 	smtpHost := "smtp.example.com" // Адрес вашего SMTP сервера
// 	smtpPort := "587"              // Порт сервера SMTP
// 	from := "support@example.com"  // Ваш email, с которого будет отправляться письмо
// 	password := "yourpassword"     // Пароль от email аккаунта
//
// 	// Создание сообщения
// 	var buf bytes.Buffer
// 	writer := multipart.NewWriter(&buf)
//
// 	// Создаем часть для текста сообщения
// 	bodyWriter, err := writer.CreatePart(map[string][]string{
// 		"Content-Type": {"text/plain; charset=UTF-8"},
// 		"Content-Transfer-Encoding": {"quoted-printable"},
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create text part: %v", err)
// 	}
//
// 	// Кодируем сообщение
// 	quotedPrintable := quotedprintable.NewWriter(bodyWriter)
// 	_, err = quotedPrintable.Write([]byte(message))
// 	if err != nil {
// 		return fmt.Errorf("failed to write quoted-printable message: %v", err)
// 	}
// 	quotedPrintable.Close()
//
// 	// Если есть файл, добавляем его как вложение
// 	if file != nil {
// 		attachmentWriter, err := writer.CreateFormFile("attachment", "support_file")
// 		if err != nil {
// 			return fmt.Errorf("failed to create file part: %v", err)
// 		}
//
// 		// Копируем файл в сообщение
// 		_, err = io.Copy(attachmentWriter, file)
// 		if err != nil {
// 			return fmt.Errorf("failed to copy file: %v", err)
// 		}
// 	}
//
// 	writer.Close()
//
// 	to := []string{"dzhumashev.arman@gmail.com"} // Email получателя
// 	subject := "Support Request from " + email // Тема письма
// 	headers := map[string]string{
// 		"From":    from,
// 		"To":      strings.Join(to, ", "),
// 		"Subject": subject,
// 		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()),
// 	}
//
// 	var messageBytes bytes.Buffer
// 	for key, value := range headers {
// 		messageBytes.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
// 	}
// 	messageBytes.WriteString("\r\n")
// 	messageBytes.Write(buf.Bytes())
//
// 	auth := smtp.PlainAuth("", from, password, smtpHost)
// 	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, messageBytes.Bytes())
// 	if err != nil {
// 		return fmt.Errorf("failed to send email: %v", err)
// 	}
//
// 	return nil
// }
