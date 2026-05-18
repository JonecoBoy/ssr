package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/jonecoboy/ssr/router"
)

func RecoverFromPanicMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.SsrRequest) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic and backtrace
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// Check if a request ID is provided
				if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
					log.Printf("Request ID: %s", reqID)
				}

				// Return HTTP 500 status
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
