package tracker

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func NewDatabase() *DB {

	db := sqlx.MustConnect("postgres", "user=khel dbname=khel sslmode=disable password=plsqlkhel")

	return &DB{db}
}

func (db *DB) AddWorkoutPlan(input WorkoutPlan) {
	db.MustExec("INSERT INTO WORKOUT_PLAN (exercise_name, repititions, sets, weights) VALUES ($1, $2, $3, $4)", input.ExerciseName, input.Repititions, input.Sets, input.Weight)
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
	err := db.Select(&plans, "SELECT exercise_name, repititions, sets, weights FROM WORKOUT_PLAN")
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
	db.MustExec("UPDATE WORKOUT_PLAN SET repititions = $1, sets = $2, weights = $3 WHERE exercise_name = $4",
		input.Repititions, input.Sets, input.Weight, input.ExerciseName)
	return nil
}

func (db *DB) Close() error {
	err := db.DB.Close()

	if err != nil {
		return err
	}
	return nil
}
