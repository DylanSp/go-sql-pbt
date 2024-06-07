package storage_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/DylanSp/go-sql-pbt/pkg/models"
	"github.com/DylanSp/go-sql-pbt/pkg/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
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

type SQLOpType string

const (
	SelectOp SQLOpType = "SELECT"
	InsertOp SQLOpType = "INSERT"
	UpdateOp SQLOpType = "UPDATE"
	DeleteOp SQLOpType = "DELETE"
)

func FuzzBasicUsage(f *testing.F) {
	dbCfg := storage.DBConfig{
		Host:   "localhost",
		Port:   "5432",
		DBName: "school",

		Username: "postgres",
		Password: "devpassword",
	}
	store, err := storage.NewStore(dbCfg)
	if err != nil {
		f.Fatal("couldn't initialize database connection: ", err)
	}

	// TODO - use f.Add() to seed corpus?
	f.Add(int64(48))

	fuzzTarget := createFuzzTarget(store)
	f.Fuzz(fuzzTarget)
}

// use a higher-order function so we don't have to recreate the store on every iteration
func createFuzzTarget(store *storage.Store) func(*testing.T, int64) {
	return func(t *testing.T, seed int64) {
		rng := rand.New(rand.NewSource(seed))

		// generate sequence of operations
		numOperations := rng.Int()
		operations := []SQLOpType{}
		for i := 0; i < numOperations; i++ {
			opType := randChoice(rng, []SQLOpType{SelectOp, InsertOp, UpdateOp, DeleteOp})
			operations = append(operations, opType)
		}

		// map functions as test oracle; behavior of SQL statements should match behavior of map
		oracle := map[uuid.UUID]string{}

		// perform operations, calling database methods and updating map, checking that results match map
		for _, opType := range operations {
			// TODO - refactor?
			// Insert/Update/Delete are pretty similar

			switch opType {
			case SelectOp:
				// randomly choose whether to query an existing record or a nonexistent record
				// either way, result from database should match test oracle
				useExistingID := randBool(rng)

				existingIDs := maps.Keys(oracle)

				var id uuid.UUID
				if useExistingID && len(existingIDs) > 0 {
					// some random existing ID
					id = randChoice(rng, existingIDs)
				} else {
					id = uuid.Must(uuid.NewRandomFromReader(rng))
				}

				// result from database
				nameFromDB, foundInDB, err := store.GetStudentByID(id)
				assert.NoError(t, err)

				// result from test oracle
				nameFromOracle, foundInOracle := oracle[id]

				// database and test oracle should match
				assert.EqualValues(t, foundInOracle, foundInDB)
				assert.EqualValues(t, nameFromOracle, nameFromDB)

			case InsertOp:
				newName := strconv.Itoa(rng.Int())

				// insert in database
				newStudent, err := store.CreateStudent(newName)
				assert.NoError(t, err)

				// insert in test oracle to keep it in sync
				oracle[newStudent.ID] = newStudent.Name

			case UpdateOp:
				// randomly choose whether to update an existing record or to update a nonexistent record
				// either way, result from database should match test oracle
				useExistingID := randBool(rng)

				existingIDs := maps.Keys(oracle)
				var id uuid.UUID
				if useExistingID && len(existingIDs) > 0 {
					// some random existing ID
					id = randChoice(rng, existingIDs)
				} else {
					id = uuid.Must(uuid.NewRandomFromReader(rng))
				}

				newName := strconv.Itoa(rng.Int())
				studentToUpdate := &models.Student{
					ID:   id,
					Name: newName,
				}

				// update database
				_, foundInDB, err := store.UpdateStudent(studentToUpdate)
				assert.NoError(t, err)

				// check if database's records match oracle's
				_, foundInOracle := oracle[id]
				assert.EqualValues(t, foundInOracle, foundInDB)

				// update value in test oracle to keep it in sync
				oracle[id] = newName

			case DeleteOp:
				// randomly choose whether to update an existing record or to update a nonexistent record
				// either way, result from database should match test oracle
				useExistingID := randBool(rng)

				existingIDs := maps.Keys(oracle)
				var id uuid.UUID
				if useExistingID && len(existingIDs) > 0 {
					// some random existing ID
					id = randChoice(rng, existingIDs)
				} else {
					id = uuid.Must(uuid.NewRandomFromReader(rng))
				}

				// delete from database
				foundInDB, err := store.DeleteStudentByID(id)
				assert.NoError(t, err)

				// check if database's records match oracle's
				_, foundInOracle := oracle[id]
				assert.EqualValues(t, foundInOracle, foundInDB)

				// delete value in test oracle to keep it in sync
				delete(oracle, id)
			}
		}
	}
}

// utility functions

func randChoice[T any](rng *rand.Rand, choices []T) T {
	index := rng.Intn(len(choices))
	return choices[index]
}

func randBool(rng *rand.Rand) bool {
	return rng.Intn(2) == 0
}
