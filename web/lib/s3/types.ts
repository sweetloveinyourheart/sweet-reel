/**
 * Type definitions for S3 operations
 */

export interface S3Config {
  accessKeyId: string;
  secretAccessKey: string;
  region: string;
  endpoint?: string; // For MinIO or custom S3-compatible endpoints
}

export interface ParsedS3Url {
  bucket: string;
  key: string;
}

