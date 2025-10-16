/**
 * S3 Library - Main exports
 */

// Main S3 client
export { S3Client, getS3Client } from './s3';

// Types
export type {
  S3Config,
  ParsedS3Url,
} from './types';

// Constants
export {
  S3_BUCKETS,
  DEFAULT_EXPIRATION_SECONDS,
  ALLOWED_VIDEO_TYPES,
  MAX_VIDEO_SIZE,
} from './constants';

// Utility functions
export {
  generateS3Key,
  getPublicS3Url,
  parseS3Url,
  validateFileType,
  validateFileSize,
  formatFileSize,
} from './utils';