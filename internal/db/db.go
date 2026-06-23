package db // Paket db untuk mengelola interaksi dengan database SQLite

import (
	"database/sql"           // Paket standar untuk SQL
	"fmt"                    // Paket untuk format input/output
	"strings"                // Paket untuk manipulasi string
	_ "modernc.org/sqlite"  // Driver SQLite yang murni ditulis dalam Go
	"golang.org/x/crypto/bcrypt" // Paket untuk enkripsi kata sandi
)

// InisialisasiDB membuat koneksi database dan menyiapkan tabel-tabel yang diperlukan
func InisialisasiDB(path string) (*sql.DB, error) {
	// Membuka koneksi ke file database SQLite
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err // Mengembalikan error jika gagal membuka file
	}

	// Daftar query untuk membuat tabel-tabel utama dalam Bahasa Indonesia
	queries := []string{
		// Tabel pengguna untuk menyimpan data akun login
		`CREATE TABLE IF NOT EXISTS pengguna (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nama_pengguna TEXT UNIQUE,
			kata_sandi TEXT,
			peran TEXT
		)`,
		// Tabel warga untuk menyimpan data kependudukan dan label kesejahteraan
		`CREATE TABLE IF NOT EXISTS warga (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			idpengguna INTEGER,
			nik TEXT UNIQUE DEFAULT '',
			no_kk TEXT DEFAULT '',
			nama_lengkap TEXT DEFAULT '',
			alamat TEXT DEFAULT '',
			rt TEXT DEFAULT '',
			rw TEXT DEFAULT '',
			kelurahan TEXT DEFAULT '',
			data_latih INTEGER DEFAULT 0,
			label_kelas TEXT DEFAULT '',
			dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(idpengguna) REFERENCES pengguna(id)
		)`,
		// Tabel data_indikator untuk menyimpan 36 nilai IM/IL per warga
		`CREATE TABLE IF NOT EXISTS data_indikator (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			warga_id INTEGER,
			indikator_id TEXT,
			nilai TEXT, -- Nilai A, B, C, atau D
			FOREIGN KEY(warga_id) REFERENCES warga(id)
		)`,
		// Tabel hasil_klasifikasi untuk menyimpan riwayat prediksi model
		`CREATE TABLE IF NOT EXISTS hasil_klasifikasi (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			warga_id INTEGER,
			nama_kelas TEXT,
			probabilitas TEXT, -- Data JSON berisi peluang 6 kelas
			dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(warga_id) REFERENCES warga(id)
		)`,
	}

	// Menjalankan setiap query pembuatan tabel
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, err // Mengembalikan error jika eksekusi query gagal
		}
	}

	// Cek apakah kolom idpengguna sudah ada di tabel warga (untuk database lama)
	var kolomAda bool
	rows, err := db.Query("PRAGMA table_info(warga)")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dflt_value interface{}
			if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err == nil {
				if name == "idpengguna" {
					kolomAda = true
					break
				}
			}
		}
	}
	if !kolomAda {
		// Tambahkan kolom idpengguna jika belum ada
		db.Exec("ALTER TABLE warga ADD COLUMN idpengguna INTEGER REFERENCES pengguna(id)")
	}

	return db, nil // Mengembalikan instance database yang siap digunakan
}

// DataLatih adalah struktur data untuk menampung satu rekaman data training
type DataLatih struct {
	ID        int
	Nama      string
	Kelas     int
	Indikator map[string]string
}

// AmbilDataLatih mengambil semua data dari tabel warga yang ditandai sebagai data_latih = 1
func AmbilDataLatih(db *sql.DB) ([]DataLatih, error) {
	// Mengambil data warga yang sudah memiliki label kelas
	rows, err := db.Query("SELECT id, nama_lengkap, label_kelas FROM warga WHERE data_latih = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk mengubah nama label menjadi ID angka (1-6)
	labelMap := map[string]int{
		"Sangat Miskin": 1, "Miskin": 2, "Hampir Miskin": 3,
		"Rentan Miskin": 4, "Pas-pasan": 5, "Pas-Pasan": 5, "Menengah ke Atas": 6,
	}

	var daftarData []DataLatih
	for rows.Next() {
		var d DataLatih
		var label string
		if err := rows.Scan(&d.ID, &d.Nama, &label); err != nil {
			return nil, err
		}
		
		d.Kelas = labelMap[label] // Konversi label string ke ID kelas
		if d.Kelas == 0 {
			fmt.Sscanf(label, "%d", &d.Kelas) // Jika label sudah berupa angka, parse langsung
		}

		d.Indikator = make(map[string]string)
		daftarData = append(daftarData, d)
	}

	// Mengambil 36 indikator untuk setiap warga yang ditemukan
	for i := range daftarData {
		irows, err := db.Query("SELECT indikator_id, nilai FROM data_indikator WHERE warga_id = ?", daftarData[i].ID)
		if err != nil {
			return nil, err
		}
		for irows.Next() {
			var id, val string
			if err := irows.Scan(&id, &val); err != nil {
				irows.Close()
				return nil, err
			}
			daftarData[i].Indikator[strings.ToUpper(id)] = val
		}
		irows.Close()
	}

	return daftarData, nil
}

// AmbilDataUji mengambil semua data dari tabel warga yang ditandai sebagai data_latih = 0 (Data Uji) dan memiliki label_kelas untuk keperluan evaluasi
func AmbilDataUji(db *sql.DB) ([]DataLatih, error) {
	// Mengambil data warga yang merupakan data uji dan sudah memiliki label kelas
	rows, err := db.Query("SELECT id, nama_lengkap, label_kelas FROM warga WHERE data_latih = 0 AND label_kelas != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk mengubah nama label menjadi ID angka (1-6)
	labelMap := map[string]int{
		"Sangat Miskin": 1, "Miskin": 2, "Hampir Miskin": 3,
		"Rentan Miskin": 4, "Pas-pasan": 5, "Pas-Pasan": 5, "Menengah ke Atas": 6,
	}

	var daftarData []DataLatih
	for rows.Next() {
		var d DataLatih
		var label string
		if err := rows.Scan(&d.ID, &d.Nama, &label); err != nil {
			return nil, err
		}
		
		d.Kelas = labelMap[label] // Konversi label string ke ID kelas
		if d.Kelas == 0 {
			fmt.Sscanf(label, "%d", &d.Kelas) // Jika label sudah berupa angka, parse langsung
		}

		d.Indikator = make(map[string]string)
		daftarData = append(daftarData, d)
	}

	// Mengambil 36 indikator untuk setiap warga yang ditemukan
	for i := range daftarData {
		irows, err := db.Query("SELECT indikator_id, nilai FROM data_indikator WHERE warga_id = ?", daftarData[i].ID)
		if err != nil {
			return nil, err
		}
		for irows.Next() {
			var id, val string
			if err := irows.Scan(&id, &val); err != nil {
				irows.Close()
				return nil, err
			}
			daftarData[i].Indikator[strings.ToUpper(id)] = val
		}
		irows.Close()
	}

	return daftarData, nil
}


// TambahWarga memasukkan data kependudukan baru ke database
func TambahWarga(db *sql.DB, nik, no_kk, nama, alamat, rt, rw, kelurahan string) (int64, error) {
	res, err := db.Exec("INSERT INTO warga (nik, no_kk, nama_lengkap, alamat, rt, rw, kelurahan) VALUES (?, ?, ?, ?, ?, ?, ?)",
		nik, no_kk, nama, alamat, rt, rw, kelurahan)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId() // Mengembalikan ID warga yang baru dibuat
}

// AmbilSemuaPengguna mengambil daftar semua akun yang ada di sistem
func AmbilSemuaPengguna(db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, nama_pengguna, peran FROM pengguna")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var daftar []map[string]interface{}
	for rows.Next() {
		var id int
		var u, r string
		rows.Scan(&id, &u, &r)
		daftar = append(daftar, map[string]interface{}{"ID": id, "Username": u, "Role": r})
	}
	return daftar, nil
}

// AmbilPenggunaBerdasarkanID mencari data pengguna berdasarkan primary key ID
func AmbilPenggunaBerdasarkanID(db *sql.DB, id interface{}) (map[string]interface{}, error) {
	var uid int
	var u, p, r string
	err := db.QueryRow("SELECT id, nama_pengguna, kata_sandi, peran FROM pengguna WHERE id = ?", id).Scan(&uid, &u, &p, &r)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"ID": uid, "Username": u, "Password": p, "Role": r}, nil
}

// TambahPengguna membuat akun pengguna baru dengan kata sandi yang di-hash
func TambahPengguna(db *sql.DB, nama_pengguna, kata_sandi, peran string) error {
	// Enkripsi kata sandi menggunakan bcrypt untuk keamanan
	hash, err := bcrypt.GenerateFromPassword([]byte(kata_sandi), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO pengguna (nama_pengguna, kata_sandi, peran) VALUES (?, ?, ?)", nama_pengguna, string(hash), peran)
	return err
}

// PerbaruiPengguna mengubah data pengguna (opsional ubah kata sandi)
func PerbaruiPengguna(db *sql.DB, id interface{}, nama_pengguna, kata_sandi, peran string) error {
	if kata_sandi != "" {
		// Jika kata sandi baru diisi, enkripsi dan simpan
		hash, err := bcrypt.GenerateFromPassword([]byte(kata_sandi), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE pengguna SET nama_pengguna = ?, kata_sandi = ?, peran = ? WHERE id = ?", nama_pengguna, string(hash), peran, id)
		return err
	}
	// Jika kata sandi kosong, hanya perbarui nama dan peran
	_, err := db.Exec("UPDATE pengguna SET nama_pengguna = ?, peran = ? WHERE id = ?", nama_pengguna, peran, id)
	return err
}

// HapusPengguna menghapus akun pengguna berdasarkan ID
func HapusPengguna(db *sql.DB, id interface{}) error {
	_, err := db.Exec("DELETE FROM pengguna WHERE id = ?", id)
	return err
}

// BenihPenggunaDefault membuat akun admin bawaan jika tabel pengguna masih kosong
func BenihPenggunaDefault(db *sql.DB) error {
	var total int
	db.QueryRow("SELECT COUNT(*) FROM pengguna").Scan(&total)
	if total == 0 {
		// Buat akun admin:admin123 sebagai akses awal
		return TambahPengguna(db, "admin", "admin123", "Admin")
	}
	return nil
}
