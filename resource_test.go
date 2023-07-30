package itsy

import "testing"

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