<p align="center">
  <img src="docs/images/ssr_logo.png" alt="SSR Logo" width="300">
</p>

# SSR - Simple Server Router

A lightweight HTTP router for Go that wraps the standard `http.ServeMux` with additional features like middleware support, route groups, and automatic parameter parsing.
THIS IS NOT A FRAMEWORK AND DON'T PRETEND TO BE.

## Installation

```bash
go get github.com/jonecoboy/ssr
```

## Features

- HTTP method routing (GET, POST, PUT, DELETE, OPTIONS, HEAD, TRACE, CONNECT)
- Middleware support
- Route groups with prefixes
- URI parameter parsing (`/users/{id}`)
- Query string parsing
- Automatic body parsing (JSON, XML, form-urlencoded)
- Custom `SsrRequest` with parsed parameters

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/jonecoboy/ssr/router"
)

func main() {
    r := router.NewRouter()

    r.GET("/hello/{name}", func(w http.ResponseWriter, r *router.SsrRequest) {
        name := r.Params.UriParams["name"]
        fmt.Fprintf(w, "Hello, %s!", name)
    }, nil)

    r.StartServer(8080)
}
```

### Route Groups

```go
r := router.NewRouter()

// Create a group with prefix "/api"
api := r.GROUP("/api", nil, nil)

api.GET("/users", listUsersHandler, nil)
api.POST("/users", createUserHandler, nil)
api.GET("/users/{id}", getUserHandler, nil)
```

### Middleware

SSR allows you to implement your own middlewares (before(pre) and after (post) request) or use the pre-existing ones. You can chain middlewares and the most left one will be the first.

For a complete list of built-in middlewares, see the [Middleware Documentation](docs/middleware.md).

```go
func loggingMiddleware(next router.Handler) router.Handler {
    return func(w http.ResponseWriter, r *router.SsrRequest) {
        fmt.Printf("%s %s\n", r.Method, r.RequestURI)
        next(w, r)
    }
}

func authMiddleware(next router.Handler) router.Handler {
    return func(w http.ResponseWriter, r *router.SsrRequest) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

// Apply middleware to a single route
r.GET("/protected", handler, []router.Middleware{authMiddleware, loggingMiddleware})

// Apply middleware to a group
api := r.GROUP("/api", []router.Middleware{loggingMiddleware}, nil)
```

### Accessing Parameters

```go
r.GET("/products/{code}", func(w http.ResponseWriter, r *router.SsrRequest) {
    params := r.GetParams()

    // URI parameters
    code := params.UriParams["code"]

    // Query string parameters (?page=1&limit=10)
    page := params.QueryString["page"]
    limit := params.QueryString["limit"]

    // Combined params (URI + query string)
    all := params.Params

    // Or get raw params map directly
    rawParams := r.GetReqParams()
}, nil)
```

### Handling POST/PUT Body

```go
r.POST("/users", func(w http.ResponseWriter, r *router.SsrRequest) {
    body, err := r.GetBody()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    username := body["username"]
    email := body["email"]

    // Process the data...
}, nil)
```

The router automatically parses request bodies based on Content-Type:
- `application/json`
- `application/xml` / `text/xml`
- `application/x-www-form-urlencoded`

## API Reference

### Router Methods

| Method | Description |
|--------|-------------|
| `NewRouter()` | Creates a new router instance |
| `GET(pattern, handler, middlewares)` | Register a GET route |
| `POST(pattern, handler, middlewares)` | Register a POST route |
| `PUT(pattern, handler, middlewares)` | Register a PUT route |
| `DELETE(pattern, handler, middlewares)` | Register a DELETE route |
| `OPTIONS(pattern, handler, middlewares)` | Register an OPTIONS route |
| `HEAD(pattern, handler, middlewares)` | Register a HEAD route |
| `TRACE(pattern, handler, middlewares)` | Register a TRACE route |
| `CONNECT(pattern, handler, middlewares)` | Register a CONNECT route |
| `GROUP(prefix, middlewares, groupMiddlewares)` | Create a route group |
| `StartServer(port)` | Start the HTTP server |

### SsrRequest

| Field | Type | Description |
|-------|------|-------------|
| `Header` | `http.Header` | Request headers |
| `Method` | `string` | HTTP method |
| `Host` | `string` | Request host |
| `RemoteAddr` | `string` | Client address |
| `UserAgent` | `string` | User-Agent header |

| Method | Return Type | Description |
|--------|-------------|-------------|
| `GetParams()` | `*SsrParamsRequest` | Get parsed parameters |
| `GetReqParams()` | `map[string]map[string]string` | Get raw params map |
| `GetBody()` | `(map[string]interface{}, error)` | Get parsed request body |

### SsrParamsRequest

| Field | Type | Description |
|-------|------|-------------|
| `UriParams` | `map[string]string` | Parameters from URI path |
| `QueryString` | `map[string]string` | Query string parameters |
| `Params` | `map[string]string` | Combined URI + query params |

## License

MIT
