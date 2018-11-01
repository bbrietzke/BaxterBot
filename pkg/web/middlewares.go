package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"golang.org/x/time/rate"
)

func loggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func rateLimitingMW(rps int64, burst int) mux.MiddlewareFunc {
	logger.Println("Rate: ", rps, "burst: ", burst)
	return func(next http.Handler) http.Handler {
		limiter := rate.NewLimiter(rate.Limit(rps), burst)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := limiter.Wait(ctx); err != nil {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
			}

			next.ServeHTTP(w, r)
		})
	}
}
