-- Drop indexes
DROP INDEX IF EXISTS idx_channel_subscriptions_subscriber_id;
DROP INDEX IF EXISTS idx_channel_subscriptions_channel_id;
DROP INDEX IF EXISTS idx_channels_handle;
DROP INDEX IF EXISTS idx_channels_owner_id;

-- Drop tables
DROP TABLE IF EXISTS channel_subscriptions;
DROP TABLE IF EXISTS channels;
