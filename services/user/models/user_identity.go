package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// UserIdentity represents a mapping between your internal user
// and an external OAuth provider (e.g. Google, GitHub).
type UserIdentity struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         uuid.UUID  `db:"user_id" json:"user_id"`
	Provider       string     `db:"provider" json:"provider"`                   // "google", "github", etc.
	ProviderUserID string     `db:"provider_user_id" json:"provider_user_id"`   // e.g., Google sub
	AccessToken    *string    `db:"access_token" json:"access_token,omitempty"` // optional
	RefreshToken   *string    `db:"refresh_token" json:"refresh_token,omitempty"`
	ExpiresAt      *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}
