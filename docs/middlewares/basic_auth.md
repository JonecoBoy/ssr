# BasicAuthMiddleware

Implements HTTP Basic Authentication for protecting routes.

## Import

```go
import "github.com/jonecoboy/ssr/middleware"
```

## Usage

```go
// Wrap your handler with BasicAuthMiddleware
protectedHandler := middleware.BasicAuthMiddleware("admin", "secret", yourHandler)
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `username` | `string` | Expected username |
| `password` | `string` | Expected password |
| `next` | `http.Handler` | Handler to call on successful auth |

## Behavior

1. Extracts credentials from `Authorization` header using `r.BasicAuth()`
2. Validates username and password match
3. If invalid:
   - Sets `WWW-Authenticate` header to prompt browser login dialog
   - Returns HTTP 401 Unauthorized
4. If valid, calls the next handler

## Response on Failure

```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Basic realm="Please enter your username and password"
```

## Notes

- Uses standard HTTP Basic Authentication
- Credentials are base64 encoded (not encrypted) - use HTTPS in production
- Works with standard `http.Handler` interface