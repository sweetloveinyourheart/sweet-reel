package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// ChannelSubscription represents a user's subscription to a channel
type ChannelSubscription struct {
	ID           uuid.UUID `db:"id" json:"id"`
	ChannelID    uuid.UUID `db:"channel_id" json:"channel_id"`
	SubscriberID uuid.UUID `db:"subscriber_id" json:"subscriber_id"`
	SubscribedAt time.Time `db:"subscribed_at" json:"subscribed_at"`
}

// ChannelSubscriptionWithDetails includes channel and subscriber information
type ChannelSubscriptionWithDetails struct {
	ChannelSubscription
	Channel    Channel `json:"channel"`
	Subscriber User    `json:"subscriber"`
}
