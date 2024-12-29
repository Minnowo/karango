package handler

import (
	"karango/components/pages"
	"karango/database"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type MainRoute struct {
	db  database.DB
	log zerolog.Logger
}

func NewMainRouteHandler(db database.DB, logger zerolog.Logger) *MainRoute {
	return &MainRoute{db: db, log: logger}
}

func (h *MainRoute) HandleEntry(w http.ResponseWriter, r *http.Request) {

	foods, err := h.db.GetAllFoods(r.Context())

	if err != nil {

		h.log.Error().Err(err).Msg("Could not get foods")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

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
		Foods:            foods,
	}).Render(r.Context(), w)
}

func (h *MainRoute) HandleRoot(w http.ResponseWriter, r *http.Request) {

	pages.Home(&pages.HomeView{
		Days: []pages.Day{
			{
				Day: "today",
				Events: []pages.Event{
					{
						Event:             "lunch",
						Time:              time.Now(),
						BG:                5.3,
						ITCR:              5.0,
						ActualTaken:       7.56,
						RecommendedAmount: 7.43,
						ISF:               3,
						BGT:               6.5,
						Foods: []pages.Food{
							{
								Name:    "apple",
								Unit:    "grams",
								Portion: 1,
								Carbs:   10,
								Protein: 10,
								Fat:     10,
								Fibre:   1,
							},
							{
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
}
