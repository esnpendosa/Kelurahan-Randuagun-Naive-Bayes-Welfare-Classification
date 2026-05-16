package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("Template_Import_Warga_filled.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("First 3 rows of Excel:")
	for i, row := range rows {
		if i > 3 { break }
		fmt.Printf("Row %d: %v\n", i, row)
	}
}
