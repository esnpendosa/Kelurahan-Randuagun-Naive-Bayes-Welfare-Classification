package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "welfare.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Read and execute SQL
	sqlContent, err := os.ReadFile("scripts/import_custom.sql")
	if err != nil {
		log.Fatal(err)
	}

	queries := strings.Split(string(sqlContent), ";")
	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		// Basic comment removal for execution
		lines := strings.Split(q, "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "--") && trimmed != "" {
				cleanLines = append(cleanLines, line)
			}
		}
		cleanQuery := strings.Join(cleanLines, "\n")
		if cleanQuery == "" {
			continue
		}

		_, err := db.Exec(cleanQuery)
		if err != nil {
			log.Printf("Error executing query: %v\nQuery Snippet: %s...", err, cleanQuery[:50])
		}
	}

	fmt.Println("SQL tables created and data inserted into raw tables.")

	// 2. Sync to main residents table for UI visibility
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := tx.Query("SELECT no_urut, nama_krt, kelas_kesejahteraan, jenis_data FROM data_warga")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	
	type rawData struct {
		noUrut, classID int
		nama, jenis string
	}
	var dataToSync []rawData
	for rows.Next() {
		var d rawData
		rows.Scan(&d.noUrut, &d.nama, &d.classID, &d.jenis)
		dataToSync = append(dataToSync, d)
	}
	rows.Close()

	classMap := map[int]string{
		1: "Sangat Miskin", 2: "Miskin", 3: "Hampir Miskin",
		4: "Rentan Miskin", 5: "Pas-Pasan", 6: "Menengah ke Atas",
	}

	count := 0
	for _, d := range dataToSync {
		nik := fmt.Sprintf("SQL-%04d", d.noUrut)
		isTraining := 0
		if d.jenis == "training" {
			isTraining = 1
		}

		_, err := tx.Exec(`
			INSERT OR REPLACE INTO residents (nik, nama_kk, is_training, class_label, alamat, dusun)
			VALUES (?, ?, ?, ?, 'Imported SQL', 'SQL Data')
		`, nik, d.nama, isTraining, classMap[d.classID])
		if err != nil {
			log.Printf("Error syncing resident %s: %v", d.nama, err)
			continue
		}

		var resID int64
		tx.QueryRow("SELECT id FROM residents WHERE nik = ?", nik).Scan(&resID)

		// Fetch indicators from data_warga for this row
		row := tx.QueryRow("SELECT * FROM data_warga WHERE no_urut = ?", d.noUrut)
		cols := make([]interface{}, 40)
		var u1, u2 int
		var s1, s2 string
		vals := make([]string, 36)
		cols[0] = &u1
		cols[1] = &s1
		cols[2] = &u2
		for i := 0; i < 36; i++ {
			cols[i+3] = &vals[i]
		}
		cols[39] = &s2

		if err := row.Scan(cols...); err == nil {
			for i := 0; i < 36; i++ {
				indicatorID := fmt.Sprintf("IM%d", i+1)
				tx.Exec("INSERT OR REPLACE INTO indicator_data (resident_id, indicator_id, value) VALUES (?, ?, ?)", 
					resID, indicatorID, vals[i])
			}
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully synced %d records to main application tables.\n", count)
}
