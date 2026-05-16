package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type TrainingRecord struct {
	Name       string            `json:"name"`
	Class      string            `json:"class"`
	Indicators map[string]string `json:"indicators"`
}

func main() {
	f, err := excelize.OpenFile("data skripsi/klasifikasi naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Perhitungan naive bayes")
	if err != nil {
		log.Fatal(err)
	}

	var records []TrainingRecord
	for i, row := range rows {
		if i < 2 || len(row) < 39 { // Skip header and short rows
			continue
		}

		record := TrainingRecord{
			Name:       row[1],
			Class:      row[2],
			Indicators: make(map[string]string),
		}

		for j := 1; j <= 36; j++ {
			id := fmt.Sprintf("IM%d", j)
			val := strings.TrimSpace(row[j+2])
			if val == "" {
				val = "A" // Default fallback
			}
			record.Indicators[id] = val
		}
		records = append(records, record)
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	os.WriteFile("thesis_training_data.json", data, 0644)
	fmt.Printf("Extracted %d records\n", len(records))
}
