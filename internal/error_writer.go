package tracker

import (
	"encoding/json"
	"net/http"
)

type ErrorWriter struct {
	ErrorMessage string
	Code         int
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	errorResponse := ErrorWriter{
		ErrorMessage: message,
		Code:         code,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

var (
	InternalServerError = func(w http.ResponseWriter, err error) {
		writeError(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	RequestBodyError = func(w http.ResponseWriter, err error) {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
	}
)
