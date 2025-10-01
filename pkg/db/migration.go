package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// PerformMigrations will automatically connect to the database and perform any unapplied migrations from the passed in
// file system. This method should be safe to run concurrently as it will acquire a database level lock before applying migrations
func PerformMigrations(fs embed.FS, dbPrefix string, connectionStr string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Global().Warn("database migrations failed to apply, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			logger.Global().Info("retrying migration")
			performMigrations(fs, dbPrefix, connectionStr)
		}
	}()
	performMigrations(fs, dbPrefix, connectionStr)
}

func performMigrations(fs embed.FS, dbPrefix string, connectionStr string) {
	prefixIsValid := validatePrefix(dbPrefix)
	if !prefixIsValid {
		err := fmt.Errorf("invalid prefix - can't do migration - %s", dbPrefix)
		panic(err)
	}

	migrationDir, err := iofs.New(fs, "migrations")
	if err != nil {
		panic(err)
	}

	connection, err := sql.Open("pgx", connectionStr)
	if err != nil {
		panic(err)
	}

	driver, err := migrateDriver.WithInstance(connection, &migrateDriver.Config{
		SchemaName:      "public",
		MigrationsTable: fmt.Sprintf("%sschema_migrations", dbPrefix),
	})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance("iofs", migrationDir, "pgx", driver)
	if err != nil {
		panic(err)
	}

	// Wait up to 60 seconds, retrying at most once a second
	for range 60 {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				// Database already migrated, return
				logger.Global().Info("no database migrations to apply")
				return
			}

			if errors.Is(err, migrate.ErrLocked) {
				// Another process is applying migrations, sleep and then try again
				time.Sleep(1 * time.Second)
				continue
			}

			if errors.Is(err, migrate.ErrLockTimeout) {
				// Another process is likely applying migrations, sleep and then try again
				time.Sleep(1 * time.Second)
				continue
			}

			if errors.Is(err, migrate.ErrDirty{}) {
				if dirty, ok := err.(migrate.ErrDirty); ok {
					if dirty.Version == 0 {
						connection, err := sql.Open("pgx", connectionStr)
						if err != nil {
							panic(err)
						}
						logger.Global().Info("dropping schema_migrations table, initial migration failed")
						_, err = connection.Exec(fmt.Sprintf("TRUNCATE TABLE %sschema_migrations;", dbPrefix))
						if err != nil {
							logger.Global().Error("failed to drop schema_migrations table", zap.Error(err))
						}
						_ = connection.Close()
					}
				}
				time.Sleep(1 * time.Second)
				continue
			}

			// Unexpected error, stop the process
			logger.Global().Error("database migrations failed to apply", zap.Error(err))
			panic(err)
		} else {
			// No error, migrations applied successfully, return
			logger.Global().Info("database migrations applied successfully")
			return
		}
	}

	logger.Global().Error("database migrations failed to apply, exceeded retry", zap.Error(err))

	err = fmt.Errorf("unable to apply database migrations, exceeded retries")
	panic(err)
}
