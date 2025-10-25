package mocks

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
)

type MockChannelRepository struct {
	mock.Mock
}

// Channel operations

func (m *MockChannelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}

func (m *MockChannelRepository) GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	args := m.Called(ctx, id)
	if channel, ok := args.Get(0).(*models.Channel); ok {
		return channel, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChannelRepository) GetChannelByHandle(ctx context.Context, handle string) (*models.Channel, error) {
	args := m.Called(ctx, handle)
	if channel, ok := args.Get(0).(*models.Channel); ok {
		return channel, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChannelRepository) GetChannelByOwnerID(ctx context.Context, ownerID uuid.UUID) (*models.Channel, error) {
	args := m.Called(ctx, ownerID)
	if channel, ok := args.Get(0).(*models.Channel); ok {
		return channel, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChannelRepository) UpdateChannel(ctx context.Context, channel *models.Channel) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}

func (m *MockChannelRepository) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChannelRepository) ListChannels(ctx context.Context, limit, offset int) ([]*models.Channel, error) {
	args := m.Called(ctx, limit, offset)
	if channels, ok := args.Get(0).([]*models.Channel); ok {
		return channels, args.Error(1)
	}
	return nil, args.Error(1)
}

// Channel statistics operations

func (m *MockChannelRepository) IncrementSubscriberCount(ctx context.Context, channelID uuid.UUID) error {
	args := m.Called(ctx, channelID)
	return args.Error(0)
}

func (m *MockChannelRepository) DecrementSubscriberCount(ctx context.Context, channelID uuid.UUID) error {
	args := m.Called(ctx, channelID)
	return args.Error(0)
}

func (m *MockChannelRepository) UpdateTotalViews(ctx context.Context, channelID uuid.UUID, totalViews int64) error {
	args := m.Called(ctx, channelID, totalViews)
	return args.Error(0)
}

func (m *MockChannelRepository) UpdateTotalVideos(ctx context.Context, channelID uuid.UUID, totalVideos int) error {
	args := m.Called(ctx, channelID, totalVideos)
	return args.Error(0)
}

func (m *MockChannelRepository) IncrementTotalVideos(ctx context.Context, channelID uuid.UUID) error {
	args := m.Called(ctx, channelID)
	return args.Error(0)
}

func (m *MockChannelRepository) DecrementTotalVideos(ctx context.Context, channelID uuid.UUID) error {
	args := m.Called(ctx, channelID)
	return args.Error(0)
}

// Channel subscription operations

func (m *MockChannelRepository) CreateSubscription(ctx context.Context, subscription *models.ChannelSubscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockChannelRepository) GetSubscription(ctx context.Context, channelID, subscriberID uuid.UUID) (*models.ChannelSubscription, error) {
	args := m.Called(ctx, channelID, subscriberID)
	if subscription, ok := args.Get(0).(*models.ChannelSubscription); ok {
		return subscription, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChannelRepository) DeleteSubscription(ctx context.Context, channelID, subscriberID uuid.UUID) error {
	args := m.Called(ctx, channelID, subscriberID)
	return args.Error(0)
}

func (m *MockChannelRepository) GetSubscriptionsByChannelID(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*models.ChannelSubscription, error) {
	args := m.Called(ctx, channelID, limit, offset)
	if subscriptions, ok := args.Get(0).([]*models.ChannelSubscription); ok {
		return subscriptions, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChannelRepository) GetSubscriptionsBySubscriberID(ctx context.Context, subscriberID uuid.UUID, limit, offset int) ([]*models.ChannelSubscription, error) {
	args := m.Called(ctx, subscriberID, limit, offset)
	if subscriptions, ok := args.Get(0).([]*models.ChannelSubscription); ok {
		return subscriptions, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChannelRepository) IsSubscribed(ctx context.Context, channelID, subscriberID uuid.UUID) (bool, error) {
	args := m.Called(ctx, channelID, subscriberID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChannelRepository) GetSubscriberCount(ctx context.Context, channelID uuid.UUID) (int64, error) {
	args := m.Called(ctx, channelID)
	return args.Get(0).(int64), args.Error(1)
}

// Ensure MockChannelRepository implements IChannelRepository
var _ repos.IChannelRepository = (*MockChannelRepository)(nil)
