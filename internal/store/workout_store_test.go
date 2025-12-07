package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE workout_entries, workouts, users RESTART IDENTITY CASCADE")

	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	err = Migrate(db, "./../../migrations")

	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body day",
				UserID:          1,
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(135.5),
						Notes:        "Warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "full body",
				Description:     "complete workout",
				UserID:          1,
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        "Keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squats",
						Sets:            4,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(150.0),
						Notes:           "Full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, len(tt.workout.Entries), len(createdWorkout.Entries))

			for i, entry := range tt.workout.Entries {
				createdEntry := createdWorkout.Entries[i]
				assert.Equal(t, entry.ExerciseName, createdEntry.ExerciseName)
				assert.Equal(t, entry.Sets, createdEntry.Sets)
				assert.Equal(t, entry.Reps, createdEntry.Reps)
				assert.Equal(t, entry.DurationSeconds, createdEntry.DurationSeconds)
				assert.Equal(t, entry.Weight, createdEntry.Weight)
				assert.Equal(t, entry.Notes, createdEntry.Notes)
				assert.Equal(t, entry.OrderIndex, createdEntry.OrderIndex)
			}

			retrieved, err := store.GetWorkoutByID(createdWorkout.ID)
			require.NoError(t, err)
			assert.Equal(t, createdWorkout, retrieved)
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}
