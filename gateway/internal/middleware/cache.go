package middleware

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/xamelllion/golang-course/gateway/internal/adapter/cache"
)

type responseRecorder struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

func CacheMiddleware(cacher cache.Cacher, ttl time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Printf("cacher error %v", err)
				next.ServeHTTP(w, r)
				return
			}

			key := "cache:endpoint:" + r.URL.Path + ":ip:" + ip
			data, isExists, err := cacher.Get(r.Context(), key)
			if err != nil {
				log.Printf("cacher error %v", err)
				next.ServeHTTP(w, r)
				return
			}

			if isExists {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(data)
				return
			}

			rec := &responseRecorder{
				ResponseWriter: w,
				body:           []byte{},
				status:         http.StatusOK,
			}
			next.ServeHTTP(rec, r)

			if rec.status == http.StatusOK {
				err := cacher.Set(r.Context(), key, rec.body, ttl)
				if err != nil {
					log.Printf("cacher set error %v", err)
				}
			}
		})
	}
}
