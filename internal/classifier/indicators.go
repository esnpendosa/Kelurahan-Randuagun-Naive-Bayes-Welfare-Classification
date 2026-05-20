package classifier

type Pilihan struct {
	Kode string
	Ket  string
}

type Indikator struct {
	ID      string
	Label   string
	Pilihan []Pilihan
	Bagian  string
}

func AmbilDaftarIndikator() []Indikator {
	return []Indikator{
		{ID: "IM1", Label: "Status Kepemilikan Rumah (IM 1)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Milik orang lain tanpa sewa"},
			{Kode: "B", Ket: "Milik orang tua"},
			{Kode: "C", Ket: "Menyewa"},
			{Kode: "D", Ket: "Milik sendiri"},
		}},
		{ID: "IM2", Label: "Dinding Rumah Terluas (IM 2)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Berdinding bambu / kayu berkualitas"},
			{Kode: "B", Ket: "Tembok tanpa plester"},
			{Kode: "C", Ket: "Papan / kayu jati"},
			{Kode: "D", Ket: "Tembok kualitas baik"},
		}},
		{ID: "IM3", Label: "Luas Lantai Bangunan (IM 3)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "< 5 m2/jiwa"},
			{Kode: "B", Ket: "6 m2/jiwa"},
			{Kode: "C", Ket: "7 m2/jiwa"},
			{Kode: "D", Ket: "> 8 m2/jiwa"},
		}},
		{ID: "IM4", Label: "Lantai Terluas (IM 4)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tanah"},
			{Kode: "B", Ket: "Plester semen / batu bata"},
			{Kode: "C", Ket: "Tegel"},
			{Kode: "D", Ket: "Keramik"},
		}},
		{ID: "IM5", Label: "Fasilitas BAB (IM 5)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Punya"},
			{Kode: "B", Ket: "Umum"},
			{Kode: "C", Ket: "Jamban bersama"},
			{Kode: "D", Ket: "Milik sendiri"},
		}},
		{ID: "IM6", Label: "Sumber Air Minum (IM 6)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Sungai / air hujan"},
			{Kode: "B", Ket: "Sumur / mata air"},
			{Kode: "C", Ket: "Ledeng eceran"},
			{Kode: "D", Ket: "PDAM / membeli air kemasan"},
		}},
		{ID: "IM7", Label: "Sumber Penerangan (IM 7)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Petromak / ublik"},
			{Kode: "B", Ket: "Listrik numpang"},
			{Kode: "C", Ket: "PLN 450 Watt"},
			{Kode: "D", Ket: "PLN 900 Watt / Lebih"},
		}},
		{ID: "IM8", Label: "Bahan Bakar Utama (IM 8)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Kayu Bakar"},
			{Kode: "B", Ket: "Arang"},
			{Kode: "C", Ket: "Gas LPG 3 Kg"},
			{Kode: "D", Ket: "Gas LPG > 3 Kg / Blue gas"},
		}},
		{ID: "IM9", Label: "Pendidikan Kepala RT (IM 9)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak sekolah / tidak tamas SD"},
			{Kode: "B", Ket: "Tamat SD / SMP"},
			{Kode: "C", Ket: "Tamat SMA"},
			{Kode: "D", Ket: "Tamat Perguruan Tinggi"},
		}},
		{ID: "IM10", Label: "Pekerjaan Utama Kepala RT (IM 10)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak punya pekerjaan"},
			{Kode: "B", Ket: "Pekerjaan bebas"},
			{Kode: "C", Ket: "Buruh / karyawan"},
			{Kode: "D", Ket: "wirausaha / pedagang besar"},
		}},
		{ID: "IM11", Label: "Penghasilan Kepala RT (1 Bulan) (IM 11)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "< Rp. 608,828,-"},
			{Kode: "B", Ket: "> Rp. 608,828,- s.d Rp. 1,500,000,-"},
			{Kode: "C", Ket: "> Rp. 1,500,000,- s.d Rp. 2,500,000,-"},
			{Kode: "D", Ket: "> Rp 2,500,000,-"},
		}},
		{ID: "IM12", Label: "Jumlah anggota RT (IM 12)", Bagian: "Kondisi Rumah", Pilihan: []Pilihan{
			{Kode: "A", Ket: "> 6 orang / 1 orang lanjut usia sebatang kara"},
			{Kode: "B", Ket: "5 orang"},
			{Kode: "C", Ket: "4 orang"},
			{Kode: "D", Ket: "1-3 orang"},
		}},
		{ID: "IM13", Label: "Jumlah KK dalam RT (IM 13)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "> 3 KK"},
			{Kode: "B", Ket: "3 KK"},
			{Kode: "C", Ket: "2 KK"},
			{Kode: "D", Ket: "1 KK"},
		}},
		{ID: "IM14", Label: "Jumlah Anggota RT yang Bekerja (IM 14)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "0 orang"},
			{Kode: "B", Ket: "1 orang"},
			{Kode: "C", Ket: "2 orang"},
			{Kode: "D", Ket: "> 3 orang"},
		}},
		{ID: "IM15", Label: "Jumlah Anggota Keluarga Masih Sekolah (IM 15)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "> 3 orang"},
			{Kode: "B", Ket: "2-3 orang"},
			{Kode: "C", Ket: "1 orang"},
			{Kode: "D", Ket: "0 orang"},
		}},
		{ID: "IM16", Label: "Anggota Keluarga Penyandang Disabilitas (IM 16)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Penyandang disabilitas multi"},
			{Kode: "B", Ket: "Penyandang disabilitas ganda"},
			{Kode: "C", Ket: "Penyandang disabilitas tunggal"},
			{Kode: "D", Ket: "Tidak ada"},
		}},
		{ID: "IM17", Label: "Anggota Keluarga yang Masuk Kategori Lanjut Usia (IM 17)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Lanjut usia terlantar"},
			{Kode: "B", Ket: "Lanjut usia tidak potensial"},
			{Kode: "C", Ket: "Lanjut usia potensial"},
			{Kode: "D", Ket: "Tidak ada"},
		}},
		{ID: "IM18", Label: "Anggota Keluarga yang Menderita Penyakit Kronis (IM 18)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Ada anggota keluarga yang menderita sakit kronis berat"},
			{Kode: "B", Ket: "Ada anggota keluarga yang menderita sakit kronis sedang"},
			{Kode: "C", Ket: "Ada anggota keluarga yang menderita sakit kronis ringan"},
			{Kode: "D", Ket: "Tidak ada"},
		}},
		{ID: "IM19", Label: "Akses Terhadap Pelayanan Kesehatan (IM 19)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak memiliki jaminan kesehatan"},
			{Kode: "B", Ket: "Memiliki jaminan kesehatan namun tidak aktif"},
			{Kode: "C", Ket: "Memiliki jaminan kesehatan namun ada tunggakan pembayaran"},
			{Kode: "D", Ket: "Memiliki jaminan kesehatan aktif"},
		}},
		{ID: "IM20", Label: "Nilai Aset yang Dimiliki yang Mudah di Jual (IM 20)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Rp. 0,- s.d Rp. 608,828,-"},
			{Kode: "B", Ket: "Rp. 608,828, s.d Rp. 1,500,000,-"},
			{Kode: "C", Ket: "Rp. 1,500,000,- s.d Rp. 5,000,000,-"},
			{Kode: "D", Ket: "Rp. 5,000,000,-"},
		}},
		{ID: "IM21", Label: "Pengeluaran Untuk Pakaian Dalam Setahun (IM 21)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak ada"},
			{Kode: "B", Ket: "1 stel untuk setiap anggota keluarga"},
			{Kode: "C", Ket: "2 stel untuk setiap anggota keluarga"},
			{Kode: "D", Ket: "> 2 stel untuk setiap anggota keluarga"},
		}},
		{ID: "IM22", Label: "Jenis Atap Terluas (IM 22)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Asbes"},
			{Kode: "B", Ket: "Seng"},
			{Kode: "C", Ket: "Genteng"},
			{Kode: "D", Ket: "Beton"},
		}},
		{ID: "IM23", Label: "Pembuangan Akhir Tinja (IM 23)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Ipal"},
			{Kode: "B", Ket: "Tangki Septik"},
		}},
		{ID: "IM24", Label: "Memiliki Bangunan Ditempat Lain (IM 24)", Bagian: "Ekonomi Keluarga", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak"},
			{Kode: "B", Ket: "Ya"},
		}},
		{ID: "IM25", Label: "Memiliki Lahan Lain (Selain yang Ditempati) (IM 25)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak"},
			{Kode: "B", Ket: "Ya"},
		}},
		{ID: "IM26", Label: "Jumlah tabung gas 5,5 kg / lebih (IM 26)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM27", Label: "Jumlah AC (IM 27)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM28", Label: "Jumlah telepon rumah (PSTN) (IM 28)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM29", Label: "Jumlah Emas/Perhiasan (IM 29)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM30", Label: "Jumlah Sepeda Motor (IM 30)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM31", Label: "Jumlah Mobil (IM 31)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM32", Label: "Jumlah lemari es/kulkas (IM 32)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM33", Label: "Jumlah televisi layar datar (IM 33)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM34", Label: "Jumlah Komputer (IM 34)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM35", Label: "Jumlah Sepeda (IM 35)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
		{ID: "IM36", Label: "Jumlah Handphone (IM 36)", Bagian: "Aset & Fasilitas", Pilihan: []Pilihan{
			{Kode: "A", Ket: "Tidak Ada = 0"},
			{Kode: "B", Ket: "Ada = 1"},
		}},
	}
}
