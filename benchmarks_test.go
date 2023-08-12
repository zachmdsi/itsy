package itsy

import (
	"fmt"
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

func setup() *Itsy {
	i := New()

	r1 := i.Register("/test/:id")
	r2 := i.Register("/test2/:id")

	r1.GET(func(ctx Context) error {
		id := ctx.GetParam("id")
		ctx.CreateField("id", id)
		return ctx.WriteHTML()
	})

	r2.GET(func(ctx Context) error {
		id := ctx.GetParam("id")
		ctx.CreateField("id", id)
		return ctx.WriteHTML()
	})

	r1.Link(r2, "related")

	return i
}

func BenchmarkServer(b *testing.B) {
	i := setup()

	// Start the server in a goroutine
	go i.Run(":8080")

	// Wait for the server to start
	http.Get("http://localhost:8080/test/1")

	b.Run("Concurrency100", func(b *testing.B) { benchmarkRequests(b, 100) })
	b.Run("Concurrency1000", func(b *testing.B) { benchmarkRequests(b, 1000) })
	b.Run("Concurrency10000", func(b *testing.B) { benchmarkRequests(b, 10000) })
}

func benchmarkRequests(b *testing.B, concurrency int) {
	b.SetParallelism(concurrency)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sendRequests()
		}
	})
}

func sendRequests() {
	// Send 100 requests in parallel
	for i := 0; i < 100; i++ {
		go func(id int) {
			http.Get(fmt.Sprintf("http://localhost:8080/test/%d", id))
			http.Get(fmt.Sprintf("http://localhost:8080/test2/%d", id))
		}(i)
	}
}