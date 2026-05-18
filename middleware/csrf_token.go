package middleware

import (
	"net/http"

	ssRouter "github.com/jonecoboy/ssr/router"
)

const csrfTokenHeader = "X-CSRF-Token"
const validToken = "valid-token" // This should be dynamically generated and stored securely

// CSRFGenerateMiddleware generates a CSRF token and includes it in the response
func CSRFGenerateMiddleware(next ssRouter.Handler) ssRouter.Handler {
	return func(w http.ResponseWriter, r *ssRouter.SsrRequest) {
		// Generate and set the CSRF token (for simplicity, using a static token here)
		w.Header().Set(csrfTokenHeader, validToken)
		next(w, r)
	}
}

// CSRFValidateMiddleware validates the CSRF token in the request
func CSRFValidateMiddleware(next ssRouter.Handler) ssRouter.Handler {
	return func(w http.ResponseWriter, r *ssRouter.SsrRequest) {
		// Validate the CSRF token
		token := r.Header.Get(csrfTokenHeader)
		if token != validToken {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
