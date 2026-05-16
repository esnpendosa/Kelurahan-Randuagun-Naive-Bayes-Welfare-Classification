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

	rows, err := db.Query("SELECT * FROM residents LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	fmt.Println("Columns in result set:", cols)

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, col := range cols {
			val := values[i]
			fmt.Printf("%s: %v\n", col, val)
		}
	}
}
