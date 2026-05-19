# NoCacheMiddleware

Prevents browsers and proxies from caching the response.

## Import

```go
import "github.com/jonecoboy/ssr/middleware"
```

## Usage

```go
// Works with http.Handler
handler := middleware.NoCacheMiddleware(yourHandler)
```

## Behavior

1. Removes ETag-related headers from the request:
   - `ETag`
   - `If-Modified-Since`
   - `If-Match`
   - `If-None-Match`
   - `If-Range`
   - `If-Unmodified-Since`

2. Sets no-cache headers on the response:
   - `Expires: Thu, 01 Jan 1970 00:00:00 UTC`
   - `Cache-Control: no-cache, no-store, no-transform, must-revalidate, private, max-age=0`
   - `Pragma: no-cache`
   - `X-Accel-Expires: 0`

3. Calls the next handler

## Response Headers

```
Expires: Thu, 01 Jan 1970 00:00:00 UTC
Cache-Control: no-cache, no-store, no-transform, must-revalidate, private, max-age=0
Pragma: no-cache
X-Accel-Expires: 0
```

## Use Cases

- Dynamic content that should never be cached
- Sensitive data endpoints
- Real-time data APIs
- Authentication/session endpoints

## Notes

- Works with standard `http.Handler` interface
- `X-Accel-Expires: 0` is for Nginx proxy caching
- Sets epoch date for `Expires` header
