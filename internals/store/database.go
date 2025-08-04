package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver for PostgreSQL
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db ,err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=workout_db port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	fmt.Printf("Connected to database:"	)
	return db , nil
}

func MigrateFS(db *sql.DB, MigratioFS fs.FS, dir string) error {
	goose.SetBaseFS(MigratioFS);
	defer func() {
		goose.SetBaseFS(nil) 
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	} 
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	fmt.Println("Migrations applied successfully")
	return nil
}
func MigrateDown(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	} 
	err = goose.Down(db, dir)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	fmt.Println("Migrations Down successfully")
	return nil
}