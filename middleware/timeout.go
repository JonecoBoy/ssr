package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/jonecoboy/ssr/router"
)

func TimeoutMiddleware(timeout time.Duration) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.SsrRequest) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r.Request = r.Request.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-ctx.Done():
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			case <-done:
			}
		})
	}
}
