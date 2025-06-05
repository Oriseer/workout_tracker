package main

import (
	"fmt"
	"log"
	"net/http"

	tracker "github.com/Oriseer/workout_tracker/internal"
)

func main() {
	db := tracker.NewDatabase()
	defer db.Close()
	server := tracker.NewWorkoutServer(db)
	fmt.Println("Starting web server on :8080")

	log.Fatal(http.ListenAndServe(":8080", server.Handler))
}
