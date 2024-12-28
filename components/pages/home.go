package pages

import (
	"fmt"
	"time"
)

type Food struct {
	Name    string  `json:"name"`
	Portion float32 `json:"portion"`
	Unit    string  `json:"unit"`
	Protein float32 `json:"protein"`
	Carbs   float32 `json:"carbs"`
	Fibre   float32 `json:"fibre"`
	Fat     float32 `json:"fat"`
}

func (f *Food) NetCarbs() float32 {
	return f.Carbs - f.Fibre
}

type Event struct {
	Event             string    `json:"event"`
	Time              time.Time `json:"time"`
	BG                float32   `json:"bg"`
	ITCR              float32   `json:"itcr"`
	ActualTaken       float32   `json:"actual_taken"`
	RecommendedAmount float32   `json:"recommended_amount"`
	ISF               float32   `json:"isf"`
	BGT               float32   `json:"bgt"`
	Foods             []Food    `json:"foods"`
}

func (e *Event) Carbs() float32 {

	s := float32(0)
	for _, i := range e.Foods {
		s += i.NetCarbs()
	}
	return s
}
func (e *Event) Fat() float32 {

	s := float32(0)
	for _, i := range e.Foods {
		s += i.Fat
	}
	return s
}
func (e *Event) Fibre() float32 {

	s := float32(0)
	for _, i := range e.Foods {
		s += i.Fibre
	}
	return s
}
func (e *Event) Protein() float32 {

	s := float32(0)
	for _, i := range e.Foods {
		s += i.Protein
	}
	return s
}

type Day struct {
	Day    string  `json:"day"`
	Events []Event `json:"events"`
}

func (d *Day) Carbs() float32 {
	s := float32(0)
	for _, i := range d.Events {
		s += i.Carbs()
	}
	return s
}
func (d *Day) Fat() float32 {
	s := float32(0)
	for _, i := range d.Events {
		s += i.Fat()
	}
	return s
}
func (d *Day) Fibre() float32 {
	s := float32(0)
	for _, i := range d.Events {
		s += i.Fibre()
	}
	return s
}
func (d *Day) Protein() float32 {
	s := float32(0)
	for _, i := range d.Events {
		s += i.Protein()
	}
	return s
}

type HomeView struct {
	Days []Day
}

func F32Str(f float32) string {

	s := fmt.Sprintf("%.2f", f)
	return s
}
