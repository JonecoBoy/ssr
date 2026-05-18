package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/jonecoboy/ssr/router"
)

func LoggingMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.SsrRequest) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
