package database

import (
	"context"
	"karango/database/migrations/pg"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGDatabase struct {
	dbbase
}

func (db *PGDatabase) DBx() *sqlx.DB {
	return db.DB
}

var migrationMap = map[int]func(db *PGDatabase){
	0: pg.Migration0001,
	1: pg.Migration0002,
}

//  map[int]func(){
//         1: func() { fmt.Println("Action 1 executed") },
//         2: func() { fmt.Println("Action 2 executed") },
//         3: func() { fmt.Println("Action 3 executed") },
//     }

func (db *PGDatabase) Migrate(ctx context.Context) error {

	version, err := db.GetVersion()

	if err != nil {
		return err
	}

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

func (db *PGDatabase) AddFood(food *Food) (int, error) {
	var foodID int
	query := `
        INSERT INTO tbl_food (food, portion, unit, protein, carb, fibre, fat)
        VALUES (:food, :portion, :unit, :protein, :carb, :fibre, :fat)
        RETURNING food_id;
    `

	// Use NamedQuery to bind struct fields to the query
	rows, err := db.NamedQuery(query, food)
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

func (db *PGDatabase) AddFoods(foods []Food) error {

	tx, err := db.Beginx()

	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return err
	}

	query := `
        INSERT INTO tbl_food (food, portion, unit, protein, carb, fibre, fat)
        VALUES (:food, :portion, :unit, :protein, :carb, :fibre, :fat);
    `

	// Iterate over the slice and execute the insert query for each `Food`
	for _, food := range foods {

		_, err := tx.NamedExec(query, food)

		if err != nil {

			log.Error().Err(err).Msg("Error executing batch insert")

			tx.Rollback()

			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {

		log.Error().Err(err).Msg("Error committing transaction")

		return err
	}

	return nil
}

func (db *PGDatabase) GetVersion() (int, error) {
	var version Version
	query := `SELECT version FROM tbl_version LIMIT 1`
	err := db.Get(&version, query)
	if err != nil {
		return -1, err
	}
	return version.Version, nil
}
func (db *PGDatabase) SetVersion(version int) error {
	var v Version
	v.Version = version
	query := `
        INSERT INTO tbl_version (version) VALUES (:version);
    `

	_, err := db.NamedQuery(query, v)
	if err != nil {
		log.Error().Err(err).Msg("Error executing insert query")
		return err
	}

	return nil
}

func (db *PGDatabase) HasFlag(flag DBFlag) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM tbl_flags WHERE flag = $1)`
	err := db.DB.Get(&exists, query, int(flag))
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SetFlag sets or unsets a flag in the database.
func (db *PGDatabase) SetFlag(flag DBFlag, on bool) error {
	if on {
		query := `INSERT INTO tbl_flags (flag) VALUES ($1) ON CONFLICT DO NOTHING`
		_, err := db.DB.Exec(query, int(flag))
		if err != nil {
			return err
		}
	} else {
		query := `DELETE FROM tbl_flags WHERE flag = $1`
		_, err := db.DB.Exec(query, int(flag))
		if err != nil {
			return err
		}
	}
	return nil
}
