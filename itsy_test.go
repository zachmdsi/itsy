package itsy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestResource struct {
	Resource
}

func TestGET(t *testing.T) {
	// Create a new Itsy instance.
	i := New()

	// Register a resource.
	r := i.Register("/", &TestResource{})

	// Register a GET handler.
	r.GET(func(c Context) error {
		return c.Response().WriteString("Hello, world")
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
	expected := "<html><body>Hello, world</body></html>"
	if string(body) != expected {
		t.Errorf("Expected %q, got %q", expected, string(body))
	}
}

func TestGETWithParams(t *testing.T) {
	// Create a new Itsy instance.
	i := New()

	// Register a resource.
	r := i.Register("/hello/:name", &TestResource{})

	// Register a GET handler.
	r.GET(func(c Context) error {
		return c.Response().WriteString("Hello, " + c.GetParam("name"))
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
	expected := "<html><body>Hello, world</body></html>"
	if string(body) != expected {
		t.Errorf("Expected %q, got %q", expected, string(body))
	}
}

func TestHTTPErrorFunction(t *testing.T) {
	i := New()

	// Create a test server.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HTTPError(http.StatusNotFound, "Not Found", w, i.Logger)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Not Found: Not Found"
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
	r := i.Register("/test", &TestResource{})
	r.GET(func(c Context) error {
		return c.Response().WriteString("Should not be called")
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
