package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

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
	
	classes := make(map[string]int)
	for i, row := range rows {
		if i < 2 { continue } // Skip header
		if len(row) > 2 {
			classes[row[2]]++
		}
	}
	
	fmt.Println("Classes found:", classes)
}
