package itsy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestResource struct {
	*BaseResource
}

func TestGETWithParams(t *testing.T) {
	// Create a new Itsy instance.
	i := New()

	// Register a resource.
	r := i.Register("/hello/:name", &TestResource{})

	// Register a GET handler.
	r.GET(func(c Context) {
		c.Response().Write([]byte("Hello, " + c.GetParam("name")))
	})

	// Create a test server.
	server := httptest.NewServer(i)
	defer server.Close()

	// Make a request.
	resp, err := http.Get(server.URL + "/hello/world")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response body.
	expected := "Hello, world"
	if string(body) != expected {
		t.Errorf("Expected %q, got %q", expected, string(body))
	}
}
