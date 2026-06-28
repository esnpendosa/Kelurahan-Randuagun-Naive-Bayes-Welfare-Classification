package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
	_ "modernc.org/sqlite"
)

func main() {
	dbConn, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Ambil dari DB
	var id int
	dbConn.QueryRow("SELECT id FROM warga WHERE nama_lengkap = 'Nanang Fathoni'").Scan(&id)
	
	dbInds := make(map[string]string)
	rows, _ := dbConn.Query("SELECT indikator_id, nilai FROM data_indikator WHERE warga_id = ?", id)
	defer rows.Close()
	for rows.Next() {
		var indId, val string
		rows.Scan(&indId, &val)
		dbInds[indId] = val
	}

	// Ambil dari Excel (Seluruh Data Warga - Nanang Fathoni di baris 2)
	f, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	excelRows, _ := f.GetRows("Seluruh Data Warga")
	excelRow := excelRows[1] // Nanang Fathoni

	fmt.Printf("%-6s | %-6s | %-6s\n", "Ind", "DB", "Excel")
	fmt.Println("--------------------")
	for idx := 1; idx <= 36; idx++ {
		indID := fmt.Sprintf("IM%d", idx)
		excelVal := excelRow[idx+2]
		dbVal := dbInds[indID]
		fmt.Printf("%-6s | %-6s | %-6s\n", indID, dbVal, excelVal)
	}
}
