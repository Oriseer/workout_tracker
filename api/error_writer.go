package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorWriter struct {
	ErrorMessage string
	Code         int
}

var (
	ErrUserName            = errors.New("username already exists")
	ErrInvalidUserDetails  = errors.New("invalid user details")
	ErrInvalidLoginDetails = errors.New("invalid username or password")
	ErrJWTToken            = errors.New("could not generate token")
	ErrInvalidToken        = errors.New("invalid input token")
	ErrInvalidExpredToken  = errors.New("invalid or expired token")
	ErrInvalidTokenClaims  = errors.New("invalid token claims")
)

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	errorResponse := ErrorWriter{
		Code:         code,
		ErrorMessage: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

var (
	InternalServerError = func(w http.ResponseWriter, err error) {
		writeError(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	StatusBadRequestServerError = func(w http.ResponseWriter, err error) {
		writeError(w, http.StatusBadRequest, "Bad Request: "+err.Error())
	}

	DatabaseError = func(w http.ResponseWriter, err error) {
		writeError(w, http.StatusInternalServerError, "Database Error: "+err.Error())
	}

	RequestBodyError = func(w http.ResponseWriter, err error) {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
	}
)
