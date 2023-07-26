package itsy

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type (
	Renderer interface {
		Render(w http.ResponseWriter, resource Resource) error
	}
	JSONRenderer struct{}
	XMLRenderer  struct{}
)

func (r *JSONRenderer) Render(w http.ResponseWriter, resource Resource) error {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(resource)
	if err != nil {
		HTTPError(http.StatusInternalServerError, w)
		return err
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		HTTPError(http.StatusInternalServerError, w)
		return err
	}
	return nil
}

func (r *XMLRenderer) Render(w http.ResponseWriter, resource Resource) error {
	w.Header().Set("Content-Type", "application/xml")
	xmlBytes, err := xml.Marshal(resource)
	if err != nil {
		HTTPError(http.StatusInternalServerError, w)
		return err
	}
	_, err = w.Write(xmlBytes)
	if err != nil {
		HTTPError(http.StatusInternalServerError, w)
		return err
	}
	return nil
}
