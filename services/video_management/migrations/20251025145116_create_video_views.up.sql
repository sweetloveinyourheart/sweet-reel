-- 1. Add view_count column to videos table
ALTER TABLE videos
ADD COLUMN view_count BIGINT DEFAULT 0;

-- 2. Create video_views table for detailed view tracking
CREATE TABLE video_views (
    id                  UUID            NOT NULL,
    video_id            UUID            NOT NULL,
    viewer_id           UUID,                           -- NULL for anonymous viewers
    viewed_at           TIMESTAMP       DEFAULT NOW(),
    watch_duration      INT,                            -- seconds watched (for analytics)
    ip_address          VARCHAR(45),                    -- IPv4/IPv6 for deduplication
    user_agent          TEXT,

    PRIMARY KEY (id),
    CONSTRAINT fk_video_view FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE
);

-- 3. Create indexes for efficient querying
CREATE INDEX idx_video_views_video_id ON video_views(video_id);
CREATE INDEX idx_video_views_viewer_id ON video_views(viewer_id) WHERE viewer_id IS NOT NULL;
CREATE INDEX idx_video_views_viewed_at ON video_views(viewed_at);
CREATE INDEX idx_videos_view_count ON videos(view_count);

-- 4. Create index on uploader_id for aggregating channel statistics
CREATE INDEX idx_videos_uploader_id ON videos(uploader_id);
