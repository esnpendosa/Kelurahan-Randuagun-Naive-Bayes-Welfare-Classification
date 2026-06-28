package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"welfare-classification/internal/classifier"
	"welfare-classification/internal/db"

	"github.com/xuri/excelize/v2"
	_ "modernc.org/sqlite"
)

func main() {
	dbConn, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Inisialisasi model Naive Bayes (Laplace Smoothing)
	modelNB := classifier.BuatModelBaru()
	daftarIndikator := classifier.AmbilDaftarIndikator()
	var namaFitur []string
	for _, ind := range daftarIndikator {
		namaFitur = append(namaFitur, ind.ID)
	}
	modelNB.DaftarFitur = namaFitur

	dataLatih, _ := db.AmbilDataLatihSplit(dbConn, 1)
	dataUji, _ := db.AmbilDataUjiSplit(dbConn, 1)

	var inputLatih []map[string]string
	var targetLatih []classifier.KelasKesejahteraan
	for _, dl := range dataLatih {
		inputLatih = append(inputLatih, dl.Indikator)
		targetLatih = append(targetLatih, classifier.KelasKesejahteraan(dl.Kelas))
	}
	modelNB.LatihModel(inputLatih, targetLatih)

	// Buka Excel untuk membaca prediksi Excel
	f, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, _ := f.GetRows("Evaluasi 1")

	fmt.Printf("%-3s | %-25s | %-12s | %-12s | %-12s | %-12s\n", "No", "Nama Warga", "Aktual", "Pred Excel", "Pred Go (S)", "Pred Go (No S)")
	fmt.Println(strings.Repeat("-", 85))

	// Model tanpa smoothing
	prior := make(map[int]float64)
	countPerClass := make(map[int]int)
	for _, dl := range dataLatih {
		countPerClass[dl.Kelas]++
	}
	for _, c := range []int{1,2,3,4,5,6} {
		prior[c] = float64(countPerClass[c]) / float64(len(dataLatih))
	}
	likelihood := make(map[int]map[string]map[string]float64)
	for _, c := range []int{1,2,3,4,5,6} {
		likelihood[c] = make(map[string]map[string]float64)
		for _, f := range namaFitur {
			likelihood[c][f] = make(map[string]float64)
		}
	}
	for _, f := range namaFitur {
		for _, c := range []int{1,2,3,4,5,6} {
			counts := make(map[string]int)
			for _, dl := range dataLatih {
				if dl.Kelas == c {
					counts[dl.Indikator[f]]++
				}
			}
			for _, val := range []string{"A", "B", "C", "D"} {
				likelihood[c][f][val] = float64(counts[val]) / float64(countPerClass[c])
			}
		}
	}

	for idx, du := range dataUji {
		// Prediksi dengan smoothing
		pS := modelNB.Prediksi(du.Indikator)
		predS := modelNB.AmbilKelasTerbaik(pS)

		// Prediksi tanpa smoothing
		bestClassNoS := -1
		bestScoreNoS := -1.0
		for _, c := range []int{1,2,3,4,5,6} {
			score := prior[c]
			for _, f := range namaFitur {
				val := du.Indikator[f]
				score *= likelihood[c][f][val]
			}
			if score > bestScoreNoS {
				bestScoreNoS = score
				bestClassNoS = c
			}
		}

		// Ambil dari Excel (Row idx+2 karena baris 1 header, index excel 1-based, index slice 0-based)
		rowNum := idx + 2
		excelPred := ""
		if rowNum <= len(rows) {
			r := rows[rowNum-1]
			if len(r) > 8 {
				excelPred = r[8]
			}
		}

		fmt.Printf("%3d | %-25s | %-12s | %-12s | %-12s | KK%d\n",
			idx+1,
			du.Nama,
			classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(du.Kelas)],
			excelPred,
			classifier.DaftarNamaKelas[predS],
			bestClassNoS,
		)
	}
}
