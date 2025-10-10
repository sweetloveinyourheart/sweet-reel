# app

Unified service launcher

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

## app api_gateway

Run as api_gateway service

```
app api_gateway [flags]
```

### Options

```
      --auth-server-url string     Auth server connection URL (default "http://auth:50070")
  -h, --help                       help for api_gateway
      --http-port int              HTTP Port to listen on (default 8080)
      --id string                  Unique identifier for this services
      --token-signing-key string   Signing key used for service to service tokens
```

### Environment Variables

- API_GATEWAY_AUTH_SERVER_URL :: `api_gateway.auth_server.url` Auth server connection URL
- API_GATEWAY_HTTP_PORT :: `api_gateway.http.port` HTTP Port to listen on
- API_GATEWAY_ID :: `api_gateway.id` Unique identifier for this services
- API_GATEWAY_SECRETS_TOKEN_SIGNING_KEY :: `api_gateway.secrets.token_signing_key` Signing key used for service to service tokens
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

## app auth

Run as auth service

```
app auth [flags]
```

### Options

```
      --google-oauth-client-id string       
      --google-oauth-client-secret string   
      --google-oauth-redirect-url string    
      --grpc-port int                       GRPC Port to listen on (default 50060)
  -h, --help                                help for auth
      --id string                           Unique identifier for this services
      --token-signing-key string            Signing key used for service to service tokens
      --user-server-url string              User server connection URL (default "http://user:50065")
```

### Environment Variables

- AUTH_GOOGLE_OAUTH_CLIENT_ID :: `auth.google.oauth.client_id` 
- AUTH_GOOGLE_OAUTH_CLIENT_SECRET :: `auth.google.oauth.client_secret` 
- AUTH_GOOGLE_OAUTH_REDIRECT_URL :: `auth.google.oauth.redirect_url` 
- AUTH_GRPC_PORT :: `auth.grpc.port` GRPC Port to listen on
- AUTH_ID :: `auth.id` Unique identifier for this services
- AUTH_SECRETS_TOKEN_SIGNING_KEY :: `auth.secrets.token_signing_key` Signing key used for service to service tokens
- AUTH_USER_SERVER_URL :: `auth.user_server.url` User server connection URL
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

## app user

Run as user service

```
app user [flags]
```

### Options

```
      --db-migrations-url string                  Database connection migrations URL
      --db-postgres-connection-max-idletime int   Max connection idle time in seconds (default 180)
      --db-postgres-connection-max-lifetime int   Max connection lifetime in seconds (default 300)
      --db-postgres-max-idle-connections int      Maximum number of idle connections (default 50)
      --db-postgres-max-open-connections int      Maximum number of connections (default 500)
      --db-postgres-timeout int                   Timeout for postgres connection (default 60)
      --db-read-url string                        Database connection readonly URL
      --db-url string                             Database connection URL
      --grpc-port int                             GRPC Port to listen on (default 50060)
  -h, --help                                      help for user
      --id string                                 Unique identifier for this services
      --token-signing-key string                  Signing key used for service to service tokens
```

### Environment Variables

- USER_DB_MIGRATIONS_URL :: `user.db.migrations.url` Database connection migrations URL
- USER_DB_POSTGRES_CONNECTION_MAX_IDLETIME :: `user.db.postgres.max_idletime` Max connection idle time in seconds
- USER_DB_POSTGRES_CONNECTION_MAX_LIFETIME :: `user.db.postgres.max_lifetime` Max connection lifetime in seconds
- USER_DB_POSTGRES_MAX_IDLE_CONNECTIONS :: `user.db.postgres.max_idle_connections` Maximum number of idle connections
- USER_DB_POSTGRES_MAX_OPEN_CONNECTIONS :: `user.db.postgres.max_open_connections` Maximum number of connections
- USER_DB_POSTGRES_TIMEOUT :: `user.db.postgres.timeout` Timeout for postgres connection
- USER_DB_READ_URL :: `user.db.read.url` Database connection readonly URL
- USER_DB_URL :: `user.db.url` Database connection URL
- USER_GRPC_PORT :: `user.grpc.port` GRPC Port to listen on
- USER_ID :: `user.id` Unique identifier for this services
- USER_SECRETS_TOKEN_SIGNING_KEY :: `user.secrets.token_signing_key` Signing key used for service to service tokens
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

## app video_processing

Run as video_processing service

```
app video_processing [flags]
```

### Options

```
      --aws_s3_access_id string          s3 access id
      --aws_s3_region string             s3 region
      --aws_s3_secret string             s3 secret
      --grpc-port int                    GRPC Port to listen on (default 50055)
  -h, --help                             help for video_processing
      --id string                        Unique identifier for this services
      --kafka-brokers string             Kafka broker addresses (comma-separated) (default "localhost:9092")
      --kafka-compression string         Compression codec (none, gzip, snappy, lz4, zstd) (default "snappy")
      --kafka-flush-bytes int32          Number of bytes to buffer before flushing (default 16384)
      --kafka-flush-frequency-ms int     Frequency of flushing messages in milliseconds (default 500)
      --kafka-flush-messages int32       Number of messages to buffer before flushing (default 100)
      --kafka-idempotent-writes          Enable idempotent writes (default true)
      --kafka-max-open-requests int32    Max open requests (default 1)
      --kafka-required-acks string       Required acknowledgments (none, leader, all) (default "all")
      --kafka-retry-backoff-ms int       Backoff time between retries in milliseconds (default 100)
      --kafka-retry-max int32            Maximum number of retries for failed requests (default 3)
      --kafka-sasl-mechanism string      SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512)
      --kafka-sasl-password string       SASL password for authentication
      --kafka-sasl-username string       SASL username for authentication
      --kafka-security-protocol string   Security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL) (default "PLAINTEXT")
      --kafka-tls-enabled                Enable TLS encryption
      --minio-url string                 MINIO URL
      --s3_bucket string                 s3 bucket
      --token-signing-key string         Signing key used for service to service tokens
```

### Environment Variables

- VIDEO_PROCESSING_AWS_S3_ACCESS_ID :: `video_processing.aws.s3.access.id` s3 access id
- VIDEO_PROCESSING_AWS_S3_REGION :: `video_processing.aws.s3.region` s3 region
- VIDEO_PROCESSING_AWS_S3_SECRET :: `video_processing.aws.s3.secret` s3 secret
- VIDEO_PROCESSING_GRPC_PORT :: `video_processing.grpc.port` GRPC Port to listen on
- VIDEO_PROCESSING_ID :: `video_processing.id` Unique identifier for this services
- VIDEO_PROCESSING_KAFKA_BROKERS :: `video_processing.kafka.brokers` Kafka broker addresses (comma-separated)
- VIDEO_PROCESSING_KAFKA_COMPRESSION :: `video_processing.kafka.compression` Compression codec (none, gzip, snappy, lz4, zstd)
- VIDEO_PROCESSING_KAFKA_FLUSH_BYTES :: `video_processing.kafka.flush_bytes` Number of bytes to buffer before flushing
- VIDEO_PROCESSING_KAFKA_FLUSH_FREQUENCY_MS :: `video_processing.kafka.flush_frequency_ms` Frequency of flushing messages in milliseconds
- VIDEO_PROCESSING_KAFKA_FLUSH_MESSAGES :: `video_processing.kafka.flush_messages` Number of messages to buffer before flushing
- VIDEO_PROCESSING_KAFKA_IDEMPOTENT_WRITES :: `video_processing.kafka.idempotent_writes` Enable idempotent writes
- VIDEO_PROCESSING_KAFKA_MAX_OPEN_REQUESTS :: `video_processing.kafka.max_open_requests` Max open requests
- VIDEO_PROCESSING_KAFKA_REQUIRED_ACKS :: `video_processing.kafka.required_acks` Required acknowledgments (none, leader, all)
- VIDEO_PROCESSING_KAFKA_RETRY_BACKOFF_MS :: `video_processing.kafka.retry_backoff_ms` Backoff time between retries in milliseconds
- VIDEO_PROCESSING_KAFKA_RETRY_MAX :: `video_processing.kafka.retry_max` Maximum number of retries for failed requests
- VIDEO_PROCESSING_KAFKA_SASL_MECHANISM :: `video_processing.kafka.sasl_mechanism` SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512)
- VIDEO_PROCESSING_KAFKA_SASL_PASSWORD :: `video_processing.kafka.sasl_password` SASL password for authentication
- VIDEO_PROCESSING_KAFKA_SASL_USERNAME :: `video_processing.kafka.sasl_username` SASL username for authentication
- VIDEO_PROCESSING_KAFKA_SECURITY_PROTOCOL :: `video_processing.kafka.security_protocol` Security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL)
- VIDEO_PROCESSING_KAFKA_TLS_ENABLED :: `video_processing.kafka.tls_enabled` Enable TLS encryption
- VIDEO_PROCESSING_MINIO_URL :: `video_processing.minio.url` MINIO URL
- VIDEO_PROCESSING_AWS_S3_BUCKET :: `video_processing.aws.s3.bucket` s3 bucket
- VIDEO_PROCESSING_SECRETS_TOKEN_SIGNING_KEY :: `video_processing.secrets.token_signing_key` Signing key used for service to service tokens
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
