# app

Sweet Real Auth Service

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


## Configuration Paths

 - /etc/sweet-reel/schema.yaml
 - $HOME/.sweet-reel/schema.yaml
 - ./schema.yaml

### Common

## Testing
```go test ./cmd/app/```
