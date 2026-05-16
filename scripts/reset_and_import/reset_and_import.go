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
	db, err := sql.Open("sqlite", "data_skripsi.db?_busy_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Cleaning database...")
	db.Exec("DELETE FROM residents")
	db.Exec("DELETE FROM indicator_data")
	db.Exec("DELETE FROM classification_results")
	db.Exec("DELETE FROM data_warga")

	fmt.Println("Importing custom SQL...")
	sqlContent, err := os.ReadFile("scripts/import_custom.sql")
	if err != nil {
		log.Fatal(err)
	}

	queries := strings.Split(string(sqlContent), ";")
	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" { continue }
		db.Exec(q)
	}

	fmt.Println("Syncing data to application tables...")
	rows, err := db.Query("SELECT no_urut, nama_krt, kelas_kesejahteraan, jenis_data FROM data_warga")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var noUrut, classID int
		var nama, jenis string
		if err := rows.Scan(&noUrut, &nama, &classID, &jenis); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		nik := fmt.Sprintf("SQL-%04d", noUrut)
		isTraining := 0
		if jenis == "training" {
			isTraining = 1
		}

		classMap := map[int]string{
			1: "Sangat Miskin", 2: "Miskin", 3: "Hampir Miskin",
			4: "Rentan Miskin", 5: "Pas-Pasan", 6: "Menengah ke Atas",
		}

		res, err := db.Exec(`
			INSERT INTO residents (nik, nama_kk, is_training, class_label, alamat, dusun, no_kk)
			VALUES (?, ?, ?, ?, 'Randuagung', 'Pusat', '')
		`, nik, nama, isTraining, classMap[classID])
		if err != nil {
			log.Printf("Insert error: %v", err)
			continue
		}

		resID, _ := res.LastInsertId()

		// Fetch indicators for this row
		row := db.QueryRow("SELECT * FROM data_warga WHERE no_urut = ?", noUrut)
		cols := make([]interface{}, 40)
		var u1, u2 int
		var s1, s2 string
		vals := make([]string, 36)
		cols[0], cols[1], cols[2] = &u1, &s1, &u2
		for i := 0; i < 36; i++ { cols[i+3] = &vals[i] }
		cols[39] = &s2

		if err := row.Scan(cols...); err == nil {
			for i := 0; i < 36; i++ {
				indicatorID := fmt.Sprintf("IM%d", i+1)
				db.Exec("INSERT INTO indicator_data (resident_id, indicator_id, value) VALUES (?, ?, ?)", 
					resID, indicatorID, vals[i])
			}
		}
		count++
	}

	fmt.Printf("Success! Imported %d records into application.\n", count)
}
