package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, password, role FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("ID | Username | Password (Hash) | Role")
	fmt.Println("--------------------------------------")
	for rows.Next() {
		var id int
		var u, p, r string
		rows.Scan(&id, &u, &p, &r)
		fmt.Printf("%d | %s | %s | %s\n", id, u, p, r)
	}
}
