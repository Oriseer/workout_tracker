package tracker

import "slices"

type InMemoryStore struct {
	store []WorkoutPlan
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{[]WorkoutPlan{}}
}

func (s *InMemoryStore) AddWorkoutPlan(input WorkoutPlan) {
	s.store = append(s.store, input)
}

func (s *InMemoryStore) DeleteWorkoutPlan(name string) {
	newSlice := make([]WorkoutPlan, len(s.store))
	for i, value := range s.store {
		if name == value.ExerciseName {
			newSlice = slices.Delete(s.store, i, 1)
		}
	}

	s.store = newSlice
}

func (s *InMemoryStore) UpdateWorkoutPlan(name string) {
	newSlice := make([]WorkoutPlan, len(s.store))
	for i, value := range s.store {
		if name == value.ExerciseName {
			newSlice = slices.Replace(s.store, i, 1)
		}
	}

	s.store = newSlice
}

func (s *InMemoryStore) GetWorkoutPlanList() []WorkoutPlan {
	return s.store
}
