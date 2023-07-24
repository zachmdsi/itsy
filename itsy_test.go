package itsy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestItsyGET tests GET request handling in Itsy.
func TestItsyGET(t *testing.T) {
	app := New()

	app.GET("/get", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("GET, itsy!"))
		return nil
	})

	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "GET, itsy!" {
		t.Errorf("Expected body: 'GET, itsy!', got '%v'", rec.Body.String())
	}
}

// TestItsyPOST tests POST request handling in Itsy.
func TestItsyPOST(t *testing.T) {
	app := New()

	app.POST("/post", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("POST, itsy!"))
		return nil
	})

	req, err := http.NewRequest("POST", "/post", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "POST, itsy!" {
		t.Errorf("Expected body: 'POST, itsy!', got '%v'", rec.Body.String())
	}
}

// TestItsyPUT tests PUT request handling in Itsy.
func TestItsyPUT(t *testing.T) {
	app := New()

	app.PUT("/put", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("PUT, itsy!"))
		return nil
	})

	req, err := http.NewRequest("PUT", "/put", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "PUT, itsy!" {
		t.Errorf("Expected body: 'PUT, itsy!', got '%v'", rec.Body.String())
	}
}

// TestItsyDELETE tests DELETE request handling in Itsy.
func TestItsyDELETE(t *testing.T) {
	app := New()

	app.DELETE("/delete", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("DELETE, itsy!"))
		return nil
	})

	req, err := http.NewRequest("DELETE", "/delete", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "DELETE, itsy!" {
		t.Errorf("Expected body: 'DELETE, itsy!', got '%v'", rec.Body.String())
	}
}

// TestItsyPATCH tests PATCH request handling in Itsy.
func TestItsyPATCH(t *testing.T) {
	app := New()

	app.PATCH("/patch", func(ctx Context) error {
		ctx.ResponseWriter().Write([]byte("PATCH, itsy!"))
		return nil
	})

	req, err := http.NewRequest("PATCH", "/patch", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rec.Code)
	}

	if rec.Body.String() != "PATCH, itsy!" {
		t.Errorf("Expected body: 'PATCH, itsy!', got '%v'", rec.Body.String())
	}
}
