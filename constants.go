package itsy

import "net/http"

const (
	// DefaultPort is the default port to listen on.
	DefaultPort = ":8080"

	// Define HTTP Methods
	GET    = http.MethodGet
	POST   = http.MethodPost
	PUT    = http.MethodPut
	PATCH  = http.MethodPatch
	DELETE = http.MethodDelete

	// Define HTTP Status Codes
	StatusOK                  = http.StatusOK                  // 200
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed    // 405
	StatusInternalServerError = http.StatusInternalServerError // 500

	// Define HTTP Header Names
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	// Define MIME Types
	MIMETextHTML = "text/html"
)

// Define a map of HTTP status codes to error messages.
var httpErrors = map[int]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusInternalServerError: "Internal Server Error",
}
