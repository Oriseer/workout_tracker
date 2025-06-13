package tracker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Oriseer/workout_tracker/api"
	"github.com/Oriseer/workout_tracker/middleware"
)

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginData struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password_hash"`
}

type WorkoutPlanStore interface {
	AddWorkoutPlan(input WorkoutPlan)
	DeleteWorkoutPlan(name string) error
	UpdateWorkoutPlan(input WorkoutPlan) error
	GetWorkoutPlanList() ([]WorkoutPlan, error)
	AddUser(userDetails UserDetails) error
	UserLogin(loginData LoginData) (LoginData, error)
}

type Token struct {
	Token string `json:"token"`
}

// workoutPlanStore is a concrete implementation of the WorkoutPlanStore interface
type WorkoutServer struct {
	store WorkoutPlanStore
	http.Handler
}

type WorkoutPlan struct {
	ExerciseName string `json:"ExerciseName" db:"exercise_name"`
	Repititions  int    `json:"Repetitions" db:"repetitions"`
	Sets         int    `json:"Sets" db:"sets"`
	Weight       int    `json:"Weight" db:"weights"`
}

func NewWorkoutServer(store WorkoutPlanStore) *WorkoutServer {
	s := new(WorkoutServer)

	s.store = store

	router := http.NewServeMux()

	// Route for storing and deleting workout plans
	router.Handle("/workout-plans/", middleware.JwtAuth(http.HandlerFunc(s.storeWorkoutHandler)))
	router.Handle("/workouts", middleware.JwtAuth(http.HandlerFunc(s.getWorkoutPlanListHandler)))
	router.Handle("/auth/register", http.HandlerFunc(s.registerUserHandler))
	router.Handle("/auth/login", http.HandlerFunc(s.loginUserHandler))
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
	byteReqBody, _ := io.ReadAll(r.Body) // Read the body to ensure it can be closed later
	r.Body.Close()
	err := json.Unmarshal(byteReqBody, v)
	return err
}

func (ws *WorkoutServer) storeWorkoutPlan(w http.ResponseWriter, plan WorkoutPlan) {
	ws.store.AddWorkoutPlan(plan)
	w.WriteHeader(http.StatusCreated)
}

func (ws *WorkoutServer) deleteWorkoutPlan(w http.ResponseWriter, name string) {
	err := ws.store.DeleteWorkoutPlan(name)
	if err != nil {
		api.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ws *WorkoutServer) updateWorkoutPlan(w http.ResponseWriter, plan WorkoutPlan) {
	err := ws.store.UpdateWorkoutPlan(plan)
	if err != nil {
		api.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ws *WorkoutServer) getWorkoutPlanListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := ws.store.GetWorkoutPlanList()
	if err != nil {
		api.InternalServerError(w, err)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (ws *WorkoutServer) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	userDetails := UserDetails{}
	jsonErr := ws.jsonDecode(r, &userDetails)
	if jsonErr != nil {
		api.StatusBadRequestServerError(w, jsonErr)
		return
	}
	validationErr := validateUserDetails(userDetails)

	if validationErr != nil {
		api.StatusBadRequestServerError(w, validationErr)
	}
	err := ws.store.AddUser(userDetails)
	if err == api.ErrUserName {
		api.StatusBadRequestServerError(w, err)
		return
	} else if err != nil {
		api.DatabaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ws *WorkoutServer) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	loginData := LoginData{}
	jsonErr := ws.jsonDecode(r, &loginData)
	if jsonErr != nil {
		api.StatusBadRequestServerError(w, jsonErr)
		return
	}
	userDetails, err := ws.store.UserLogin(loginData)
	if err == api.ErrInvalidLoginDetails {
		api.StatusBadRequestServerError(w, err)
		return
	} else if err != nil {
		api.StatusBadRequestServerError(w, err)
		return
	}

	tokenStr, err := JwtGenerator(userDetails)

	if err != nil {
		api.StatusBadRequestServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Token{tokenStr})

}

func validateUserDetails(userDetails UserDetails) error {
	if userDetails.Username == "" || userDetails.Password == "" || userDetails.Email == "" {
		return api.ErrInvalidUserDetails
	}
	return nil
}
