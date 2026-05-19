# UrlChangeMiddleware

Rewrites URL paths by replacing a prefix with a new value.

## Import

```go
import "github.com/jonecoboy/ssr/middleware"
```

## Usage

```go
// Create middleware that rewrites /old-api/* to /api/*
rewriter := middleware.UrlChangeMiddleware("/old-api", "/api")

// Apply to handler
handler := rewriter(yourHandler)
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `old` | `string` | URL prefix to match and replace |
| `new` | `string` | Replacement prefix |

## Behavior

1. Checks if request path starts with `old` prefix
2. If it matches, replaces the first occurrence with `new`
3. Calls the next handler with the modified path

## Examples

### API Version Migration

```go
// Redirect old API version to new
rewriter := middleware.UrlChangeMiddleware("/v1/api", "/v2/api")
// /v1/api/users -> /v2/api/users
```

### Path Alias

```go
// Create shorter alias
rewriter := middleware.UrlChangeMiddleware("/u", "/users")
// /u/123 -> /users/123
```

### Legacy URL Support

```go
// Support old URL structure
rewriter := middleware.UrlChangeMiddleware("/legacy", "/new-system")
// /legacy/data -> /new-system/data
```

## Notes

- Only replaces the first occurrence of `old` in the path
- Uses `strings.HasPrefix` for matching
- Works with standard `http.Handler` interface
- Useful for URL migrations without breaking existing links
