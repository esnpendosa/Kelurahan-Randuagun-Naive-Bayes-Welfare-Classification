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

	var total, training, testing int
	db.QueryRow("SELECT COUNT(*) FROM residents").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM residents WHERE is_training = 1").Scan(&training)
	db.QueryRow("SELECT COUNT(*) FROM residents WHERE is_training = 0").Scan(&testing)

	fmt.Printf("Total: %d\nTraining: %d\nTesting: %d\n", total, training, testing)
}
