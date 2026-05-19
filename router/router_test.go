package router

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRouter(t *testing.T) {
	// Create a new router
	nr := NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *SsrRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	postHandler := func(w http.ResponseWriter, r *SsrRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	putHandler := func(w http.ResponseWriter, r *SsrRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}
	deleteHandler := func(w http.ResponseWriter, r *SsrRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler
	nr.GET("/hello/{id}", helloHandler, []Middleware{})
	nr.POST("/post/{id}", postHandler, []Middleware{})
	nr.PUT("/put/{id}", putHandler, []Middleware{})
	nr.DELETE("/delete/{id}", deleteHandler, []Middleware{})

	tests := []struct {
		method     string
		url        string
		wantStatus int
	}{
		{"GET", "/hello/123", http.StatusOK},
		{"POST", "/post/321", http.StatusOK},
		{"PUT", "/put/321", http.StatusOK},
		{"DELETE", "/delete/321", http.StatusOK},
		{"POST", "/hello/123", http.StatusMethodNotAllowed},
		{"GET", "/post/321", http.StatusMethodNotAllowed},
		{"DELETE", "/put/321", http.StatusMethodNotAllowed},
		{"PUT", "/delete/321", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestRouterParsesUriAndQueryParams(t *testing.T) {
	nr := NewRouter()

	var gotParams *SsrParamsRequest
	nr.GET("/products/{code}", func(w http.ResponseWriter, r *SsrRequest) {
		gotParams = r.GetParams()
		w.WriteHeader(http.StatusOK)
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/sku-123?joneco=2&abc=456", nil)
	rr := httptest.NewRecorder()

	nr.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got status %v, want %v", rr.Code, http.StatusOK)
	}

	if gotParams == nil {
		t.Fatal("expected params to be populated")
	}

	if gotParams.UriParams["code"] != "sku-123" {
		t.Fatalf("got uri param %q, want %q", gotParams.UriParams["code"], "sku-123")
	}

	if gotParams.QueryString["joneco"] != "2" {
		t.Fatalf("got query param %q, want %q", gotParams.QueryString["joneco"], "2")
	}

	if gotParams.QueryString["abc"] != "456" {
		t.Fatalf("got query param %q, want %q", gotParams.QueryString["abc"], "456")
	}

	if gotParams.Params["code"] != "sku-123" {
		t.Fatalf("got merged param %q, want %q", gotParams.Params["code"], "sku-123")
	}

	if gotParams.Params["joneco"] != "2" {
		t.Fatalf("got merged param %q, want %q", gotParams.Params["joneco"], "2")
	}
}

func TestStartServer(t *testing.T) {
	nr := NewRouter()

	nr.GET("/health", func(w http.ResponseWriter, r *SsrRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}, nil)

	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Start server in goroutine
	go nr.StartServer(port)

	// Wait for server to be ready
	serverURL := fmt.Sprintf("http://localhost:%d", port)
	var resp *http.Response
	for i := 0; i < 50; i++ {
		resp, err = http.Get(serverURL + "/health")
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if err != nil {
		t.Fatalf("server did not start: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %v, want %v", resp.StatusCode, http.StatusOK)
	}
}
