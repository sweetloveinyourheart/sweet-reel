/**
 * Utility functions for S3 operations
 */

import type { ParsedS3Url } from './types';

/**
 * Helper function to generate a unique S3 key for a file
 */
export function generateS3Key(
  prefix: string,
  filename: string,
  userId?: string
): string {
  const timestamp = Date.now();
  const randomString = Math.random().toString(36).substring(2, 15);
  const sanitizedFilename = filename.replace(/[^a-zA-Z0-9.-]/g, '_');
  
  if (userId) {
    return `${prefix}/${userId}/${timestamp}-${randomString}-${sanitizedFilename}`;
  }
  
  return `${prefix}/${timestamp}-${randomString}-${sanitizedFilename}`;
}

/**
 * Get the public URL for an S3 object (for publicly accessible buckets)
 */
export function getPublicS3Url(
  bucket: string,
  key: string,
  region?: string,
  endpoint?: string
): string {
  if (endpoint) {
    // For MinIO or custom endpoints
    return `${endpoint}/${bucket}/${key}`;
  }
  
  const s3Region = region || process.env.AWS_REGION || 'us-east-1';
  return `https://${bucket}.s3.${s3Region}.amazonaws.com/${key}`;
}

/**
 * Parse S3 URL to extract bucket and key
 */
export function parseS3Url(url: string): ParsedS3Url | null {
  // Match AWS S3 URL format: https://bucket.s3.region.amazonaws.com/key
  const s3Match = url.match(/https:\/\/([^.]+)\.s3\.[^.]+\.amazonaws\.com\/(.+)/);
  if (s3Match) {
    return {
      bucket: s3Match[1],
      key: s3Match[2],
    };
  }

  // Match MinIO or custom endpoint format: http(s)://endpoint/bucket/key
  const customMatch = url.match(/https?:\/\/[^/]+\/([^/]+)\/(.+)/);
  if (customMatch) {
    return {
      bucket: customMatch[1],
      key: customMatch[2],
    };
  }

  return null;
}

/**
 * Validate file type
 */
export function validateFileType(
  file: File,
  allowedTypes: string[]
): boolean {
  return allowedTypes.some(type => {
    if (type.endsWith('/*')) {
      const prefix = type.slice(0, -2);
      return file.type.startsWith(prefix);
    }
    return file.type === type;
  });
}

/**
 * Validate file size (maxSize in bytes)
 */
export function validateFileSize(file: File, maxSize: number): boolean {
  return file.size <= maxSize;
}

/**
 * Format file size for display
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}
