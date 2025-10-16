# Sweet Reel

A microservices-based video processing and management platform built with Go. Sweet Reel handles video uploads, OAuth authentication, and automated video transcoding with FFmpeg.

## Architecture

Sweet Reel consists of five microservices:

- **API Gateway** (HTTP:8080) - REST API gateway for client requests
- **Auth Service** (gRPC:50070) - OAuth authentication (Google, GitHub, etc.)
- **User Service** (gRPC:50065) - User profile management
- **Video Management** (gRPC:50060) - Video metadata and presigned URL generation
- **Video Processing** (gRPC:50055) - FFmpeg-based video transcoding and segmentation

**Infrastructure:**
- PostgreSQL (databases: `user`, `video-management`)
- Redis (session/caching)
- Apache Kafka (event streaming for video processing)
- MinIO (S3-compatible object storage)

## Prerequisites

- **Go** 1.25.1 or higher
- **Docker** and **Docker Compose**
- **Make**
- **direnv** (optional, for auto-loading environment variables)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/sweetloveinyourheart/sweet-reel.git
cd sweet-reel
```

### 2. Setup Environment Variables

The project uses a `.env` file for configuration. A sample is already included:

```bash
cat .env
```

**Key variables:**
- `COMPOSE_PROJECT_NAME=srl` - Docker Compose project name
- Container images (Alpine, Postgres, Redis, Kafka versions)
- `AUTH_GOOGLE_OAUTH_CLIENT_ID` and `AUTH_GOOGLE_OAUTH_CLIENT_SECRET` - Configure your own Google OAuth credentials

**Using direnv (recommended):**

```bash
# Install direnv (macOS)
brew install direnv

# Add hook to your shell (~/.zshrc)
echo 'eval "$(direnv hook zsh)"' >> ~/.zshrc
source ~/.zshrc

# Create .envrc to auto-load .env
echo 'dotenv' > .envrc
direnv allow
```

### 3. Install Go Dependencies

```bash
make go-deps
```

This installs:
- `golangci-lint` - Linter
- `goimports` - Import formatter
- Protocol buffer tools (via `go.mod`)

## Running the Project

### Quick Start (Docker Compose)

Start all services with Docker Compose:

```bash
make compose-up
```

This launches:
- Infrastructure services (Postgres, Redis, Kafka, MinIO)
- All microservices (api_gateway, auth, user, video_management, video_processing)

**Access points:**
- API Gateway: http://localhost:8080
- MinIO Console: http://localhost:9001 (credentials: `sweetreel` / `sweetreel4000`)
- PostgreSQL: `localhost:5432` (user: `root_admin`, password: `admin@123`)

### Stop Services

```bash
make compose-down
```

### Building Docker Images

Build all containers:

```bash
make build
```

This generates:
- `srl:latest` - Standard services binary
- `srl-ffmpeg:latest` - Video processing service with FFmpeg

## Development

### Code Generation

Generate gRPC/Protocol Buffer code:

```bash
make gen
```

This runs:
- `buf generate` for proto files
- `goimports` for formatting

### Linting

```bash
make lint
```

### Testing

Run all unit tests:

```bash
make test
```

Run tests with verbose output:

```bash
make test-verbose
```

Run tests with coverage:

```bash
make test-coverage
```

Coverage reports are saved to `tests/logs/cov-*/`.

**Run service-specific tests:**

```bash
make ut-auth                # Auth service tests
make ut-user                # User service tests
make ut-video_management    # Video management tests
make ut-video_processing    # Video processing tests
```

### Running Individual Services Locally

Build the main binary:

```bash
make app-build
```

Run a specific service:

```bash
# Ensure dependencies are running (DB, Kafka, etc.)
make compose-up

# Run service (requires appropriate env vars)
SWEET_REEL_SERVICE=api_gateway \
API_GATEWAY_HTTP_PORT=8080 \
API_GATEWAY_SECRETS_TOKEN_SIGNING_KEY=secr3t_k3y \
./cmd/app/app
```

Check service health:

```bash
./cmd/app/app check http localhost:8080/api/v1/health
```

## API Endpoints

### Authentication

**POST** `/api/v1/oauth` - OAuth login (Google)
```json
{
  "provider": "google",
  "access_token": "<google_access_token>"
}
```

**GET** `/api/v1/auth/refresh-token` - Refresh JWT token

### Video Management (Protected)

**POST** `/api/v1/videos/presigned-url` - Generate presigned upload URL
```json
{
  "title": "My Video",
  "description": "Video description",
  "file_name": "video.mp4"
}
```

## Project Structure

```
.
├── cmd/
│   ├── app/              # Unified service launcher
│   └── standalone/       # Standalone service binaries
├── services/             # Service implementations
│   ├── api_gateway/
│   ├── auth/
│   ├── user/
│   ├── video_management/
│   └── video_processing/
├── proto/                # Protocol buffer definitions
├── pkg/                  # Shared packages
│   ├── ffmpeg/          # FFmpeg transcoding utilities
│   ├── kafka/           # Kafka client
│   ├── s3/              # S3/MinIO client
│   └── db/              # Database utilities
├── dockerfiles/          # Docker Compose configuration
└── scripts/              # Build and utility scripts
```

## Makefile Targets

Run `make help` to see all available targets:

```bash
make help
```

**Common targets:**
- `make build` - Build Docker images
- `make compose-up` - Start all services
- `make compose-down` - Stop all services
- `make gen` - Generate code from proto files
- `make lint` - Run linters
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage
- `make go-deps` - Install Go dependencies

## Troubleshooting

### Port Conflicts

If ports are already in use, modify them in `dockerfiles/docker-compose.yml`:
- PostgreSQL: 5432
- Redis: 6379
- Kafka: 9092
- MinIO: 9000, 9001
- API Gateway: 8080

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Update `.env` with your `CLIENT_ID` and `CLIENT_SECRET`
6. Set redirect URL: `http://localhost:3000/auth/callback/google`

### MinIO Access

Default credentials are in `.env`:
- Access Key: `sweetreel`
- Secret Key: `sweetreel4000`

Buckets are auto-created:
- `video-uploaded` - Raw uploaded videos
- `video-processed` - Transcoded videos

### Database Migrations

Migrations run automatically on service startup. To run manually:

```bash
# Access the database container
docker exec -it srl_database psql -U root_admin -d video-management

# Or from host
psql "postgres://root_admin:admin@123@localhost:5432/video-management?sslmode=disable"
```

### Logs

View service logs:

```bash
# All services
docker compose -f dockerfiles/docker-compose.yml logs -f

# Specific service
docker compose -f dockerfiles/docker-compose.yml logs -f api_gateway
```

## License

See [LICENSE](LICENSE) file for details.