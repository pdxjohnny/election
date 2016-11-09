package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gizak/termui"
)

const (
	PrezURL = "https://opb-data.s3.amazonaws.com/2016-11-08-election/presidential_results.json"
	Width   = 50
	Height  = 3
)

var (
	PartyColor = map[string]termui.Attribute{
		"GOP": termui.ColorRed,
		"Dem": termui.ColorBlue,
		"Lib": termui.ColorYellow,
		"Gre": termui.ColorGreen,
	}
)

type PrezResults struct {
	US State `json:"US"`
}

type State struct {
	Reporting  float64     `json:"reporting"`
	Called     bool        `json:"called"`
	OfficeName string      `json:"officename"`
	National   bool        `json:"national"`
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	CandidateID   string  `json:"candidateid"`
	Electtotal    int     `json:"electtotal"`
	VoteTotal     int     `json:"vote_total"`
	ElectWon      int     `json:"electwon"`
	VotePercent   float64 `json:"vote_percent"`
	Name          string  `json:"name"`
	IsWinner      bool    `json:"is_winner"`
	DelegateCount int     `json:"delegatecount"`
	Party         string  `json:"party"`
	IsIncumbent   bool    `json:"is_incumbent"`
}

func GetPrezResults() (error, *PrezResults) {
	resp, err := http.Get(PrezURL)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	p := &PrezResults{}
	err = dec.Decode(p)
	if err != nil {
		return err, nil
	}

	return nil, p
}

func main() {
	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	draw := func() {
		Y := 0
		err, p := GetPrezResults()
		if err != nil {
			log.Fatal(err)
		}

		b := make([]termui.Bufferer, 0)

		// Show how far the election is from done
		g := termui.NewGauge()
		g.Percent = int(math.Ceil(p.US.Reporting))
		g.Width = Width
		g.Height = Height
		g.Y = Y
		g.BorderLabel = p.US.OfficeName
		g.BarColor = termui.ColorGreen
		g.BorderFg = termui.ColorWhite
		g.BorderLabelFg = termui.ColorCyan
		Y += Height
		b = append(b, g)

		// Add all the canidates
		for _, c := range p.US.Candidates {
			g := termui.NewGauge()
			g.Percent = int(math.Ceil(c.VotePercent))
			g.Width = Width
			g.Height = Height
			g.Y = Y
			g.BorderLabel = c.Name
			g.BarColor = PartyColor[c.Party]
			g.BorderFg = termui.ColorWhite
			g.BorderLabelFg = termui.ColorCyan
			Y += Height
			b = append(b, g)
		}

		termui.Render(b...)
	}

	draw()

	termui.Handle("/sys/kbd/q", func(e termui.Event) {
		termui.StopLoop()
	})

	go func() {
		c := time.Tick(10 * time.Second)
		for range c {
			draw()
		}
	}()

	termui.Loop()
}
