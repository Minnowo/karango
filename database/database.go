package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// DB is interface for accessing and manipulating data in database.
type DB interface {
	// DBx is the underlying sqlx.DB
	DBx() *sqlx.DB

	// Migrate runs migrations for this database
	Migrate(ctx context.Context) error

	// Inserts a new food into the datbase
	AddFood(ctx context.Context, food *Food) (int, error)
	AddFoods(ctx context.Context, foods []Food) error
	GetAllFoods(ctx context.Context) ([]Food, error)

	GetVersion(ctx context.Context) (int, error)
	SetVersion(ctx context.Context, version DBVersion) error

	HasFlag(ctx context.Context, flag DBFlag) (bool, error)
	SetFlag(ctx context.Context, flag DBFlag, on bool) error

	withTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error
	base() *dbbase
}

type dbbase struct {
	*sqlx.DB
}

func (db *dbbase) base() *dbbase {
	return db
}

func (db *dbbase) withTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {

	log := zerolog.Ctx(ctx)

	tx, err := db.BeginTxx(ctx, nil)

	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	err = fn(tx)

	if err != nil {

		log.Error().Err(err).Msg("error running transaction function")

		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			log.Error().Err(rollbackErr).Msg("error during rollback")
		}

		return errors.WithStack(err)
	}

	err = tx.Commit()

	if err != nil {
		log.Error().Err(err).Msg("error during commit")
	}

	return err
}

func Connect(ctx context.Context, db DBType, connection string) (DB, error) {

	switch db {

	case MOCK:

		con, err := OpenMockDB(ctx, connection)

		return con, err

	case POSTGRES:

		con, err := OpenPGDatabase(ctx, connection)

		return con, err
	}

	panic("Database type is not supported!")
}
