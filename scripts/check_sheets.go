package main

import (
	"fmt"
	"log"
	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("klasifikasi naive bayes tambahan data.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i, name := range f.GetSheetList() {
		fmt.Printf("Sheet %d: %s\n", i+1, name)
	}
}
