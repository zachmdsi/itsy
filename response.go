package itsy

type (
	// Response describes an HTTP response.
	Response struct {
		StatusCode int               // The HTTP status code.
		Headers    map[string]string // The HTTP headers.
		Body       string            // The HTTP body.
	}
)
