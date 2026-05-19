# LoggingMiddleware

Logs HTTP request information including method, path, and duration.

## Import

```go
import "github.com/jonecoboy/ssr/middleware"
```

## Usage

```go
r.GET("/api/users", handler, []router.Middleware{
    middleware.LoggingMiddleware,
})
```

## Behavior

1. Logs when request starts: `Started GET /api/users`
2. Calls the next handler
3. Logs when request completes: `Completed /api/users in 1.234ms`

## Output Example

```
2024/01/15 10:30:00 Started GET /api/users
2024/01/15 10:30:00 Completed /api/users in 1.234567ms
```

## Notes

- Uses Go's standard `log` package
- Measures time using `time.Since()`
- Works with `router.Handler` signature