package middleware

import (
	"net/http"
	"social-network/pkg/logger"
)

func MiddlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLogger.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		logger.InfoLogger.Println("Executing middlewareOne again")
	})
}
