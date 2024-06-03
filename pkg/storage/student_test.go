package storage_test

import (
	"testing"

	"github.com/DylanSp/go-sql-pbt/pkg/models"
	"github.com/DylanSp/go-sql-pbt/pkg/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasicUsage(t *testing.T) {
	t.Run("Run all storage methods with a new student", func(t *testing.T) {
		dbCfg := storage.DBConfig{
			Host:   "localhost",
			Port:   "5432",
			DBName: "school",

			Username: "postgres",
			Password: "devpassword",
		}
		store, err := storage.NewStore(dbCfg)
		require.NoError(t, err)

		// CreateStudent()
		initialStudentName := "Alice"
		createdStudent, err := store.CreateStudent(initialStudentName)
		require.NoError(t, err)
		require.NotEqualValues(t, uuid.Nil, createdStudent.ID)
		require.EqualValues(t, "Alice", createdStudent.Name)

		// GetStudentByID()
		fetchedStudent, found, err := store.GetStudentByID(createdStudent.ID)
		require.NoError(t, err)
		require.True(t, found)
		require.EqualValues(t, "Alice", fetchedStudent.Name)

		// UpdateStudent()
		studentWithNameChanged := &models.Student{
			ID:   createdStudent.ID,
			Name: "Bob",
		}
		updatedStudent, found, err := store.UpdateStudent(studentWithNameChanged)
		require.NoError(t, err)
		require.True(t, found)
		require.EqualValues(t, "Bob", updatedStudent.Name)

		// DeleteStudent()
		found, err = store.DeleteStudentByID(updatedStudent.ID)
		require.NoError(t, err)
		require.True(t, found)

		// check that student was deleted
		_, found, err = store.GetStudentByID(updatedStudent.ID)
		require.NoError(t, err)
		require.False(t, found)
	})

	t.Run("Check that storage methods behave correctly when not finding students", func(t *testing.T) {
		dbCfg := storage.DBConfig{
			Host:   "localhost",
			Port:   "5432",
			DBName: "school",

			Username: "postgres",
			Password: "devpassword",
		}
		store, err := storage.NewStore(dbCfg)
		require.NoError(t, err)

		_, found, err := store.GetStudentByID(uuid.New())
		require.NoError(t, err)
		require.False(t, found)

		nonexistentStudent := &models.Student{
			ID: uuid.New(),
		}
		_, found, err = store.UpdateStudent(nonexistentStudent)
		require.NoError(t, err)
		require.False(t, found)

		found, err = store.DeleteStudentByID(uuid.New())
		require.NoError(t, err)
		require.False(t, found)
	})
}
