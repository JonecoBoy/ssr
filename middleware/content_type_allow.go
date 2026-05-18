package middleware

import (
	"net/http"
	"strings"

	"github.com/jonecoboy/ssr/router"
)

func AllowedContentTypeMiddleware(allowedTypes []string) router.Middleware {
	allowed := make(map[string]struct{})
	for _, contentType := range allowedTypes {
		allowed[contentType] = struct{}{}
	}

	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.SsrRequest) {
			contentType := r.Header.Get("Content-Type")
			for allowedType := range allowed {
				if strings.HasPrefix(contentType, allowedType) {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		})
	}
}
