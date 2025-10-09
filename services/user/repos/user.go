package repos

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	CreateIdentity(ctx context.Context, identity *models.UserIdentity) error
	GetIdentity(ctx context.Context, provider, providerUserID string) (*models.UserIdentity, error)
	GetUserWithIdentity(ctx context.Context, provider, providerUserID string) (*models.User, error)
}

type UserRepository struct {
	Tx db.DbOrTx
}

func NewUserRepository(tx db.DbOrTx) IUserRepository {
	return &UserRepository{
		Tx: tx,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, picture)
		VALUES ($1, $2, $3, $4)`

	_, err := r.Tx.Exec(ctx, query,
		user.ID, user.Email,
		user.Name, user.Picture)
	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, created_at, updated_at
		FROM users WHERE id = $1`

	user := &models.User{}
	err := r.Tx.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, created_at, updated_at
		FROM users WHERE email = $1`

	user := &models.User{}
	err := r.Tx.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateIdentity(ctx context.Context, identity *models.UserIdentity) error {
	query := `
		INSERT INTO user_identities (id, user_id, provider, provider_user_id, access_token, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.Tx.Exec(ctx, query,
		identity.ID, identity.UserID,
		identity.Provider, identity.ProviderUserID,
		identity.AccessToken, identity.RefreshToken,
		identity.ExpiresAt)
	return err
}

func (r *UserRepository) GetIdentity(ctx context.Context, provider, providerUserID string) (*models.UserIdentity, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, access_token, refresh_token, expires_at, created_at, updated_at
		FROM user_identities WHERE provider = $1 and provider_user_id = $2`

	userIdentity := &models.UserIdentity{}
	err := r.Tx.QueryRow(ctx, query, provider, providerUserID).Scan(
		&userIdentity.ID, &userIdentity.UserID, &userIdentity.Provider, &userIdentity.ProviderUserID,
		&userIdentity.AccessToken, &userIdentity.RefreshToken, &userIdentity.ExpiresAt,
		&userIdentity.CreatedAt, &userIdentity.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return userIdentity, nil
}

func (r *UserRepository) GetUserWithIdentity(ctx context.Context, provider, providerUserID string) (*models.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.picture, u.created_at, u.updated_at
        FROM users u
        JOIN user_identities ui ON u.id = ui.user_id
        WHERE ui.provider = $1 AND ui.provider_user_id = $2`

	user := &models.User{}
	err := r.Tx.QueryRow(ctx, query, provider, providerUserID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}
