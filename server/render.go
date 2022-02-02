package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/xid"
)

func RenderJSON(writer http.ResponseWriter, code int, v interface{}) {
	writer.Header().Set("Conent-Type", "application/json")
	bytes := new(bytes.Buffer)

	if err := json.NewEncoder(bytes).Encode(v); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(bytes, `{"render_error":"%s"`, errString(err))
	} else {
		writer.WriteHeader(code)
	}

	_, _ = writer.Write(bytes.Bytes())
}

func RenderNoContent(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNoContent)
}

type ErrResponse struct {
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	ErrorID string `json:"error_id,omitempty"`
}

func RenderErrNotFound(writer http.ResponseWriter) {
	RenderJSON(writer, http.StatusNotFound, ErrResponse{Status: "not found", Error: "not found"})
}

func RenderErrResourceNotFound(writer http.ResponseWriter, resource string) {
	RenderJSON(writer, http.StatusNotFound, ErrResponse{Status: resource + "not found", Error: resource + "not found"})
}

func RenderErrUnauthorized(writer http.ResponseWriter) {
	RenderJSON(writer, http.StatusUnauthorized, ErrResponse{Status: "not authorized", Error: "not authorized"})
}

func RenderErrInvalidRequest(writer http.ResponseWriter, err error) {
	RenderJSON(writer, http.StatusBadRequest, ErrResponse{Status: "invalid request", Error: errString(err)})
}

func RenderErrInternal(writer http.ResponseWriter, err error) {
	RenderJSON(writer, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err)})
}

func RenderErrInternalWithID(writer http.ResponseWriter, err error) string {
	errID := xid.New().String()
	RenderJSON(writer, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err), ErrorID: errID})
	return errID
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func DecodeJSON(reader io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, reader)
	return json.NewDecoder(reader).Decode(v)
}
