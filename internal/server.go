package tracker

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type WorkoutPlanStore interface {
	AddWorkoutPlan(input WorkoutPlan)
	DeleteWorkoutPlan(name string) error
	UpdateWorkoutPlan(input WorkoutPlan) error
	GetWorkoutPlanList() ([]WorkoutPlan, error)
}

// workoutPlanStore is a concrete implementation of the WorkoutPlanStore interface
type WorkoutServer struct {
	store WorkoutPlanStore
	http.Handler
}

type WorkoutPlan struct {
	ExerciseName string `json:"ExerciseName" db:"exercise_name"`
	Repititions  int    `json:"Repititions" db:"repititions"`
	Sets         int    `json:"Sets" db:"sets"`
	Weight       int    `json:"Weight" db:"weights"`
}

func NewWorkoutServer(store WorkoutPlanStore) *WorkoutServer {
	s := new(WorkoutServer)

	s.store = store

	router := http.NewServeMux()

	// Route for storing and deleting workout plans
	router.Handle("/workout-plans/", http.HandlerFunc(s.storeWorkoutHandler))
	router.Handle("/workouts", http.HandlerFunc(s.getWorkoutPlanListHandler))
	s.Handler = router

	return s
}

func (ws *WorkoutServer) storeWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/workout-plans/")
	workoutPlan := WorkoutPlan{}

	if r.Body != nil {
		jsonErr := ws.jsonDecode(r, &workoutPlan)
		if jsonErr != nil {
			log.Println("Error decoding JSON:", jsonErr)
		}
	}

	switch r.Method {
	case http.MethodPost:
		ws.storeWorkoutPlan(w, workoutPlan)
	case http.MethodDelete:
		ws.deleteWorkoutPlan(w, name)
	case http.MethodPut:
		ws.updateWorkoutPlan(w, workoutPlan)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ws *WorkoutServer) jsonDecode(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}

func (ws *WorkoutServer) storeWorkoutPlan(w http.ResponseWriter, plan WorkoutPlan) {
	ws.store.AddWorkoutPlan(plan)
	w.WriteHeader(http.StatusCreated)
}

func (ws *WorkoutServer) deleteWorkoutPlan(w http.ResponseWriter, name string) {
	err := ws.store.DeleteWorkoutPlan(name)
	if err != nil {
		InternalServerError(w, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ws *WorkoutServer) updateWorkoutPlan(w http.ResponseWriter, plan WorkoutPlan) {
	err := ws.store.UpdateWorkoutPlan(plan)
	if err != nil {
		InternalServerError(w, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ws *WorkoutServer) getWorkoutPlanListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := ws.store.GetWorkoutPlanList()
	if err != nil {
		InternalServerError(w, err)
	}
	json.NewEncoder(w).Encode(list)
}
