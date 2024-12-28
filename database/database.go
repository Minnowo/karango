package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// DB is interface for accessing and manipulating data in database.
type DB interface {
	// DBx is the underlying sqlx.DB
	DBx() *sqlx.DB

	// Migrate runs migrations for this database
	Migrate(ctx context.Context) error
}

type dbbase struct {
	*sqlx.DB
}

type DBType int

const (
	POSTGRES DBType = iota
)

func DBTypeFromStr(s string) DBType {

	switch s {
	case "postgres":
		return POSTGRES
	}

	return POSTGRES
}

func Connect(ctx context.Context, db DBType, connection string) (DB, error) {

	switch db {

	case POSTGRES:

		con, err := OpenPGDatabase(ctx, connection)

		return con, err
	}

	panic("Database type is not supported!")
}
