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

	rows, err := db.Query("SELECT nik, nama_kk, is_training FROM residents ORDER BY id DESC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Recent Residents:")
	for rows.Next() {
		var nik, nama string
		var isT int
		rows.Scan(&nik, &nama, &isT)
		fmt.Printf("NIK: [%s], Nama: [%s], Training: %d\n", nik, nama, isT)
	}
}
