-- Remove view_count column
ALTER TABLE video_manifests
DROP COLUMN IF EXISTS quality;
