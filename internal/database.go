package tracker

import (
	"fmt"
	"os"

	"github.com/Oriseer/workout_tracker/api"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sqlx.DB
}

func NewDatabase() *DB {

	godotenv.Load()

	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")

	db := sqlx.MustConnect("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", db_user, db_name, db_password))
	return &DB{db}
}

func (db *DB) AddWorkoutPlan(input WorkoutPlan) {
	db.MustExec("INSERT INTO WORKOUT_PLAN (exercise_name, repetitions, sets, weights) VALUES ($1, $2, $3, $4)", input.ExerciseName, input.Repititions, input.Sets, input.Weight)
}

func (db *DB) DeleteWorkoutPlan(name string) error {

	workoutPlan := WorkoutPlan{}
	// Select the workout plan to ensure it exists before deleting
	err := db.Get(&workoutPlan, "SELECT exercise_name FROM WORKOUT_PLAN WHERE exercise_name = $1", name)
	if err != nil {
		return err
	}
	db.MustExec("DELETE FROM WORKOUT_PLAN WHERE exercise_name = $1", name)
	return nil
}

func (db *DB) GetWorkoutPlanList() ([]WorkoutPlan, error) {
	var plans []WorkoutPlan
	err := db.Select(&plans, "SELECT exercise_name, repetitions, sets, weights FROM WORKOUT_PLAN")
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (db *DB) UpdateWorkoutPlan(input WorkoutPlan) error {
	workoutPlan := WorkoutPlan{}
	// Select the workout plan to ensure it exists before deleting
	err := db.Get(&workoutPlan, "SELECT exercise_name FROM WORKOUT_PLAN WHERE exercise_name = $1", input.ExerciseName)
	if err != nil {
		return err
	}
	db.MustExec("UPDATE WORKOUT_PLAN SET repetitions = $1, sets = $2, weights = $3 WHERE exercise_name = $4",
		input.Repititions, input.Sets, input.Weight, input.ExerciseName)
	return nil
}

func (db *DB) AddUser(userDetails UserDetails) error {
	user := UserDetails{}

	// Check if the user already exists
	err := db.Get(&user, "SELECT username FROM USERS WHERE username = $1", userDetails.Username)

	if err == nil {
		return api.ErrUserName // User already exists
	}

	pass_hash, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO USERS (username, password_hash, email) VALUES ($1, $2, $3)", userDetails.Username, pass_hash, userDetails.Email)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UserLogin(loginData LoginData) (LoginData, error) {
	userDetails := LoginData{}
	// check if user valid
	err := db.Get(&userDetails, "SELECT username, password_hash FROM USERS WHERE username = $1", loginData.Username)

	if err != nil {
		// Return sql.ErrNoRow
		return LoginData{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(loginData.Password)); err != nil {
		return LoginData{}, api.ErrInvalidLoginDetails
	}

	return userDetails, nil

}

func (db *DB) Close() error {
	err := db.DB.Close()

	if err != nil {
		return err
	}
	return nil
}
