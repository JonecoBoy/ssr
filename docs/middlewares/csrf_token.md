# CSRF Token Middlewares

Provides Cross-Site Request Forgery (CSRF) protection with two middlewares: one for generating tokens and one for validating them.

## Import

```go
import "github.com/jonecoboy/ssr/middleware"
```

## CSRFGenerateMiddleware

Generates a CSRF token and includes it in the response header.

### Usage

```go
r.GET("/form", formHandler, []router.Middleware{
    middleware.CSRFGenerateMiddleware,
})
```

### Behavior

1. Sets `X-CSRF-Token` header in the response
2. Calls the next handler

## CSRFValidateMiddleware

Validates the CSRF token in incoming requests.

### Usage

```go
r.POST("/submit", submitHandler, []router.Middleware{
    middleware.CSRFValidateMiddleware,
})
```

### Behavior

1. Reads `X-CSRF-Token` header from request
2. If token is invalid or missing:
   - Returns HTTP 403 Forbidden
3. If valid, calls the next handler

## Typical Flow

```go
// GET route generates token
r.GET("/form", showForm, []router.Middleware{
    middleware.CSRFGenerateMiddleware,
})

// POST route validates token
r.POST("/form", handleSubmit, []router.Middleware{
    middleware.CSRFValidateMiddleware,
})
```

## Client Usage

```javascript
// Get token from response header
const token = response.headers.get('X-CSRF-Token');

// Include in subsequent requests
fetch('/form', {
    method: 'POST',
    headers: {
        'X-CSRF-Token': token
    }
});
```

## Notes

- Header name: `X-CSRF-Token`
- Current implementation uses a static token for demonstration
- In production, implement dynamic token generation and secure storage
