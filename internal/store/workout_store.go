package store

import "database/sql"

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

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	transaction, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(transaction)

	query := `
			INSERT INTO workouts(title, description, duration_minutes, calories_burned)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

	err = transaction.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		entryQuery := `
				INSERT INTO workout_entries(workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id
			`

		err = transaction.QueryRow(entryQuery, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)

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

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}

	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`

	err := pg.db.QueryRow(query, id).Scan(&workout)

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	updateQuery := `UPDATE workouts SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4 WHERE id = $5`

	_, err := pg.db.Exec(updateQuery, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	deleteQuery := `DELETE FROM workouts WHERE id = $1`

	_, err := pg.db.Exec(deleteQuery, id)

	if err != nil {
		return err
	}

	return nil
}
