package main

import (
	"fmt"
	"log"
	"strings"

	"welfare-classification/internal/db"

	"github.com/xuri/excelize/v2"
	_ "modernc.org/sqlite"
)

func main() {
	// Use InisialisasiDB to ensure schema is initialized/migrated
	sqlDB, err := db.InisialisasiDB("data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Start a transaction
	tx, err := sqlDB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resetting database tables...")
	tx.Exec("DELETE FROM data_indikator")
	tx.Exec("DELETE FROM hasil_klasifikasi")
	tx.Exec("DELETE FROM warga")

	f, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer f.Close()

	// Read all residents from "Seluruh Data Warga"
	rows, err := f.GetRows("Seluruh Data Warga")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	classMap := map[string]string{
		"1": "Sangat Miskin",
		"2": "Miskin",
		"3": "Hampir Miskin",
		"4": "Rentan Miskin",
		"5": "Pas-pasan",
		"6": "Menengah ke Atas",
	}

	fmt.Printf("Importing %d residents...\n", len(rows)-1)
	insertedNamesMap := make(map[string]int64)

	stmtWarga, err := tx.Prepare(`
		INSERT INTO warga (nik, no_kk, nama_lengkap, alamat, label_kelas, data_latih, data_latih_2)
		VALUES (?, ?, ?, ?, ?, 0, 0)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmtWarga.Close()

	stmtInd, err := tx.Prepare("INSERT INTO data_indikator (warga_id, indikator_id, nilai) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmtInd.Close()

	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if len(row) < 3 || strings.TrimSpace(row[1]) == "" {
			continue
		}

		name := strings.TrimSpace(row[1])
		classCode := strings.TrimSpace(row[2])
		className := classMap[classCode]
		if className == "" {
			className = classCode // fallback
		}

		// Generate NIK/KK dummy
		nik := fmt.Sprintf("350801%010d", i)
		kk := fmt.Sprintf("350801%010d", i+10000)
		alamat := "Dusun Randuagung RT 01 RW 01"

		res, err := stmtWarga.Exec(nik, kk, name, alamat, className)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Row %d Error inserting resident: %v", i, err)
		}

		wargaID, _ := res.LastInsertId()
		insertedNamesMap[name] = wargaID

		// Insert 36 indicators
		for colIdx := 3; colIdx < len(row) && colIdx < 39; colIdx++ {
			indID := fmt.Sprintf("IM%d", colIdx-2)
			val := strings.ToUpper(strings.TrimSpace(row[colIdx]))
			if val != "" {
				_, err = stmtInd.Exec(wargaID, indID, val)
				if err != nil {
					tx.Rollback()
					log.Fatalf("Error inserting indicator: %v", err)
				}
			}
		}
	}

	fmt.Println("Resident profiles and indicators successfully prepared.")

	// Set split 1 data latih
	t1Rows, err := f.GetRows("Data Training 1")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	t1Count := 0
	for i, row := range t1Rows {
		if i == 0 || len(row) < 2 || strings.TrimSpace(row[1]) == "" {
			continue
		}
		name := strings.TrimSpace(row[1])
		wargaID, exists := insertedNamesMap[name]
		if exists {
			tx.Exec("UPDATE warga SET data_latih = 1 WHERE id = ?", wargaID)
			t1Count++
		} else {
			fmt.Printf("Warning: Training 1 name %q not found in all residents!\n", name)
		}
	}
	fmt.Printf("Prepared %d residents as training data for Split 1.\n", t1Count)

	// Set split 2 data latih
	t2Rows, err := f.GetRows("Data Training 2")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	t2Count := 0
	for i, row := range t2Rows {
		if i == 0 || len(row) < 2 || strings.TrimSpace(row[1]) == "" {
			continue
		}
		name := strings.TrimSpace(row[1])
		wargaID, exists := insertedNamesMap[name]
		if exists {
			tx.Exec("UPDATE warga SET data_latih_2 = 1 WHERE id = ?", wargaID)
			t2Count++
		} else {
			fmt.Printf("Warning: Training 2 name %q not found in all residents!\n", name)
		}
	}
	fmt.Printf("Prepared %d residents as training data for Split 2.\n", t2Count)

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All done! Data synchronized successfully via transaction.")
}
