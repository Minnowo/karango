package database

type Food struct {
	FoodID  int     `db:"food_id"`
	Food    string  `db:"food"`
	Portion float64 `db:"portion"`
	Unit    string  `db:"unit"`
	Protein float64 `db:"protein"`
	Carb    float64 `db:"carb"`
	Fibre   float64 `db:"fibre"`
	Fat     float64 `db:"fat"`
}

type Version struct {
	Version int `db:"version"`
}

type Flag struct {
	Flag int `db:"flag"`
}
