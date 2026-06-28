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
	dbConn, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Ambil nama dari DB (uji split 1)
	rowsDB, _ := dbConn.Query("SELECT nama_lengkap FROM warga WHERE data_latih = 0 AND label_kelas != '' ORDER BY id ASC")
	defer rowsDB.Close()
	var namesDB []string
	for rowsDB.Next() {
		var name string
		rowsDB.Scan(&name)
		namesDB = append(namesDB, name)
	}

	// Ambil nama dari Excel (Data Uji 1)
	f, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rowsExcel, _ := f.GetRows("Data Uji 1")
	var namesExcel []string
	for i, r := range rowsExcel {
		if i == 0 {
			continue
		}
		if len(r) > 1 && strings.TrimSpace(r[1]) != "" {
			namesExcel = append(namesExcel, strings.TrimSpace(r[1]))
		}
	}

	fmt.Printf("DB Uji Count: %d, Excel Uji Count: %d\n", len(namesDB), len(namesExcel))

	// Cek perbedaan
	fmt.Println("\nNames in DB but NOT in Excel:")
	for _, n := range namesDB {
		found := false
		for _, ne := range namesExcel {
			if n == ne {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("- %q\n", n)
		}
	}

	fmt.Println("\nNames in Excel but NOT in DB:")
	for _, ne := range namesExcel {
		found := false
		for _, n := range namesDB {
			if n == ne {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("- %q\n", ne)
		}
	}
}
