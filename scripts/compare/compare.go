package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
	_ "modernc.org/sqlite"

	"welfare-classification/internal/classifier"
	"welfare-classification/internal/db"
)

func main() {
	// Open DB
	dbConn, err := sql.Open("sqlite", "data_skripsi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Open Excel
	excelFile, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer excelFile.Close()

	evalRows, err := excelFile.GetRows("Evaluasi 1")
	if err != nil {
		log.Fatal(err)
	}

	// Features
	namaFitur := make([]string, 36)
	for i := 1; i <= 36; i++ {
		namaFitur[i-1] = fmt.Sprintf("IM%d", i)
	}

	// Train model
	dataLatih, err := db.AmbilDataLatihSplit(dbConn, 1)
	if err != nil {
		log.Fatal(err)
	}

	model := classifier.BuatModelBaru()
	model.DaftarFitur = namaFitur

	var in []map[string]string
	var tg []classifier.KelasKesejahteraan
	for _, dl := range dataLatih {
		in = append(in, dl.Indikator)
		tg = append(tg, classifier.KelasKesejahteraan(dl.Kelas))
	}
	model.LatihModel(in, tg)

	// Test data
	dataUji, err := db.AmbilDataUjiSplit(dbConn, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-20s | %-12s | %-12s | %-25s | %-25s\n", "Nama", "Excel Pred", "Go Pred", "Excel Probs (KK1, KK2)", "Go Probs (KK1, KK2)")
	fmt.Println(strings.Repeat("-", 110))

	for _, du := range dataUji {
		// Go prediction
		pGo := model.Prediksi(du.Indikator)
		predGo := model.AmbilKelasTerbaik(pGo)

		// Excel prediction & probs
		var excelPred string
		var excelKK1, excelKK2 string
		for _, r := range evalRows {
			if len(r) > 9 && strings.EqualFold(strings.TrimSpace(r[1]), strings.TrimSpace(du.Nama)) {
				excelPred = r[9]
				excelKK1 = r[2]
				excelKK2 = r[3]
				break
			}
		}

		fmt.Printf("%-20s | %-12s | %-12s | (%-10s, %-10s) | (%.4E, %.4E)\n",
			du.Nama, excelPred, classifier.DaftarNamaKelas[predGo], excelKK1, excelKK2, pGo[classifier.SangatMiskin], pGo[classifier.Miskin])
	}
}
