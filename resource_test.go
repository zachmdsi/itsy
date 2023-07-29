package itsy

import "testing"

func TestLink(t *testing.T) {
	r := &BaseResource{}
	err := r.Link("https://example.com", "Example", Attr{Key: "rel", Value: "example"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(r.Links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(r.Links))
	}

	link := r.Links[0]
	if link.Href != "https://example.com" {
		t.Errorf("Expected href to be 'https://example.com', got '%s'", link.Href)
	}
	if link.Prompt != "Example" {
		t.Errorf("Expected prompt to be 'Example', got '%s'", link.Prompt)
	}
	if link.Rel != "example" {
		t.Errorf("Expected rel to be 'example', got '%s'", link.Rel)
	}
}

func TestForm(t *testing.T) {
	r := &BaseResource{}
	form := NewForm("testForm", "/submit", "POST", NewFormField("name", "value"))
	r.Form(form)

	if len(r.Forms) != 1 {
		t.Errorf("Expected 1 form, got %d", len(r.Forms))
	}

	f := r.Forms[0]
	if f.Name != "testForm" {
		t.Errorf("Expected form name to be 'testForm', got '%s'", f.Name)
	}
	if f.Href != "/submit" {
		t.Errorf("Expected form href to be '/submit', got '%s'", f.Href)
	}
	if f.Method != "POST" {
		t.Errorf("Expected form method to be 'POST', got '%s'", f.Method)
	}
	if len(f.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(f.Fields))
	}
	if f.Fields[0].Name != "name" {
		t.Errorf("Expected field name to be 'name', got '%s'", f.Fields[0].Name)
	}
	if f.Fields[0].Value != "value" {
		t.Errorf("Expected field value to be 'value', got '%s'", f.Fields[0].Value)
	}
}

func TestRenderBase(t *testing.T) {
	r := &BaseResource{}
	output := r.RenderBase(nil)
	if output == "" {
		t.Errorf("Expected output to not be empty")
	}
}

func TestNewForm(t *testing.T) {
	form := NewForm("testForm", "/submit", "POST", NewFormField("name", "value"))
	if form.Name != "testForm" {
		t.Errorf("Expected form name to be 'testForm', got '%s'", form.Name)
	}
	if form.Href != "/submit" {
		t.Errorf("Expected form href to be '/submit', got '%s'", form.Href)
	}
	if form.Method != "POST" {
		t.Errorf("Expected form method to be 'POST', got '%s'", form.Method)
	}
	if len(form.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(form.Fields))
	}
	if form.Fields[0].Name != "name" {
		t.Errorf("Expected field name to be 'name', got '%s'", form.Fields[0].Name)
	}
	if form.Fields[0].Value != "value" {
		t.Errorf("Expected field value to be 'value', got '%s'", form.Fields[0].Value)
	}
}

func TestNewFormField(t *testing.T) {
	field := NewFormField("name", "value")
	if field.Name != "name" {
		t.Errorf("Expected field name to be 'name', got '%s'", field.Name)
	}
	if field.Value != "value" {
		t.Errorf("Expected field value to be 'value', got '%s'", field.Value)
	}
}

func TestParseLink(t *testing.T) {
	tag := A("/test/link", "Test Link", Attr{"rel", "test"}, Attr{"name", "test"}, Attr{"render", "test"}, Attr{"prompt", "Test 1"})
	link, err := ParseLink(tag)
	if err != nil {
		t.Fatalf("ParseLink: %v", err)
	}

	if link.Href != "/test/link" {
		t.Errorf("got href %q, want %q", link.Href, "/test/link")
	}

	if link.Prompt != "Test Link" {
		t.Errorf("got prompt %q, want %q", link.Prompt, "Test Link")
	}
}

func TestAddLink(t *testing.T) {
	// Create a new BaseResource.
	resource := &BaseResource{}

	// Add the link to the resource.
	err := resource.Link("/test/link", "Test Link", Data("rel", "test"), Data("name", "test"), Data("render", "test"), Data("prompt", "Test 1"))
	if err != nil {
		t.Fatalf("AddLink: %v", err)
	}

	// Check that the link was added.
	if len(resource.Links) != 1 {
		t.Errorf("got %d links, want 1", len(resource.Links))
	}

	// Check the link's attributes.
	link := resource.Links[0]
	if link.Href != "/test/link" {
		t.Errorf("got href %q, want %q", link.Href, "/test/link")
	}
	if link.Prompt != "Test Link" {
		t.Errorf("got prompt %q, want %q", link.Prompt, "Test Link")
	}
}