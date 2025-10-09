package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// User represents the main user entity in your system.
// It’s your internal identity record.
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Picture   string    `db:"picture" json:"picture"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
