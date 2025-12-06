package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/DavidGudovic/api_exercise/internal/store"
	"github.com/DavidGudovic/api_exercise/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

// NewWorkoutHandler Constructor
func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

// HandleGetWorkoutByID GET /workouts/{workoutID}
func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workout"})
		return
	}

	_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

// HandleCreateWorkout POST /workouts
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to decode workout"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create workout"})
		return
	}

	_ = utils.WriteJson(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

// HandleUpdateWorkout PUT /workouts/{workoutID}
func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)

	workout := store.Workout{
		ID: int(workoutID),
	}

	err = json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to decode workout"})
		return
	}

	err = wh.workoutStore.UpdateWorkout(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update workout"})
	}

	_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

// HandleDeleteWorkout DELETE /workouts/{workoutID}
func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete workout"})
		return
	}

	_ = utils.WriteJson(w, http.StatusNoContent, nil)
}

// HandleGetAllWorkouts GET /workouts
func (wh *WorkoutHandler) HandleGetAllWorkouts(w http.ResponseWriter, _ *http.Request) {
	workouts, err := wh.workoutStore.GetAllWorkouts()

	if errors.Is(err, sql.ErrNoRows) {
		_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workouts": []store.Workout{}})
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workouts"})
		return
	}

	_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workouts": workouts})
}
