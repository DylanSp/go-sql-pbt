package models

import "github.com/google/uuid"

type Student struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
