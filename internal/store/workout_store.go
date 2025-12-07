package store

import (
	"database/sql"
	"errors"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int) error
	GetAllWorkouts() ([]*Workout, error)
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

func (pg *PostgresWorkoutStore) GetAllWorkouts() ([]*Workout, error) {
	var workouts []*Workout

	rows, err := pg.db.Query(`SELECT id, title, description, duration_minutes, calories_burned FROM workouts`)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		workout := &Workout{}
		err = rows.Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

		if err != nil {
			return nil, err
		}

		err = pg.populateEntriesForWorkout(workout)

		if err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	return workouts, nil
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	transaction, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}

	defer func() { _ = transaction.Rollback() }()

	query := `
			INSERT INTO workouts(title, description, duration_minutes, calories_burned)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

	err = transaction.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

	if err != nil {
		return nil, err
	}

	for index := range workout.Entries {
		entryQuery := `
				INSERT INTO workout_entries(workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id
			`

		err = transaction.QueryRow(
			entryQuery,
			workout.ID,
			workout.Entries[index].ExerciseName,
			workout.Entries[index].Sets,
			workout.Entries[index].Reps,
			workout.Entries[index].DurationSeconds,
			workout.Entries[index].Weight,
			workout.Entries[index].Notes,
			workout.Entries[index].OrderIndex,
		).Scan(&workout.Entries[index].ID)

		if err != nil {
			return nil, err
		}
	}

	err = transaction.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int) (*Workout, error) {
	workout := &Workout{}

	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`

	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	err = pg.populateEntriesForWorkout(workout)

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	transaction, err := pg.db.Begin()

	if err != nil {
		return err
	}

	defer func() { _ = transaction.Rollback() }()

	updateQuery := `UPDATE workouts SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4 WHERE id = $5`

	result, err := transaction.Exec(updateQuery, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	updateEntryQuery := `UPDATE workout_entries SET exercise_name = $1, sets = $2, reps = $3, duration_seconds = $4, weight = $5, notes = $6, order_index = $7 WHERE id = $8 AND workout_id = $9`

	for _, entry := range workout.Entries {
		_, err = transaction.Exec(updateEntryQuery, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex, entry.ID, workout.ID)

		if err != nil {
			return err
		}
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int) error {
	deleteQuery := `DELETE FROM workouts WHERE id = $1`

	_, err := pg.db.Exec(deleteQuery, id)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresWorkoutStore) populateEntriesForWorkout(workout *Workout) error {
	rows, err := pg.db.Query(`SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index FROM workout_entries WHERE workout_id = $1 ORDER BY order_index`, workout.ID)

	if err != nil {
		return err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		entry := WorkoutEntry{}
		err = rows.Scan(&entry.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)

		if err != nil {
			return err
		}

		workout.Entries = append(workout.Entries, entry)
	}

	return nil
}
