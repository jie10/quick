// Package basicauth provides HTTP Basic Authentication middleware for Go web servers.
// This file contains tests for the BasicAuth middleware implementation.
package basicauth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeffotoni/quick/internal/concat"
)
// TestBasicAuth_Success tests successful authentication with valid credentials.
//
// Test steps:
// 1. Creates middleware with valid credentials ("admin", "1234")
// 2. Sets up a test handler that returns 200 OK
// 3. Sends request with correct Basic Auth header
// 4. Verifies:
//    - Status code is 200 OK
//    - Request proceeds to the handler
//
// Example:
// TestBasicAuth_Success test with credentials success
// TestBasicAuth_Success(t *testing.T)

func TestBasicAuth_Success(t *testing.T) {
	// valid users
	username := "admin"
	password := "1234"

	// Create middleware
	middleware := BasicAuth(username, password)

	// Create test handler
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Create a mock request with correct credentials
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	auth := concat.String("Basic ", base64.StdEncoding.EncodeToString([]byte(concat.String(username, ":", password))))
	req.Header.Set("Authorization", auth)

	// Simulate response
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Validate HTTP status
	if rec.Code != http.StatusOK {
		t.Errorf("Esperado %d, mas recebeu %d", http.StatusOK, rec.Code)
	}
}

// TestBasicAuth_Failure tests authentication failure with invalid credentials.

// Test steps:
// 1. Creates middleware with credentials ("admin", "1234")
// 2. Sends request with invalid credentials ("wronguser:wrongpass")
// 3. Verifies:
//    - Status code is 401 Unauthorized
//    - Request doesn't reach the handler
//
// TestBasicAuth_Failure test with invalid credentials
// TestBasicAuth_Failure(t *testing.T)
func TestBasicAuth_Failure(t *testing.T) {
	middleware := BasicAuth("admin", "1234")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a mock request with invalid credentials
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("wronguser:wrongpass"))
	req.Header.Set("Authorization", auth)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Validate that returns 401 Unauthorized
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Esperado %d, mas recebeu %d", http.StatusUnauthorized, rec.Code)
	}
}

// Test failed authentication without credentials
func TestBasicAuth_NoCredentials(t *testing.T) {
	middleware := BasicAuth("admin", "1234")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create request without Authorization header
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Validate that returns 401 Unauthorized
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Esperado %d, mas recebeu %d", http.StatusUnauthorized, rec.Code)
	}
}
