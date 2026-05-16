package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Record struct {
	Indicators map[string]string `json:"indicators"`
	Class      int               `json:"class"`
}

func main() {
	indicators := []string{
		"status_rumah", "dinding", "luas_lantai", "lantai", "atap", "fasilitas_bab",
		"pembuangan_tinja", "sumber_air", "penerangan", "bahan_bakar", "pendidikan_kk",
		"pekerjaan_kk", "penghasilan_kk", "jumlah_art", "jumlah_kk_rumah", "art_bekerja",
		"art_sekolah", "disabilitas", "lansia", "penyakit_kronis", "akses_kesehatan",
		"nilai_aset", "pengeluaran_pakaian", "bangunan_lain", "lahan_lain", "gas_5_5kg",
		"ac", "pstn", "emas", "motor", "mobil", "kulkas", "tv_flat", "komputer", "sepeda", "hp",
	}

	options := map[string][]string{
		"status_rumah":     {"Milik Sendiri", "Kontrak/Sewa", "Bebas Sewa", "Dinas", "Lainnya"},
		"dinding":          {"Tembok", "Kayu/Papan", "Bambu/Anyaman", "Seng", "Lainnya"},
		"luas_lantai":      {"< 20 m2", "20-50 m2", "50-100 m2", "> 100 m2"},
		"lantai":           {"Keramik/Marmer", "Semen/Ubin", "Kayu/Papan", "Tanah", "Lainnya"},
		"atap":             {"Genteng", "Seng/Asbes", "Beton", "Bambu/Rumbia", "Lainnya"},
		"fasilitas_bab":    {"Sendiri", "Bersama/Umum", "Tidak Ada"},
		"pembuangan_tinja": {"Tangki Septik", "Lubang Tanah", "Kolam/Sawah", "Sungai/Laut", "Lainnya"},
		"sumber_air":       {"Air Kemasan/Isi Ulang", "Leding/PAM", "Sumur Bor/Pompa", "Sumur Terlindung", "Mata Air", "Air Hujan"},
		"penerangan":       {"PLN Berlangganan", "PLN Non-Berlangganan", "Non-PLN", "Bukan Listrik"},
		"bahan_bakar":      {"Listrik/Gas > 3kg", "Gas 3kg", "Minyak Tanah", "Kayu Bakar", "Lainnya"},
		"pendidikan_kk":    {"Tidak Sekolah", "SD/Sederajat", "SMP/Sederajat", "SMA/Sederajat", "Diploma/Sarjana"},
		"pekerjaan_kk":     {"Tidak Bekerja", "Petani", "Buruh", "Pedagang", "Karyawan Swasta", "PNS/TNI/Polri", "Lainnya"},
		"penghasilan_kk":   {"< 500rb", "500rb - 1.5jt", "1.5jt - 3jt", "3jt - 5jt", "> 5jt"},
		"jumlah_art":       {"1-2 orang", "3-4 orang", "5-6 orang", "> 6 orang"},
		"jumlah_kk_rumah":  {"1 KK", "2 KK", "> 2 KK"},
		"art_bekerja":      {"Tidak Ada", "1 orang", "2 orang", "> 2 orang"},
		"art_sekolah":      {"Tidak Ada", "1 orang", "2 orang", "> 2 orang"},
		"disabilitas":      {"Ada", "Tidak Ada"},
		"lansia":           {"Ada", "Tidak Ada"},
		"penyakit_kronis":  {"Ada", "Tidak Ada"},
		"akses_kesehatan":  {"Sangat Mudah", "Mudah", "Sulit", "Sangat Sulit"},
		"nilai_aset":       {"< 500rb", "500rb - 5jt", "5jt - 20jt", "> 20jt"},
		"pengeluaran_pakaian": {"< 500rb", "500rb - 1.5jt", "> 1.5jt"},
		"bangunan_lain":    {"Ya", "Tidak"},
		"lahan_lain":       {"Ya", "Tidak"},
		"gas_5_5kg":        {"Ya", "Tidak"},
		"ac":               {"Ya", "Tidak"},
		"pstn":             {"Ya", "Tidak"},
		"emas":             {"Ya", "Tidak"},
		"motor":            {"Ada", "Tidak Ada"},
		"mobil":            {"Ada", "Tidak Ada"},
		"kulkas":           {"Ada", "Tidak Ada"},
		"tv_flat":          {"Ada", "Tidak Ada"},
		"komputer":         {"Ada", "Tidak Ada"},
		"sepeda":           {"Ada", "Tidak Ada"},
		"hp":               {"Ada", "Tidak Ada"},
	}

	var data []Record
	for i := 0; i < 200; i++ {
		class := rand.Intn(6) + 1
		rec := Record{
			Indicators: make(map[string]string),
			Class:      class,
		}
		
		for _, ind := range indicators {
			opts := options[ind]
			// Biased selection based on class for better "dummy" training
			idx := rand.Intn(len(opts))
			if class <= 2 { // Poor
				if ind == "penghasilan_kk" { idx = 0 }
				if ind == "pendidikan_kk" { idx = rand.Intn(2) }
			} else if class >= 5 { // Rich
				if ind == "penghasilan_kk" { idx = len(opts)-1 }
				if ind == "pendidikan_kk" { idx = len(opts)-1 }
			}
			rec.Indicators[ind] = opts[idx]
		}
		data = append(data, rec)
	}

	file, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile("dummy_data.json", file, 0644)
	fmt.Println("Generated 200 dummy records in dummy_data.json")
}
