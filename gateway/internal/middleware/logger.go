package middleware

import (
	"log"
	"net"
	"net/http"
)

func LoggerMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Printf("rate limiter error %v", err)
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("%v %v %v", r.Method, r.URL.Path, ip)

			next.ServeHTTP(w, r)
		})
	}
}
