package main

import (
	"context"
	"net/http"
	"time"

	"karango/assets"
	"karango/components/pages"
	"karango/database"
	"karango/logging"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func test() {

	// basic database test
	db, err := database.OpenPGDatabase(context.Background(), "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
		return
	}

    log.Info().Msg("Connected to the database")

	err = db.Migrate(context.Background())

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
		return
	}

    log.Info().Msg("Migration successful")

	// we are hard coding for postgres right now
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	tx.Commit()

    log.Info().Msg("Test data added successfully")

}

func main() {

	logging.InitFromEnv()

	test()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home().Render(r.Context(), w)
	})

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

	err := srv.ListenAndServe()

	log.Fatal().Err(err).Msg("Site is dead")
}
