package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/DylanSp/go-sql-pbt/pkg/models"
	"github.com/google/uuid"
)

func (s *Store) CreateStudent(name string) (*models.Student, error) {
	query := `
		INSERT INTO students (
			name
		) VALUES (
			:name
		) RETURNING
			id,
			name
		;
	`

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	returnedStudent := &models.Student{}
	args := map[string]any{
		"name": name,
	}

	err = stmt.Get(returnedStudent, args)
	if err != nil {
		return nil, err
	}

	return returnedStudent, nil
}

func (s *Store) GetStudentByID(id uuid.UUID) (student *models.Student, found bool, err error) {
	query := `
		SELECT id,
			name
		FROM students
		WHERE id = :id
	`

	// query := `
	// 	SELECT id,
	// 		name
	// 	FROM students
	// 	ORDER BY id ASC
	// 	LIMIT 1
	// `

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return nil, false, err
	}

	returnedStudent := &models.Student{}
	args := map[string]any{
		"id": id,
	}

	err = stmt.Get(returnedStudent, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	return returnedStudent, true, nil
}

func (s *Store) UpdateStudent(student *models.Student) (updatedStudent *models.Student, found bool, err error) {
	query := `
		UPDATE students
		SET
			name = :name
		WHERE id = :id
		RETURNING
			id,
			name
		;
	`

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return nil, false, err
	}

	updatedStudent = &models.Student{}

	err = stmt.Get(updatedStudent, student)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	return updatedStudent, true, nil
}

func (s *Store) DeleteStudentByID(id uuid.UUID) (found bool, err error) {
	query := `
		DELETE
		FROM students
		WHERE id = :id
	`
	args := map[string]any{
		"id": id,
	}

	result, err := s.db.NamedExec(query, args)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	fmt.Printf("Affected %v rows\n", rowsAffected)
	return (rowsAffected > 0), nil
}
