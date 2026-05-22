package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/jonecoboy/ssr/pagination"
)

// Response represents a standard API response.
type Response struct {
	Success bool        `json:"success" xml:"success"`
	Data    interface{} `json:"data" xml:"data"`
}

// PaginatedResponse represents a paginated API response.
type PaginatedResponse struct {
	Success bool            `json:"success" xml:"success"`
	Data    interface{}     `json:"data" xml:"data"`
	Meta    pagination.Meta `json:"meta" xml:"meta"`
}

// ErrorResponse represents an error API response.
type ErrorResponse struct {
	Success bool   `json:"success" xml:"success"`
	Error   string `json:"error" xml:"error"`
}

// JSON writes a successful JSON response with the given data.
func JSON(w http.ResponseWriter, data interface{}) {
	JSONWithStatus(w, http.StatusOK, data)
}

// JSONWithStatus writes a JSON response with a custom status code.
func JSONWithStatus(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// JSONPaginated writes a paginated JSON response.
func JSONPaginated(w http.ResponseWriter, data interface{}, p pagination.Pagination, total int64) {
	JSONPaginatedWithStatus(w, http.StatusOK, data, p, total)
}

// JSONPaginatedWithStatus writes a paginated JSON response with a custom status code.
func JSONPaginatedWithStatus(w http.ResponseWriter, status int, data interface{}, p pagination.Pagination, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(PaginatedResponse{
		Success: true,
		Data:    data,
		Meta:    p.GetMeta(total),
	})
}

// JSONError writes an error JSON response.
func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// XML writes a successful XML response with the given data.
func XML(w http.ResponseWriter, data interface{}) {
	XMLWithStatus(w, http.StatusOK, data)
}

// XMLWithStatus writes an XML response with a custom status code.
func XMLWithStatus(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// XMLPaginated writes a paginated XML response.
func XMLPaginated(w http.ResponseWriter, data interface{}, p pagination.Pagination, total int64) {
	XMLPaginatedWithStatus(w, http.StatusOK, data, p, total)
}

// XMLPaginatedWithStatus writes a paginated XML response with a custom status code.
func XMLPaginatedWithStatus(w http.ResponseWriter, status int, data interface{}, p pagination.Pagination, total int64) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(PaginatedResponse{
		Success: true,
		Data:    data,
		Meta:    p.GetMeta(total),
	})
}

// XMLError writes an error XML response.
func XMLError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// Raw writes raw data with a custom content type.
func Raw(w http.ResponseWriter, contentType string, data []byte) {
	RawWithStatus(w, http.StatusOK, contentType, data)
}

// RawWithStatus writes raw data with a custom status code and content type.
func RawWithStatus(w http.ResponseWriter, status int, contentType string, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write(data)
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Created writes a 201 Created JSON response.
func Created(w http.ResponseWriter, data interface{}) {
	JSONWithStatus(w, http.StatusCreated, data)
}
