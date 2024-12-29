package database

import (
	"context"
	"fmt"
)

var (
	ErrDefaultDataCreate error = fmt.Errorf("default data exists")
)

func CreateDefaultData(ctx context.Context, db DB) error {

	defaultDataExists, err := db.HasFlag(ctx, FLAG_DEFAULT_DATA_CREATED)

	if err != nil {
		return fmt.Errorf("could not check for IS_DEFAULT_DATA_CREATE: %w", err)
	}

	if defaultDataExists {
		return ErrDefaultDataCreate
	}

	db.SetFlag(ctx, FLAG_DEFAULT_DATA_CREATED, true)

	testData := []Food{
		{FoodID: 1, Food: "tim bar", Portion: 1, Unit: "bar", Protein: 6, Carb: 22, Fibre: 4, Fat: 15},
		{FoodID: 2, Food: "english muffin - wonder", Portion: 57, Unit: "g", Protein: 5, Carb: 25, Fibre: 1, Fat: 1.5},
		{FoodID: 3, Food: "spaghettini - great value", Portion: 85, Unit: "g", Protein: 12, Carb: 63, Fibre: 3, Fat: 1},
		{FoodID: 4, Food: "egg - medium", Portion: 1, Unit: "egg", Protein: 5.5, Carb: 0.5, Fibre: 0, Fat: 4.5},
		{FoodID: 5, Food: "whey - ON Gold Standard - Natural Vanilla", Portion: 32, Unit: "g", Protein: 24, Carb: 5, Fibre: 0, Fat: 1},
		{FoodID: 6, Food: "cassein - ON Gold Standard - Natural Chocolate", Portion: 38, Unit: "g", Protein: 24, Carb: 8, Fibre: 1, Fat: 1.5},
		{FoodID: 7, Food: "silk - oat yeah", Portion: 250, Unit: "ml", Protein: 1, Carb: 6, Fibre: 1, Fat: 2.5},
		{FoodID: 8, Food: "oatmeal - quaker - large flake", Portion: 30, Unit: "g", Protein: 4, Carb: 20, Fibre: 3, Fat: 2},
		{FoodID: 9, Food: "blueberries - frozen", Portion: 80, Unit: "g", Protein: 0.5, Carb: 10, Fibre: 2, Fat: 0},
		{FoodID: 10, Food: "rice - jasmine", Portion: 50, Unit: "g", Protein: 4, Carb: 38, Fibre: 2, Fat: 0},
		{FoodID: 12, Food: "sugar cube", Portion: 1, Unit: "cube", Protein: 0, Carb: 4, Fibre: 0, Fat: 0},
		{FoodID: 13, Food: "fudge bar", Portion: 1, Unit: "bar", Protein: 2, Carb: 18, Fibre: 1, Fat: 2},
		{FoodID: 14, Food: "float bar", Portion: 1, Unit: "bar", Protein: 0.2, Carb: 10, Fibre: 0, Fat: 2},
		{FoodID: 15, Food: "english muffin - no name", Portion: 57, Unit: "g", Protein: 5, Carb: 26, Fibre: 1, Fat: 1},
		{FoodID: 16, Food: "grapes", Portion: 100, Unit: "g", Protein: 0.72, Carb: 18.1, Fibre: 0.9, Fat: 0.16},
		// Continue adding the rest of the test data here...
	}

	err = db.AddFoods(ctx, testData)

	if err != nil {
		return fmt.Errorf("could not insert foods: %w", err)
	}

	return nil
}
