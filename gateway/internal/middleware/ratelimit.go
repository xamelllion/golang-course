package middleware

import (
	"log"
	"net"
	"net/http"

	"github.com/xamelllion/golang-course/gateway/internal/adapter/ratelimiter"
)

type Middleware func(next http.Handler) http.Handler

func RateLimitMiddleware(ratelimiter ratelimiter.RateLimiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Printf("rate limiter error %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			allow, err := ratelimiter.Allow(r.Context(), ip)
			if err != nil {
				log.Printf("rate limiter error %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allow {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
