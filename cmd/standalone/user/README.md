# app

Sweet Real User Service

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


## Configuration Paths

 - /etc/sweet-reel/schema.yaml
 - $HOME/.sweet-reel/schema.yaml
 - ./schema.yaml

### Common

## Testing
```go test ./cmd/app/```
