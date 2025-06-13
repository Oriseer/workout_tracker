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
	userAdded    int
	userLogged   int
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

func (s *StubWorkoutPlanStore) AddUser(userDetails UserDetails) error {
	s.userAdded++
	return nil
}

func (s *StubWorkoutPlanStore) UserLogin(loginData LoginData) (LoginData, error) {
	s.userLogged++
	return LoginData{}, nil
}

func TestStoreWorkoutPlan(t *testing.T) {
	userDetails := LoginData{
		Username: "test",
		Password: "pass",
	}

	token, _ := JwtGenerator(userDetails)

	t.Run("successfully adds a workout plan", func(t *testing.T) {
		store := &StubWorkoutPlanStore{}
		reqBody := []byte(`{"exerciseName": "pushup", "repetitions": 10, "sets": 3, "weight": 20}`)
		request, _ := http.NewRequest(http.MethodPost, "/workout-plans/", bytes.NewBuffer(reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
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
			0,
			0,
		}
		request, _ := http.NewRequest(http.MethodDelete, "/workout-plans/pushup", nil)
		request.Header.Set("Authorization", "Bearer "+token)
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
			0,
			0,
		}
		request, _ := http.NewRequest(http.MethodPut, "/workout-plans/", nil)
		request.Header.Set("Authorization", "Bearer "+token)
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
	userDetails := LoginData{
		Username: "test",
		Password: "pass",
	}

	token, _ := JwtGenerator(userDetails)

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
		0,
		0,
	}
	server := NewWorkoutServer(store)
	request, _ := http.NewRequest(http.MethodGet, "/workouts", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	jsonResponse := fmt.Sprintf(`[{"ExerciseName":"%s","Repetitions":%d,"Sets":%d,"Weight":%d},{"ExerciseName":"%s","Repetitions":%d,"Sets":%d,"Weight":%d}]`,
		workoutplan[0].ExerciseName, workoutplan[0].Repititions, workoutplan[0].Sets, workoutplan[0].Weight,
		workoutplan[1].ExerciseName, workoutplan[1].Repititions, workoutplan[1].Sets, workoutplan[1].Weight)

	if jsonResponse != strings.TrimSpace(response.Body.String()) {
		t.Errorf("Expected workout plans %v, got %v", jsonResponse, response.Body.String())
	}

	AssertResponseStatus(t, http.StatusOK, response.Code)

}

func TestUserRegistration(t *testing.T) {
	userDetails := LoginData{
		Username: "test",
		Password: "pass",
	}

	token, _ := JwtGenerator(userDetails)
	t.Run("successfully registers a user", func(t *testing.T) {
		store := &StubWorkoutPlanStore{
			nil,
			nil,
			nil,
			0,
			0,
		}

		server := NewWorkoutServer(store)
		//reqBody := []byte(`{"username": "testuser", "password": "testpass", "email": "test@gmail.com"}`)
		reqBody := []byte(`{"username": "testuser", "password": "testpass", "email": "test@gmail.com"}`)
		request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseStatus(t, http.StatusCreated, response.Code)
		if store.userAdded != 1 {
			t.Error("Expected user to be added, but it was not.")
		}
	})
}

func TestUserLogin(t *testing.T) {
	userDetails := LoginData{
		Username: "test",
		Password: "pass",
	}

	token, _ := JwtGenerator(userDetails)
	t.Run("successfully, login", func(t *testing.T) {

		store := &StubWorkoutPlanStore{
			nil,
			nil,
			nil,
			0,
			0,
		}

		server := NewWorkoutServer(store)
		reqBody := []byte(`{"username": "testuser", "password": "testpass"}`)
		request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseStatus(t, http.StatusOK, response.Code)
		if store.userLogged != 1 {
			t.Error("Expected user to be logged in, but it was not.")
		}
	})
}

func AssertResponseStatus(t *testing.T, expected, got int) {
	t.Helper()
	if expected != got {
		t.Errorf("Expected status code %d, got %d", expected, got)
	}
}
