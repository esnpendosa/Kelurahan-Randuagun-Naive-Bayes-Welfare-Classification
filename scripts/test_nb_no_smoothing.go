package main

import (
	"database/sql"
	"fmt"
	"log"

	"welfare-classification/internal/classifier"
	"welfare-classification/internal/db"

	_ "modernc.org/sqlite"
)

func main() {
	dbConn, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	dataLatih, _ := db.AmbilDataLatihSplit(dbConn, 1)
	dataUji, _ := db.AmbilDataUjiSplit(dbConn, 1)

	daftarIndikator := classifier.AmbilDaftarIndikator()
	var namaFitur []string
	for _, ind := range daftarIndikator {
		namaFitur = append(namaFitur, ind.ID)
	}

	classes := []int{1, 2, 3, 4, 5, 6}

	// Prior
	prior := make(map[int]float64)
	countPerClass := make(map[int]int)
	for _, dl := range dataLatih {
		countPerClass[dl.Kelas]++
	}
	for _, c := range classes {
		prior[c] = float64(countPerClass[c]) / float64(len(dataLatih))
	}

	// Likelihood (No Smoothing)
	// likelihood[class][feature][value]
	likelihood := make(map[int]map[string]map[string]float64)
	for _, c := range classes {
		likelihood[c] = make(map[string]map[string]float64)
		for _, f := range namaFitur {
			likelihood[c][f] = make(map[string]float64)
		}
	}

	for _, f := range namaFitur {
		for _, c := range classes {
			// Hitung count untuk setiap nilai
			counts := make(map[string]int)
			totalClass := countPerClass[c]
			for _, dl := range dataLatih {
				if dl.Kelas == c {
					counts[dl.Indikator[f]]++
				}
			}
			// Tanpa smoothing
			// Kita ambil semua possible values (A, B, C, D)
			for _, val := range []string{"A", "B", "C", "D"} {
				likelihood[c][f][val] = float64(counts[val]) / float64(totalClass)
			}
		}
	}

	// Evaluasi Data Uji
	benar := 0
	matriks := make(map[int]map[int]int)
	for _, c := range classes {
		matriks[c] = make(map[int]int)
	}

	for _, du := range dataUji {
		// Prediksi
		bestClass := -1
		bestScore := -1.0

		for _, c := range classes {
			score := prior[c]
			for _, f := range namaFitur {
				val := du.Indikator[f]
				score *= likelihood[c][f][val]
			}
			if score > bestScore {
				bestScore = score
				bestClass = c
			}
		}

		matriks[du.Kelas][bestClass]++
		if du.Kelas == bestClass {
			benar++
		}
	}

	fmt.Printf("Akurasi Tanpa Smoothing: %.4f (Benar: %d / %d)\n", float64(benar)/float64(len(dataUji)), benar, len(dataUji))
	fmt.Println("Confusion Matrix:")
	for _, row := range classes {
		fmt.Printf("KK%-2d: ", row)
		for _, col := range classes {
			fmt.Printf("%d ", matriks[row][col])
		}
		fmt.Println()
	}
}
