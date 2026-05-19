# SSR Middlewares

SSR provides several built-in middlewares to handle common tasks. You can also create your own custom middlewares.

## Available Middlewares

| Middleware | Description | Documentation |
|------------|-------------|---------------|
| `LoggingMiddleware` | Logs request start/end times and duration | [logging.md](middlewares/logging.md) |
| `RecoverFromPanicMiddleware` | Recovers from panics and returns HTTP 500 | [recover_from_panic.md](middlewares/recover_from_panic.md) |
| `BasicAuthMiddleware` | HTTP Basic Authentication | [basic_auth.md](middlewares/basic_auth.md) |
| `CSRFGenerateMiddleware` | Generates CSRF tokens | [csrf_token.md](middlewares/csrf_token.md) |
| `CSRFValidateMiddleware` | Validates CSRF tokens | [csrf_token.md](middlewares/csrf_token.md) |
| `NoCacheMiddleware` | Prevents response caching | [no_cache.md](middlewares/no_cache.md) |
| `UrlChangeMiddleware` | Rewrites URL paths | [url_change.md](middlewares/url_change.md) |

## Using Middlewares

### Single Route

```go
import "github.com/jonecoboy/ssr/middleware"

r.GET("/api/data", handler, []router.Middleware{
    middleware.LoggingMiddleware,
    middleware.RecoverFromPanicMiddleware,
})
```

### Route Group

```go
api := r.GROUP("/api", []router.Middleware{
    middleware.LoggingMiddleware,
}, nil)
```

## Middleware Execution Order

Middlewares are executed from left to right in the declaration order. The first middleware in the slice wraps the second, which wraps the third, and so on.

```go
// Execution order: Logging -> Auth -> Handler -> Auth -> Logging
r.GET("/", handler, []router.Middleware{loggingMiddleware, authMiddleware})
```

## Creating Custom Middlewares

```go
func MyCustomMiddleware(next router.Handler) router.Handler {
    return func(w http.ResponseWriter, r *router.SsrRequest) {
        // Before request
        fmt.Println("Before")

        next(w, r)

        // After request
        fmt.Println("After")
    }
}
```