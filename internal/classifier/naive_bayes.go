package classifier // Paket classifier untuk menangani logika klasifikasi

import (
	"math" // Mengimpor paket math untuk operasi matematika seperti logaritma
)

// KelasKesejahteraan merepresentasikan 6 tingkat kesejahteraan warga
type KelasKesejahteraan int

// Konstanta untuk setiap tingkat kesejahteraan (1 sampai 6)
const (
	SangatMiskin KelasKesejahteraan = iota + 1 // Kelas 1: Sangat Miskin
	Miskin                                     // Kelas 2: Miskin
	HampirMiskin                               // Kelas 3: Hampir Miskin
	RentanMiskin                               // Kelas 4: Rentan Miskin
	PasPasan                                   // Kelas 5: Pas-pasan
	MenengahKeAtas                             // Kelas 6: Menengah ke Atas
)

// DaftarNamaKelas memetakan ID kelas ke nama string dalam Bahasa Indonesia
var DaftarNamaKelas = map[KelasKesejahteraan]string{
	SangatMiskin:   "Sangat Miskin",   // Pemetaan Kelas 1
	Miskin:         "Miskin",          // Pemetaan Kelas 2
	HampirMiskin:   "Hampir Miskin",   // Pemetaan Kelas 3
	RentanMiskin:   "Rentan Miskin",   // Pemetaan Kelas 4
	PasPasan:       "Pas-pasan",       // Pemetaan Kelas 5
	MenengahKeAtas: "Menengah ke Atas", // Pemetaan Kelas 6
}

// KlasifikasiNaiveBayes adalah struktur data untuk model Naive Bayes
type KlasifikasiNaiveBayes struct {
	SemuaKelas      []KelasKesejahteraan                               // Daftar semua kelas yang tersedia
	PeluangPrior    map[KelasKesejahteraan]float64                     // P(C): Peluang awal setiap kelas
	PeluangLikelihood map[KelasKesejahteraan]map[string]map[string]float64 // P(X|C): Peluang fitur muncul di kelas tertentu
	DaftarFitur     []string                                           // Nama-nama fitur (IM1 - IM36) yang digunakan
}

// BuatModelBaru menginisialisasi model KlasifikasiNaiveBayes yang baru
func BuatModelBaru() *KlasifikasiNaiveBayes {
	return &KlasifikasiNaiveBayes{
		// Mengisi daftar semua kelas dari 1 sampai 6
		SemuaKelas: []KelasKesejahteraan{SangatMiskin, Miskin, HampirMiskin, RentanMiskin, PasPasan, MenengahKeAtas},
		// Inisialisasi map untuk menyimpan nilai peluang
		PeluangPrior:  make(map[KelasKesejahteraan]float64),
		PeluangLikelihood: make(map[KelasKesejahteraan]map[string]map[string]float64),
	}
}

// LatihModel menghitung prior dan likelihood dengan teknik Laplace Smoothing
func (nb *KlasifikasiNaiveBayes) LatihModel(data []map[string]string, target []KelasKesejahteraan) {
	jumlahTotal := float64(len(target)) // Menghitung total jumlah data training
	hitungPerKelas := make(map[KelasKesejahteraan]int) // Map untuk menghitung jumlah data per kelas
	
	// Menghitung frekuensi kemunculan setiap kelas dalam data training
	for _, t := range target {
		hitungPerKelas[t]++
	}

	// Menghitung Peluang Prior P(C) untuk setiap kelas
	for _, c := range nb.SemuaKelas {
		// Prior = jumlah data di kelas C / jumlah total data
		nb.PeluangPrior[c] = float64(hitungPerKelas[c]) / jumlahTotal
		// Inisialisasi map likelihood untuk kelas tersebut
		nb.PeluangLikelihood[c] = make(map[string]map[string]float64)
	}

	// Menghitung Peluang Likelihood P(X|C) untuk setiap fitur
	for _, fitur := range nb.DaftarFitur {
		// Mencari nilai unik untuk fitur ini (biasanya A, B, C, D) untuk keperluan Laplace Smoothing
		nilaiUnik := make(map[string]bool)
		for _, baris := range data {
			nilaiUnik[baris[fitur]] = true
		}
		jumlahV := float64(len(nilaiUnik)) // |V|: Jumlah kategori nilai dalam fitur

		// Hitung likelihood untuk setiap kelas
		for _, c := range nb.SemuaKelas {
			nb.PeluangLikelihood[c][fitur] = make(map[string]float64)
			
			// Menghitung berapa kali nilai fitur tertentu muncul di kelas C
			jumlahMuncul := make(map[string]int)
			totalDiKelas := 0
			for i, baris := range data {
				if target[i] == c {
					jumlahMuncul[baris[fitur]]++
					totalDiKelas++
				}
			}

			// Menerapkan Laplace Smoothing: P(Xi=v|Ck) = (jumlah+1) / (total_di_kelas + |V|)
			for v := range nilaiUnik {
				nb.PeluangLikelihood[c][fitur][v] = (float64(jumlahMuncul[v]) + 1) / (float64(totalDiKelas) + jumlahV)
			}
		}
	}
}

// Prediksi menghitung probabilitas setiap kelas untuk data input baru
func (nb *KlasifikasiNaiveBayes) Prediksi(input map[string]string) map[KelasKesejahteraan]float64 {
	hasilPeluang := make(map[KelasKesejahteraan]float64) // Map untuk menyimpan hasil akhir probabilitas

	for _, c := range nb.SemuaKelas {
		// Menggunakan logaritma untuk menghindari "underflow" (angka yang terlalu kecil hingga menjadi nol)
		p := math.Log(nb.PeluangPrior[c]) // Mulai dengan log P(C)
		if nb.PeluangPrior[c] == 0 {
			p = -1e10 // Gunakan angka sangat kecil jika prior 0
		}

		for _, fitur := range nb.DaftarFitur {
			nilai := input[fitur] // Ambil nilai fitur dari input
			
			// Pastikan map likelihood untuk kelas dan fitur ini sudah terinisialisasi
			if nb.PeluangLikelihood[c] != nil && nb.PeluangLikelihood[c][fitur] != nil {
				if l, ada := nb.PeluangLikelihood[c][fitur][nilai]; ada && l > 0 {
					p += math.Log(l)
				} else {
					p += math.Log(1e-6) // Penanganan nilai yang tidak ada di data training
				}
			} else {
				p += math.Log(1e-6)
			}
		}
		// Kembalikan nilai log ke bentuk eksponensial (probabilitas asli)
		hasilPeluang[c] = math.Exp(p)
	}

	return hasilPeluang // Kembalikan map probabilitas untuk setiap kelas
}

// AmbilKelasTerbaik mencari kelas dengan probabilitas tertinggi
func (nb *KlasifikasiNaiveBayes) AmbilKelasTerbaik(peluang map[KelasKesejahteraan]float64) KelasKesejahteraan {
	var kelasTerbaik KelasKesejahteraan // Variabel penampung kelas dengan skor tertinggi
	var peluangMaks float64 // Variabel penampung nilai probabilitas tertinggi
	for c, p := range peluang {
		if p > peluangMaks {
			peluangMaks = p
			kelasTerbaik = c
		}
	}
	return kelasTerbaik // Kembalikan kelas yang paling mungkin (hasil argmax)
}
