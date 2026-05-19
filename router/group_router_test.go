package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupRouter(t *testing.T) {
	// Create a new router
	nr := NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *SsrRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Create a group with a prefix
	group := nr.GROUP("/api", nil, nil)

	// Register routes within the group
	group.GET("/hello", helloHandler, nil)
	group.POST("/hello", helloHandler, nil)
	group.PUT("/hello", helloHandler, nil)
	group.DELETE("/hello", helloHandler, nil)

	tests := []struct {
		method     string
		url        string
		wantStatus int
		wantBody   string
	}{
		{"GET", "/api/hello", http.StatusOK, "Hello, World!"},
		{"POST", "/api/hello", http.StatusOK, "Hello, World!"},
		{"PUT", "/api/hello", http.StatusOK, "Hello, World!"},
		{"DELETE", "/api/hello", http.StatusOK, "Hello, World!"},
		{"GET", "/hello", http.StatusNotFound, "404 page not found\n"},
		{"POST", "/hello", http.StatusNotFound, "404 page not found\n"},
		{"PUT", "/hello", http.StatusNotFound, "404 page not found\n"},
		{"DELETE", "/hello", http.StatusNotFound, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}

			if rr.Body.String() != tt.wantBody {
				t.Errorf("got body %v, want %v", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestGroupRouterParsesUriAndQueryParams(t *testing.T) {
	nr := NewRouter()
	group := nr.GROUP("/products", nil, nil)

	var gotParams *SsrParamsRequest
	group.GET("/{code}", func(w http.ResponseWriter, r *SsrRequest) {
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

	if gotParams.Params["abc"] != "456" {
		t.Fatalf("got merged param %q, want %q", gotParams.Params["abc"], "456")
	}
}
