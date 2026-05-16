package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"welfare-classification/internal/db"

	_ "modernc.org/sqlite"
)

type TrainingData struct {
	Name       string            `json:"name"`
	Class      string            `json:"class"`
	Indicators map[string]string `json:"indicators"`
}

func main() {
	// Membaca file JSON data training
	file, err := os.ReadFile("thesis_training_data.json")
	if err != nil {
		log.Fatal(err)
	}

	var data []TrainingData
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Inisialisasi koneksi database menggunakan fungsi baru
	database, err := db.InisialisasiDB("data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Memulai transaksi database agar proses sinkronisasi cepat dan aman
	tx, err := database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Menyiapkan statement SQL untuk memasukkan data warga (ditambah NIK, No KK, dan Alamat dummy)
	stmtResident, err := tx.Prepare("INSERT OR IGNORE INTO warga (nik, no_kk, nama_lengkap, alamat, data_latih, label_kelas) VALUES (?, ?, ?, ?, 1, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtResident.Close()

	// Menyiapkan statement SQL untuk memasukkan data indikator
	stmtIndicator, err := tx.Prepare("INSERT INTO data_indikator (warga_id, indikator_id, nilai) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIndicator.Close()

	fmt.Printf("Sinkronisasi %d data sedang berjalan...\n", len(data))

	for i, record := range data {
		// Buat NIK dan KK dummy berdasarkan indeks agar unik
		nikDummy := fmt.Sprintf("350801%010d", i+1)
		kkDummy := fmt.Sprintf("350801%010d", i+10001)
		alamatDummy := "Dusun Randuagun RT 01 RW 01"

		// Simpan data warga
		res, err := stmtResident.Exec(nikDummy, kkDummy, record.Name, alamatDummy, record.Class)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Error saat memasukkan warga ke-%d: %v", i, err)
		}

		idWarga, _ := res.LastInsertId()

		// Simpan 36 indikator untuk warga tersebut
		for id, nilai := range record.Indicators {
			_, err = stmtIndicator.Exec(idWarga, id, nilai)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Error saat memasukkan indikator %s untuk ID %d: %v", id, idWarga, err)
			}
		}

		if (i+1)%100 == 0 {
			fmt.Printf("Sudah memproses %d/%d data\n", i+1, len(data))
		}
	}

	// Commit semua perubahan ke database
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sinkronisasi database berhasil diselesaikan!")
}
