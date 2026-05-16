package main

import (
	"log"
	"welfare-classification/internal/db"
	_ "modernc.org/sqlite"
)

func main() {
	// Inisialisasi koneksi database menggunakan fungsi baru
	database, err := db.InisialisasiDB("data_skripsi.db")
	if err != nil {
		log.Fatal(err) // Berhenti jika gagal membuka database
	}
	defer database.Close() // Pastikan database ditutup saat selesai

	// Daftar data penduduk untuk uji coba
	penduduk := []struct {
		nik    string
		nokk   string
		nama   string
		alamat string
		dusun  string
	}{
		{"3508011001000001", "3508010101010001", "Slamet Raharjo", "RT 01 RW 02", "Randu Baros"},
		{"3508011001000002", "3508010101010002", "M. Yusuf", "RT 05 RW 01", "Krajan"},
		{"3508011001000003", "3508010101010003", "Siti Aminah", "RT 02 RW 02", "Randu Baros"},
		{"3508011001000004", "3508010101010004", "Bambang Sugiantoro", "RT 03 RW 01", "Krajan"},
		{"3508011001000005", "3508010101010005", "Lilik Handayani", "RT 01 RW 01", "Krajan"},
	}

	// Memasukkan setiap data ke dalam tabel warga (nama tabel baru)
	for _, p := range penduduk {
		_, err := database.Exec("INSERT INTO warga (nik, no_kk, nama_lengkap, alamat, rt, rw, dusun, data_latih) VALUES (?, ?, ?, ?, ?, ?, ?, 0)",
			p.nik, p.nokk, p.nama, p.alamat, "01", "01", p.dusun)
		if err != nil {
			log.Printf("Gagal memasukkan data %s: %v", p.nama, err)
		}
	}

	log.Println("Berhasil menambahkan 5 data warga untuk uji coba.")
}
