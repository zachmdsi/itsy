package itsy

import (
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	itsy := &Itsy{}
	router := newRouter(itsy)

	// Define a handler function for testing
	handler := func(ctx Context) error { return nil }

	// Add routes
	router.addRoute(http.MethodGet, "/test", handler)
	router.addRoute(http.MethodPost, "/test", handler)
	router.addRoute(http.MethodGet, "/test/:param", handler)

	// Test that routes have been added correctly
	testNode := router.index.children["test"]
	if testNode == nil {
		t.Fatal("Expected /test route to exist, but it doesn't")
	}

	if testNode.handlers[http.MethodGet] == nil {
		t.Fatal("Expected GET /test route to exist, but it doesn't")
	}

	if testNode.handlers[http.MethodPost] == nil {
		t.Fatal("Expected POST /test route to exist, but it doesn't")
	}

	paramNode := testNode.children[":param"]
	if paramNode == nil {
		t.Fatal("Expected /test/:param route to exist, but it doesn't")
	}

	if paramNode.handlers[http.MethodGet] == nil {
		t.Fatal("Expected GET /test/:param route to exist, but it doesn't")
	}

	if paramNode.param != "param" {
		t.Fatalf("Expected parameter name to be 'param', got '%s'", paramNode.param)
	}
}
