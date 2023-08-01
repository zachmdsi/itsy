package itsy

import "net/http"

type (
	// Response describes an HTTP response.
	Response struct {
		itsy 	  *Itsy                // The main framework instance.
		Writer     http.ResponseWriter // The HTTP response writer.
		StatusCode int                 // The HTTP status code.
		Headers    map[string]string   // The HTTP headers.
		Body       string              // The HTTP body.
	}
)

// NewResponse creates a new response instance.
func NewResponse(res http.ResponseWriter, i *Itsy) *Response {
	return &Response{
		Writer: res,
		itsy: i,
	}
}

// Write writes the response body.
func (r *Response) Write(b []byte) (n int, err error) {
	r.WriteHeader(StatusOK)
	n, err = r.Writer.Write(b)
	return
}

// WriteHeader writes the response header.
func (r *Response) WriteHeader(code int) {
	r.StatusCode = code
	r.Writer.WriteHeader(r.StatusCode)
}
