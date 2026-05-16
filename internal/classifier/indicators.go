package classifier // Paket classifier untuk menangani struktur data indikator

// Indikator merepresentasikan satu poin pertanyaan/kriteria penilaian
type Indikator struct {
	ID      string   // ID unik indikator (contoh: IM1, IM2)
	Label   string   // Teks pertanyaan atau label kriteria
	Pilihan []string // Daftar opsi jawaban (contoh: A, B, C, D)
	Bagian  string   // Kategori/Section indikator (contoh: Kondisi Rumah)
}

// AmbilDaftarIndikator mengembalikan 36 indikator kesejahteraan yang digunakan dalam sistem
func AmbilDaftarIndikator() []Indikator {
	return []Indikator{
		// Indikator Bagian 1: Kondisi Rumah
		{ID: "IM1", Label: "Status Kepemilikan Rumah", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM2", Label: "Dinding Rumah Terluas", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM3", Label: "Luas Lantai Bangunan", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM4", Label: "Lantai Terluas", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM5", Label: "Fasilitas BAB", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM6", Label: "Sumber Air Minum", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM7", Label: "Sumber Penerangan", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM8", Label: "Bahan Bakar Utama", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM9", Label: "Pendidikan Kepala RT", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM10", Label: "Pekerjaan Utama Kepala RT", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM11", Label: "Penghasilan Kepala RT", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},
		{ID: "IM12", Label: "Jumlah anggota RT", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Kondisi Rumah"},

		// Indikator Bagian 2: Ekonomi Keluarga
		{ID: "IM13", Label: "Jumlah KK dalam RT", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM14", Label: "Jumlah Anggota RT yang Bekerja", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM15", Label: "Jumlah Anggota Keluarga Masih Sekolah", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM16", Label: "Anggota Keluarga Penyandang Disabilitas", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM17", Label: "Anggota Keluarga yang Masuk Kategori Lanjut Usia", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM18", Label: "Anggota Keluarga yang Menderita Penyakit Kronis", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM19", Label: "Akses Terhadap Pelayanan Kesehatan", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM20", Label: "Nilai Aset yang Dimiliki yang Mudah di Jual", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM21", Label: "Pengeluaran Untuk Pakaian Dalam Setahun", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM22", Label: "Jenis Atap Terluas", Pilihan: []string{"A", "B", "C", "D"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM23", Label: "Pembuangan Akhir Tinja", Pilihan: []string{"A", "B"}, Bagian: "Ekonomi Keluarga"},
		{ID: "IM24", Label: "Memiliki Bangunan Ditempat Lain", Pilihan: []string{"A", "B"}, Bagian: "Ekonomi Keluarga"},

		// Indikator Bagian 3: Aset & Fasilitas
		{ID: "IM25", Label: "Memiliki Lahan Lain", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM26", Label: "Jumlah tabung gas 5,5 kg / lebih", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM27", Label: "Jumlah AC", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM28", Label: "Jumlah telepon rumah", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM29", Label: "Jumlah Emas/Perhiasan", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM30", Label: "Jumlah Sepeda Motor", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM31", Label: "Jumlah Mobil", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM32", Label: "Jumlah lemari es/kulkas", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM33", Label: "Jumlah televisi layar datar", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM34", Label: "Jumlah Komputer", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM35", Label: "Jumlah Sepeda", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
		{ID: "IM36", Label: "Jumlah Handphone", Pilihan: []string{"A", "B"}, Bagian: "Aset & Fasilitas"},
	}
}
