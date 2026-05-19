# RecoverFromPanicMiddleware

Recovers from panics in handlers and returns an HTTP 500 Internal Server Error response.

## Import

```go
import "github.com/jonecoboy/ssr/middleware"
```

## Usage

```go
r.GET("/api/risky", handler, []router.Middleware{
    middleware.RecoverFromPanicMiddleware,
})
```

## Behavior

1. Wraps the handler in a `defer` with `recover()`
2. If a panic occurs:
   - Logs the panic message and stack trace
   - Logs the `X-Request-ID` header if present
   - Returns HTTP 500 Internal Server Error
3. If no panic occurs, request proceeds normally

## Output Example

When a panic occurs:
```
2024/01/15 10:30:00 Panic: runtime error: index out of range
goroutine 1 [running]:
runtime/debug.Stack()
    /usr/local/go/src/runtime/debug/stack.go:24 +0x65
...
2024/01/15 10:30:00 Request ID: abc-123
```

## Notes

- Always place this middleware early in the chain to catch panics from subsequent middlewares
- Uses `runtime/debug.Stack()` for detailed stack traces
- Supports request ID tracking via `X-Request-ID` header