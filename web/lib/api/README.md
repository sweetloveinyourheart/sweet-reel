# API Client Library

A unified API client for both client-side and server-side usage with automatic NextAuth token injection.

## Architecture

```
lib/api/
├── core/           # Core API client functionality
│   ├── client.ts   # ApiClient class
│   ├── config.ts   # API configuration
│   ├── errors.ts   # Error handling
│   └── types.ts    # TypeScript types
├── client.ts       # Client-side hook (useApiClient)
├── server.ts       # Server-side function (getServerApiClient)
└── index.ts        # Main exports
```

## Usage

### Client-Side (React Components)

Use the `useApiClient()` hook in client components:

```tsx
"use client"

import { useApiClient } from "@/lib/api"

function MyComponent() {
  const api = useApiClient()

  const fetchVideos = async () => {
    try {
      const videos = await api.get("/videos")
      console.log(videos)
    } catch (error) {
      console.error("Failed to fetch videos:", error)
    }
  }

  return <button onClick={fetchVideos}>Fetch Videos</button>
}
```

### Server-Side (Server Components, Actions, Route Handlers)

Use the `getServerApiClient()` function in server contexts:

#### Server Component

```tsx
import { getServerApiClient } from "@/lib/api/server"

export default async function VideosPage() {
  const api = await getServerApiClient()
  const videos = await api.get("/videos")

  return (
    <div>
      {videos.map(video => (
        <div key={video.id}>{video.title}</div>
      ))}
    </div>
  )
}
```

#### Server Action

```tsx
"use server"

import { getServerApiClient } from "@/lib/api/server"

export async function deleteVideo(id: string) {
  const api = await getServerApiClient()
  await api.delete(`/videos/${id}`)
  revalidatePath("/videos")
}
```

#### Route Handler

```tsx
import { getServerApiClient } from "@/lib/api/server"
import { NextRequest } from "next/server"

export async function GET(request: NextRequest) {
  const api = await getServerApiClient()
  const data = await api.get("/videos")
  return Response.json(data)
}

export async function POST(request: NextRequest) {
  const api = await getServerApiClient()
  const body = await request.json()
  const result = await api.post("/videos", body)
  return Response.json(result)
}
```

## API Methods

The API client provides the following methods:

### HTTP Methods

```typescript
// GET request
await api.get<T>(endpoint, options?)

// POST request
await api.post<T>(endpoint, body?, options?)

// PUT request
await api.put<T>(endpoint, body?, options?)

// PATCH request
await api.patch<T>(endpoint, body?, options?)

// DELETE request
await api.delete<T>(endpoint, options?)
```

### File Upload

```typescript
await api.uploadFile(
  endpoint,
  file,
  onProgress?: (progress: number) => void,
  additionalData?: Record<string, string>
)
```

### Request Options

```typescript
interface RequestOptions {
  params?: Record<string, string | number | boolean>  // Query parameters
  token?: string                                      // Override auth token
  timeout?: number                                    // Request timeout
  retry?: number                                      // Retry attempts
  skipInterceptors?: boolean                          // Skip interceptors
  headers?: HeadersInit                               // Additional headers
  // ... other fetch options
}
```

### Examples

#### With Query Parameters

```typescript
const videos = await api.get("/videos", {
  params: {
    page: 1,
    limit: 10,
    category: "music"
  }
})
// GET /videos?page=1&limit=10&category=music
```

#### With Custom Headers

```typescript
const result = await api.post("/videos", videoData, {
  headers: {
    "X-Custom-Header": "value"
  }
})
```

#### File Upload with Progress

```typescript
"use client"

import { useApiClient } from "@/lib/api"

function UploadComponent() {
  const api = useApiClient()
  const [progress, setProgress] = useState(0)

  const handleUpload = async (file: File) => {
    await api.uploadFile(
      "/videos/upload",
      file,
      (progress) => setProgress(progress),
      { title: "My Video", description: "..." }
    )
  }

  return <input type="file" onChange={(e) => handleUpload(e.target.files[0])} />
}
```

## Authentication

Both client and server functions automatically inject the NextAuth access token:

- **Client-side**: Token is retrieved from the NextAuth session via `useSession()` hook
- **Server-side**: Token is retrieved from the NextAuth session via `auth()` function

No manual token management is required!

## Error Handling

```typescript
import { ApiError } from "@/lib/api"

try {
  const data = await api.get("/videos")
} catch (error) {
  if (error instanceof ApiError) {
    console.error(`API Error: ${error.message}`)
    console.error(`Status: ${error.status}`)
    console.error(`Response:`, error.response)
  }
}
```

## Interceptors

You can add custom interceptors to the API client:

```typescript
import { apiClient } from "@/lib/api/core/client"

// Request interceptor
apiClient.addRequestInterceptor((config, url) => {
  console.log(`Request: ${config.method} ${url}`)
  return config
})

// Response interceptor
apiClient.addResponseInterceptor((response) => {
  console.log(`Response: ${response.status}`)
  return response
})

// Error interceptor
apiClient.addErrorInterceptor((error) => {
  console.error(`Error: ${error.message}`)
  return error
})
```

## Configuration

Configure the API base URL in your environment variables:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

Default configuration is in `core/config.ts`:

```typescript
export const API_CONFIG = {
  baseUrl: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1",
  timeout: 30000,
}
```

## Migration Guide

### From Old Service Pattern

**Before:**
```typescript
import { VideoService } from "@/lib/api/services"

const videos = await VideoService.getVideos()
```

**After (Client-side):**
```typescript
"use client"
import { useApiClient } from "@/lib/api"

function MyComponent() {
  const api = useApiClient()
  const videos = await api.get("/videos")
}
```

**After (Server-side):**
```typescript
import { getServerApiClient } from "@/lib/api/server"

async function MyServerComponent() {
  const api = await getServerApiClient()
  const videos = await api.get("/videos")
}
```

### From ApiProvider Pattern

The `ApiProvider` component is no longer needed! Token injection happens automatically:

**Before:**
```tsx
<SessionProvider>
  <ApiProvider>
    <YourApp />
  </ApiProvider>
</SessionProvider>
```

**After:**
```tsx
<SessionProvider>
  <YourApp />
</SessionProvider>
```

## Best Practices

1. **Use the hook in client components**: Always use `useApiClient()` in components marked with `"use client"`
2. **Use the function in server contexts**: Always use `getServerApiClient()` in server components, actions, and route handlers
3. **Handle errors**: Always wrap API calls in try-catch blocks
4. **Type your responses**: Use TypeScript generics for type-safe responses
5. **Avoid mixing contexts**: Don't try to use `useApiClient()` in server components or `getServerApiClient()` in client components

## Troubleshooting

### "useSession must be wrapped in SessionProvider"

Make sure your app is wrapped in `<SessionProvider>` (already configured in the root layout).

### "Cannot use 'use client' directive in server component"

Make sure you're importing from the correct module:
- Client components: `import { useApiClient } from "@/lib/api"`
- Server components: `import { getServerApiClient } from "@/lib/api/server"`

### Unauthorized (401) errors

Check that:
1. The user is logged in (has a valid session)
2. The access token is not expired
3. The API endpoint requires authentication
