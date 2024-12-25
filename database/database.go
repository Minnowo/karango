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
