package itsy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestItsyGET tests GET request handling in Itsy.
func TestItsyGET(t *testing.T) {
	itsy := New()

	itsy.GET("/get", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("GET, itsy!"))
		return nil
	})

	req := httptest.NewRequest("GET", "/get", nil)
	w := httptest.NewRecorder()

	itsy.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "GET, itsy!" {
		t.Errorf("Expected body: 'GET, itsy!', got '%v'", w.Body.String())
	}
}

// TestItsyPOST tests POST request handling in Itsy.
func TestItsyPOST(t *testing.T) {
	itsy := New()

	itsy.POST("/post", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("POST, itsy!"))
		return nil
	})

	req := httptest.NewRequest("POST", "/post", nil)
	w := httptest.NewRecorder()
	itsy.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "POST, itsy!" {
		t.Errorf("Expected body: 'POST, itsy!', got '%v'", w.Body.String())
	}
}

// TestItsyPUT tests PUT request handling in Itsy.
func TestItsyPUT(t *testing.T) {
	itsy := New()

	itsy.PUT("/put", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("PUT, itsy!"))
		return nil
	})

	req := httptest.NewRequest("PUT", "/put", nil)
	w := httptest.NewRecorder()
	itsy.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "PUT, itsy!" {
		t.Errorf("Expected body: 'PUT, itsy!', got '%v'", w.Body.String())
	}
}

// TestItsyDELETE tests DELETE request handling in Itsy.
func TestItsyDELETE(t *testing.T) {
	itsy := New()

	itsy.DELETE("/delete", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("DELETE, itsy!"))
		return nil
	})

	req := httptest.NewRequest("DELETE", "/delete", nil)
	w := httptest.NewRecorder()
	itsy.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "DELETE, itsy!" {
		t.Errorf("Expected body: 'DELETE, itsy!', got '%v'", w.Body.String())
	}
}

// TestItsyPATCH tests PATCH request handling in Itsy.
func TestItsyPATCH(t *testing.T) {
	itsy := New()

	itsy.PATCH("/patch", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("PATCH, itsy!"))
		return nil
	})

	req := httptest.NewRequest("PATCH", "/patch", nil)
	w := httptest.NewRecorder()
	itsy.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", w.Code)
	}

	if w.Body.String() != "PATCH, itsy!" {
		t.Errorf("Expected body: 'PATCH, itsy!', got '%v'", w.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	itsy := New()
	itsy.Logger = zaptest.NewLogger(t)

	middlewareCalled := false
	handlerCalled := false

	itsy.Use(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			middlewareCalled = true
			ctx.Logger().Info("Before handler")
			err := next(ctx)
			ctx.Logger().Info("After handler")
			return err
		}
	})

	itsy.GET("/", func(ctx Context) error {
		handlerCalled = true
		ctx.Logger().Info("Handler")
		return nil
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	itsy.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !middlewareCalled {
		t.Error("Expected middleware to be called")
	}

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}
