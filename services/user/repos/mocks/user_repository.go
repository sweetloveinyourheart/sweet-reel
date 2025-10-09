package mocks

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateIdentity(ctx context.Context, identity *models.UserIdentity) error {
	args := m.Called(ctx, identity)
	return args.Error(0)
}

func (m *MockUserRepository) GetIdentity(ctx context.Context, provider, providerUserID string) (*models.UserIdentity, error) {
	args := m.Called(ctx, provider, providerUserID)
	if identity, ok := args.Get(0).(*models.UserIdentity); ok {
		return identity, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserWithIdentity(ctx context.Context, provider, providerUserID string) (*models.User, error) {
	args := m.Called(ctx, provider, providerUserID)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
