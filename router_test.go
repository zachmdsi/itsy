package itsy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRouter tests the routing functionality.
func TestRouter(t *testing.T) {
	itsy := New()
	router := NewRouter(itsy)
	router.Handle("GET", "/hello", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("Hello, world!"))
		return nil
	})

	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "Hello, world!" {
		t.Errorf("Expected body 'Hello, world!', got '%v'", rec.Body.String())
	}
}

// TestRouterNoRoute tests that the router returns a 404 for an unknown route.
func TestRouterNoRoute(t *testing.T) {
	itsy := New()
	router := NewRouter(itsy)

	req, err := http.NewRequest("GET", "/unknown", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %v", rec.Code)
	}
}

// TestRouterMethodNotAllowed tests that the router returns a 405 for an unsupported method.
func TestRouterMethodNotAllowed(t *testing.T) {
	itsy := New()
	router := NewRouter(itsy)
	router.Handle("GET", "/onlyget", func(ctx Context) error { return nil })

	req, err := http.NewRequest("POST", "/onlyget", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed, got %v", rec.Code)
	}
}

// TestRouterParam tests that the router correctly matches parameterized routes.
func TestRouterParam(t *testing.T) {
	itsy := New()
	router := NewRouter(itsy)
	router.Handle("GET", "/hello/:name", func(ctx Context) error {
		name := strings.TrimPrefix(ctx.Request().URL.Path, "/hello/")
		ctx.ResponseWriter().Write([]byte("Hello, " + name + "!"))
		return nil
	})

	req, err := http.NewRequest("GET", "/hello/world", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "Hello, world!" {
		t.Errorf("Expected body 'Hello, world!', got '%v'", rec.Body.String())
	}
}
