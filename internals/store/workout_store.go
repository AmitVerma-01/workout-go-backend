package store

import (
	"database/sql"
	"fmt"
)

type Workout struct {
	ID              int64            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"` // in seconds
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"` // in seconds
	WeightKg        *float64 `json:"weight"`           // in kilograms
	Notes           *string  `json:"notes"`
	OrderIndex      int      `json:"order_index"` // to maintain the order of entries
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{
		db: db,
	}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(id int64, workout *Workout) error
	DeleteWorkout(id int64) error
	GetWorkouts() ([]Workout, error)
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	if workout.Title == "" {
		return nil, fmt.Errorf("workout title is required")
	}
	if workout.DurationMinutes < 0 {
		return nil, fmt.Errorf("workout duration cannot be negative")
	}

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		 INSERT INTO workouts (title, description, duration_minutes, calories_burned)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, title, description
	`
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(
		&workout.ID,
		&workout.Title,
		&workout.Description)

	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		entryQuery := `		 
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds,
			weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, exercise_name
		`
		err = tx.QueryRow(entryQuery, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps,
			entry.DurationSeconds, entry.WeightKg, entry.Notes, entry.OrderIndex).Scan(
			&entry.ID,
			&entry.ExerciseName)

		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`
	workout := &Workout{}
	err := pg.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.Title,
		&workout.Description,
		&workout.DurationMinutes,
		&workout.CaloriesBurned)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No workout found
		}
		return nil, err
	}

	entryQuery := `
		SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
		FROM workout_entries
		WHERE workout_id = $1
	`
	rows, err := pg.db.Query(entryQuery, id)
	if err == sql.ErrNoRows {
		return nil, nil // No entries found
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		entry := WorkoutEntry{}
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.WeightKg,
			&entry.Notes,
			&entry.OrderIndex)

		if err != nil {
			return nil, err
		}

		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(id int64, workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existingWorkout, err := pg.GetWorkoutByID(id)
	if err != nil {
		return fmt.Errorf("failed to get existing workout: %v", err)
	}
	if existingWorkout == nil {
		return fmt.Errorf("workout with ID %d not found", id)
	}

	query := `UPDATE workouts
		SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
		WHERE id = $5
		RETURNING id, title, description, duration_minutes, calories_burned
	`
	res, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes,
		workout.CaloriesBurned, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("workout with ID %d not found", id)
	}

	tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, id)
	for _, entry := range workout.Entries {
		entryQuery := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds,
			weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		_, err = tx.Exec(entryQuery, id, entry.ExerciseName, entry.Sets, entry.Reps,
			entry.DurationSeconds, entry.WeightKg, entry.Notes, entry.OrderIndex)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows // No workout found
		}
		return err
	}
	return nil
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, id)
	if err != nil {
		return err
	}

	res, err := tx.Exec(`DELETE FROM workouts WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("workout with ID %d not found", id)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresWorkoutStore) GetWorkouts() ([]Workout, error) { // TODO: change to apiWorkout, error) {
	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
	`
	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query workouts: %w", err)
	}
	defer rows.Close()
	var workouts []Workout
	for rows.Next() {
		workout := Workout{}
		err = rows.Scan(
			&workout.ID,
			&workout.Title,
			&workout.Description,
			&workout.DurationMinutes,
			&workout.CaloriesBurned)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout: %w", err)
		}
		// Fetch entries for each workout
		entryQuery := `			SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
			FROM workout_entries
			WHERE workout_id = $1
		`
		entryRows, err := pg.db.Query(entryQuery, workout.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query workout entries: %w", err)
		}
		defer entryRows.Close()
		for entryRows.Next() {
			entry := WorkoutEntry{}
			err = entryRows.Scan(
				&entry.ID,
				&entry.ExerciseName,
				&entry.Sets,
				&entry.Reps,
				&entry.DurationSeconds,
				&entry.WeightKg,
				&entry.Notes,
				&entry.OrderIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to scan workout entry: %w", err)
			}
			workout.Entries = append(workout.Entries, entry)
		}
		if err = entryRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over workout entries: %w", err)
		}
		workouts = append(workouts, workout)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over workouts: %w", err)
	}
	return workouts, nil
}
