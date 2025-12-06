package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("db:open %w", err)
	}

	fmt.Println("Database connection established")

	return db, nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("migration error: goose:setdialect %w", err)
	}

	err = goose.Up(db, dir)

	if err != nil {
		return fmt.Errorf("migration error: goose:up %w", err)
	}

	fmt.Println("Database migration completed")

	return nil
}
