package database

import (
	"context"
	"embed"
	"fmt"
	"path"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

//go:embed migrations/*
var migrationFiles embed.FS

type DBVersion int

const (
	DBVERSION_NONE DBVersion = -1
	DBVERSION_0    DBVersion = 0
)

func (v DBVersion) String() string {
	return strconv.Itoa(int(v))
}

type migrationFunc func(ctx context.Context, db *dbbase) error

// migration represents a database schema migration
type migration struct {
	fromVersion   DBVersion
	toVersion     DBVersion
	migrationFunc migrationFunc
}

// newFuncMigration creates a new migration from a function.
func newFuncMigration(fromVersion, toVersion DBVersion, migrationFunc migrationFunc) migration {
	return migration{
		fromVersion:   fromVersion,
		toVersion:     toVersion,
		migrationFunc: migrationFunc,
	}
}

// newFileMigration creates a new migration from a file.
func newFileMigration(fromVersion, toVersion DBVersion, filename string) migration {

	return newFuncMigration(fromVersion, toVersion, func(ctx context.Context, db *dbbase) error {

		return db.withTx(ctx, func(tx *sqlx.Tx) error {

			migrationSQL, err := migrationFiles.ReadFile(path.Join("migrations", filename+".up.sql"))

			if err != nil {
				return fmt.Errorf("failed to read migration file: %w", err)
			}

			if _, err := tx.Exec(string(migrationSQL)); err != nil {
				return fmt.Errorf("failed to execute migration %s to %s: %w", fromVersion, toVersion, err)
			}
			return nil
		})
	})
}

// runMigrations runs the given migrations.
func runMigrations(ctx context.Context, db DB, version DBVersion, migrations []migration) error {

	log := zerolog.Ctx(ctx)

	log.Info().Int("version", int(version)).Msg("Running database migrations")

	defer func(v *DBVersion) {
		log.Info().Int("version", int(version)).Msg("Finished database migrations")
	}(&version)

	for _, migration := range migrations {

		if version != migration.fromVersion {
			continue
		}

		log.Info().
			Int("from", int(migration.fromVersion)).
			Int("to", int(migration.toVersion)).
			Msg("Migrating database")

		if err := migration.migrationFunc(ctx, db.base()); err != nil {
			return fmt.Errorf("failed to run migration from %s to %s: %w", migration.fromVersion, migration.toVersion, err)
		}

		version = migration.toVersion

		if err := db.SetVersion(ctx, version); err != nil {
			return fmt.Errorf("failed to store database version %s from %s to %s: %w", version.String(), migration.fromVersion, migration.toVersion, err)
		}
	}

	return nil
}
