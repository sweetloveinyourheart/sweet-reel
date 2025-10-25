-- 1. Channels table
CREATE TABLE channels (
    id                  UUID            NOT NULL,
    owner_id            UUID            NOT NULL,
    name                VARCHAR(255)    NOT NULL,
    handle              VARCHAR(100)    NOT NULL,       -- e.g. @username
    description         TEXT,
    banner_url          TEXT,
    subscriber_count    INTEGER         DEFAULT 0,      -- cached count from channel_subscriptions
    total_views         BIGINT          DEFAULT 0,      -- cached aggregate from video_management service
    total_videos        INTEGER         DEFAULT 0,      -- cached count of videos
    created_at          TIMESTAMP       DEFAULT NOW(),
    updated_at          TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_channel_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (handle),
    UNIQUE (owner_id)                                   -- one channel per user
);

-- 2. Channel Subscriptions table
CREATE TABLE channel_subscriptions (
    id                  UUID            NOT NULL,
    channel_id          UUID            NOT NULL,
    subscriber_id       UUID            NOT NULL,
    subscribed_at       TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_subscription_channel FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_subscriber FOREIGN KEY (subscriber_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (channel_id, subscriber_id)
);

-- 3. Create indexes for better query performance
CREATE INDEX idx_channels_owner_id ON channels(owner_id);
CREATE INDEX idx_channels_handle ON channels(handle);
CREATE INDEX idx_channel_subscriptions_channel_id ON channel_subscriptions(channel_id);
CREATE INDEX idx_channel_subscriptions_subscriber_id ON channel_subscriptions(subscriber_id);
