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

	var id int
	var nik, name string
	err = db.QueryRow("SELECT id, nik, nama_kk FROM residents WHERE id = 1422").Scan(&id, &nik, &name)
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Printf("ID: %d, NIK: %s, Name: [%s]\n", id, nik, name)
	}
}
