package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=workouts_db sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("db:open %w", err)
	}

	fmt.Println("Database connection established")

	return db, nil
}
