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

	rows, err := db.Query("SELECT no_urut, nama_krt, kelas_kesejahteraan FROM data_warga LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Data Warga Content:")
	for rows.Next() {
		var no int
		var nama string
		var class int
		rows.Scan(&no, &nama, &class)
		fmt.Printf("No: %d, Nama: [%s], Class: %d\n", no, nama, class)
	}
}
