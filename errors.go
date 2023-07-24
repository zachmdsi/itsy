package itsy

import (
	"net/http"
	"strconv"
)

// Define a map of HTTP status codes to error messages.
var errors = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusMethodNotAllowed:    "Method Not Allowed",
	http.StatusInternalServerError: "Internal Server Error",
}

// HTTPError writes an error message to the response given a status code
func HTTPError(statusCode int, w http.ResponseWriter, r *http.Request) {
	// Get the status text from the map. If not found, default to "Internal Server Error".
	statusText, ok := errors[statusCode]
	if !ok {
		statusCode = http.StatusInternalServerError
		statusText = errors[statusCode]
	}

	// Write the status code and message to the response.
	w.WriteHeader(statusCode)
	w.Write([]byte(strconv.Itoa(statusCode) + " " + statusText))
}
