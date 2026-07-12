package classifier // Paket classifier untuk menangani logika klasifikasi



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

// LatihModel menghitung prior dan likelihood tanpa Laplace Smoothing
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
		// Mencari nilai unik untuk fitur ini (biasanya A, B, C, D)
		nilaiUnik := make(map[string]bool)
		for _, baris := range data {
			nilaiUnik[baris[fitur]] = true
		}

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

			// Tanpa Laplace Smoothing: P(Xi=v|Ck) = jumlahMuncul / totalDiKelas
			for v := range nilaiUnik {
				if totalDiKelas > 0 {
					nb.PeluangLikelihood[c][fitur][v] = float64(jumlahMuncul[v]) / float64(totalDiKelas)
				} else {
					nb.PeluangLikelihood[c][fitur][v] = 0.0
				}
			}
		}
	}
}

// Prediksi menghitung probabilitas setiap kelas untuk data input baru
func (nb *KlasifikasiNaiveBayes) Prediksi(input map[string]string) map[KelasKesejahteraan]float64 {
	hasilPeluang := make(map[KelasKesejahteraan]float64) // Map untuk menyimpan hasil akhir probabilitas

	for _, c := range nb.SemuaKelas {
		p := nb.PeluangPrior[c]
		if p == 0 {
			hasilPeluang[c] = 0
			continue
		}

		for _, fitur := range nb.DaftarFitur {
			nilai := input[fitur] // Ambil nilai fitur dari input
			
			// Pastikan map likelihood untuk kelas dan fitur ini sudah terinisialisasi
			var l float64
			if nb.PeluangLikelihood[c] != nil && nb.PeluangLikelihood[c][fitur] != nil {
				l = nb.PeluangLikelihood[c][fitur][nilai] // nilainya default 0.0 jika tidak ada
			}
			
			p *= l
			if p == 0 {
				break
			}
		}
		hasilPeluang[c] = p
	}

	return hasilPeluang // Kembalikan map probabilitas untuk setiap kelas
}

// PrediksiLaplace menghitung probabilitas dengan Laplace Smoothing untuk menghindari peluang 0
func (nb *KlasifikasiNaiveBayes) PrediksiLaplace(input map[string]string, data []map[string]string, target []KelasKesejahteraan) map[KelasKesejahteraan]float64 {
	hasilPeluang := make(map[KelasKesejahteraan]float64)
	jumlahTotal := float64(len(target))
	hitungPerKelas := make(map[KelasKesejahteraan]int)
	for _, t := range target {
		hitungPerKelas[t]++
	}

	for _, c := range nb.SemuaKelas {
		// Prior dengan Laplace Smoothing: (count(c) + 1) / (N + K)
		p := float64(hitungPerKelas[c] + 1) / (jumlahTotal + float64(len(nb.SemuaKelas)))

		for _, fitur := range nb.DaftarFitur {
			nilai := input[fitur]
			jumlahMuncul := 0
			totalDiKelas := 0
			nilaiUnik := make(map[string]bool)
			for i, baris := range data {
				nilaiUnik[baris[fitur]] = true
				if target[i] == c {
					totalDiKelas++
					if baris[fitur] == nilai {
						jumlahMuncul++
					}
				}
			}

			v := len(nilaiUnik)
			if v == 0 {
				v = 4
			}

			// Likelihood dengan Laplace Smoothing: (count(fitur=nilai|c) + 1) / (count(c) + V)
			l := float64(jumlahMuncul + 1) / float64(totalDiKelas + v)
			p *= l
		}
		hasilPeluang[c] = p
	}
	return hasilPeluang
}

// AmbilKelasTerbaik mencari kelas dengan probabilitas tertinggi secara deterministik
func (nb *KlasifikasiNaiveBayes) AmbilKelasTerbaik(peluang map[KelasKesejahteraan]float64) KelasKesejahteraan {
	var kelasTerbaik KelasKesejahteraan = 1 // default
	var peluangMaks float64 = -1.0          // Mulai dari -1.0 agar kelas dengan peluang 0 bisa terpilih jika semuanya 0
	for _, c := range nb.SemuaKelas {
		p := peluang[c]
		if p > peluangMaks {
			peluangMaks = p
			kelasTerbaik = c
		}
	}
	return kelasTerbaik // Kembalikan kelas yang paling mungkin (hasil argmax)
}
