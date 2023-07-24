package itsy

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

type MockResource struct {
	BaseResource	
}
func TestRenderResource(t *testing.T) {
	logger := zap.NewExample()
	baseCtx := BaseContext{
		request:        &http.Request{},
		responseWriter: httptest.NewRecorder(),
		logger:         logger,
	}

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
