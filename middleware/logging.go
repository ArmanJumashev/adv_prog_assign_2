// logging.go
package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Оборачиваем ResponseWriter, чтобы получить статус код ответа
		lrw := &LoggedResponseWriter{ResponseWriter: w}
		// Обрабатываем запрос
		next.ServeHTTP(lrw, r)
		// Логируем информацию о запросе
		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"url":        r.URL.Path,
			"status_code": lrw.statusCode,
			"duration":    time.Since(start),
		}).Info("Запрос обработан")
	})
}

// Обертка для ResponseWriter для логирования статус кода
type LoggedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *LoggedResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
