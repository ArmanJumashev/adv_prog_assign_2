// ratelimit.go
package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
	"github.com/sirupsen/logrus"
)

// Middleware для rate limiting
func RateLimitMiddleware(next http.Handler, limiter *rate.Limiter, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			logger.Warn("Too many requests")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
