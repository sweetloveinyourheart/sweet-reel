# S3 Library

A TypeScript library for interacting with S3-compatible storage (AWS S3 or MinIO) in the Sweet Reel web application.

## Features

- **Singleton Pattern**: Single instance for managing S3 operations throughout the application
- **Type-Safe**: Full TypeScript support with comprehensive type definitions
- **Modular Architecture**: Separated into logical modules (client, types, constants, utilities)
- **Progress Tracking**: Built-in upload progress monitoring
- **Backend-Driven**: Uses presigned URLs from backend API for security

## File Structure

```
lib/s3/
├── index.ts          # Main exports
├── s3.ts            # S3Client singleton class
├── types.ts         # TypeScript type definitions
├── constants.ts     # S3-related constants
├── utils.ts         # Utility functions
├── examples.ts      # Usage examples
└── README.md        # This file
```

## Usage

### Basic Setup

```typescript
import { getS3Client } from '@/lib/s3';

// Get the singleton instance
const s3Client = getS3Client();
```

### Upload a File

```typescript
import { getS3Client, generateS3Key } from '@/lib/s3';

async function uploadVideo(file: File, userId: string) {
  // Generate a unique key for the file
  const key = generateS3Key('videos', file.name, userId);
  
  // Get presigned URL from your backend API
  const presignedUrl = await fetch('/api/s3/presigned-upload', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ bucket: 'video-uploaded', key }),
  }).then(res => res.json()).then(data => data.url);
  
  // Upload with progress tracking
  const s3Client = getS3Client();
  await s3Client.upload(presignedUrl, file, {
    contentType: file.type,
    onProgress: (progress) => {
      console.log(`Upload progress: ${progress.toFixed(2)}%`);
    }
  });
  
  return key;
}
```

### Download a File

```typescript
import { getS3Client } from '@/lib/s3';

async function downloadVideo(key: string) {
  // Get presigned URL from your backend API
  const presignedUrl = await fetch('/api/s3/presigned-download', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ bucket: 'video-processed', key }),
  }).then(res => res.json()).then(data => data.url);
  
  const s3Client = getS3Client();
  const blob = await s3Client.download(presignedUrl);
  
  // Create a download link
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'video.mp4';
  a.click();
  URL.revokeObjectURL(url);
}
```

### Utility Functions

```typescript
import {
  generateS3Key,
  validateFileType,
  validateFileSize,
  formatFileSize,
  parseS3Url,
  getPublicS3Url,
  ALLOWED_VIDEO_TYPES,
  MAX_VIDEO_SIZE,
} from '@/lib/s3';

// Generate unique S3 key
const key = generateS3Key('videos', 'my-video.mp4', 'user-123');
// Result: "videos/user-123/1697472000000-abc123def-my-video.mp4"

// Validate file type
const isValidType = validateFileType(file, ALLOWED_VIDEO_TYPES);

// Validate file size
const isValidSize = validateFileSize(file, MAX_VIDEO_SIZE);

// Format file size for display
const sizeText = formatFileSize(file.size);
// Result: "15.5 MB"

// Parse S3 URL
const parsed = parseS3Url('https://mybucket.s3.us-east-1.amazonaws.com/path/to/file.mp4');
// Result: { bucket: 'mybucket', key: 'path/to/file.mp4' }

// Get public S3 URL
const publicUrl = getPublicS3Url('mybucket', 'path/to/file.mp4', 'us-east-1');
```

### Custom Configuration

```typescript
import { S3Client } from '@/lib/s3';

// Get instance with custom config
const s3Client = S3Client.getInstance({
  accessKeyId: 'custom-key',
  secretAccessKey: 'custom-secret',
  region: 'us-west-2',
  endpoint: 'https://custom-endpoint.com',
});

// Update configuration
s3Client.updateConfig({
  region: 'eu-west-1',
});

// Get current configuration
const config = s3Client.getConfig();
```

## API Reference

### S3Client

#### Methods

- `getInstance(config?: S3Config): S3Client` - Get singleton instance
- `resetInstance(): void` - Reset singleton (useful for testing)
- `getConfig(): S3Config` - Get current configuration
- `updateConfig(config: Partial<S3Config>): void` - Update configuration
- `upload(presignedUrl, file, options?): Promise<void>` - Upload file using presigned URL
- `download(presignedUrl): Promise<Blob>` - Download file using presigned URL

### Constants

- `S3_BUCKETS` - Predefined bucket names
  - `VIDEO_UPLOADED`
  - `VIDEO_PROCESSED`
- `DEFAULT_EXPIRATION_SECONDS` - 600 (10 minutes)
- `ALLOWED_VIDEO_TYPES` - Array of allowed video MIME types
- `MAX_VIDEO_SIZE` - 500 MB in bytes

### Utility Functions

- `generateS3Key(prefix, filename, userId?): string` - Generate unique S3 key
- `getPublicS3Url(bucket, key, region?, endpoint?): string` - Get public S3 URL
- `parseS3Url(url): ParsedS3Url | null` - Parse S3 URL
- `validateFileType(file, allowedTypes): boolean` - Validate file type
- `validateFileSize(file, maxSize): boolean` - Validate file size
- `formatFileSize(bytes): string` - Format bytes for display

## Environment Variables

The library uses the following environment variables:

```env
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
AWS_ENDPOINT_URL=http://localhost:9000  # Optional, for MinIO
```

## Testing

```typescript
import { S3Client } from '@/lib/s3';

// Reset singleton between tests
afterEach(() => {
  S3Client.resetInstance();
});

test('should upload file', async () => {
  const s3Client = S3Client.getInstance({
    accessKeyId: 'test-key',
    secretAccessKey: 'test-secret',
    region: 'us-east-1',
  });
  
  // Your test logic here
});
```

## License

MIT
