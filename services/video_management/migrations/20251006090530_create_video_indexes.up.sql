-- Index on videos.uploader_id for fast queries by uploader
CREATE INDEX idx_videos_uploader_id ON videos (uploader_id);

-- Index on videos.status for filtering by status
CREATE INDEX idx_videos_status ON videos (status);

-- Composite index on videos.uploader_id and created_at for paginated queries by uploader
CREATE INDEX idx_videos_uploader_created_at ON videos (uploader_id, created_at DESC);

-- Index on videos.created_at for general ordering and time-based queries
CREATE INDEX idx_videos_created_at ON videos (created_at DESC);

-- Index on videos.updated_at for recently updated content
CREATE INDEX idx_videos_updated_at ON videos (updated_at DESC);

-- Index on video_manifests.video_id for fast lookups of manifests by video
CREATE INDEX idx_video_manifests_video_id ON video_manifests (video_id);

-- Index on video_variants.video_id for fast lookups of variants by video
CREATE INDEX idx_video_variants_video_id ON video_variants (video_id);

-- Composite index on video_variants.video_id and quality for specific quality lookups
CREATE INDEX idx_video_variants_video_quality ON video_variants (video_id, quality);

-- Index on video_thumbnails.video_id for fast lookups of thumbnails by video
CREATE INDEX idx_video_thumbnails_video_id ON video_thumbnails (video_id);

-- Index on video_thumbnails.created_at for ordering thumbnails
CREATE INDEX idx_video_thumbnails_created_at ON video_thumbnails (created_at);