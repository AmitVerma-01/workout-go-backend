package api

import (
	"encoding/json"
	"fmt"
	"go_beginner/internals/store"
	"go_beginner/utils"
	"log"
	"net/http"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error:: Reading workout ID: %v", err)
		http.Error(w, fmt.Sprintf("Invalid workout ID: %v", err), http.StatusBadRequest)
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error:: Getting workout by ID: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get workout: %v", err), http.StatusInternalServerError)
		return
	}
	if workout == nil {
		wh.logger.Printf("Error:: Workout not found for ID: %d", workoutID)
		http.Error(w, fmt.Sprintf("Workout not found for ID: %d", workoutID), http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"workout": workout,
	})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	// Logic to create a workout
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error:: Decoding workout request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error:: Creating workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create workout: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"workout": createdWorkout,
	})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error:: Reading workout ID: %v", err)
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error:: Getting existing workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get existing workout: %v", err), http.StatusInternalServerError)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"` // in seconds
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("Error:: Decoding update workout request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkout(workoutID, existingWorkout)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workout: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"workout": existingWorkout,
	})
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}
	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err != nil {
		wh.logger.Printf("Error:: Deleting workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to delete workout: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message": fmt.Sprintf("Workout with ID %d deleted successfully", workoutID),
	})
}

func (wh *WorkoutHandler) HandleGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := wh.workoutStore.GetWorkouts()
	if err != nil {
		wh.logger.Printf("Error:: Getting all workouts: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get all workouts: %v", err), http.StatusInternalServerError)
		return
	}
	if len(workouts) == 0 {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"message": "No workouts found",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"workouts": workouts,
	})
}
