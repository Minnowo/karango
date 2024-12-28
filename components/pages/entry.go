package pages

import (
	"fmt"
	"time"
)

const HTML_TIME_FMT string = "2006-01-02T15:04"

type EntryView struct {
	Time time.Time

	BGL     float32
	ITCR    float32
	AIT     float32
	RIA     float32
	Portion float32

	BGLIncrement     float32
	ITCRIncrement    float32
	AITIncrement     float32
	RIAIncrement     float32
	PortionIncrement float32
}

func FmtFloat(f float32) string {
	s := fmt.Sprintf("%.2f", f)
	return s
}
