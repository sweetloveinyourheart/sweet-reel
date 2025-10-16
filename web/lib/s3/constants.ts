
// Constants
export const S3_BUCKETS = {
  VIDEO_UPLOADED: 'video-uploaded',
  VIDEO_PROCESSED: 'video-processed',
} as const;

export const DEFAULT_EXPIRATION_SECONDS = 600; // 10 minutes

export const ALLOWED_VIDEO_TYPES = [
  'video/mp4',
  'video/quicktime',
  'video/x-msvideo',
  'video/x-matroska',
  'video/webm',
];

export const MAX_VIDEO_SIZE = 500 * 1024 * 1024; // 500 MB
