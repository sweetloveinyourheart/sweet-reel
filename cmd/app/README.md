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
      --grpc-port int              GRPC Port to listen on (default 50055)
  -h, --help                       help for video_processing
      --id string                  Unique identifier for this services
      --token-signing-key string   Signing key used for service to service tokens
```

### Environment Variables

- VIDEO_PROCESSING_GRPC_PORT :: `video_processing.grpc.port` GRPC Port to listen on
- VIDEO_PROCESSING_ID :: `video.processing.id` Unique identifier for this services
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
