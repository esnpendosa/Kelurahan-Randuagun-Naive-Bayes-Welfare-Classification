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

	// Print formula for Row 2 to 10 in sheet "Evaluasi 1"
	// Column I is index 8 (Hasil)
	// Column J is index 9 (Prediksi)
	for r := 2; r <= 15; r++ {
		cellHasil := fmt.Sprintf("I%d", r)
		cellPred := fmt.Sprintf("J%d", r)
		valHasil, _ := f.GetCellValue("Evaluasi 1", cellHasil)
		formHasil, _ := f.GetCellFormula("Evaluasi 1", cellHasil)
		valPred, _ := f.GetCellValue("Evaluasi 1", cellPred)
		formPred, _ := f.GetCellFormula("Evaluasi 1", cellPred)
		
		// Print scores too (B=KK1, C=KK2, D=KK3, E=KK4, F=KK5, G=KK6)
		vB, _ := f.GetCellValue("Evaluasi 1", fmt.Sprintf("B%d", r))
		vC, _ := f.GetCellValue("Evaluasi 1", fmt.Sprintf("C%d", r))
		vD, _ := f.GetCellValue("Evaluasi 1", fmt.Sprintf("D%d", r))
		vE, _ := f.GetCellValue("Evaluasi 1", fmt.Sprintf("E%d", r))
		vF, _ := f.GetCellValue("Evaluasi 1", fmt.Sprintf("F%d", r))
		vG, _ := f.GetCellValue("Evaluasi 1", fmt.Sprintf("G%d", r))
		
		fmt.Printf("Row %2d | KK1: %s | KK2: %s | KK3: %s | KK4: %s | KK5: %s | KK6: %s |\n", r, vB, vC, vD, vE, vF, vG)
		fmt.Printf("       | Hasil Val: %-5s | Formula: %s\n", valHasil, formHasil)
		fmt.Printf("       | Pred  Val: %-5s | Formula: %s\n\n", valPred, formPred)
	}
}
