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

	rows, err := db.Query("PRAGMA table_info(residents)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Residents Table Columns:")
	for rows.Next() {
		var cid int
		var name, dtype string
		var notnull, pk int
		var dflt_value interface{}
		rows.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk)
		fmt.Printf("- %s (%s)\n", name, dtype)
	}
}
