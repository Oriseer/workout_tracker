package main

import (
	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db := sqlx.MustConnect("postgres", "user=khel dbname=khel sslmode=disable password=plsqlkhel")

	tx := db.MustBegin()

	tx.MustExec("INSERT INTO EXERCISES (id, exercise_name, description, category) VALUES ($1, $2, $3, $4)", 1, "pushup", "bodyweight exercise", "strength")
	tx.MustExec("INSERT INTO EXERCISES (id, exercise_name, description, category) VALUES ($1, $2, $3, $4)", 2, "pullup", "bodyweight exercise", "strength")
	tx.MustExec("INSERT INTO EXERCISES (id, exercise_name, description, category) VALUES ($1, $2, $3, $4)", 3, "curlup", "bodyweight exercise", "strength")
	tx.MustExec("INSERT INTO EXERCISES (id, exercise_name, description, category) VALUES ($1, $2, $3, $4)", 4, "bench press", "barbell", "strength")
	tx.Commit()
}
