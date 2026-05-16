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

	sheets := f.GetSheetList()
	for _, name := range sheets {
		fmt.Printf("--- Sheet: %s ---\n", name)
		rows, err := f.GetRows(name)
		if err != nil {
			fmt.Printf("Error reading sheet %s: %v\n", name, err)
			continue
		}
		for i := 0; i < 5 && i < len(rows); i++ {
			fmt.Println(rows[i])
		}
	}
}
