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

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "Hello, world!" {
		t.Errorf("Expected body 'Hello, world!', got '%v'", w.Body.String())
	}
}

// TestRouterNoRoute tests that the router returns a 404 for an unknown route.
func TestRouterNoRoute(t *testing.T) {
	itsy := New()
	router := NewRouter(itsy)

	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %v", w.Code)
	}
}

// TestRouterMethodNotAllowed tests that the router returns a 405 for an unsupported method.
func TestRouterMethodNotAllowed(t *testing.T) {
	itsy := New()
	router := NewRouter(itsy)
	router.Handle("GET", "/onlyget", func(ctx Context) error { return nil })

	req := httptest.NewRequest("POST", "/onlyget", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed, got %v", w.Code)
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

	req := httptest.NewRequest("GET", "/hello/world", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "Hello, world!" {
		t.Errorf("Expected body 'Hello, world!', got '%v'", w.Body.String())
	}
}
