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

	// Inisialisasi model Naive Bayes
	modelNB := classifier.BuatModelBaru()
	daftarIndikator := classifier.AmbilDaftarIndikator()
	var namaFitur []string
	for _, ind := range daftarIndikator {
		namaFitur = append(namaFitur, ind.ID)
	}
	modelNB.DaftarFitur = namaFitur

	testFold := func(fold int) {
		fmt.Printf("\n=== EVALUASI FOLD %d ===\n", fold)
		dataLatih, err := db.AmbilDataLatihSplit(dbConn, fold)
		if err != nil {
			log.Fatalf("Error data latih fold %d: %v", fold, err)
		}
		dataUji, err := db.AmbilDataUjiSplit(dbConn, fold)
		if err != nil {
			log.Fatalf("Error data uji fold %d: %v", fold, err)
		}

		fmt.Printf("Data Latih: %d, Data Uji: %d\n", len(dataLatih), len(dataUji))

		// Latih model
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

		for _, du := range dataUji {
			p := modelNB.Prediksi(du.Indikator)
			pred := modelNB.AmbilKelasTerbaik(p)
			aktual := classifier.KelasKesejahteraan(du.Kelas)
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

		fmt.Println("\nMetrik Per Kelas:")
		var totP, totR float64
		for _, k := range modelNB.SemuaKelas {
			tp := float64(matriks[k][k])
			var fp, fn float64
			for _, ac := range modelNB.SemuaKelas {
				if ac != k {
					fp += float64(matriks[ac][k])
				}
			}
			for _, pc := range modelNB.SemuaKelas {
				if pc != k {
					fn += float64(matriks[k][pc])
				}
			}
			pr := 0.0
			if tp+fp > 0 {
				pr = tp / (tp + fp)
			}
			rc := 0.0
			if tp+fn > 0 {
				rc = tp / (tp + fn)
			}
			f1 := 0.0
			if pr+rc > 0 {
				f1 = 2 * pr * rc / (pr + rc)
			}
			totP += pr
			totR += rc
			fmt.Printf("Kelas KK%-2d (%-18s) | Precision: %.4f | Recall: %.4f | F1-Score: %.4f\n",
				k, classifier.DaftarNamaKelas[k], pr, rc, f1)
		}
		macroP := totP / 6.0
		macroR := totR / 6.0
		macroF1 := 0.0
		if macroP+macroR > 0 {
			macroF1 = 2 * macroP * macroR / (macroP + macroR)
		}
		fmt.Printf("\nMacro Average Precision: %.4f\n", macroP)
		fmt.Printf("Macro Average Recall   : %.4f\n", macroR)
		fmt.Printf("Macro Average F1-Score : %.4f\n", macroF1)
	}

	testFold(1)
	testFold(2)
}
