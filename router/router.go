package router

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ServeMux struct {
	*http.ServeMux
}

type Handler func(w http.ResponseWriter, r *SsrRequest)

func (h Handler) ServeHTTP(writer http.ResponseWriter, request *SsrRequest) {
	// Call the original golang http handler
	h(writer, request)
}

type Middlewares []Middleware
type Middleware func(Handler) Handler

func NewRouter() *ServeMux {
	mux := &ServeMux{
		ServeMux: http.NewServeMux(),
	}
	return mux
}

type GenericXML struct {
	XMLName xml.Name
	Content string       `xml:",chardata"`
	Nodes   []GenericXML `xml:",any"`
}

type SsrRequest struct {
	*http.Request
	Header        http.Header
	ContentLength int64
	Form          *url.Values
	PostForm      *url.Values
	MultipartForm *multipart.Form
	RemoteAddr    string
	Method        string
	RequestURI    string
	ctx           context.Context
	tls           *tls.ConnectionState
	UserAgent     string
	Proto         string
	Host          string
	pattern       string
	params        *SsrParamsRequest
	body          interface{}
	ValidatedData map[string]string
}

type SsrParamsRequest struct {
	QueryString map[string]string
	UriParams   map[string]string
	Params      map[string]string
}

func (nr *SsrRequest) GetBody() (map[string]interface{}, error) {
	// Check if the body exists
	if nr.body == nil {
		return nil, fmt.Errorf("body is empty or not initialized")
	}

	// Try to assert the body to a map
	body, ok := nr.body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse body as map[string]interface{}")
	}

	return body, nil
}

func (nr *SsrRequest) GetParams() *SsrParamsRequest {
	return nr.params
}

func (nr *SsrRequest) GetReqParams() map[string]map[string]string {
	qs := parseQueryString(nr.Request)
	uriParams := parseUriParams(nr.Request, nr.pattern)

	params := make(map[string]map[string]string)
	params["params"] = make(map[string]string)
	params["queryString"] = qs
	params["uriParams"] = uriParams

	for key, value := range qs {
		params["params"][key] = value
	}
	for key, value := range uriParams {
		params["params"][key] = value
	}

	return params
}

func (mux *ServeMux) GET(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodGet+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
		}
		finalHandler(w, ssrRequest)
	}))
}

func (mux *ServeMux) POST(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodPost+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Restore the body for potential reuse
		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

		// Parse body into a unified map
		parsedBody := make(map[string]interface{})
		contentType := r.Header.Get("Content-Type")

		switch {
		case contentType == "application/json":
			// Parse JSON
			if err := json.Unmarshal(bodyBytes, &parsedBody); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
		case contentType == "application/xml" || contentType == "text/xml":
			// Parse the XML into a generic tree structure
			var root GenericXML
			if err := xml.Unmarshal(bodyBytes, &root); err != nil {
				http.Error(w, "Invalid XML", http.StatusBadRequest)
				return
			}

			// Convert the XML tree to a map
			parsedBody = xmlToMap(root)

		case contentType == "application/x-www-form-urlencoded":
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form data", http.StatusBadRequest)
				return
			}
			for key, values := range r.PostForm {
				// Add form data to the map (use the first value for simplicity)
				if len(values) > 0 {
					parsedBody[key] = values[0]
				}
			}
		default:
			// For unsupported content types, treat as raw text and try to parse
			rawBody := string(bodyBytes)
			parsedMap, err := parseRawBody(rawBody)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to parse raw body: %v", err), http.StatusBadRequest)
				return
			}
			for key, value := range parsedMap {
				parsedBody[key] = value
			}
			parsedBody["rawBody"] = string(bodyBytes)
		}

		// Set up request parameters
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		// Create the custom ssrRequest
		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
			body:          parsedBody, // Store the unified map
		}

		finalHandler(w, ssrRequest)
	}))
}

func (mux *ServeMux) PUT(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodPut+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Restore the body for potential reuse
		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

		// Parse body into a unified map
		parsedBody := make(map[string]interface{})
		contentType := r.Header.Get("Content-Type")

		switch {
		case contentType == "application/json":
			// Parse JSON
			if err := json.Unmarshal(bodyBytes, &parsedBody); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
		case contentType == "application/xml" || contentType == "text/xml":
			// Parse the XML into a generic tree structure
			var root GenericXML
			if err := xml.Unmarshal(bodyBytes, &root); err != nil {
				http.Error(w, "Invalid XML", http.StatusBadRequest)
				return
			}

			// Convert the XML tree to a map
			parsedBody = xmlToMap(root)

		case contentType == "application/x-www-form-urlencoded":
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form data", http.StatusBadRequest)
				return
			}
			for key, values := range r.PostForm {
				// Add form data to the map (use the first value for simplicity)
				if len(values) > 0 {
					parsedBody[key] = values[0]
				}
			}
		default:
			// For unsupported content types, treat as raw text and try to parse
			rawBody := string(bodyBytes)
			parsedMap, err := parseRawBody(rawBody)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to parse raw body: %v", err), http.StatusBadRequest)
				return
			}
			for key, value := range parsedMap {
				parsedBody[key] = value
			}
			parsedBody["rawBody"] = string(bodyBytes)
		}

		// Set up request parameters
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		// Create the custom ssrRequest
		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
			body:          parsedBody, // Store the unified map
		}

		finalHandler(w, ssrRequest)
	}))
}

func (mux *ServeMux) DELETE(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodDelete+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
		}
		finalHandler(w, ssrRequest)
	}))
}

func xmlToMap(node GenericXML) map[string]interface{} {
	result := make(map[string]interface{})

	// If the node has no children, it's a leaf node
	if len(node.Nodes) == 0 {
		result[node.XMLName.Local] = node.Content
		return result
	}

	// Recursively process child nodes
	children := map[string]interface{}{}
	for _, child := range node.Nodes {
		childMap := xmlToMap(child)
		for key, value := range childMap {
			// Handle duplicate keys by appending to a slice
			if existing, found := children[key]; found {
				switch v := existing.(type) {
				case []interface{}:
					children[key] = append(v, value)
				default:
					children[key] = []interface{}{v, value}
				}
			} else {
				children[key] = value
			}
		}
	}

	// Add the processed children to the current node
	result[node.XMLName.Local] = children
	return result
}

func (mux *ServeMux) TRACE(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodTrace+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodTrace {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
		}
		finalHandler(w, ssrRequest)
	}))
}

func (mux *ServeMux) OPTIONS(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodOptions+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
		}
		finalHandler(w, ssrRequest)
	}))
}

func (mux *ServeMux) HEAD(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodHead+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
		}
		finalHandler(w, ssrRequest)
	}))
}

func (mux *ServeMux) CONNECT(pattern string, handler Handler, middlewares []Middleware) {
	finalHandler := applyMiddlewares(handler, middlewares...)
	mux.ServeMux.Handle(http.MethodConnect+" "+pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reqParams := getReqParams(r, pattern)
		params := &SsrParamsRequest{
			QueryString: reqParams["queryString"],
			UriParams:   reqParams["uriParams"],
			Params:      reqParams["params"],
		}

		ssrRequest := &SsrRequest{
			Request:       r,
			Header:        r.Header,
			Form:          &r.Form,
			Method:        r.Method,
			PostForm:      &r.PostForm,
			ctx:           r.Context(),
			ContentLength: r.ContentLength,
			tls:           r.TLS,
			Proto:         r.Proto,
			Host:          r.Host,
			pattern:       pattern,
			params:        params,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
		}
		finalHandler(w, ssrRequest)
	}))
}

func getReqParams(r *http.Request, pattern string) map[string]map[string]string {
	qs := parseQueryString(r)
	uriParams := parseUriParams(r, pattern)

	params := make(map[string]map[string]string)
	params["params"] = make(map[string]string)
	params["queryString"] = qs
	params["uriParams"] = uriParams

	for key, value := range qs {
		params["params"][key] = value
	}
	for key, value := range uriParams {
		params["params"][key] = value
	}

	return params
}

func parseQueryString(r *http.Request) map[string]string {
	rawQS := r.URL.Query()
	qs := make(map[string]string)
	for key, values := range rawQS {
		if len(values) > 0 {
			qs[key] = values[0]
		}
	}
	return qs
}

func parseRawBody(rawBody string) (map[string]string, error) {
	// Define valid separators
	separators := []rune{',', ';', '|'}
	result := make(map[string]string)

	var keyBuilder, valueBuilder strings.Builder
	isEscaped := false
	isParsingKey := true

	for _, char := range rawBody {
		switch {
		case isEscaped:
			// Append the escaped character
			if isParsingKey {
				keyBuilder.WriteRune(char)
			} else {
				valueBuilder.WriteRune(char)
			}
			isEscaped = false
		case char == '\\':
			// Handle escape sequences
			isEscaped = true
		case contains(separators, char):
			// Handle separators
			if isParsingKey {
				return nil, errors.New("missing key-value separator '=' in raw body")
			}
			// Save the current key-value pair
			key := strings.TrimSpace(keyBuilder.String())
			value := strings.TrimSpace(valueBuilder.String())
			if key != "" {
				result[key] = value
			}
			// Reset builders for the next pair
			keyBuilder.Reset()
			valueBuilder.Reset()
			isParsingKey = true
		case char == '=' && isParsingKey:
			// Switch to parsing the value
			isParsingKey = false
		default:
			// Append to the appropriate builder
			if isParsingKey {
				keyBuilder.WriteRune(char)
			} else {
				valueBuilder.WriteRune(char)
			}
		}
	}

	// Add the final key-value pair if any
	key := strings.TrimSpace(keyBuilder.String())
	value := strings.TrimSpace(valueBuilder.String())
	if key != "" {
		result[key] = value
	}

	return result, nil
}

// Helper function to check if a rune is in a slice
func contains(slice []rune, char rune) bool {
	for _, v := range slice {
		if v == char {
			return true
		}
	}
	return false
}

func parseUriParams(r *http.Request, pattern string) map[string]string {
	patternParts := splitRoutePath(extractRoutePath(pattern))
	requestParts := splitRoutePath(r.URL.Path)

	params := make(map[string]string)
	if len(patternParts) != len(requestParts) {
		return params
	}

	for idx, patternPart := range patternParts {
		if isRouteParam(patternPart) {
			paramName := strings.TrimSuffix(strings.TrimPrefix(patternPart, "{"), "}")
			if paramName != "" {
				params[paramName] = requestParts[idx]
			}
			continue
		}

		if patternPart != requestParts[idx] {
			return map[string]string{}
		}
	}

	return params
}

func extractRoutePath(pattern string) string {
	parts := strings.SplitN(strings.TrimSpace(pattern), " ", 2)
	if len(parts) == 2 && strings.HasPrefix(parts[1], "/") {
		return parts[1]
	}

	return pattern
}

func splitRoutePath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}

	return strings.Split(trimmed, "/")
}

func isRouteParam(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
}

func (r *SsrRequest) SetContext(ctx context.Context) {
	r.Request = r.Request.WithContext(ctx)
}

// variadic so can be any size of array
func applyMiddlewares(h Handler, middlewares ...Middleware) Handler {
	// in this normal order will be last middleware first!
	//for _, middleware := range middlewares {
	//	h = middleware(h)
	//}
	// in this order will load from the left to the right in the declaration
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func (mux *ServeMux) StartServer(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
