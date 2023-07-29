package itsy

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestBaseContext(t *testing.T) {
	i := &Itsy{Logger: zap.NewNop()}
	formData := url.Values{
		"formKey": []string{"formValue"},
	}
	r := httptest.NewRequest("POST", "http://example.com", strings.NewReader(formData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	c := i.newBaseContext(r, w)

	if c.Request() != r {
		t.Errorf("Expected request to be %v, got %v", r, c.Request())
	}

	if c.ResponseWriter() != w {
		t.Errorf("Expected response writer to be %v, got %v", w, c.ResponseWriter())
	}

	paramKey := "key"
	paramValue := "value"
	c.SetParam(paramKey, paramValue)
	if val, ok := c.Params()[paramKey]; !ok || val != paramValue {
		t.Errorf("Expected param '%s' to be '%s', got '%s'", paramKey, paramValue, val)
	}

	err := c.ParseForm()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val, ok := c.formValues["formKey"]; !ok || val != "formValue" {
		t.Errorf("Expected form value 'formKey' to be 'formValue', got '%s'", val)
	}
}
