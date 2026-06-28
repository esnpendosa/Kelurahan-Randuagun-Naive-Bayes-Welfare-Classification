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

	rows, err := f.GetRows("IM (Indikator miskin)")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Print first 40 rows of sheet 'IM (Indikator miskin)':")
	for i, r := range rows {
		if i > 40 {
			break
		}
		fmt.Printf("Row %2d: ", i+1)
		for _, val := range r {
			fmt.Printf("[%s] ", val)
		}
		fmt.Println()
	}
}
