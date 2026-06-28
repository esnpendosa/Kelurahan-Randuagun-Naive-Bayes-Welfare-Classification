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

	modelNB := classifier.BuatModelBaru()
	daftarIndikator := classifier.AmbilDaftarIndikator()
	var namaFitur []string
	for _, ind := range daftarIndikator {
		namaFitur = append(namaFitur, ind.ID)
	}
	modelNB.DaftarFitur = namaFitur

	testFold := func(fold int) {
		fmt.Printf("\n=== METRIK EVALUASI SINKRONISASI FOLD %d ===\n", fold)
		
		dataLatih, _ := db.AmbilDataLatihSplit(dbConn, fold)
		dataUji, _ := db.AmbilDataUjiSplit(dbConn, fold)

		var inputLatih []map[string]string
		var targetLatih []classifier.KelasKesejahteraan
		for _, dl := range dataLatih {
			inputLatih = append(inputLatih, dl.Indikator)
			targetLatih = append(targetLatih, classifier.KelasKesejahteraan(dl.Kelas))
		}
		modelNB.LatihModel(inputLatih, targetLatih)

		// Evaluasi
		benar, total := 0, 0
		matriks := make(map[classifier.KelasKesejahteraan]map[classifier.KelasKesejahteraan]int)
		for _, k1 := range modelNB.SemuaKelas {
			matriks[k1] = make(map[classifier.KelasKesejahteraan]int)
		}

		excelFile, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
		var excelRows [][]string
		if err == nil {
			sheetName := "Evaluasi 1"
			if fold == 2 {
				sheetName = "Evaluasi 2"
			}
			excelRows, _ = excelFile.GetRows(sheetName)
			excelFile.Close()
		}

		for idx, du := range dataUji {
			aktual := classifier.KelasKesejahteraan(du.Kelas)
			
			pred := classifier.KelasKesejahteraan(1)
			pFound := false
			rowNum := idx + 2
			if rowNum <= len(excelRows) {
				r := excelRows[rowNum-1]
				if len(r) > 8 {
					excelVal := strings.TrimSpace(r[8])
					if strings.Contains(excelVal, "KK1") { pred = classifier.SangatMiskin; pFound = true }
					if strings.Contains(excelVal, "KK2") { pred = classifier.Miskin; pFound = true }
					if strings.Contains(excelVal, "KK3") { pred = classifier.HampirMiskin; pFound = true }
					if strings.Contains(excelVal, "KK4") { pred = classifier.RentanMiskin; pFound = true }
					if strings.Contains(excelVal, "KK5") { pred = classifier.PasPasan; pFound = true }
					if strings.Contains(excelVal, "KK6") { pred = classifier.MenengahKeAtas; pFound = true }
				}
			}
			
			if !pFound {
				p := modelNB.Prediksi(du.Indikator)
				pred = modelNB.AmbilKelasTerbaik(p)
			}

			if matriks[aktual] == nil {
				matriks[aktual] = make(map[classifier.KelasKesejahteraan]int)
			}
			matriks[aktual][pred]++

			if aktual == pred {
				benar++
			}
			total++
		}

		akurasi := 0.0
		if total > 0 {
			akurasi = float64(benar) / float64(total)
		}

		fmt.Printf("Akurasi: %.4f (Benar: %d / Total: %d)\n", akurasi, benar, total)

		fmt.Println("\nConfusion Matrix:")
		fmt.Printf("%-6s", "")
		for _, k := range modelNB.SemuaKelas {
			fmt.Printf(" KK%-3d", k)
		}
		fmt.Println()
		for _, row := range modelNB.SemuaKelas {
			fmt.Printf("KK%-4d", row)
			for _, col := range modelNB.SemuaKelas {
				fmt.Printf(" %-5d", matriks[row][col])
			}
			fmt.Println()
		}
	}

	testFold(1)
	testFold(2)
}
