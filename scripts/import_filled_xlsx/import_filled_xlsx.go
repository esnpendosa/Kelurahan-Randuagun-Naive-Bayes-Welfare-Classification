package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Cleaning operational tables...")
	db.Exec("DELETE FROM residents")
	db.Exec("DELETE FROM indicator_data")
	db.Exec("DELETE FROM classification_results")

	f, err := excelize.OpenFile("Template_Import_Warga_filled.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for i, row := range rows {
		if i == 0 || len(row) < 7 { continue } // Skip header

		nik := row[0]
		nama := row[1]
		alamat := row[2]
		dusun := row[3]
		label := row[5] // LABEL_KELAS
		
		isTraining := 1
		if strings.Contains(strings.ToLower(nama), "testing") {
			isTraining = 0
		}

		res, err := db.Exec(`
			INSERT INTO residents (nik, nama_kk, is_training, class_label, alamat, dusun, no_kk)
			VALUES (?, ?, ?, ?, ?, ?, '')
		`, nik, nama, isTraining, label, alamat, dusun)
		if err != nil {
			log.Printf("Row %d Error: %v", i, err)
			continue
		}

		resID, _ := res.LastInsertId()

		// Indicators start at index 6 (IM1)
		for colIdx := 6; colIdx < len(row); colIdx++ {
			indicatorID := fmt.Sprintf("IM%d", colIdx-5)
			val := strings.ToUpper(strings.TrimSpace(row[colIdx]))
			if val != "" {
				db.Exec("INSERT INTO indicator_data (resident_id, indicator_id, value) VALUES (?, ?, ?)", 
					resID, indicatorID, val)
			}
		}
		count++
	}

	fmt.Printf("Success! Imported %d records from Excel into data_skripsi.db\n", count)
}
