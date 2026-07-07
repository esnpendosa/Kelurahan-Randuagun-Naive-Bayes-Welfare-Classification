package main // Program utama untuk menjalankan script membaca file Excel

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2" // Menggunakan pustaka excelize untuk membaca berkas spreadsheet (.xlsx)
)

func main() {
	// Membuka berkas Excel hasil perhitungan manual skripsi
	f, err := excelize.OpenFile("data skripsi/klasifikasi naive bayes.xlsx")
	if err != nil {
		log.Fatal(err) // Hentikan script jika berkas Excel tidak ditemukan atau gagal dibuka
	}
	defer f.Close() // Pastikan berkas Excel ditutup secara aman setelah selesai digunakan

	// Membaca seluruh baris data dari sheet bernama "Perhitungan naive bayes"
	rows, err := f.GetRows("Perhitungan naive bayes")
	if err != nil {
		log.Fatal(err) // Hentikan script jika sheet tidak ditemukan
	}
	
	// Membuat peta (map) untuk menghitung distribusi atau jumlah warga per kelas kesejahteraan
	classes := make(map[string]int)
	for i, row := range rows {
		if i < 2 { continue } // Melewati baris ke-1 dan ke-2 karena merupakan header kolom
		if len(row) > 2 {
			classes[row[2]]++ // Menghitung kemunculan kelas kesejahteraan (kolom indeks ke-2)
		}
	}
	
	// Menampilkan statistik jumlah kelas kesejahteraan warga yang ditemukan di Excel ke terminal
	fmt.Println("Classes found:", classes)
}
