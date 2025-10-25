-- Drop indexes
DROP INDEX IF EXISTS idx_videos_uploader_id;
DROP INDEX IF EXISTS idx_videos_view_count;
DROP INDEX IF EXISTS idx_video_views_viewed_at;
DROP INDEX IF EXISTS idx_video_views_viewer_id;
DROP INDEX IF EXISTS idx_video_views_video_id;

-- Drop table
DROP TABLE IF EXISTS video_views;

-- Remove view_count column
ALTER TABLE videos
DROP COLUMN IF EXISTS view_count;
