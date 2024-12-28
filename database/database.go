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

	// Inserts a new food into the datbase
	AddFood(food *Food) (int, error)
	AddFoods(foods []Food) error

	GetVersion() (int, error)
	SetVersion(version int) error

	HasFlag(flag DBFlag) (bool, error)
	SetFlag(flag DBFlag, on bool) error
}

type dbbase struct {
	*sqlx.DB
}

func Connect(ctx context.Context, db DBType, connection string) (DB, error) {

	switch db {

	case POSTGRES:

		con, err := OpenPGDatabase(ctx, connection)

		return con, err
	}

	panic("Database type is not supported!")
}
