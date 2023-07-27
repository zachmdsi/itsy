package itsy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestResource struct {
	BaseResource
}

func (t *TestResource) Render() string {
	return "<h1>Hello, world!</h1>"
}

func TestItsy(t *testing.T) {
	// Create a new Itsy instance.
	itsy := New()

	// Add a test resource.
	itsy.AddResource(http.MethodGet,"/test", &TestResource{})

	// Create a test HTTP server.
	server := httptest.NewServer(itsy)
	defer server.Close()

	// Make a GET request to the test resource.
	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}

	// Check the status code.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Check the response body.
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}
	if got, want := string(body), "<h1>Hello, world!</h1>"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
}
