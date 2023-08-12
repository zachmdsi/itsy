package itsy

import (
	"io"
	"net/http"
)

type (
	// Response describes an HTTP response.
	Response struct {
		itsy       *Itsy               // The main framework instance.
		Writer     http.ResponseWriter // The HTTP response writer.
		StatusCode int                 // The HTTP status code.
	}
	// responseWriterWrapper is a wrapper around the response writer.
	responseWriterWrapper struct {
		writer     io.Writer
		statusCode int
		original   http.ResponseWriter
	}
)

// NewResponse creates a new response instance.
func NewResponse(res http.ResponseWriter, i *Itsy) *Response {
	return &Response{
		itsy:       i,
		Writer:     res,
		StatusCode: -1,
	}
}

// Write writes the response body.
func (r *Response) Write(b []byte) (n int, err error) {
	// Write the header if it hasn't been written yet.
	if r.StatusCode == -1 {
		r.WriteHeader(StatusOK)
	}
	n, err = r.Writer.Write(b)
	return
}

// WriteHeader writes the response header.
func (r *Response) WriteHeader(code int) {
	// Don't write the header if it has already been written.
	if r.StatusCode != -1 {
		return
	}
	r.StatusCode = code
}

// Write writes to the response.
func (res *responseWriterWrapper) Write(b []byte) (n int, err error) {
	return res.writer.Write(b)
}

// Header gets the header of the response.
func (res *responseWriterWrapper) Header() http.Header {
	return res.original.Header()
}

// WriteHeader writes the header of the response.
func (res *responseWriterWrapper) WriteHeader(code int) {
	res.statusCode = code
	res.original.WriteHeader(code)
}
