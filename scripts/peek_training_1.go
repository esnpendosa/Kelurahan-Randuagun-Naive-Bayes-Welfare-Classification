package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Training 1+Uji 1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Print first 25 rows of sheet 'Training 1+Uji 1':")
	for i, r := range rows {
		if i > 25 {
			break
		}
		fmt.Printf("Row %2d: ", i+1)
		for colIdx, val := range r {
			if colIdx < 10 { // limit columns
				fmt.Printf("[%s] ", val)
			}
		}
		fmt.Println()
	}
}
