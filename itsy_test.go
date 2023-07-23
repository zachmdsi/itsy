package itsy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestItsySuccess tests a successful request.
func TestItsySuccess(t *testing.T) {
	app := New()

	app.GET("/hello", func(ctx *Context) {
		ctx.ResponseWriter.Write([]byte("Hello, itsy!"))
	})

	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "Hello, itsy!" {
		t.Errorf("Expected body: 'Hello, itsy!', got %v", rec.Body.String())
	}
}

// TestItsyNotFound tests a request for a non-existent route.
func TestItsyNotFound(t *testing.T) {
	app := New()

	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %v", rec.Code)
	}

	if rec.Body.String() != "404 Not Found" {
		t.Errorf("Expected body: '404 Not Found', got %v", rec.Body.String())
	}
}