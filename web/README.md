# Sweet Reel - Web Frontend

The web frontend for Sweet Reel, a video processing and management platform built with Next.js.

## Overview

This is a modern Next.js application that provides the user interface for Sweet Reel's video processing platform. It features:

- **Authentication**: OAuth integration with Google using NextAuth.js
- **Video Upload**: Direct upload capabilities to S3-compatible storage
- **Video Management**: Browse, view, and manage uploaded videos
- **User Profiles**: User profile management and viewing
- **Responsive UI**: Built with Tailwind CSS and Radix UI components

## Tech Stack

- **Framework**: Next.js (latest)
- **Authentication**: NextAuth.js with Google OAuth
- **UI Components**: Radix UI, Tailwind CSS
- **State Management**: React Hooks
- **TypeScript**: Full TypeScript support
- **Storage**: S3-compatible object storage integration

## Project Structure

```
web/
├── app/                    # Next.js App Router pages
│   ├── api/               # API routes
│   ├── auth/              # Authentication routes
│   ├── profile/           # User profile pages
│   ├── signin/            # Sign in page
│   ├── upload/            # Video upload page
│   └── video/             # Video viewing pages
├── components/            # Reusable React components
├── hooks/                 # Custom React hooks
├── lib/                   # Utility libraries
│   ├── api/              # API client and configuration
│   └── s3/               # S3 storage utilities
├── types/                 # TypeScript type definitions
└── auth.ts               # NextAuth configuration
```

## Getting Started

### Prerequisites

- Node.js >= 20.0.0
- pnpm package manager
- Backend services running (see main project README)

### 1. Install dependencies

From the `web` directory:

```bash
pnpm install
```

### 2. Configure environment variables

Create a `.env.local` file in the `web` directory with the following variables:

```env
# NextAuth Configuration
AUTH_SECRET=your-secret-key-here
AUTH_GOOGLE_ID=your-google-oauth-client-id
AUTH_GOOGLE_SECRET=your-google-oauth-client-secret

# API Gateway URL
NEXT_PUBLIC_API_URL=http://localhost:8080

# S3/MinIO Configuration
NEXT_PUBLIC_S3_ENDPOINT=http://localhost:9000
NEXT_PUBLIC_S3_BUCKET=video-uploaded
```

### 3. Set up Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:3000/auth/callback/google`
6. Copy the Client ID and Client Secret to your `.env.local`

### 4. Start the development server

```bash
pnpm run dev
```

The application will be available at [http://localhost:3000](http://localhost:3000)

### 5. Build for production

```bash
pnpm run build
pnpm run start
```

## Features

### Authentication
- Google OAuth integration via NextAuth.js
- JWT-based session management
- Automatic token refresh
- Protected routes with middleware

### Video Upload
- Direct upload to S3-compatible storage
- Progress tracking
- File validation
- Multipart upload support

### Video Management
- Browse uploaded videos
- Video playback
- Video metadata display
- User-specific video filtering

### User Interface
- Responsive design with Tailwind CSS
- Accessible components with Radix UI
- Dark/light mode support
- Modern, clean interface

## API Integration

The frontend communicates with the backend services through the API Gateway:

- **Auth Service**: User authentication and token management
- **User Service**: User profile and account management
- **Video Management Service**: Video metadata and listing
- **S3 Storage**: Direct video upload and retrieval

## Development

### Available Scripts

- `pnpm dev` - Start development server
- `pnpm build` - Build for production
- `pnpm start` - Start production server

### Code Organization

- **Components**: Reusable UI components following atomic design principles
- **Hooks**: Custom React hooks for common functionality
- **API Client**: Type-safe API client with error handling
- **Types**: Centralized TypeScript definitions

## Contributing

This is part of the Sweet Reel project. Please refer to the main project repository for contribution guidelines.

## Related Documentation

- [NextAuth.js Documentation](https://authjs.dev)
- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Radix UI Documentation](https://www.radix-ui.com/docs)

## License

ISC
