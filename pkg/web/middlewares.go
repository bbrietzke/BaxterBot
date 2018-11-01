package web

import (
	"context"
	"net/http"
	"time"
)

func loggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func allowLimitsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "TooManyRequests", http.StatusTooManyRequests)
		}

		next.ServeHTTP(w, r)
	})
}

func waitLimitsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := limiter.Wait(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		}

		next.ServeHTTP(w, r)
	})
}
