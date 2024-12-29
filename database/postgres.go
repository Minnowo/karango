package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPGDatabase(ctx context.Context, connString string) (pgDB *PGDatabase, err error) {

	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Second)

	pgDB = &PGDatabase{
		dbbase: dbbase{
			DB: db,
		},
	}
	return pgDB, err
}

var postgresMigrations = []migration{
	newFileMigration(DBVERSION_NONE, DBVERSION_0, "pg/0001_system"),
}

type PGDatabase struct {
	dbbase
}

func (db *PGDatabase) DBx() *sqlx.DB {
	return db.DB
}

func (db *PGDatabase) Migrate(ctx context.Context) error {

	version, err := db.GetVersion(ctx)

	if err != nil {

		pgErr, ok := err.(*pgconn.PgError)

		if !ok || pgErr.Code != "42P01" {
			return err
		}

		version = -1
	}

	err = runMigrations(ctx, db, DBVersion(version), postgresMigrations)

	if err != nil {
		return err
	}

	return nil
}

func (db *PGDatabase) GetVersion(ctx context.Context) (int, error) {
	var version Version
	query := `SELECT version FROM tbl_version LIMIT 1`
	err := db.GetContext(ctx, &version, query)
	if err != nil {
		return -1, err
	}
	return version.Version, nil
}

func (db *PGDatabase) SetVersion(ctx context.Context, version DBVersion) error {

	_, err := db.ExecContext(ctx, "UPDATE tbl_version SET version = $1", int(version))

	if err != nil {
		return fmt.Errorf("could not update database version: %s", err)
	}

	return nil
}

func (db *PGDatabase) AddFood(ctx context.Context, food *Food) (int, error) {
	var foodID int
	query := `
        INSERT INTO tbl_food (food, portion, unit, protein, carb, fibre, fat)
        VALUES (:food, :portion, :unit, :protein, :carb, :fibre, :fat)
        RETURNING food_id;
    `

	rows, err := db.NamedQueryContext(ctx, query, food)
	if err != nil {
		log.Error().Err(err).Msg("Error executing insert query")
		return 0, err
	}
	defer rows.Close()

	// Retrieve the auto-generated ID
	if rows.Next() {
		if err := rows.Scan(&foodID); err != nil {
			log.Error().Err(err).Msg("Error scanning food_id")
			return 0, err
		}
	}

	return foodID, nil
}

func (db *PGDatabase) AddFoods(ctx context.Context, foods []Food) error {

	return db.withTx(ctx, func(tx *sqlx.Tx) error {

		query := `
        INSERT INTO tbl_food (food, portion, unit, protein, carb, fibre, fat)
        VALUES (:food, :portion, :unit, :protein, :carb, :fibre, :fat);
    `
		for _, food := range foods {

			_, err := tx.NamedExec(query, food)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *PGDatabase) HasFlag(ctx context.Context, flag DBFlag) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM tbl_flags WHERE flag = $1)`
	err := db.DB.GetContext(ctx, &exists, query, int(flag))
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SetFlag sets or unsets a flag in the database.
func (db *PGDatabase) SetFlag(ctx context.Context, flag DBFlag, on bool) error {
	if on {
		query := `INSERT INTO tbl_flags (flag) VALUES ($1) ON CONFLICT DO NOTHING`
		_, err := db.DB.ExecContext(ctx, query, int(flag))
		if err != nil {
			return err
		}
	} else {
		query := `DELETE FROM tbl_flags WHERE flag = $1`
		_, err := db.DB.ExecContext(ctx, query, int(flag))
		if err != nil {
			return err
		}
	}
	return nil
}
func (db *PGDatabase) GetAllFoods(ctx context.Context) ([]Food, error) {

	var foods []Food
	query := `SELECT food_id, food, portion, unit, protein, carb, fibre, fat FROM tbl_food`
	err := db.SelectContext(ctx, &foods, query)
	if err != nil {
		return nil, err
	}
	return foods, nil
}
