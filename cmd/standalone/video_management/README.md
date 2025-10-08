# app

Sweet Real Video Management

```
app [flags]
```

### Options

```
      --config string              config file (default is $HOME/.sweet-reel/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
  -h, --help                       help for app
      --log-level string           log level to use (default "info")
  -s, --service string             which service to run
```

### Environment Variables

- SWEET_REEL_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- SWEET_REEL_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- SWEET_REEL_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LOG_LEVEL :: `log.level` log level to use
- SWEET_REEL_SERVICE :: `service` which service to run
```

## app check

Health check commands

### Synopsis

Commands for running health checks

### Options

```
  -h, --help   help for check
```

### Environment Variables

```

### Options inherited from parent commands

```
      --config string              config file (default is $HOME/.sweet-reel/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
      --log-level string           log level to use (default "info")
  -s, --service string             which service to run
```

### Environment Variables inherited from parent commands

- SWEET_REEL_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- SWEET_REEL_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- SWEET_REEL_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LOG_LEVEL :: `log.level` log level to use
- SWEET_REEL_SERVICE :: `service` which service to run
```

## app video_management

Run as video_management service

```
app video_management [flags]
```

### Options

```
      --aws_s3_access_id string                   s3 access id
      --aws_s3_region string                      s3 region
      --aws_s3_secret string                      s3 secret
      --db-migrations-url string                  Database connection migrations URL
      --db-postgres-connection-max-idletime int   Max connection idle time in seconds (default 180)
      --db-postgres-connection-max-lifetime int   Max connection lifetime in seconds (default 300)
      --db-postgres-max-idle-connections int      Maximum number of idle connections (default 50)
      --db-postgres-max-open-connections int      Maximum number of connections (default 500)
      --db-postgres-timeout int                   Timeout for postgres connection (default 60)
      --db-read-url string                        Database connection readonly URL
      --db-url string                             Database connection URL
      --grpc-port int                             GRPC Port to listen on (default 50060)
  -h, --help                                      help for video_management
      --id string                                 Unique identifier for this services
      --minio-url string                          MINIO URL
      --s3_bucket string                          s3 bucket
      --token-signing-key string                  Signing key used for service to service tokens
```

### Environment Variables

- VIDEO_MANAGEMENT_AWS_S3_ACCESS_ID :: `video_management.aws.s3.access.id` s3 access id
- VIDEO_MANAGEMENT_AWS_S3_REGION :: `video_management.aws.s3.region` s3 region
- VIDEO_MANAGEMENT_AWS_S3_SECRET :: `video_management.aws.s3.secret` s3 secret
- VIDEO_MANAGEMENT_DB_MIGRATIONS_URL :: `video_management.db.migrations.url` Database connection migrations URL
- VIDEO_MANAGEMENT_DB_POSTGRES_CONNECTION_MAX_IDLETIME :: `video_management.db.postgres.max_idletime` Max connection idle time in seconds
- VIDEO_MANAGEMENT_DB_POSTGRES_CONNECTION_MAX_LIFETIME :: `video_management.db.postgres.max_lifetime` Max connection lifetime in seconds
- VIDEO_MANAGEMENT_DB_POSTGRES_MAX_IDLE_CONNECTIONS :: `video_management.db.postgres.max_idle_connections` Maximum number of idle connections
- VIDEO_MANAGEMENT_DB_POSTGRES_MAX_OPEN_CONNECTIONS :: `video_management.db.postgres.max_open_connections` Maximum number of connections
- VIDEO_MANAGEMENT_DB_POSTGRES_TIMEOUT :: `video_management.db.postgres.timeout` Timeout for postgres connection
- VIDEO_MANAGEMENT_DB_READ_URL :: `video_management.db.read.url` Database connection readonly URL
- VIDEO_MANAGEMENT_DB_URL :: `video_management.db.url` Database connection URL
- VIDEO_MANAGEMENT_GRPC_PORT :: `video_management.grpc.port` GRPC Port to listen on
- VIDEO_MANAGEMENT_ID :: `video_management.id` Unique identifier for this services
- VIDEO_MANAGEMENT_MINIO_URL :: `video_management.minio.url` MINIO URL
- VIDEO_MANAGEMENT_AWS_S3_BUCKET :: `video_management.aws.s3.bucket` s3 bucket
- VIDEO_MANAGEMENT_SECRETS_TOKEN_SIGNING_KEY :: `video_management.secrets.token_signing_key` Signing key used for service to service tokens
```

### Options inherited from parent commands

```
      --config string              config file (default is $HOME/.sweet-reel/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
      --log-level string           log level to use (default "info")
  -s, --service string             which service to run
```

### Environment Variables inherited from parent commands

- SWEET_REEL_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- SWEET_REEL_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- SWEET_REEL_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LOG_LEVEL :: `log.level` log level to use
- SWEET_REEL_SERVICE :: `service` which service to run
```


## Configuration Paths

 - /etc/sweet-reel/schema.yaml
 - $HOME/.sweet-reel/schema.yaml
 - ./schema.yaml

### Common

## Testing
```go test ./cmd/app/```
