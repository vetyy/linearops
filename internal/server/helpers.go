package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleSuccess(w http.ResponseWriter, jsonBody interface{}) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(jsonBody); err != nil {
		handleInternalError(w, err, "failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Error("could not write JSON response", err)
	}
}

func handleInternalError(w http.ResponseWriter, err error, format string, args ...interface{}) {
	handleError(w, http.StatusInternalServerError, err, format, args...)
}

func handleError(w http.ResponseWriter, status int, err error, format string, args ...interface{}) {
	if err != nil {
		args = append(args, err.Error())
		format += ": %v"
		log.Errorf(format, args...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encoded structure is known, any error must be caused by writing the response
	resp := struct {
		Error string `json:"error"`
	}{
		Error: fmt.Sprintf(format, args...),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("could not write JSON response", err)
	}
}
