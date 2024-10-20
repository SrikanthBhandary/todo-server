package router

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // Call the next handler in the chain
		log.Printf("Host: %s URL: %s Method %s Time: %v", r.RemoteAddr, r.URL, r.Method, time.Since(start))
	})
}
