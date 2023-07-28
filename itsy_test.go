package itsy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNew tests the New function.
func TestNew(t *testing.T) {
	itsy := New()
	if itsy == nil {
		t.Errorf("New returned nil")
	}
}

// TestAdd tests the Add function.
func TestAdd(t *testing.T) {
	itsy := New()
	resource := &TestResource{}
	itsy.Add(http.MethodGet, "/test", resource)
	if _, ok := itsy.resources["/test"]; !ok {
		t.Errorf("resource not added")
	}
}

// TestHandleResource tests the handleResource function.
func TestHandleResource(t *testing.T) {
	itsy := New()
	resource := &TestResource{}
	itsy.Add(http.MethodGet, "/test", resource)
	server := httptest.NewServer(itsy)
	defer server.Close()
	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

// TestHTTPError tests the HTTPError function.
func TestHTTPError(t *testing.T) {
	itsy := New()
	server := httptest.NewServer(itsy)
	defer server.Close()
	resp, err := http.Get(server.URL + "/nonexistent")
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("got status %v, want %v", resp.StatusCode, http.StatusNotFound)
	}
}

// TestResource is a test resource.
type TestResource struct {
	BaseResource
}

// Render renders the test resource.
func (t *TestResource) Render(ctx Context) string {
	return "<h1>Hello, world!</h1>"
}
