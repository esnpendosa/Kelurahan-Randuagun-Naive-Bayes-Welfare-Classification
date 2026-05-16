package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "welfare.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nik, no_kk, nama_kk FROM residents LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var nik, no_kk, nama_kk string
		rows.Scan(&id, &nik, &no_kk, &nama_kk)
		fmt.Printf("ID: %d | NIK: %s | NoKK: %s | Nama: %s\n", id, nik, no_kk, nama_kk)
	}
}
