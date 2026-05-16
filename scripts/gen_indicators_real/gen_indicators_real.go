package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Indicator struct {
	ID      string
	Label   string
	Options map[string]string
}

func main() {
	f, err := excelize.OpenFile("data skripsi/klasifikasi naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("IM (Indikator miskin)")
	if err != nil {
		log.Fatal(err)
	}

	var indicators []Indicator
	var currentInd *Indicator

	reID := regexp.MustCompile(`\(IM\s+(\d+)\)`)

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		firstCol := strings.TrimSpace(row[0])
		
		// Check if it's a new indicator
		match := reID.FindStringSubmatch(firstCol)
		if match != nil {
			if currentInd != nil {
				indicators = append(indicators, *currentInd)
			}
			id := "IM" + match[1]
			label := strings.TrimSpace(strings.Split(firstCol, "(")[0])
			currentInd = &Indicator{
				ID:      id,
				Label:   label,
				Options: make(map[string]string),
			}
		}

		// Check for options in current row
		if currentInd != nil && len(row) >= 3 {
			optCode := strings.TrimSpace(row[1])
			optLabel := strings.TrimSpace(row[2])
			if optCode != "" && optLabel != "" {
				currentInd.Options[optCode] = optLabel
			}
		}
	}

	if currentInd != nil {
		indicators = append(indicators, *currentInd)
	}

	// Write to file
	out, err := os.Create("internal/classifier/indicators.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	fmt.Fprintln(out, "package classifier")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "func GetIndicators() []Indicator {")
	fmt.Fprintln(out, "\treturn []Indicator{")
	for _, ind := range indicators {
		fmt.Fprintf(out, "\t\t{ID: \"%s\", Label: \"%s\", Options: []string{", ind.ID, ind.Label)
		opts := []string{"A", "B", "C", "D"}
		var existingOpts []string
		for _, o := range opts {
			if _, ok := ind.Options[o]; ok {
				existingOpts = append(existingOpts, o)
			}
		}
		if len(existingOpts) == 0 {
			for k := range ind.Options {
				existingOpts = append(existingOpts, k)
			}
		}
		
		for i, o := range existingOpts {
			fmt.Fprintf(out, "\"%s\"", o)
			if i < len(existingOpts)-1 {
				fmt.Fprint(out, ", ")
			}
		}
		
		section := "Lainnya"
		idNum := 0
		fmt.Sscanf(ind.ID, "IM%d", &idNum)
		if idNum <= 12 {
			section = "Kondisi Rumah"
		} else if idNum <= 24 {
			section = "Ekonomi Keluarga"
		} else {
			section = "Aset & Fasilitas"
		}

		fmt.Fprintf(out, "}, Section: \"%s\"},\n", section)
	}
	fmt.Fprintln(out, "\t}")
	fmt.Fprintln(out, "}")
}
