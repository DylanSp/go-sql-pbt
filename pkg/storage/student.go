package storage

import (
	"database/sql"
	"errors"

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

	returnedStudent := models.Student{}
	args := map[string]any{
		"name": name,
	}

	err = stmt.Get(&returnedStudent, args)
	if err != nil {
		return nil, err
	}

	return &returnedStudent, nil
}

func (s *Store) GetStudentByID(id uuid.UUID) (student *models.Student, found bool, err error) {
	query := `
		SELECT id,
			name
		FROM students
		WHERE id = :id
	`

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return nil, false, err
	}

	returnedStudent := models.Student{}
	args := map[string]any{
		"id": id,
	}

	err = stmt.Get(&returnedStudent, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &returnedStudent, true, nil
}
