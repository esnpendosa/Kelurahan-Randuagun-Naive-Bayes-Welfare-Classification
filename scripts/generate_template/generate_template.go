package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"NIK", "NAMA", "ALAMAT", "DUSUN"}
	for i := 1; i <= 36; i++ {
		headers = append(headers, fmt.Sprintf("IM%d", i))
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Example Row
	example := []string{"1234567890", "Contoh Nama", "Jl. Mawar No. 1", "Pusat"}
	for i := 1; i <= 36; i++ {
		example = append(example, "A") // Example value
	}

	for i, v := range example {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, v)
	}

	if err := f.SaveAs("Template_Import_Warga.xlsx"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Template Excel 'Template_Import_Warga.xlsx' berhasil dibuat.")
}
