package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonecoboy/ssr/pagination"
)

type testData struct {
	ID   int    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := testData{ID: 1, Name: "Test"}

	JSON(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("JSON() status = %v, want %v", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("JSON() Content-Type = %v, want application/json", ct)
	}

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success {
		t.Error("JSON() Success = false, want true")
	}
}

func TestJSONWithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	data := testData{ID: 1, Name: "Test"}

	JSONWithStatus(w, http.StatusCreated, data)

	if w.Code != http.StatusCreated {
		t.Errorf("JSONWithStatus() status = %v, want %v", w.Code, http.StatusCreated)
	}
}

func TestJSONPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	data := []testData{{ID: 1, Name: "Test"}}
	p := pagination.New(2, 10)

	JSONPaginated(w, data, p, 55)

	if w.Code != http.StatusOK {
		t.Errorf("JSONPaginated() status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success {
		t.Error("JSONPaginated() Success = false, want true")
	}
	if resp.Meta.Page != 2 {
		t.Errorf("JSONPaginated() Meta.Page = %v, want 2", resp.Meta.Page)
	}
	if resp.Meta.Limit != 10 {
		t.Errorf("JSONPaginated() Meta.Limit = %v, want 10", resp.Meta.Limit)
	}
	if resp.Meta.Total != 55 {
		t.Errorf("JSONPaginated() Meta.Total = %v, want 55", resp.Meta.Total)
	}
	if resp.Meta.TotalPages != 6 {
		t.Errorf("JSONPaginated() Meta.TotalPages = %v, want 6", resp.Meta.TotalPages)
	}
}

func TestJSONError(t *testing.T) {
	w := httptest.NewRecorder()

	JSONError(w, http.StatusBadRequest, "invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("JSONError() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var resp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Success {
		t.Error("JSONError() Success = true, want false")
	}
	if resp.Error != "invalid input" {
		t.Errorf("JSONError() Error = %v, want 'invalid input'", resp.Error)
	}
}

func TestXML(t *testing.T) {
	w := httptest.NewRecorder()
	data := testData{ID: 1, Name: "Test"}

	XML(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("XML() status = %v, want %v", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/xml" {
		t.Errorf("XML() Content-Type = %v, want application/xml", ct)
	}
	if !strings.Contains(w.Body.String(), "<?xml") {
		t.Error("XML() response should contain XML header")
	}
}

func TestXMLPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	data := []testData{{ID: 1, Name: "Test"}}
	p := pagination.New(1, 10)

	XMLPaginated(w, data, p, 100)

	if w.Code != http.StatusOK {
		t.Errorf("XMLPaginated() status = %v, want %v", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "<meta>") {
		t.Error("XMLPaginated() response should contain meta element")
	}
}

func TestXMLError(t *testing.T) {
	w := httptest.NewRecorder()

	XMLError(w, http.StatusInternalServerError, "server error")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("XMLError() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<error>server error</error>") {
		t.Errorf("XMLError() response should contain error message, got: %s", body)
	}
}

func TestRaw(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("plain text content")

	Raw(w, "text/plain", data)

	if w.Code != http.StatusOK {
		t.Errorf("Raw() status = %v, want %v", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/plain" {
		t.Errorf("Raw() Content-Type = %v, want text/plain", ct)
	}
	if w.Body.String() != "plain text content" {
		t.Errorf("Raw() body = %v, want 'plain text content'", w.Body.String())
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	NoContent(w)

	if w.Code != http.StatusNoContent {
		t.Errorf("NoContent() status = %v, want %v", w.Code, http.StatusNoContent)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	data := testData{ID: 1, Name: "New"}

	Created(w, data)

	if w.Code != http.StatusCreated {
		t.Errorf("Created() status = %v, want %v", w.Code, http.StatusCreated)
	}
}

func TestXMLResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	data := testData{ID: 1, Name: "Test"}

	XML(w, data)

	// Skip XML header for parsing
	body := w.Body.String()
	body = strings.TrimPrefix(body, xml.Header)

	var resp Response
	err := xml.Unmarshal([]byte(body), &resp)
	if err != nil {
		t.Errorf("XML() response should be valid XML: %v", err)
	}
	if !resp.Success {
		t.Error("XML() Success = false, want true")
	}
}
