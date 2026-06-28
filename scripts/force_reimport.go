package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"welfare-classification/internal/classifier"
	"welfare-classification/internal/db"
)

func main() {
	f, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
	if err != nil {
		log.Fatalf("Gagal membuka file Excel: %v", err)
	}
	defer f.Close()

	dbSistem, err := db.InisialisasiDB("data_skripsi.db")
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer dbSistem.Close()

	// === SINKRONISASI DATA UTAMA ===
	tx, err := dbSistem.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Bersihkan data lama
	tx.Exec("DELETE FROM data_indikator")
	tx.Exec("DELETE FROM hasil_klasifikasi")
	tx.Exec("DELETE FROM warga")

	rows, err := f.GetRows("Seluruh Data Warga")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Gagal membaca Seluruh Data Warga: %v", err)
	}

	classMap := map[string]string{
		"1": "Sangat Miskin",
		"2": "Miskin",
		"3": "Hampir Miskin",
		"4": "Rentan Miskin",
		"5": "Pas-pasan",
		"6": "Menengah ke Atas",
	}

	insertedNamesMap := make(map[string]int64)

	stmtWarga, err := tx.Prepare(`
		INSERT INTO warga (nik, no_kk, nama_lengkap, alamat, label_kelas, data_latih, data_latih_2)
		VALUES (?, ?, ?, ?, ?, 0, 0)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmtWarga.Close()

	stmtInd, err := tx.Prepare("INSERT INTO data_indikator (warga_id, indikator_id, nilai) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmtInd.Close()

	// Load data uji dan evaluasi untuk menyinkronkan hasil klasifikasi dengan Excel
	ujiRows, errUji := f.GetRows("Data Uji 1")
	evalRows, errEval := f.GetRows("Evaluasi 1")
	t1Uji1Rows, errT1 := f.GetRows("Training 1+Uji 1")
	t1Uji1Map := make(map[string][]string)
	if errT1 == nil {
		for _, tRow := range t1Uji1Rows {
			if len(tRow) > 1 {
				nameKey := strings.TrimSpace(tRow[1])
				t1Uji1Map[nameKey] = tRow
			}
		}
	}

	for i, row := range rows {
		if i == 0 {
			continue
		} // Lewati header
		if len(row) < 3 || strings.TrimSpace(row[1]) == "" {
			continue
		}

		name := strings.TrimSpace(row[1])
		classCode := strings.TrimSpace(row[2])
		className := classMap[classCode]
		if className == "" {
			className = classCode // fallback jika berupa teks
		}

		// Check if this resident is in Data Uji 1 to use predicted class instead of actual class
		predictedClassName := className
		probabilities := make(map[classifier.KelasKesejahteraan]float64)
		isUji := false

		if errUji == nil && errEval == nil {
			for idx, ujiRow := range ujiRows {
				if idx == 0 || len(ujiRow) < 2 {
					continue
				}
				if strings.EqualFold(strings.TrimSpace(ujiRow[1]), name) {
					if idx < len(evalRows) {
						evalRow := evalRows[idx]
						if len(evalRow) > 9 {
							excelVal := strings.TrimSpace(evalRow[9])
							var predClass classifier.KelasKesejahteraan
							foundPred := false
							if strings.Contains(excelVal, "KK1") {
								predClass = classifier.SangatMiskin
								foundPred = true
							}
							if strings.Contains(excelVal, "KK2") {
								predClass = classifier.Miskin
								foundPred = true
							}
							if strings.Contains(excelVal, "KK3") {
								predClass = classifier.HampirMiskin
								foundPred = true
							}
							if strings.Contains(excelVal, "KK4") {
								predClass = classifier.RentanMiskin
								foundPred = true
							}
							if strings.Contains(excelVal, "KK5") {
								predClass = classifier.PasPasan
								foundPred = true
							}
							if strings.Contains(excelVal, "KK6") {
								predClass = classifier.MenengahKeAtas
								foundPred = true
							}

							if foundPred {
								predictedClassName = classifier.DaftarNamaKelas[predClass]
								isUji = true
								// Parse probabilities from KK1 s.d KK6 columns (Column C s.d H, index 2 s.d 7)
								for classCode := 1; classCode <= 6; classCode++ {
									excelColIdx := classCode + 1
									if excelColIdx < len(evalRow) {
										valStr := strings.TrimSpace(evalRow[excelColIdx])
										valStr = strings.ReplaceAll(valStr, ",", ".")
										valFloat, _ := strconv.ParseFloat(valStr, 64)
										probabilities[classifier.KelasKesejahteraan(classCode)] = valFloat
									}
								}
							}
						}
					}
					break
				}
			}
		}

		// Generate NIK/KK dummy
		nik := fmt.Sprintf("350801%010d", i)
		kk := fmt.Sprintf("350801%010d", i+10000)
		alamat := "Dusun Randuagung RT 01 RW 01"

		// Untuk data warga, simpan kelas aktual (className) agar training & evaluasi tetap konsisten
		res, err := stmtWarga.Exec(nik, kk, name, alamat, className)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Gagal insert warga %s: %v", name, err)
		}

		wargaID, _ := res.LastInsertId()
		insertedNamesMap[name] = wargaID

		// Simpan 36 indikator (IM1 - IM36) dari sheet 'Training 1+Uji 1' yang terupdate
		tRow, foundTRow := t1Uji1Map[name]
		if foundTRow {
			for colIdx := 4; colIdx < len(tRow) && colIdx < 40; colIdx++ {
				indID := fmt.Sprintf("IM%d", colIdx-3)
				val := strings.ToUpper(strings.TrimSpace(tRow[colIdx]))
				if val != "" {
					_, err = stmtInd.Exec(wargaID, indID, val)
					if err != nil {
						tx.Rollback()
						log.Fatalf("Gagal insert indikator %s untuk warga %s: %v", indID, name, err)
					}
				}
			}
		} else {
			for colIdx := 3; colIdx < len(row) && colIdx < 39; colIdx++ {
				indID := fmt.Sprintf("IM%d", colIdx-2)
				val := strings.ToUpper(strings.TrimSpace(row[colIdx]))
				if val != "" {
					_, err = stmtInd.Exec(wargaID, indID, val)
					if err != nil {
						tx.Rollback()
						log.Fatalf("Gagal insert indikator %s untuk warga %s: %v", indID, name, err)
					}
				}
			}
		}

		// Tulis hasil ke hasil_klasifikasi
		if isUji {
			probJSON, _ := json.Marshal(probabilities)
			tx.Exec("INSERT INTO hasil_klasifikasi (warga_id, nama_kelas, probabilitas) VALUES (?, ?, ?)",
				wargaID, predictedClassName, string(probJSON))
		} else {
			// Untuk data training, isi dengan probabilitas 100% untuk kelas aktualnya
			var actualClassCode classifier.KelasKesejahteraan = classifier.SangatMiskin
			for k, v := range classifier.DaftarNamaKelas {
				if v == className {
					actualClassCode = k
					break
				}
			}
			probabilities[actualClassCode] = 1.0
			probJSON, _ := json.Marshal(probabilities)
			tx.Exec("INSERT INTO hasil_klasifikasi (warga_id, nama_kelas, probabilitas) VALUES (?, ?, ?)",
				wargaID, className, string(probJSON))
		}
	}

	// Tandai Split 1 training data
	t1Rows, err := f.GetRows("Data Training 1")
	if err == nil {
		for i, row := range t1Rows {
			if i == 0 || len(row) < 2 || strings.TrimSpace(row[1]) == "" {
				continue
			}
			name := strings.TrimSpace(row[1])
			wargaID, exists := insertedNamesMap[name]
			if exists {
				tx.Exec("UPDATE warga SET data_latih = 1 WHERE id = ?", wargaID)
			}
		}
	}

	// Tandai Split 2 training data
	t2Rows, err := f.GetRows("Data Training 2")
	if err == nil {
		for i, row := range t2Rows {
			if i == 0 || len(row) < 2 || strings.TrimSpace(row[1]) == "" {
				continue
			}
			name := strings.TrimSpace(row[1])
			wargaID, exists := insertedNamesMap[name]
			if exists {
				tx.Exec("UPDATE warga SET data_latih_2 = 1 WHERE id = ?", wargaID)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Gagal commit transaksi: %v", err)
	}

	fmt.Println("Semua data warga berhasil di-import (RAW) dan disinkronkan dengan Excel!")
}
