package tracker

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StubWorkoutPlanStore struct {
	workoutCalls []int
	workouts     map[string]string
	workoutPlans []WorkoutPlan
}

func (s *StubWorkoutPlanStore) AddWorkoutPlan(input WorkoutPlan) {
	s.workoutCalls = append(s.workoutCalls, 1)
}

func (s *StubWorkoutPlanStore) DeleteWorkoutPlan(name string) error {
	delete(s.workouts, name)
	return nil
}

func (s *StubWorkoutPlanStore) UpdateWorkoutPlan(input WorkoutPlan) error {
	s.workouts["pushup"] = "updated workout plan"
	return nil
}

func (s *StubWorkoutPlanStore) GetWorkoutPlanList() ([]WorkoutPlan, error) {
	return s.workoutPlans, nil
}

func TestStoreWorkoutPlan(t *testing.T) {
	t.Run("successfully adds a workout plan", func(t *testing.T) {
		store := &StubWorkoutPlanStore{}
		reqBody := []byte(`{"exerciseName": "pushup", "repititions": 10, "sets": 3, "weight": 20}`)
		request, _ := http.NewRequest(http.MethodPost, "/workout-plans/", bytes.NewBuffer(reqBody))
		response := httptest.NewRecorder()

		server := NewWorkoutServer(store)
		server.ServeHTTP(response, request)

		if store.workoutCalls == nil {
			t.Error("Expected workout plan to be stored, but it was not.")
		}
		AssertResponseStatus(t, http.StatusCreated, response.Code)
	})

	t.Run("successfully delete workout plans", func(t *testing.T) {
		store := &StubWorkoutPlanStore{
			nil,
			map[string]string{
				"pushup": "10 3 20",
				"pullup": "5 2 20",
			},
			nil,
		}
		request, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		response := httptest.NewRecorder()

		server := NewWorkoutServer(store)
		server.ServeHTTP(response, request)

		if _, exists := store.workouts["pushup"]; exists {
			t.Errorf("Expected workout plan 'pushup' to be deleted, but it still exists.")
		}
		AssertResponseStatus(t, http.StatusNoContent, response.Code)
	})

	t.Run("successfully updated workout plan", func(t *testing.T) {
		store := &StubWorkoutPlanStore{
			nil,
			map[string]string{
				"pushup": "10 3 20",
				"pullup": "5 2 20",
			},
			nil,
		}
		request, _ := http.NewRequest(http.MethodPut, "/workout-plans/", nil)
		response := httptest.NewRecorder()

		server := NewWorkoutServer(store)
		server.ServeHTTP(response, request)

		if store.workouts["pushup"] != "updated workout plan" {
			t.Errorf("Expected workout plan 'pushup' to be updated, but it was not.")
		}
		AssertResponseStatus(t, http.StatusNoContent, response.Code)
	})
}

func TestGetWorkoutPlanList(t *testing.T) {

	workoutplan := []WorkoutPlan{
		{
			ExerciseName: "pushup",
			Repititions:  10,
			Sets:         3,
			Weight:       20,
		},
		{
			ExerciseName: "pullup",
			Repititions:  5,
			Sets:         2,
			Weight:       10,
		},
	}
	store := &StubWorkoutPlanStore{
		nil,
		make(map[string]string),
		workoutplan,
	}
	server := NewWorkoutServer(store)
	request, _ := http.NewRequest(http.MethodGet, "/workouts", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	jsonResponse := fmt.Sprintf(`[{"ExerciseName":"%s","Repititions":%d,"Sets":%d,"Weight":%d},{"ExerciseName":"%s","Repititions":%d,"Sets":%d,"Weight":%d}]`,
		workoutplan[0].ExerciseName, workoutplan[0].Repititions, workoutplan[0].Sets, workoutplan[0].Weight,
		workoutplan[1].ExerciseName, workoutplan[1].Repititions, workoutplan[1].Sets, workoutplan[1].Weight)

	if jsonResponse != strings.TrimSpace(response.Body.String()) {
		t.Errorf("Expected workout plans %v, got %v", jsonResponse, response.Body.String())
	}

	AssertResponseStatus(t, http.StatusOK, response.Code)

}

func AssertResponseStatus(t *testing.T, expected, got int) {
	t.Helper()
	if expected != got {
		t.Errorf("Expected status code %d, got %d", expected, got)
	}
}
