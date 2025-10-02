# app

Sweet Real Video Processing

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
- VIDEO_PROCESSING_ID :: `video.processing.id` Unique identifier for this services
- VIDEO_PROCESSING_KAFKA_BROKERS :: `video.processing.kafka.brokers` Kafka broker addresses (comma-separated)
- VIDEO_PROCESSING_KAFKA_COMPRESSION :: `video.processing.kafka.compression` Compression codec (none, gzip, snappy, lz4, zstd)
- VIDEO_PROCESSING_KAFKA_FLUSH_BYTES :: `video.processing.kafka.flush_bytes` Number of bytes to buffer before flushing
- VIDEO_PROCESSING_KAFKA_FLUSH_FREQUENCY_MS :: `video.processing.kafka.flush_frequency_ms` Frequency of flushing messages in milliseconds
- VIDEO_PROCESSING_KAFKA_FLUSH_MESSAGES :: `video.processing.kafka.flush_messages` Number of messages to buffer before flushing
- VIDEO_PROCESSING_KAFKA_IDEMPOTENT_WRITES :: `video.processing.kafka.idempotent_writes` Enable idempotent writes
- VIDEO_PROCESSING_KAFKA_REQUIRED_ACKS :: `video.processing.kafka.required_acks` Required acknowledgments (none, leader, all)
- VIDEO_PROCESSING_KAFKA_RETRY_BACKOFF_MS :: `video.processing.kafka.retry_backoff_ms` Backoff time between retries in milliseconds
- VIDEO_PROCESSING_KAFKA_RETRY_MAX :: `video.processing.kafka.retry_max` Maximum number of retries for failed requests
- VIDEO_PROCESSING_KAFKA_SASL_MECHANISM :: `video.processing.kafka.sasl_mechanism` SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512)
- VIDEO_PROCESSING_KAFKA_SASL_PASSWORD :: `video.processing.kafka.sasl_password` SASL password for authentication
- VIDEO_PROCESSING_KAFKA_SASL_USERNAME :: `video.processing.kafka.sasl_username` SASL username for authentication
- VIDEO_PROCESSING_KAFKA_SECURITY_PROTOCOL :: `video.processing.kafka.security_protocol` Security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL)
- VIDEO_PROCESSING_KAFKA_TLS_ENABLED :: `video.processing.kafka.tls_enabled` Enable TLS encryption
- VIDEO_PROCESSING_MINIO_URL :: `video_processing.minio.url` MINIO URL
- VIDEO_PROCESSING_AWS_S3_BUCKET :: `video_processing.aws.s3.bucket` s3 bucket
- VIDEO_PROCESSING_SECRETS_TOKEN_SIGNING_KEY :: `video.processing.secrets.token_signing_key` Signing key used for service to service tokens
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
