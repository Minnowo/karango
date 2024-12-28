package database

import "fmt"

var (
	ErrDefaultDataCreate error = fmt.Errorf("default data exists")
)

func CreateDefaultData(db DB) error {

	defaultDataExists, err := db.HasFlag(FLAG_DEFAULT_DATA_CREATED)

	if err != nil {
		return fmt.Errorf("could not check for IS_DEFAULT_DATA_CREATE: %w", err)
	}

	if defaultDataExists {
		return ErrDefaultDataCreate
	}

	db.SetFlag(FLAG_DEFAULT_DATA_CREATED, true)

	testData := []Food{
		{Food: "tim bar", Portion: 1, Unit: "bar", Protein: 6, Carb: 22, Fibre: 4, Fat: 15},
		{Food: "english muffin - wonder", Portion: 57, Unit: "g", Protein: 5, Carb: 25, Fibre: 1, Fat: 1.5},
		{Food: "spaghettini - great value", Portion: 85, Unit: "g", Protein: 12, Carb: 63, Fibre: 3, Fat: 1},
		{Food: "egg - medium", Portion: 1, Unit: "egg", Protein: 5.5, Carb: 0.5, Fibre: 0, Fat: 4.5},
		{Food: "whey - ON Gold Standard - Natural Vanilla", Portion: 32, Unit: "g", Protein: 24, Carb: 5, Fibre: 0, Fat: 1},
		{Food: "cassein - ON Gold Standard - Natural Chocolate", Portion: 38, Unit: "g", Protein: 24, Carb: 8, Fibre: 1, Fat: 1.5},
		{Food: "silk - oat yeah", Portion: 250, Unit: "ml", Protein: 1, Carb: 6, Fibre: 1, Fat: 2.5},
		{Food: "oatmeal - quaker - large flake", Portion: 30, Unit: "g", Protein: 4, Carb: 20, Fibre: 3, Fat: 2},
		{Food: "blueberries - frozen", Portion: 80, Unit: "g", Protein: 0.5, Carb: 10, Fibre: 2, Fat: 0},
		{Food: "rice - jasmine", Portion: 50, Unit: "g", Protein: 4, Carb: 38, Fibre: 2, Fat: 0},
		{Food: "sugar cube", Portion: 1, Unit: "cube", Protein: 0, Carb: 4, Fibre: 0, Fat: 0},
		{Food: "fudge bar", Portion: 1, Unit: "bar", Protein: 2, Carb: 18, Fibre: 1, Fat: 2},
		{Food: "float bar", Portion: 1, Unit: "bar", Protein: 0.2, Carb: 10, Fibre: 0, Fat: 2},
		{Food: "english muffin - no name", Portion: 57, Unit: "g", Protein: 5, Carb: 26, Fibre: 1, Fat: 1},
		{Food: "grapes", Portion: 100, Unit: "g", Protein: 0.72, Carb: 18.1, Fibre: 0.9, Fat: 0.16},
		// Continue adding the rest of the test data here...
	}

	err = db.AddFoods(testData)

	if err != nil {
		return fmt.Errorf("could not insert foods: %w", err)
	}

	return nil
}
