package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

func OpenMockDB(ctx context.Context, connString string) (pgDB *MockDatabase, err error) {
	pgDB = &MockDatabase{
		version: -1,
		foods:   make([]Food, 0),
		flags:   make(map[DBFlag]bool),
	}
	return pgDB, err
}

type MockDatabase struct {
	dbbase
	sync.RWMutex
	version DBVersion
	foods   []Food
	flags   map[DBFlag]bool
}

func (db *MockDatabase) DBx() *sqlx.DB {
	return db.DB
}

func (db *MockDatabase) Migrate(ctx context.Context) error {

	db.SetFlag(ctx, FLAG_DEFAULT_DATA_CREATED, false)

	return nil
}

func (db *MockDatabase) GetVersion(ctx context.Context) (int, error) {
	db.RLock()
	defer db.RUnlock()
	return int(db.version), nil
}

func (db *MockDatabase) SetVersion(ctx context.Context, version DBVersion) error {

	db.Lock()
	defer db.Unlock()

	db.version = version

	return nil
}

func (db *MockDatabase) AddFood(ctx context.Context, food *Food) (int, error) {

	db.Lock()
	defer db.Unlock()

	db.foods = append(db.foods, *food)

	return len(db.foods), nil
}

func (db *MockDatabase) AddFoods(ctx context.Context, foods []Food) error {

	for _, f := range foods {

		_, err := db.AddFood(ctx, &f)

		if err != nil {
			return err
		}
	}
	return nil
}

func (db *MockDatabase) HasFlag(ctx context.Context, flag DBFlag) (bool, error) {

	db.RLock()
	defer db.RUnlock()

	f, ok := db.flags[flag]

	if !ok {
		return false, fmt.Errorf("TABLE flags do not exist")
	}

	return f, nil
}

// SetFlag sets or unsets a flag in the database.
func (db *MockDatabase) SetFlag(ctx context.Context, flag DBFlag, on bool) error {

	db.Lock()
	defer db.Unlock()

	db.flags[flag] = on

	return nil
}
func (db *MockDatabase) GetAllFoods(ctx context.Context) ([]Food, error) {

	db.RLock()
	defer db.RUnlock()

	foods := make([]Food, len(db.foods))

	for i := range len(db.foods) {
		foods[i] = db.foods[i]
	}

	return foods, nil
}
