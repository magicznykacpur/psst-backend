package api

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	errorResponseBytes, err := json.Marshal(errorResponse{Message: message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't marshal error response"))
		return
	}

	w.WriteHeader(status)
	w.Write(errorResponseBytes)
}

func respondWithJSON(w http.ResponseWriter, status int, body interface{}) {
	responseBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't marshal response body"))
		return
	}

	w.WriteHeader(status)
	w.Write(responseBytes)
}
