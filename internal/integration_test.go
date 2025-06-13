package tracker_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Oriseer/workout_tracker/api"
	tracker "github.com/Oriseer/workout_tracker/internal"
)

func TestIntegration(t *testing.T) {
	db := tracker.NewDatabase()

	server := tracker.NewWorkoutServer(db)

	reqBody := []byte(`{"username": "testuser", "password": "testpass", "email": "test@gmail.com"}`)
	request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	tracker.AssertResponseStatus(t, http.StatusCreated, response.Code)

	token := tracker.Token{}

	t.Run("User Login, token should be generated", func(t *testing.T) {
		reqBody := []byte(`{"username": "testuser", "password": "testpass"}`)
		request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		// Token save to global variable
		json.NewDecoder(response.Body).Decode(&token)

		tracker.AssertResponseStatus(t, http.StatusOK, response.Code)

	})

	t.Run("add workout plan with correct token", func(t *testing.T) {
		reqBody := []byte(`{"exerciseName": "pushup", "repetitions": 11, "sets": 2, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPost, "/workout-plans/", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token.Token)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, req)

		if response.Body.String() != "" {
			t.Errorf("Expected no response, got %q", response.Body.String())
		}

		tracker.AssertResponseStatus(t, http.StatusCreated, response.Code)
	})
	t.Run("add workout plan with incorrect token", func(t *testing.T) {
		reqBody := []byte(`{"exerciseName": "pushup", "repetitions": 11, "sets": 2, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPost, "/workout-plans/", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "dummy")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, req)

		body := api.ErrorWriter{}

		json.NewDecoder(response.Body).Decode(&body)

		expectedError := "Bad Request: " + api.ErrInvalidToken.Error()

		if expectedError != body.ErrorMessage {
			t.Errorf("Expected error message %q, got %q", expectedError, body.ErrorMessage)
		}

		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Get Workout Plan List with correct token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/workouts", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		responseBody := response.Body.String()
		requiredResponse := `[{"ExerciseName":"pushup","Repetitions":11,"Sets":2,"Weight":20}]`
		tracker.AssertResponseStatus(t, http.StatusOK, response.Code)

		if responseBody == "" {
			t.Errorf("Expected non-empty response body, got empty")
		}

		if requiredResponse != strings.TrimSpace(responseBody) {
			t.Errorf("Expected response body '%s', got '%s'", requiredResponse, responseBody)
		}
	})
	t.Run("Get Workout Plan List with incorrect token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/workouts", nil)
		req.Header.Set("Authorization", "dummy")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		body := api.ErrorWriter{}

		json.NewDecoder(response.Body).Decode(&body)

		expectedError := "Bad Request: " + api.ErrInvalidToken.Error()

		if expectedError != body.ErrorMessage {
			t.Errorf("Expected error message %q, got %q", expectedError, body.ErrorMessage)
		}
		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Update Workout Plan with correct token", func(t *testing.T) {

		reqBody := []byte(`{"exerciseName": "pushup", "repetitions": 8, "sets": 4, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPut, "/workout-plans/", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token.Token)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		if response.Body.String() != "" {
			t.Errorf("Got response %v, expected no response", response.Body.String())
		}

		tracker.AssertResponseStatus(t, http.StatusNoContent, response.Code)

	})
	t.Run("Update Workout Plan with incorrect token", func(t *testing.T) {

		reqBody := []byte(`{"exerciseName": "pushup", "repetitions": 8, "sets": 4, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPut, "/workout-plans/", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "dummy")

		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		body := api.ErrorWriter{}

		json.NewDecoder(response.Body).Decode(&body)

		expectedError := "Bad Request: " + api.ErrInvalidToken.Error()

		if expectedError != body.ErrorMessage {
			t.Errorf("Expected error message %q, got %q", expectedError, body.ErrorMessage)
		}
		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)

	})

	t.Run("Delete Workout Plan with correct token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusNoContent, response.Code)

	})
	t.Run("Delete Workout Plan with incorrect token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		req.Header.Set("Authorization", "dummy")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		body := api.ErrorWriter{}

		json.NewDecoder(response.Body).Decode(&body)

		expectedError := "Bad Request: " + api.ErrInvalidToken.Error()

		if expectedError != body.ErrorMessage {
			t.Errorf("Expected error message %q, got %q", expectedError, body.ErrorMessage)
		}
		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)

	})

	t.Run("Update Workout Plan, workout not exists", func(t *testing.T) {

		reqBody := []byte(`{"exerciseName": "pushup", "repetitions": 8, "sets": 3, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPut, "/workout-plans/", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token.Token)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusInternalServerError, response.Code)

		responseBody := api.ErrorWriter{}
		json.NewDecoder(response.Body).Decode(&responseBody)
		requiredError := "Internal Server Error: sql: no rows in result set"

		assertErrorMessage(t, responseBody, requiredError)

	})

	t.Run("Delete Workout Plan, workout not exists", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusInternalServerError, response.Code)

		responseBody := api.ErrorWriter{}
		json.NewDecoder(response.Body).Decode(&responseBody)
		requiredError := "Internal Server Error: sql: no rows in result set"

		assertErrorMessage(t, responseBody, requiredError)

	})

	t.Run("Add new user with existing username", func(t *testing.T) {

		reqBody := []byte(`{"username": "testuser", "password": "testpass", "email": "test@gmail.com"}`)
		request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		responseBody := api.ErrorWriter{}

		json.NewDecoder(response.Body).Decode(&responseBody)

		expectedError := fmt.Sprintf("Bad Request: %s", api.ErrUserName.Error())

		if responseBody.ErrorMessage != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, responseBody.ErrorMessage)
		}

		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Add new user with invalid request body, missing field", func(t *testing.T) {
		reqBody := []byte(`{"username": "testuser", "password": "testpass"}`) // email is missing
		request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseBody := api.ErrorWriter{}
		json.NewDecoder(response.Body).Decode(&responseBody)

		expectedError := "Bad Request: " + api.ErrInvalidUserDetails.Error()

		if expectedError != responseBody.ErrorMessage {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, responseBody.ErrorMessage)
		}

		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)
	})

	t.Run("User login, incorrect password", func(t *testing.T) {

		reqBody := []byte(`{"username": "testuser", "password": "incorrect"}`)
		request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseBody := api.ErrorWriter{}
		json.NewDecoder(response.Body).Decode(&responseBody)

		expectedError := "Bad Request: " + api.ErrInvalidLoginDetails.Error()

		if responseBody.ErrorMessage != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, responseBody.ErrorMessage)
		}

		tracker.AssertResponseStatus(t, http.StatusBadRequest, response.Code)
	})

}

func assertErrorMessage(t testing.TB, responseBody api.ErrorWriter, expectedError string) {
	t.Helper()
	if responseBody.ErrorMessage != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, responseBody.ErrorMessage)
	}
}
