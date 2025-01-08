package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

//     "online-shop/data_loader"

	"online-shop/db"
	"online-shop/routes"
	"online-shop/middleware"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

func main() {
	// Настройка логирования в файл
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла лога: %v", err)
	}
	defer logFile.Close()

	// Настройка логера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{}) // Формат JSON для логов
	logger.SetLevel(logrus.InfoLevel)             // Уровень логирования
	logger.SetOutput(logFile)                     // Вывод логов в файл

	// Подключение к базе данных
	logger.Info("Подключение к базе данных...")
	database := db.Connect()

	defer func() {
		logger.Info("Закрытие соединения с базой данных.")
		database.Close()
	}()

	// Настройка маршрутов
	logger.Info("Настройка маршрутов...")
	router := routes.SetupRoutes(database)

	// Добавляем промежуточное ПО для обработки ошибок, логирования и rate limiting
	limiter := rate.NewLimiter(2, 5) // Ограничение: 2 запроса в секунду, максимум 5
	router.Use(func(next http.Handler) http.Handler {
		return middleware.ErrorHandlingMiddleware(
			middleware.RateLimitMiddleware(
				middleware.LoggingMiddleware(next, logger), limiter, logger), logger)
	})

	// Обработка статических файлов
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./client"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Отдача index.html для запроса: ", r.URL.Path)
		http.ServeFile(w, r, "./client/index.html")
	})

	// Настройка graceful shutdown
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Канал для получения сигналов от ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Запуск сервера
	go func() {
		logger.WithField("port", port).Info("Запуск сервера...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Не удалось запустить сервер")
		}
	}()

	// Ожидание сигнала о завершении работы
	<-stop
	logger.Info("Получен сигнал о завершении работы")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Остановка сервера...")
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Ошибка при остановке сервера")
	} else {
		logger.Info("Сервер остановлен успешно")
	}
}
