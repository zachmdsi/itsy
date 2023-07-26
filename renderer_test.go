package itsy

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockResource struct {
	BaseResource
}

func TestRenderResourceJSON(t *testing.T) {
	// Create a new Itsy instance to get a new BaseContext
	itsyInstance := New()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	respRec := httptest.NewRecorder()

	baseCtx := itsyInstance.newBaseContext(req, respRec)

	testRes := &MockResource{}

	// Call RenderResource and check if it returns an error.
	err := baseCtx.RenderResource(testRes)
	if err != nil {
		t.Errorf("RenderResource returned an error: %v", err)
	}

	// Check the response.
	resp := baseCtx.responseWriter.(*httptest.ResponseRecorder).Result()

	// Check if the Content-Type header is set to application/json.
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, got %v", contentType)
	}

	// Check if the response body is what we expect.
	expectedBody, _ := json.Marshal(testRes)
	actualBody, _ := io.ReadAll(resp.Body)
	if !reflect.DeepEqual(expectedBody, actualBody) {
		t.Errorf("Response body doesn't match. Expected %v, got %v", string(expectedBody), string(actualBody))
	}
}

func TestRenderResourceXML(t *testing.T) {
	// Create a new Itsy instance to get a new BaseContext
	itsyInstance := New()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/xml")
	respRec := httptest.NewRecorder()

	baseCtx := itsyInstance.newBaseContext(req, respRec)

	testRes := &MockResource{}

	// Call RenderResource and check if it returns an error.
	err := baseCtx.RenderResource(testRes)
	if err != nil {
		t.Errorf("RenderResource returned an error: %v", err)
	}

	// Check the response.
	resp := baseCtx.responseWriter.(*httptest.ResponseRecorder).Result()

	// Check if the Content-Type header is set to application/xml.
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/xml" {
		t.Errorf("Expected Content-Type to be application/xml, got %v", contentType)
	}

	// Check if the response body is what we expect.
	expectedBody, _ := xml.Marshal(testRes)
	actualBody, _ := io.ReadAll(resp.Body)
	if !reflect.DeepEqual(expectedBody, actualBody) {
		t.Errorf("Response body doesn't match. Expected %v, got %v", string(expectedBody), string(actualBody))
	}
}
