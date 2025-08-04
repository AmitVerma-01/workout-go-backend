package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func setupTestDB (t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres_test port=5433")
	if err != nil {
		t.Fatalf("opening test db %v", err)
	}
	
	err = Migrate(db ,"../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db %v", err)
	}

	_ , err = db.Exec("TRUNCATE TABLE workouts, workout_entries CASCADE")
	if err != nil {
		t.Fatalf("truncating test db %v", err)
	}

	return db
}

func TestWorkoutCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name string
		workout *Workout
		wantErr bool  
	}{	
		 {
			name : "valid workout",
			workout: &Workout{
				Title:           "Morning Run ",
				Description:     "A refreshing morning run",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{ExerciseName: "Running",  Reps:  IntPtr(300), Sets: 1, WeightKg: Float64Ptr(210.0)},
				},
			}, 
			wantErr: false,
		 },
		 {
			name: "invalid workout - missing title",
			workout: &Workout{
				Description:     "A refreshing morning run",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{ExerciseName: "Running", DurationSeconds: IntPtr(30), Sets: 1},
				},
			},
			wantErr: true,
		 },
		 {
			name: "invalid workout - negative duration",
			workout: &Workout{
				Title:           "Evening Walk",
				Description:     "A relaxing evening walk",
				DurationMinutes: -20,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{ExerciseName: "Walking", DurationSeconds: IntPtr(-200), Sets: 2},
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

			retrieved, err := store.GetWorkoutByID(createdWorkout.ID)
			require.NoError(t, err)
			 
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i, entry := range tt.workout.Entries {
				assert.Equal(t, entry.ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, entry.Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, entry.Reps, retrieved.Entries[i].Reps)
				assert.Equal(t, entry.DurationSeconds, retrieved.Entries[i].DurationSeconds)
				assert.Equal(t, entry.WeightKg, retrieved.Entries[i].WeightKg)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}