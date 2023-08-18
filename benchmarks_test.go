package itsy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	// Create a new instance.
	i := New()

	// Register a resource.
	r := i.Register("/")

	// Register a GET handler.
	r.GET(func(c Context) error {
		return c.WriteString("Hello, world")
	})

	// Create a request.
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Make a request.
	for n := 0; n < b.N; n++ {
		i.ServeHTTP(rr, req)
	}
}
