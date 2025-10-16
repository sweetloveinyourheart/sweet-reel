/**
 * S3 Client Library for Sweet Reel
 * 
 * This library provides a singleton client for interacting with S3-compatible storage
 * (AWS S3 or MinIO) from the client-side.
 */

import type { S3Config } from './types';

/**
 * S3Client - Singleton class for managing S3 operations
 */
export class S3Client {
  private static instance: S3Client | null = null;
  private config: S3Config;

  private constructor(config?: S3Config) {
    this.config = config || this.getDefaultConfig();
  }

  /**
   * Get the singleton instance of S3Client
   */
  public static getInstance(config?: S3Config): S3Client {
    if (!S3Client.instance) {
      S3Client.instance = new S3Client(config);
    }
    return S3Client.instance;
  }

  /**
   * Reset the singleton instance (useful for testing)
   */
  public static resetInstance(): void {
    S3Client.instance = null;
  }

  /**
   * Get S3 configuration from environment variables
   */
  private getDefaultConfig(): S3Config {
    return {
      accessKeyId: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY_ID || '',
      secretAccessKey: process.env.NEXT_PUBLIC_AWS_SECRET_ACCESS_KEY || '',
      region: process.env.NEXT_PUBLIC_AWS_REGION || 'us-east-1',
      endpoint: process.env.NEXT_PUBLIC_AWS_ENDPOINT_URL,
    };
  }

  /**
   * Get the current configuration
   */
  public getConfig(): S3Config {
    return { ...this.config };
  }

  /**
   * Update the configuration
   */
  public updateConfig(config: Partial<S3Config>): void {
    this.config = { ...this.config, ...config };
  }

  /**
   * Upload a file to S3 using a presigned URL
   * The presigned URL should be obtained from your backend API
   */
  public async upload(
    presignedUrl: string,
    file: File | Blob,
    options?: {
      contentType?: string;
      acl?: string;
      onProgress?: (progress: number) => void;
    }
  ): Promise<void> {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      // Track upload progress
      if (options?.onProgress) {
        xhr.upload.addEventListener('progress', (event) => {
          if (event.lengthComputable) {
            const progress = (event.loaded / event.total) * 100;
            options.onProgress?.(progress);
          }
        });
      }

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve();
        } else {
          reject(new Error(`Upload failed with status ${xhr.status}`));
        }
      });

      xhr.addEventListener('error', () => {
        reject(new Error('Upload failed'));
      });

      xhr.addEventListener('abort', () => {
        reject(new Error('Upload aborted'));
      });

      xhr.open('PUT', presignedUrl);

      if (options?.acl) {
        xhr.setRequestHeader('x-amz-acl', options.acl);
      }
      
      if (options?.contentType) {
        xhr.setRequestHeader('Content-Type', options.contentType);
      }

      xhr.send(file);
    });
  }

  /**
   * Download a file from S3 using a presigned URL
   * The presigned URL should be obtained from your backend API
   */
  public async download(presignedUrl: string): Promise<Blob> {
    const response = await fetch(presignedUrl);

    if (!response.ok) {
      throw new Error(`Download failed with status ${response.status}`);
    }

    return response.blob();
  }
}

/**
 * Get the default S3Client instance
 */
export function getS3Client(config?: S3Config): S3Client {
  return S3Client.getInstance(config);
}
