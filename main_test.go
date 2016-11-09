package main

import (
	"log"
	"testing"

	"github.com/gizak/termui"
)

func TestGetPrezResults(t *testing.T) {
	err, p := GetPrezResults()
	if err != nil {
		t.Fatal(err)
	}

	b := make([]termui.Bufferer, 0)

	for _, c := range p.US.Candidates {
		g := termui.NewGauge()
		g.Percent = int(c.VotePercent)
		g.Width = 50
		g.Height = 3
		g.BorderLabel = c.Name
		g.BarColor = PartyColor[c.Party]
		g.BorderFg = termui.ColorWhite
		g.BorderLabelFg = termui.ColorCyan
		log.Println(g)
		b = append(b, g)
	}

	log.Println(b)
}
