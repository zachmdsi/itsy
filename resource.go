package itsy

type (
	// Resource is an interface that should be implemented by types that represent themselves as hypermedia resources.
	Resource interface {
		GetLinks() []Link
		GetForms() []Form
		GetEmbeds() []Embed
		GetTemplates() []Template
		GetActions() []Action
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
	// BaseResource is a base implementation of the Resource interface.
	BaseResource struct {
		Links     []Link
		Forms     []Form
		Embeds    []Embed
		Templates []Template
		Actions   []Action
	}
)

func (b *BaseResource) GetLinks() []Link         { return b.Links }
func (b *BaseResource) GetForms() []Form         { return b.Forms }
func (b *BaseResource) GetEmbeds() []Embed       { return b.Embeds }
func (b *BaseResource) GetTemplates() []Template { return b.Templates }
func (b *BaseResource) GetActions() []Action     { return b.Actions }
