package itsy

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestResourceLinking(t *testing.T) {
	i := New()

	// Define a dummy handler for GET requests
	dummyHandler := func(c Context) error {
		return c.WriteHTML()
	}

	// Register a new resource
	resource1 := i.Register("/resource1")
	resource1.GET(dummyHandler)

	// Register another resource to link to
	resource2 := i.Register("/resource2")
	resource2.GET(dummyHandler)

	// Link resource1 to resource2
	err := resource1.Link("/resource2", "related")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Make an HTTP request to resource1
	req := httptest.NewRequest("GET", "/resource1", nil)
	rr := httptest.NewRecorder()
	i.ServeHTTP(rr, req)

	// Check if the link to resource2 is correctly rendered in the response
	if !strings.Contains(rr.Body.String(), "<a href=\"/resource2\" rel=\"related\"></a>") {
		t.Fatalf("expected link to /resource2, got %v", rr.Body.String())
	}

	// Create a resource that doesn't exist in Itsy's resources
	newBaseResource("/resource3", i)

	// Test linking resource1 to a non-existing resource
	err = resource1.Link("/resource3", "nonexistent")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "resource does not exist" {
		t.Fatalf("expected resource does not exist, got %v", err)
	}
}

func TestMultipleResourceLinking(t *testing.T) {
	itsy := New()

	// Define a dummy handler for GET requests
	dummyHandler := func(c Context) error {
		return c.WriteHTML()
	}

	// Register primary resource
	primaryResource := itsy.Register("/primary")
	primaryResource.GET(dummyHandler)

	// Register multiple resources to link to
	linkedResources := []Resource{}
	for i := 1; i <= 3; i++ {
		resourcePath := "/linked" + strconv.Itoa(i)
		linkedResource := itsy.Register(resourcePath)
		linkedResource.GET(dummyHandler)
		linkedResources = append(linkedResources, linkedResource)

		// Link primary resource to each linked resource
		err := primaryResource.Link(resourcePath, "related"+strconv.Itoa(i))
		if err != nil {
			t.Fatalf("expected no error while linking %s, got %v", resourcePath, err)
		}
	}

	// Make an HTTP request to primary resource
	req := httptest.NewRequest("GET", "/primary", nil)
	rr := httptest.NewRecorder()
	itsy.ServeHTTP(rr, req)

	// Check if the links to all linked resources are correctly rendered in the response
	for i, linkedResource := range linkedResources {
		expectedLink := "<a href=\"" + linkedResource.Path() + "\" rel=\"related" + strconv.Itoa(i+1) + "\"></a>"
		if !strings.Contains(rr.Body.String(), expectedLink) {
			t.Fatalf("expected link to %s, but it was not found in the response", linkedResource.Path())
		}
	}
}

func TestParameterizedResourceLinking(t *testing.T) {
	// Create a new Itsy instance
	i := New()

	// Register a primary resource with a parameterized route
	primaryResource := i.Register("/primary/:id")
	primaryResource.GET(func(c Context) error {
		return c.WriteHTML()
	})

	// Register a linked resource with a parameterized route
	linkedResource := i.Register("/linked/:id")
	linkedResource.GET(func(c Context) error {
		return c.WriteHTML()
	})

	// Link the primary resource to the linked resource
	err := primaryResource.Link("/linked/:id", "related")
	if err != nil {
		t.Fatalf("Failed to link resources: %v", err)
	}

	// Send a GET request to the primary resource with a parameter
	req, err := http.NewRequest(GET, "/primary/123", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response
	recorder := httptest.NewRecorder()
	i.ServeHTTP(recorder, req)

	response := recorder.Body.String()

	expectedLink := "<a href=\"/linked/123\" rel=\"related\"></a>"

	if !strings.Contains(response, expectedLink) {
		t.Fatalf("Expected response to contain link to '/linked/123', but got: %s", response)
	}
}

func TestMultipleParameterResourceLinking(t *testing.T) {
	// Create a new Itsy instance
	i := New()

	// Register a primary resource with multiple parameters
	productResource := i.Register("/products/:category/:id")
	productResource.GET(func(c Context) error {
		return c.WriteHTML()
	})

	// Register a linked resource also with multiple parameters
	reviewResource := i.Register("/reviews/:category/:id")
	reviewResource.GET(func(c Context) error {
		return c.WriteHTML()
	})

	// Link the product resource to the review resource
	err := productResource.Link("/reviews/:category/:id", "view review")
	if err != nil {
		t.Fatalf("Failed to link resources: %v", err)
	}

	// Send a GET request to the product resource with category and ID parameters
	req, err := http.NewRequest(GET, "/products/electronics/123", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response
	recorder := httptest.NewRecorder()
	i.ServeHTTP(recorder, req)

	response := recorder.Body.String()

	expectedLink := "<a href=\"/reviews/electronics/123\" rel=\"view review\"></a>"

	if !strings.Contains(response, expectedLink) {
		t.Fatalf("Expected response to contain link to '%s', but got: %s", expectedLink, response)
	}
}
