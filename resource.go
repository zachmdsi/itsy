package itsy

import (
	"bytes"
	"errors"
	"strings"
	"text/template"
)

type (
	// Resource is an interface that should be implemented by types that represent themselves as hypermedia resources.
	Resource interface {
		GetLinks() []Link                       // GetLinks returns a slice of links that describe the resource.
		GetForms() []Form                       // GetForms returns a slice of forms that describe the resource.
		GetEmbeds() []Embed                     // GetEmbeds returns a slice of embedded resources.
		GetTemplates() []Template               // GetTemplates returns a slice of URL templates that clients can use to construct URLs to resources.
		GetActions() []Action                   // GetActions returns a slice of actions that can be invoked by the client.

		AddLink(tag *Tag) error           // AddLink adds a link to the resource.

		RenderBase(Context) string              // RenderBase renders the base context for the resource.
		Render(Context) string                  // Render renders the resource.
	}
	// Link represents link from one resource to another.
	Link struct {
		Rel    string // Rel identifies the relationship between the linked resource and the current context.
		Href   string // Href holds the URL of the linked resource.
		Prompt string // Prompt is a human-readable description of the link.
		Name   string // Name can be used as a secondary key for selecting link elements.
		Render string // Render hints how the linked resource should be rendered in a given media type.
	}
	// Form represents an HTML form with fields for parameters and an action URL.
	Form struct {
		Name   string      // Name serves as the form's identifier.
		Href   string      // Href holds the URL where the form data will be sent upon submission.
		Method string      // Method defines the HTTP method used to submit the form.
		Type   string      // Type specifies how the form-data should be encoded when submitting it to the server.
		Fields []FormField // Fields is a slice of all the input fields in the forms.
	}
	// FormField represents a field in an HTML form.
	FormField struct {
		Name  string // Name is the name of the input field.
		Value string // Value is the default value of the input field.
	}
	// Embed represents a resource that is embedded within another resource.
	Embed struct {
		Rel      string   // Rel describes the relationship of the embedded resource to the parent resource.
		Resource Resource // Resource is the embedded resource.
	}
	// Template represents a URL template that clients can use to construct URLs to resources.
	Template struct {
		Name string // Name serves as the template's identifier.
		Href string // Href holds the URL template string.
	}
	// Action represents the server-side operation that can be invoked by the client.
	Action struct {
		Name   string      // Name serves as the action's identifier.
		Href   string      // Href is the URL where the request will be sent upon invocation.
		Method string      // Method defines the HTTP method used to submit the form.
		Type   string      // Type specifies how the action data should be encoded when submitting it to the server.
		Fields []FormField // Fields is a slice of all the input fields in the action.
	}
	// Tag represents an HTML tag.
	Tag struct {
		name  string
		attrs []Attr
		text  string
	}
	// Attr represents an HTML attribute.
	Attr struct {
		Key   string
		Value string
	}
	// BaseResource is a base implementation of the Resource interface.
	BaseResource struct {
		Links     []Link
		Forms     []Form
		Embeds    []Embed
		Templates []Template
		Actions   []Action
	}
)

var resourceTemplate = template.Must(template.New("resource").Parse(`
<div>
	<h2>Links</h2>
	{{range .Links}}
	<a href="{{.Href}}">{{.Prompt}}</a>
	{{end}}
	<h2>Forms</h2>
	{{range .Forms}}
	<form action="{{.Href}}" method="{{.Method}}">
	{{range .Fields}}
		<label>{{.Name}}: <input type="text" name="{{.Name}}" value="{{.Value}}"></label>
	{{end}}
	<input type="submit" value="Submit">
	</form>
	{{end}}
	<h2>Actions</h2>
	{{range .Actions}}
	<form action="{{.Href}}" method="{{.Method}}">
	{{range .Fields}}
		<label>{{.Name}}: <input type="text" name="{{.Name}}" value="{{.Value}}"></label>
	{{end}}
	<input type="submit" value="Submit">
	</form>
	{{end}}
	<h2>Embeds</h2>
	{{range .Embeds}}
	<div>
		<h3>{{.Rel}}</h3>
		{{.Resource.Render}}
	</div>
	{{end}}
	<h2>Templates</h2>
	{{range .Templates}}
	<div>
		<h3>{{.Name}}</h3>
		<code>{{.Href}}</code>
	</div>
	{{end}}
</div>
`))

// GetLinks returns a slice of links that describe the resource.
func (b *BaseResource) GetLinks() []Link         { return b.Links }

// GetForms returns a slice of forms that describe the resource.
func (b *BaseResource) GetForms() []Form         { return b.Forms }

// GetEmbeds returns a slice of embedded resources.
func (b *BaseResource) GetEmbeds() []Embed       { return b.Embeds }

// GetTemplates returns a slice of URL templates that clients can use to construct URLs to resources.
func (b *BaseResource) GetTemplates() []Template { return b.Templates }

// GetActions returns a slice of actions that can be invoked by the client.
func (b *BaseResource) GetActions() []Action     { return b.Actions }

// Render is a no-op implementation of the Render method.
func (b *BaseResource) Render(Context) string    { return "" }

// AddLink adds a link to the resource.
func (b *BaseResource) AddLink(tag *Tag) error {
	// Parse the link from the tag.
	link, err := ParseLink(tag)
	if err != nil {
		return err
	}

	// Add the link to the resource.
	b.Links = append(b.Links, link)

	return nil
}

// ParseLink parses a link from a Tag.
func ParseLink(tag *Tag) (Link, error) {
	// Check that the tag is an anchor element.
	if tag.name != "a" {
		return Link{}, errors.New("tag is not an anchor element")
	}

	// Create a new Link.
	link := Link{}

	// Parse the tag's attributes.
	for _, attr := range tag.attrs {
		switch attr.Key {
		case "href":
			link.Href = attr.Value
		case "rel":
			link.Rel = attr.Value
		case "name":
			link.Name = attr.Value
		case "render":
			link.Render = attr.Value
		case "prompt":
			link.Prompt = attr.Value
		}
	}

	// Set the link's prompt to the tag's text.
	link.Prompt = tag.text

	return link, nil
}

// A is a helper function for generating HTML anchor elements.
func A(href string, text string, attrs ...Attr) *Tag {
	// Create a new anchor element.
    a := &Tag{
        name: "a",
        attrs: []Attr{
			NewAttr("href", href),
		},
		text: text,
    }

	// Add the attributes to the anchor element.
	for _, attr := range attrs {
		a.Set(attr.Key, attr.Value)
	}

    return a
}

// Data is a helper function for generating HTML data attributes.
func Data(key, value string) Attr {
	return Attr{Key: "data-" + key, Value: value}
}

// RenderBase renders the base context for the resource using the resource template.
func (b *BaseResource) RenderBase(Context) string {
	var buf bytes.Buffer
	resourceTemplate.Execute(&buf, b)
	return buf.String()
}

// Set sets the value of an attribute.
func (t* Tag) Set(key, val string) {
	t.attrs = append(t.attrs, NewAttr(key, val))
}

// String returns the string representation of the tag.
func (t *Tag) String() string {
	var b strings.Builder
	b.WriteString("<")
	b.WriteString(t.name)
	for _, attr := range t.attrs {
		b.WriteString(" ")
		b.WriteString(attr.Key)
		b.WriteString("=\"")
		b.WriteString(attr.Value)
		b.WriteString("\"")
	}
	b.WriteString(">")
	return b.String()
}

// NewAttr creates a new attribute.
func NewAttr(key, val string) Attr {
	return Attr{
		Key:   key,
		Value: val,
	}
}