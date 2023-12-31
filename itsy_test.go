package itsy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGET(t *testing.T) {
	// Create a new Itsy instance.
	i := New()

	// Register a resource.
	r := i.Register("/")

	// Register a GET handler.
	r.GET(func(c Context) error {
		return c.WriteString("Hello, world")
	})

	// Create a test server.
	server := httptest.NewServer(i)
	defer server.Close()

	// Make a request.
	resp, err := http.Get(server.URL)
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

func TestGETWithParams(t *testing.T) {
	// Create a new Itsy instance.
	i := New()

	// Register a resource.
	r := i.Register("/hello/:name")

	// Register a GET handler.
	r.GET(func(c Context) error {
		return c.WriteString("Hello, " + c.GetParamValue("name"))
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

func TestResourceDoesNotExist(t *testing.T) {
	i := New()

	// Create a test server.
	server := httptest.NewServer(i)
	defer server.Close()

	// Make a request to an unregistered resource.
	resp, err := http.Get(server.URL + "/doesnotexist")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Not Found: Resource does not exist"
	if string(body) != expected {
		t.Errorf("Expected %q, got %q", expected, string(body))
	}
}

func TestUnsupportedMethod(t *testing.T) {
	i := New()

	// Register a resource but don't add a POST handler.
	r := i.Register("/test")
	r.GET(func(c Context) error {
		return c.WriteString("Should not be called")
	})

	// Create a test server.
	server := httptest.NewServer(i)
	defer server.Close()

	// Make a POST request.
	resp, err := http.Post(server.URL+"/test", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Method Not Allowed: Handler does not exist for the request method"
	if string(body) != expected {
		t.Errorf("Expected %q, got %q", expected, string(body))
	}
}
