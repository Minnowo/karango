package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"karango/assets"
	"karango/components/pages"
	"karango/database"
	"karango/logging"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	DBCON_ENV = "DATABASE_CONN"
	DBVENDOR  = "DATABASE_VENDOR"
)

func DBConnect() {

	conn := os.Getenv(DBCON_ENV)
	vendor := database.DBTypeFromStr(os.Getenv(DBVENDOR))

	log.Info().
		Int("vendor", int(vendor)).
		Str("str", conn).
		Msg("Got databse connection string from env")

	db, err := database.Connect(
		context.Background(),
		vendor,
		conn,
	)

	if err != nil {

		log.Error().Err(err).Msg("Cannot get database connection")
		return
	}

	db.Migrate(context.Background())
}

func main() {

	logging.InitFromEnv()

	err := godotenv.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
		return
	}

	DBConnect()

	r := mux.NewRouter()

	r.HandleFunc("/entry", func(w http.ResponseWriter, r *http.Request) {

		pages.EntryPage(&pages.EntryView{
			Time:             time.Now(),
			BGL:              20,
			ITCR:             1,
			AIT:              1,
			RIA:              1,
			Portion:          23,
			BGLIncrement:     0.1,
			ITCRIncrement:    0.5,
			AITIncrement:     0.5,
			RIAIncrement:     0,
			PortionIncrement: 1,
		}).Render(r.Context(), w)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		pages.Home(&pages.HomeView{
			Days: []pages.Day{
				pages.Day{
					Day: "today",
					Events: []pages.Event{
						pages.Event{
							Event:             "lunch",
							Time:              time.Now(),
							BG:                5.3,
							ITCR:              5.0,
							ActualTaken:       7.56,
							RecommendedAmount: 7.43,
							ISF:               3,
							BGT:               6.5,
							Foods: []pages.Food{
								pages.Food{
									Name:    "apple",
									Unit:    "grams",
									Portion: 1,
									Carbs:   10,
									Protein: 10,
									Fat:     10,
									Fibre:   1,
								},
								pages.Food{
									Name:    "pear",
									Unit:    "grams",
									Portion: 1,
									Carbs:   10,
									Protein: 10,
									Fat:     10,
									Fibre:   1,
								},
							},
						},
					},
				},
			},
		}).Render(r.Context(), w)
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

	err = srv.ListenAndServe()

	log.Fatal().Err(err).Msg("Site is dead")
}
