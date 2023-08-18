package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

// sendHTTPError writes an error message to the response given a status code
func (i *Itsy) sendHTTPError(statusCode int, message string, res http.ResponseWriter, logger *zap.Logger) {
	statusText, ok := httpErrors[statusCode]
	if !ok {
		statusText = httpErrors[StatusInternalServerError]
	}

	errorMessage := statusText + ": " + message
	logger.Error("HTTP Error", zap.Int("status", statusCode), zap.String("message", errorMessage))

	res.WriteHeader(statusCode)
	res.Write([]byte(errorMessage))
}
