package middleware

import (
	"net/http"

	"github.com/jonecoboy/ssr/router"
	"github.com/microcosm-cc/bluemonday"
)

func XSSSanitizeMiddleware() router.Middleware {
	policy := bluemonday.UGCPolicy()

	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.SsrRequest) {
			// Sanitize query parameters
			query := r.URL.Query()
			for key, values := range query {
				for i, value := range values {
					query[key][i] = policy.Sanitize(value)
				}
			}
			r.URL.RawQuery = query.Encode()

			// Sanitize form parameters
			if err := r.ParseForm(); err == nil {
				for key, values := range *r.PostForm {
					for i, value := range values {
						(*r.PostForm)[key][i] = policy.Sanitize(value)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
