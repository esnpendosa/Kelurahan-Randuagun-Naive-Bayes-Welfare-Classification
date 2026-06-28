package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var totalWarga int
	db.QueryRow("SELECT COUNT(*) FROM warga").Scan(&totalWarga)
	fmt.Printf("Total Warga: %d\n", totalWarga)

	var l1, u1 int
	db.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih = 1").Scan(&l1)
	db.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih = 0 AND label_kelas != ''").Scan(&u1)
	fmt.Printf("Dataset 1: Latih = %d, Uji = %d\n", l1, u1)

	var l2, u2 int
	db.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih_2 = 1").Scan(&l2)
	db.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih_2 = 0 AND label_kelas != ''").Scan(&u2)
	fmt.Printf("Dataset 2: Latih = %d, Uji = %d\n", l2, u2)

	// Tampilkan sampel data
	fmt.Println("\nSampel Warga:")
	rows, err := db.Query("SELECT id, nama_lengkap, data_latih, data_latih_2, label_kelas FROM warga LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, dl, dl2 int
		var nama, kelas string
		rows.Scan(&id, &nama, &dl, &dl2, &kelas)
		fmt.Printf("ID: %2d | Nama: %-25s | Latih 1: %d | Latih 2: %d | Kelas: %s\n", id, nama, dl, dl2, kelas)
	}
}
