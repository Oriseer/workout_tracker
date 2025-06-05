package tracker_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tracker "github.com/Oriseer/workout_tracker/internal"
)

func TestIntegration(t *testing.T) {
	db := tracker.NewDatabase()

	server := tracker.NewWorkoutServer(db)

	reqBody := []byte(`{"exerciseName": "pushup", "repititions": 11, "sets": 2, "weight": 20}`)
	req, _ := http.NewRequest(http.MethodPost, "/workout-plans/", bytes.NewBuffer(reqBody))
	response := httptest.NewRecorder()
	server.ServeHTTP(response, req)

	tracker.AssertResponseStatus(t, http.StatusCreated, response.Code)
	t.Run("Get Workout Plan List", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/workouts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		responseBody := response.Body.String()
		requiredResponse := `[{"ExerciseName":"pushup","Repititions":11,"Sets":2,"Weight":20}]`
		tracker.AssertResponseStatus(t, http.StatusOK, response.Code)

		if responseBody == "" {
			t.Errorf("Expected non-empty response body, got empty")
		}

		if requiredResponse != strings.TrimSpace(responseBody) {
			t.Errorf("Expected response body '%s', got '%s'", requiredResponse, responseBody)
		}
	})

	t.Run("Update Workout Plan", func(t *testing.T) {

		reqBody := []byte(`{"exerciseName": "pushup", "repititions": 8, "sets": 4, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPut, "/workout-plans/", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusNoContent, response.Code)

	})

	t.Run("Delete Workout Plan", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusNoContent, response.Code)

	})

	t.Run("Update Workout Plan, workout not exists", func(t *testing.T) {

		reqBody := []byte(`{"exerciseName": "pushup", "repititions": 8, "sets": 3, "weight": 20}`)
		req, _ := http.NewRequest(http.MethodPut, "/workout-plans/", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusInternalServerError, response.Code)

		responseBody := tracker.ErrorWriter{}
		json.NewDecoder(response.Body).Decode(&responseBody)
		requiredError := "Internal Server Error: sql: no rows in result set"

		assertErrorMessage(t, responseBody, requiredError)

	})

	t.Run("Delete Workout Plan, workout not exists", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		tracker.AssertResponseStatus(t, http.StatusInternalServerError, response.Code)

		responseBody := tracker.ErrorWriter{}
		json.NewDecoder(response.Body).Decode(&responseBody)
		requiredError := "Internal Server Error: sql: no rows in result set"

		assertErrorMessage(t, responseBody, requiredError)

	})

}

func assertErrorMessage(t testing.TB, responseBody tracker.ErrorWriter, expectedError string) {
	t.Helper()
	if responseBody.ErrorMessage != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, responseBody.ErrorMessage)
	}
}
