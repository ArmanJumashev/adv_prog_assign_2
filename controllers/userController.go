package controllers

import (
	"context"
	//     "bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	// 	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"

	//      "net/smtp"
	//      "mime/multipart"
	//     "mime/quotedprintable"

	//     "path/filepath"
	"online-shop/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
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

		query += strings.Join(setClauses, ", ") + " WHERE email = $" + strconv.Itoa(len(params)+1)
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

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	File    string `json:"file"`
}

func encodeFileToBase64(filePath string) (string, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileData), nil
}

// func SendEmail(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	var req SendEmailRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
// 		return
// 	}
//
// 	credentialsFile := `D:\edu\shop_adv\controllers\credentials.json`
// 	ctx := context.Background()
// 	srv, err := getGmailService(ctx, credentialsFile)
// 	if err != nil {
// 		http.Error(w, "Ошибка создания сервиса Gmail", http.StatusInternalServerError)
// 		log.Printf("Ошибка создания сервиса Gmail: %v", err)
// 		return
// 	}
//
// 	var fileBase64 string
// 	if req.File != "" {
// 		fileBase64, err = encodeFileToBase64(req.File)
// 		if err != nil {
// 			http.Error(w, "Ошибка кодирования файла", http.StatusInternalServerError)
// 			return
// 		}
// 	}
//
// 	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", req.To, req.Subject, req.Body)
//
// 	// Формируем MIME-мультипарт
// 	var buf bytes.Buffer
// 	writer := multipart.NewWriter(&buf)
//
// 	// Добавляем текстовую часть
// 	part, err := writer.CreatePart(textPlainHeader())
// 	if err != nil {
// 		http.Error(w, "Ошибка создания части сообщения", http.StatusInternalServerError)
// 		return
// 	}
// 	part.Write([]byte(message))
//
// 	// Если файл прикреплен, добавляем его как вложение
// 	if fileBase64 != "" {
// 		attachmentPart, err := writer.CreatePart(fileAttachmentHeader())
// 		if err != nil {
// 			http.Error(w, "Ошибка создания части вложения", http.StatusInternalServerError)
// 			return
// 		}
// 		attachmentPart.Write([]byte(fileBase64))
// 	}
//
// 	// Закрываем writer, чтобы сформировать окончательный MIME-формат
// 	writer.Close()
//
// 	// Создаем сообщение с MIME-контентом
// 	msg := &gmail.Message{
// 		Raw: encodeWeb64String(buf.Bytes()),
// 	}
//
// 	// Отправляем письмо через Gmail API
// 	_, err = srv.Users.Messages.Send("me", msg).Do()
// 	if err != nil {
// 		http.Error(w, "Ошибка отправки письма", http.StatusInternalServerError)
// 		log.Printf("Ошибка отправки письма: %v", err)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Письмо успешно отправлено"))
// }

func SendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	credentialsFile := `credentials.json`

	ctx := context.Background()
	srv, err := getGmailService(ctx, credentialsFile)
	if err != nil {
		http.Error(w, "Ошибка создания сервиса Gmail", http.StatusInternalServerError)
		log.Printf("Ошибка создания сервиса Gmail: %v", err)
		return
	}

	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", req.To, req.Subject, req.Body)
	msg := &gmail.Message{
		Raw: encodeWeb64String([]byte(message)),
	}

	_, err = srv.Users.Messages.Send("me", msg).Do()
	if err != nil {
		http.Error(w, "Ошибка отправки письма", http.StatusInternalServerError)
		log.Printf("Ошибка отправки письма: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Письмо успешно отправлено"))
}

// textPlainHeader создает заголовки для текстовой части
func textPlainHeader() textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "text/plain; charset=UTF-8")
	return header
}

// fileAttachmentHeader создает заголовки для файла вложения
func fileAttachmentHeader() textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", "attachment; filename=\"attachment.txt\"")
	header.Set("Content-Transfer-Encoding", "base64")
	return header
}

func createAttachment(fileBase64 string) *gmail.MessagePart {
	return &gmail.MessagePart{
		Filename: "attachment.txt",           // Укажите нужное имя файла
		MimeType: "application/octet-stream", // Тип MIME
		Body: &gmail.MessagePartBody{
			Data: fileBase64,
		},
	}
}

func getGmailService(ctx context.Context, credentialsFile string) (*gmail.Service, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл учетных данных: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать конфигурацию клиента: %v", err)
	}

	client := getClient(ctx, config)
	return gmail.NewService(ctx, option.WithHTTPClient(client))
}

func encodeWeb64String(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(ctx, tok)
}

//func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
//	tokenFile := "token.json"
//	tok, err := tokenFromFile(tokenFile)
//
//	tok = getTokenFromWeb(config)
//	saveToken(tokenFile, tok)
//
//	// Создаем TokenSource, который автоматически обновляет токен при его истечении
//	tokenSource := config.TokenSource(ctx, tok)
//
//	// Обновляем токен, если он истек
//	newToken, err := tokenSource.Token()
//	if err != nil {
//		log.Fatalf("Ошибка обновления токена: %v", err)
//	}
//
//	// Если токен изменился, сохраняем его
//	if newToken.AccessToken != tok.AccessToken {
//		saveToken(tokenFile, newToken)
//	}
//
//	return oauth2.NewClient(ctx, tokenSource)
//}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по ссылке для авторизации: \n%v\n", authURL)

	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Ошибка получения кода авторизации", http.StatusBadRequest)
			return
		}

		tok, err := config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка обмена кода на токен: %v", err), http.StatusInternalServerError)
			return
		}

		saveToken("token.json", tok)

		w.Write([]byte("token success"))
	})

	port := "8081"
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	fmt.Println("Ожидаем на localhost:" + port)
	select {}
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Не удалось создать файл токена: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
