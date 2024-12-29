package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"karango/assets"
	"karango/database"
	"karango/handler"
	"karango/logging"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	DBCON_ENV = "DATABASE_CONN"
	DBVENDOR  = "DATABASE_VENDOR"
)

func DBConnect() database.DB {

	conn := os.Getenv(DBCON_ENV)
	vendor := database.DBTypeFromStr(os.Getenv(DBVENDOR))

	log.Debug().
		Int("vendor", int(vendor)).
		Str("str", conn).
		Msg("Got databse connection string from env")

	ctx := log.Logger.WithContext(context.Background())

	db, err := database.Connect(
		ctx,
		vendor,
		conn,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get database connection")
		return nil
	}

	err = db.Migrate(ctx)

	if err != nil {
		log.Fatal().Err(err).Msg("could not migrate")
		return nil
	}

	err = database.CreateDefaultData(ctx, db)

	if err != nil {
		if err != database.ErrDefaultDataCreate {
			log.Warn().Err(err).Msg("Could not create default data")
		} else {
			log.Info().Msg("Default data already exists")
		}
	}

	foods, err := db.GetAllFoods(ctx)

	if err != nil {
		log.Error().Err(err).Msg("error selecting all foods")
	} else {
		for _, f := range foods {
			log.Debug().Str("food", f.Food).Msg("got food")
		}
	}

	return db
}

func main() {

	logging.InitFromEnv()

	err := godotenv.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
		return
	}

	conn := DBConnect()

	mainHandler := handler.NewMainRouteHandler(conn, log.Logger)

	r := mux.NewRouter()
	r.HandleFunc("/entry", mainHandler.HandleEntry)
	r.HandleFunc("/", mainHandler.HandleRoot)

	assets.Register(r)

	const addr = "0.0.0.0:6732"

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Str("address", addr).Msg("Site is running")

	err = srv.ListenAndServe()

	log.Fatal().Err(err).Msg("Site is dead")
}
