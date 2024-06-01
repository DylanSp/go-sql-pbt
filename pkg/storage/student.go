package storage

import (
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

func (s *Store) GetStudentByID(id uuid.UUID) (*models.Student, error) {
	query := `
		SELECT id,
			name
		FROM students
		WHERE id = :id
	`

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	returnedStudent := models.Student{}
	args := map[string]any{
		"id": id,
	}

	err = stmt.Get(&returnedStudent, args)
	if err != nil {
		return nil, err
	}

	return &returnedStudent, nil
}
