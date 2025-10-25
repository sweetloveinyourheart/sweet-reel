package repos

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
)

type IChannelRepository interface {
	// Channel operations
	CreateChannel(ctx context.Context, channel *models.Channel) error
	GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error)
	GetChannelByHandle(ctx context.Context, handle string) (*models.Channel, error)
	GetChannelByOwnerID(ctx context.Context, ownerID uuid.UUID) (*models.Channel, error)
	UpdateChannel(ctx context.Context, channel *models.Channel) error
	DeleteChannel(ctx context.Context, id uuid.UUID) error
	ListChannels(ctx context.Context, limit, offset int) ([]*models.Channel, error)

	// Channel statistics operations
	IncrementSubscriberCount(ctx context.Context, channelID uuid.UUID) error
	DecrementSubscriberCount(ctx context.Context, channelID uuid.UUID) error
	UpdateTotalViews(ctx context.Context, channelID uuid.UUID, totalViews int64) error
	UpdateTotalVideos(ctx context.Context, channelID uuid.UUID, totalVideos int) error
	IncrementTotalVideos(ctx context.Context, channelID uuid.UUID) error
	DecrementTotalVideos(ctx context.Context, channelID uuid.UUID) error

	// Channel subscription operations
	CreateSubscription(ctx context.Context, subscription *models.ChannelSubscription) error
	GetSubscription(ctx context.Context, channelID, subscriberID uuid.UUID) (*models.ChannelSubscription, error)
	DeleteSubscription(ctx context.Context, channelID, subscriberID uuid.UUID) error
	GetSubscriptionsByChannelID(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*models.ChannelSubscription, error)
	GetSubscriptionsBySubscriberID(ctx context.Context, subscriberID uuid.UUID, limit, offset int) ([]*models.ChannelSubscription, error)
	IsSubscribed(ctx context.Context, channelID, subscriberID uuid.UUID) (bool, error)
	GetSubscriberCount(ctx context.Context, channelID uuid.UUID) (int64, error)
}

type ChannelRepository struct {
	Tx db.DbOrTx
}

func NewChannelRepository(tx db.DbOrTx) IChannelRepository {
	return &ChannelRepository{
		Tx: tx,
	}
}

// Channel operations

func (r *ChannelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	query := `
		INSERT INTO channels (id, owner_id, name, handle, description, banner_url, subscriber_count, total_views, total_videos)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.Tx.Exec(ctx, query,
		channel.ID, channel.OwnerID, channel.Name, channel.Handle,
		channel.Description, channel.BannerURL, channel.SubscriberCount,
		channel.TotalViews, channel.TotalVideos)
	return err
}

func (r *ChannelRepository) GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	query := `
		SELECT id, owner_id, name, handle, description, banner_url, subscriber_count, total_views, total_videos, created_at, updated_at
		FROM channels WHERE id = $1`

	channel := &models.Channel{}
	err := r.Tx.QueryRow(ctx, query, id).Scan(
		&channel.ID, &channel.OwnerID, &channel.Name, &channel.Handle,
		&channel.Description, &channel.BannerURL, &channel.SubscriberCount,
		&channel.TotalViews, &channel.TotalVideos,
		&channel.CreatedAt, &channel.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (r *ChannelRepository) GetChannelByHandle(ctx context.Context, handle string) (*models.Channel, error) {
	query := `
		SELECT id, owner_id, name, handle, description, banner_url, subscriber_count, total_views, total_videos, created_at, updated_at
		FROM channels WHERE handle = $1`

	channel := &models.Channel{}
	err := r.Tx.QueryRow(ctx, query, handle).Scan(
		&channel.ID, &channel.OwnerID, &channel.Name, &channel.Handle,
		&channel.Description, &channel.BannerURL, &channel.SubscriberCount,
		&channel.TotalViews, &channel.TotalVideos,
		&channel.CreatedAt, &channel.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (r *ChannelRepository) GetChannelByOwnerID(ctx context.Context, ownerID uuid.UUID) (*models.Channel, error) {
	query := `
		SELECT id, owner_id, name, handle, description, banner_url, subscriber_count, total_views, total_videos, created_at, updated_at
		FROM channels WHERE owner_id = $1`

	channel := &models.Channel{}
	err := r.Tx.QueryRow(ctx, query, ownerID).Scan(
		&channel.ID, &channel.OwnerID, &channel.Name, &channel.Handle,
		&channel.Description, &channel.BannerURL, &channel.SubscriberCount,
		&channel.TotalViews, &channel.TotalVideos,
		&channel.CreatedAt, &channel.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (r *ChannelRepository) UpdateChannel(ctx context.Context, channel *models.Channel) error {
	query := `
		UPDATE channels
		SET name = $2, handle = $3, description = $4, banner_url = $5,
		    subscriber_count = $6, total_views = $7, total_videos = $8, updated_at = NOW()
		WHERE id = $1`

	_, err := r.Tx.Exec(ctx, query,
		channel.ID, channel.Name, channel.Handle,
		channel.Description, channel.BannerURL, channel.SubscriberCount,
		channel.TotalViews, channel.TotalVideos)
	return err
}

func (r *ChannelRepository) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM channels WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id)
	return err
}

func (r *ChannelRepository) ListChannels(ctx context.Context, limit, offset int) ([]*models.Channel, error) {
	query := `
		SELECT id, owner_id, name, handle, description, banner_url, subscriber_count, total_views, total_videos, created_at, updated_at
		FROM channels
		ORDER BY subscriber_count DESC, created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.Tx.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(
			&channel.ID, &channel.OwnerID, &channel.Name, &channel.Handle,
			&channel.Description, &channel.BannerURL, &channel.SubscriberCount,
			&channel.TotalViews, &channel.TotalVideos,
			&channel.CreatedAt, &channel.UpdatedAt)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, rows.Err()
}

// Channel statistics operations

func (r *ChannelRepository) IncrementSubscriberCount(ctx context.Context, channelID uuid.UUID) error {
	query := `UPDATE channels SET subscriber_count = subscriber_count + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, channelID)
	return err
}

func (r *ChannelRepository) DecrementSubscriberCount(ctx context.Context, channelID uuid.UUID) error {
	query := `UPDATE channels SET subscriber_count = GREATEST(subscriber_count - 1, 0), updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, channelID)
	return err
}

func (r *ChannelRepository) UpdateTotalViews(ctx context.Context, channelID uuid.UUID, totalViews int64) error {
	query := `UPDATE channels SET total_views = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, channelID, totalViews)
	return err
}

func (r *ChannelRepository) UpdateTotalVideos(ctx context.Context, channelID uuid.UUID, totalVideos int) error {
	query := `UPDATE channels SET total_videos = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, channelID, totalVideos)
	return err
}

func (r *ChannelRepository) IncrementTotalVideos(ctx context.Context, channelID uuid.UUID) error {
	query := `UPDATE channels SET total_videos = total_videos + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, channelID)
	return err
}

func (r *ChannelRepository) DecrementTotalVideos(ctx context.Context, channelID uuid.UUID) error {
	query := `UPDATE channels SET total_videos = GREATEST(total_videos - 1, 0), updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, channelID)
	return err
}

// Channel subscription operations

func (r *ChannelRepository) CreateSubscription(ctx context.Context, subscription *models.ChannelSubscription) error {
	query := `
		INSERT INTO channel_subscriptions (id, channel_id, subscriber_id)
		VALUES ($1, $2, $3)`

	_, err := r.Tx.Exec(ctx, query,
		subscription.ID, subscription.ChannelID, subscription.SubscriberID)
	return err
}

func (r *ChannelRepository) GetSubscription(ctx context.Context, channelID, subscriberID uuid.UUID) (*models.ChannelSubscription, error) {
	query := `
		SELECT id, channel_id, subscriber_id, subscribed_at
		FROM channel_subscriptions
		WHERE channel_id = $1 AND subscriber_id = $2`

	subscription := &models.ChannelSubscription{}
	err := r.Tx.QueryRow(ctx, query, channelID, subscriberID).Scan(
		&subscription.ID, &subscription.ChannelID, &subscription.SubscriberID,
		&subscription.SubscribedAt)

	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (r *ChannelRepository) DeleteSubscription(ctx context.Context, channelID, subscriberID uuid.UUID) error {
	query := `DELETE FROM channel_subscriptions WHERE channel_id = $1 AND subscriber_id = $2`
	_, err := r.Tx.Exec(ctx, query, channelID, subscriberID)
	return err
}

func (r *ChannelRepository) GetSubscriptionsByChannelID(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*models.ChannelSubscription, error) {
	query := `
		SELECT id, channel_id, subscriber_id, subscribed_at
		FROM channel_subscriptions
		WHERE channel_id = $1
		ORDER BY subscribed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.Tx.Query(ctx, query, channelID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.ChannelSubscription
	for rows.Next() {
		subscription := &models.ChannelSubscription{}
		err := rows.Scan(
			&subscription.ID, &subscription.ChannelID, &subscription.SubscriberID,
			&subscription.SubscribedAt)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, rows.Err()
}

func (r *ChannelRepository) GetSubscriptionsBySubscriberID(ctx context.Context, subscriberID uuid.UUID, limit, offset int) ([]*models.ChannelSubscription, error) {
	query := `
		SELECT id, channel_id, subscriber_id, subscribed_at
		FROM channel_subscriptions
		WHERE subscriber_id = $1
		ORDER BY subscribed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.Tx.Query(ctx, query, subscriberID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.ChannelSubscription
	for rows.Next() {
		subscription := &models.ChannelSubscription{}
		err := rows.Scan(
			&subscription.ID, &subscription.ChannelID, &subscription.SubscriberID,
			&subscription.SubscribedAt)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, rows.Err()
}

func (r *ChannelRepository) IsSubscribed(ctx context.Context, channelID, subscriberID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM channel_subscriptions WHERE channel_id = $1 AND subscriber_id = $2)`

	var exists bool
	err := r.Tx.QueryRow(ctx, query, channelID, subscriberID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ChannelRepository) GetSubscriberCount(ctx context.Context, channelID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM channel_subscriptions WHERE channel_id = $1`

	var count int64
	err := r.Tx.QueryRow(ctx, query, channelID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
