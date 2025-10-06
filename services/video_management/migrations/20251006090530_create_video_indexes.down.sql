-- Drop indexes in reverse order

-- Drop video_thumbnails indexes
DROP INDEX IF EXISTS idx_video_thumbnails_created_at;
DROP INDEX IF EXISTS idx_video_thumbnails_video_id;

-- Drop video_variants indexes
DROP INDEX IF EXISTS idx_video_variants_video_quality;
DROP INDEX IF EXISTS idx_video_variants_video_id;

-- Drop video_manifests indexes
DROP INDEX IF EXISTS idx_video_manifests_video_id;

-- Drop videos indexes
DROP INDEX IF EXISTS idx_videos_updated_at;
DROP INDEX IF EXISTS idx_videos_created_at;
DROP INDEX IF EXISTS idx_videos_uploader_created_at;
DROP INDEX IF EXISTS idx_videos_status;
DROP INDEX IF EXISTS idx_videos_uploader_id;