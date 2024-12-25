package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGDatabase struct {
	dbbase
}

func (db *PGDatabase) DBx() *sqlx.DB {
	return db.DB
}

func (db *PGDatabase) Migrate(ctx context.Context) error {

	var schema = `
CREATE TABLE IF NOT EXISTS person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE IF NOT EXISTS place (
    country text,
    city text NULL,
    telcode integer
)`

	db.MustExec(schema)

	return nil
}

func OpenPGDatabase(ctx context.Context, connString string) (pgDB *PGDatabase, err error) {
	// Open database and start transaction
	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Second)

	pgDB = &PGDatabase{dbbase: dbbase{db}}
	return pgDB, err
}
