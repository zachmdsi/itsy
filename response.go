package itsy

import (
	"errors"
	"net/http"
)

type (
	// Response describes an HTTP response.
	Response struct {
		itsy       *Itsy               // The main framework instance.
		Writer     http.ResponseWriter // The HTTP response writer.
		StatusCode int                 // The HTTP status code.
	}
)

// NewResponse creates a new response instance.
func NewResponse(res http.ResponseWriter, i *Itsy) *Response {
	return &Response{
		Writer:     res,
		itsy:       i,
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

// WriteString writes a string to the response.
func (r *Response) WriteString(s string) error {
	data := []byte(s)
	written, err := r.Write(data)
	if err != nil {
		return err
	}

	if written != len(s) {
		HTTPError(StatusInternalServerError, "Response length mismatch", r.Writer, r.itsy.Logger)
		return errors.New("Response length mismatch")
	}

	return nil
}
