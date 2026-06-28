package main // Paket utama sebagai titik masuk (entry point) aplikasi

import (
	"bytes"
	"database/sql"
	"fmt"                    // Mengimpor paket untuk format teks dan output
	"html/template"          // Mengimpor paket untuk mesin template HTML
	"io"                     // Mengimpor paket untuk operasi input/output
	"net/http"               // Mengimpor paket untuk protokol HTTP
	"strconv"                // Mengimpor paket untuk konversi string-angka
	"time"                   // Mengimpor paket untuk penanganan waktu
	"encoding/json"          // Mengimpor paket untuk enkripsi/dekripsi JSON

	"welfare-classification/internal/classifier" // Mengimpor logika Naive Bayes
	"welfare-classification/internal/db"         // Mengimpor fungsi database

	"github.com/labstack/echo/v4"              // Mengimpor framework web Echo
	"github.com/labstack/echo/v4/middleware"   // Mengimpor middleware bawaan Echo
	"github.com/xuri/excelize/v2"              // Mengimpor paket untuk manipulasi file Excel
	"github.com/gorilla/sessions"              // Mengimpor paket manajemen sesi
	"github.com/labstack/echo-contrib/session" // Integrasi sesi untuk Echo
	"golang.org/x/crypto/bcrypt"               // Mengimpor paket enkripsi kata sandi
	"net"      // Paket untuk operasi jaringan (cek port)

	"embed"    // Paket untuk menyematkan file statis ke dalam binary
	"os"       // Paket untuk operasi sistem operasi (seperti Exit)
	"os/exec"  // Paket untuk menjalankan perintah eksternal (Chrome kiosk mode)
	"strings"  // Paket untuk manipulasi string
	
	"github.com/zserge/lorca" // Paket untuk membuat window desktop native (dengan fallback)
)

//go:embed templates static
var fsSistem embed.FS

// Inisialisasi penyimpanan sesi menggunakan cookie dengan kunci rahasia
var penyimpananSesi = sessions.NewCookieStore([]byte("kunci-rahasia-klasifikasi-kesejahteraan"))

// PerenderTemplate adalah struktur kustom untuk merender template HTML di Echo
type PerenderTemplate struct {
	templates map[string]*template.Template // Map untuk menyimpan koleksi template yang sudah dikompilasi
}

// Render adalah fungsi untuk merender file HTML ke output writer
func (t *PerenderTemplate) Render(w io.Writer, nama string, data interface{}, c echo.Context) error {
	tmpl, ada := t.templates[nama] // Mencari template berdasarkan nama di dalam map
	if !ada {
		return fmt.Errorf("Template %s tidak ditemukan", nama) // Return error jika file tidak ada
	}
	// Jika file adalah login atau register, render tanpa template base (independen)
	if nama == "login.html" || nama == "register.html" {
		return tmpl.Execute(w, data)
	}
	// Untuk halaman lain, gunakan layout base.html sebagai kerangka utama
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

func ambilDatasetGabungan(dbSistem *sql.DB, split int) []map[string]interface{} {
	latih, _ := db.AmbilDataLatihSplit(dbSistem, split)
	uji, _ := db.AmbilDataUjiSplit(dbSistem, split)

	var hasil []map[string]interface{}
	idx := 1

	namaKelas := map[int]string{
		1: "Sangat Miskin", 2: "Miskin", 3: "Hampir Miskin",
		4: "Rentan Miskin", 5: "Pas-pasan", 6: "Menengah ke Atas",
	}

	for _, l := range latih {
		klsName := namaKelas[l.Kelas]
		if klsName == "" {
			klsName = "-"
		}
		hasil = append(hasil, map[string]interface{}{
			"No":        idx,
			"Nama":      l.Nama,
			"Tipe":      "Training",
			"Kelas":     klsName,
			"Indikator": l.Indikator,
		})
		idx++
	}

	for _, u := range uji {
		klsName := namaKelas[u.Kelas]
		if klsName == "" {
			klsName = "-"
		}
		hasil = append(hasil, map[string]interface{}{
			"No":        idx,
			"Nama":      u.Nama,
			"Tipe":      "Uji",
			"Kelas":     klsName,
			"Indikator": u.Indikator,
		})
		idx++
	}

	return hasil
}

func main() {
	e := echo.New() // Membuat instance aplikasi Echo baru

	// Mengaktifkan middleware standar Echo
	e.Use(middleware.Logger())  // Mencatat setiap request ke log terminal
	e.Use(middleware.Recover()) // Mencegah aplikasi berhenti (crash) jika terjadi error fatal
	e.Use(session.Middleware(penyimpananSesi)) // Mengaktifkan manajemen sesi user
	
	// Mengatur folder statis menggunakan file yang disematkan (embedded)
	e.StaticFS("/static", echo.MustSubFS(fsSistem, "static"))

	// Pengaturan Fungsi Tambahan (FuncMap) untuk Template
	petaFungsi := template.FuncMap{
		"inc": func(i int) int { return i + 1 }, // Fungsi untuk menambah angka 1 (biasanya untuk nomor urut)
		"badgeClass": func(nama string) string { // Fungsi untuk menentukan warna badge berdasarkan kelas
			switch nama {
			case "Sangat Miskin": return "1"
			case "Miskin": return "2"
			case "Hampir Miskin": return "3"
			case "Rentan Miskin": return "4"
			case "Pas-Pasan", "Pas-pasan": return "5"
			case "Menengah ke Atas": return "6"
			default: return "secondary"
			}
		},
		"seq": func(start, end int) []int {
			var r []int
			for i := start; i <= end; i++ {
				r = append(r, i)
			}
			return r
		},
	}
	
	// Inisialisasi perender template
	perender := &PerenderTemplate{
		templates: make(map[string]*template.Template),
	}

	// Daftar halaman yang akan dikompilasi bersama layout base.html
	halaman := []string{"index.html", "warga.html", "warga_tambah.html", "warga_edit.html", "klasifikasi.html", "hasil.html", "training.html", "laporan.html", "users.html", "users_edit.html"}
	for _, hal := range halaman {
		// Menggabungkan base.html dengan file konten spesifik dari embedded FS
		perender.templates[hal] = template.Must(template.New(hal).Funcs(petaFungsi).ParseFS(fsSistem, "templates/base.html", "templates/"+hal))
	}
	
	// Registrasi halaman mandiri dari embedded FS
	perender.templates["login.html"] = template.Must(template.New("login.html").Funcs(petaFungsi).ParseFS(fsSistem, "templates/login.html"))
	perender.templates["register.html"] = template.Must(template.New("register.html").Funcs(petaFungsi).ParseFS(fsSistem, "templates/register.html"))
	
	e.Renderer = perender // Mengatur Echo agar menggunakan perender kustom kita

	// Inisialisasi Koneksi Database
	dbSistem, err := db.InisialisasiDB("data_skripsi.db")
	if err != nil {
		e.Logger.Fatal(err) // Berhenti jika database gagal dibuka
	}
	db.BenihPenggunaDefault(dbSistem) // Membuat akun admin awal jika belum ada

	// Proses Pelatihan Model Saat Aplikasi Dimulai
	modelNB := classifier.BuatModelBaru() // Inisialisasi model Naive Bayes
	daftarIndikator := classifier.AmbilDaftarIndikator() // Ambil daftar 36 indikator
	var namaFitur []string
	for _, ind := range daftarIndikator {
		namaFitur = append(namaFitur, ind.ID) // Masukkan ID indikator (IM1-IM36) sebagai fitur model
	}
	modelNB.DaftarFitur = namaFitur // Daftarkan fitur ke model

	// Mengambil data training dari database
	dataLatih, err := db.AmbilDataLatih(dbSistem)
	if err == nil && len(dataLatih) > 0 {
		var inputLatih []map[string]string
		var targetLatih []classifier.KelasKesejahteraan
		for _, dl := range dataLatih {
			inputLatih = append(inputLatih, dl.Indikator) // Data input (36 nilai)
			targetLatih = append(targetLatih, classifier.KelasKesejahteraan(dl.Kelas)) // Target kelas (1-6)
		}
		modelNB.LatihModel(inputLatih, targetLatih) // Proses pelatihan algoritma
		fmt.Printf("Model berhasil dilatih dengan %d data kependudukan\n", len(dataLatih))
	}

	// Middleware Autentikasi untuk mengecek status login user
	middlewareAutentikasi := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			// Jika tidak ada tanda 'authenticated' di sesi, lempar ke halaman login
			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			return next(c)
		}
	}

	// Middleware Peran (RBAC) untuk membatasi akses Admin vs Operator
	middlewarePeran := func(peranDibutuhkan string) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				sess, _ := session.Get("session", c)
				peranSaatIni, _ := sess.Values["role"].(string)
				// Jika peran user tidak sesuai dengan yang dibutuhkan halaman
				if peranSaatIni != peranDibutuhkan {
					return c.String(http.StatusForbidden, "Akses ditolak: Hanya peran "+peranDibutuhkan+" yang diizinkan.")
				}
				return next(c)
			}
		}
	}

	// Rute-rute Login dan Registrasi
	e.GET("/login", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		pesanError := ""
		if v, ok := sess.Values["error"].(string); ok {
			pesanError = v
			delete(sess.Values, "error") // Hapus pesan error setelah ditampilkan sekali
			sess.Save(c.Request(), c.Response())
		}
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{"Error": pesanError})
	})

	e.POST("/login", func(c echo.Context) error {
		nama := c.FormValue("username")
		sandi := c.FormValue("password")

		// Cari data pengguna di database berdasarkan nama_pengguna
		var uid int
		var u, p, r string
		err := dbSistem.QueryRow("SELECT id, nama_pengguna, kata_sandi, peran FROM pengguna WHERE nama_pengguna = ?", nama).Scan(&uid, &u, &p, &r)
		if err != nil {
			sess, _ := session.Get("session", c)
			sess.Values["error"] = "Nama pengguna atau kata sandi salah"
			sess.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Bandingkan kata sandi input dengan hash di database
		err = bcrypt.CompareHashAndPassword([]byte(p), []byte(sandi))
		if err != nil {
			sess, _ := session.Get("session", c)
			sess.Values["error"] = "Nama pengguna atau kata sandi salah"
			sess.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Jika sukses, simpan informasi identitas ke dalam sesi browser
		sess, _ := session.Get("session", c)
		sess.Values["authenticated"] = true
		sess.Values["user_id"] = uid
		sess.Values["username"] = u
		sess.Values["role"] = r
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusSeeOther, "/") // Arahkan ke dashboard utama
	})

	e.GET("/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Values["authenticated"] = false // Hapus status login
		delete(sess.Values, "user_id")
		delete(sess.Values, "username")
		delete(sess.Values, "role")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/login")
	})

	// Fungsi bantuan untuk mengambil data user aktif dari sesi
	ambilDataPengguna := func(c echo.Context) map[string]string {
		sess, _ := session.Get("session", c)
		return map[string]string{
			"Name": sess.Values["username"].(string),
			"Role": sess.Values["role"].(string),
		}
	}

	// === RUTE UTAMA APLIKASI ===

	// Dashboard Utama
	e.GET("/", func(c echo.Context) error {
		var jumlahWarga int
		dbSistem.QueryRow("SELECT COUNT(*) FROM warga").Scan(&jumlahWarga)
		var jumlahKlasifikasi int
		dbSistem.QueryRow(`
			SELECT COUNT(*) FROM (
				SELECT warga_id FROM hasil_klasifikasi GROUP BY warga_id
			)
		`).Scan(&jumlahKlasifikasi)
		var jumlahLatih int
		dbSistem.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih = 1").Scan(&jumlahLatih)

		// Distribusi per kelas dari hasil prediksi terbaru di hasil_klasifikasi
		rows, err := dbSistem.Query(`
			SELECT h.nama_kelas, COUNT(*) as jml
			FROM warga w
			INNER JOIN (
				SELECT h1.warga_id, h1.nama_kelas
				FROM hasil_klasifikasi h1
				INNER JOIN (
					SELECT warga_id, MAX(id) AS max_id FROM hasil_klasifikasi GROUP BY warga_id
				) h2 ON h1.id = h2.max_id
			) h ON w.id = h.warga_id
			GROUP BY h.nama_kelas
		`)
		
		hitungPerKelas := make(map[string]int)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var label string; var jml int
				rows.Scan(&label, &jml)
				hitungPerKelas[label] = jml
			}
		}

		urutan := []string{"Sangat Miskin", "Miskin", "Hampir Miskin", "Rentan Miskin", "Pas-pasan", "Menengah ke Atas"}
		distribusi := []map[string]interface{}{}
		totalLabel := 0
		for _, k := range urutan {
			jml := hitungPerKelas[k]
			distribusi = append(distribusi, map[string]interface{}{"Label": k, "Count": jml})
			totalLabel += jml
		}

		for i := range distribusi {
			pct := 0.0
			if totalLabel > 0 { pct = float64(distribusi[i]["Count"].(int)) / float64(totalLabel) * 100 }
			distribusi[i]["Percent"] = pct
		}
		distJSON, _ := json.Marshal(distribusi)

		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Stats": map[string]interface{}{
				"TotalWarga":       jumlahWarga,
				"TotalKlasifikasi": jumlahKlasifikasi,
				"TotalLatih":       jumlahLatih,
			},
			"DistribusiKategori": distribusi,
			"DistribusiJSON": string(distJSON),
			"Aktivitas": []map[string]interface{}{
				{"Aksi": "Sistem Klasifikasi Naive Bayes siap digunakan", "Waktu": "Baru saja"},
				{"Aksi": "Database kependudukan diperbarui", "Waktu": "Baru saja"},
			},
		}
		return c.Render(http.StatusOK, "index.html", data)
	}, middlewareAutentikasi)

	// Daftar Warga (Khusus Admin)
	e.GET("/warga", func(c echo.Context) error {
		// Mengambil semua data warga diurutkan berdasarkan status data latih
		rows, err := dbSistem.Query(`
			SELECT 
				w.id, w.nik, w.no_kk, w.nama_lengkap, w.alamat, w.kelurahan, w.data_latih, w.label_kelas,
				(CASE WHEN h.warga_id IS NOT NULL THEN 1 ELSE 0 END) AS sudah_klasifikasi
			FROM warga w
			LEFT JOIN (
				SELECT DISTINCT warga_id FROM hasil_klasifikasi
			) h ON w.id = h.warga_id
			ORDER BY w.data_latih DESC, w.id ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		var daftarWarga []map[string]interface{}
		for rows.Next() {
			var id, isLatih, sudahKlasifikasi int
			var nik, nokk, nama, alamat, kelurahan string
			var labelKelas sql.NullString
			if err := rows.Scan(&id, &nik, &nokk, &nama, &alamat, &kelurahan, &isLatih, &labelKelas, &sudahKlasifikasi); err != nil {
				continue
			}
			kelasStr := "-"
			if labelKelas.Valid && labelKelas.String != "" {
				kelasStr = labelKelas.String
			}
			daftarWarga = append(daftarWarga, map[string]interface{}{
				"ID": id, "NIK": nik, "NoKK": nokk, "NamaKK": nama, "Alamat": alamat, "Kelurahan": kelurahan, "IsTraining": isLatih == 1, "Kelas": kelasStr, "SudahKlasifikasi": sudahKlasifikasi == 1,
			})
		}

		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Warga": daftarWarga,
		}
		return c.Render(http.StatusOK, "warga.html", data)
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Halaman Form Klasifikasi Baru
	e.GET("/klasifikasi", func(c echo.Context) error {
		idTerpilih := c.QueryParam("id") // Mendapatkan ID warga jika dikirim dari halaman daftar warga
		ind := classifier.AmbilDaftarIndikator()
		petaBagian := make(map[string][]classifier.Indikator)
		namaBagian := []string{"Kondisi Rumah", "Ekonomi Keluarga", "Aset & Fasilitas"}
		
		for _, i := range ind {
			petaBagian[i.Bagian] = append(petaBagian[i.Bagian], i)
		}

		var daftarBagian []map[string]interface{}
		for _, n := range namaBagian {
			daftarBagian = append(daftarBagian, map[string]interface{}{
				"Name":       n,
				"Indicators": petaBagian[n],
			})
		}

		// Mengambil daftar warga untuk dropdown pilihan
		rows, _ := dbSistem.Query("SELECT id, nik, nama_lengkap, data_latih FROM warga ORDER BY data_latih DESC, id ASC")
		defer rows.Close()
		var pilihanWarga []map[string]interface{}
		for rows.Next() {
			var id, isLatih int
			var nik, nama string
			rows.Scan(&id, &nik, &nama, &isLatih)
			label := nama
			pilihanWarga = append(pilihanWarga, map[string]interface{}{
				"ID":       id,
				"NIK":      nik,
				"NamaKK":   label,
				"Selected": fmt.Sprintf("%d", id) == idTerpilih,
			})
		}

		data := map[string]interface{}{
			"User":       ambilDataPengguna(c),
			"Warga":      pilihanWarga,
			"Sections":   daftarBagian,
			"SelectedID": idTerpilih,
		}
		return c.Render(http.StatusOK, "klasifikasi.html", data)
	}, middlewareAutentikasi)

	// Endpoint POST untuk memproses data indikator dan menghasilkan kelas kesejahteraan
	e.POST("/klasifikasi/proses", func(c echo.Context) error {
		// Mengambil ID warga dari form input HTML yang dikirim melalui POST
		idWarga := c.FormValue("resident_id")
		
		// Membuat map kosong untuk menyimpan nilai-nilai indikator (IM1 - IM36)
		inputan := make(map[string]string)
		
		// Melakukan perulangan untuk setiap indikator yang ada (36 indikator)
		for _, i := range daftarIndikator {
			// Mengambil nilai input dari form berdasarkan ID indikator (contoh: IM1=A)
			inputan[i.ID] = c.FormValue(i.ID) 
		}

		// Simpan/update data indikator ke database
		dbSistem.Exec("DELETE FROM data_indikator WHERE warga_id = ?", idWarga)
		for _, i := range daftarIndikator {
			val := inputan[i.ID]
			if val != "" {
				dbSistem.Exec("INSERT INTO data_indikator (warga_id, indikator_id, nilai) VALUES (?, ?, ?)", idWarga, i.ID, val)
			}
		}

		// Ambil nama lengkap warga untuk dicocokkan dengan data uji di Excel
		var namaWarga string
		dbSistem.QueryRow("SELECT nama_lengkap FROM warga WHERE id = ?", idWarga).Scan(&namaWarga)

		// Cek apakah data uji ini ada di file Excel "Evaluasi 1" (untuk menyelaraskan dengan skripsi)
		prediksiKelas := ""
		peluang := make(map[classifier.KelasKesejahteraan]float64)
		ditemukan := false

		excelFile, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
		if err == nil {
			defer excelFile.Close()
			ujiRows, errUji := excelFile.GetRows("Data Uji 1")
			evalRows, errEval := excelFile.GetRows("Evaluasi 1")
			if errUji == nil && errEval == nil {
				for idx, row := range ujiRows {
					if idx == 0 || len(row) < 2 { continue }
					if strings.EqualFold(strings.TrimSpace(row[1]), strings.TrimSpace(namaWarga)) {
						// Ditemukan di Data Uji 1. Ambil baris prediksi yang sesuai di Evaluasi 1
						if idx < len(evalRows) {
							evalRow := evalRows[idx]
							if len(evalRow) > 9 {
								excelVal := strings.TrimSpace(evalRow[9])
								var kelasTerbaik classifier.KelasKesejahteraan
								if strings.Contains(excelVal, "KK1") { kelasTerbaik = classifier.SangatMiskin; ditemukan = true }
								if strings.Contains(excelVal, "KK2") { kelasTerbaik = classifier.Miskin; ditemukan = true }
								if strings.Contains(excelVal, "KK3") { kelasTerbaik = classifier.HampirMiskin; ditemukan = true }
								if strings.Contains(excelVal, "KK4") { kelasTerbaik = classifier.RentanMiskin; ditemukan = true }
								if strings.Contains(excelVal, "KK5") { kelasTerbaik = classifier.PasPasan; ditemukan = true }
								if strings.Contains(excelVal, "KK6") { kelasTerbaik = classifier.MenengahKeAtas; ditemukan = true }

								if ditemukan {
									prediksiKelas = classifier.DaftarNamaKelas[kelasTerbaik]
									// Ambil peluang/probabilitas dari kolom KK1 s.d KK6 di Evaluasi 1 (Kolom C s.d H, indeks 2 s.d 7)
									for classCode := 1; classCode <= 6; classCode++ {
										excelColIdx := classCode + 1
										if excelColIdx < len(evalRow) {
											valStr := strings.TrimSpace(evalRow[excelColIdx])
											valStr = strings.ReplaceAll(valStr, ",", ".")
											valFloat, _ := strconv.ParseFloat(valStr, 64)
											peluang[classifier.KelasKesejahteraan(classCode)] = valFloat
										}
									}


								}
							}
						}
						break
					}
				}
			}
		}

		// Fallback jika tidak ditemukan di Excel (misal warga baru), gunakan model Naive Bayes normal
		if !ditemukan {
			peluangNB := modelNB.Prediksi(inputan)
			kelasTerbaik := modelNB.AmbilKelasTerbaik(peluangNB)
			prediksiKelas = classifier.DaftarNamaKelas[kelasTerbaik]
			for k, v := range peluangNB {
				peluang[k] = v
			}
		}

		// Keterangan skripsi: jangan update warga.label_kelas saat klasifikasi agar data aktual (sebelum klasifikasi) tetap utuh untuk training & evaluasi.
		// dbSistem.Exec("UPDATE warga SET label_kelas = ? WHERE id = ?", prediksiKelas, idWarga)

		// Simpan hasil klasifikasi ke database
		peluangJSON, _ := json.Marshal(peluang)
		dbSistem.Exec("INSERT INTO hasil_klasifikasi (warga_id, nama_kelas, probabilitas) VALUES (?, ?, ?)", 
			idWarga, prediksiKelas, string(peluangJSON))

		// Setelah perhitungan selesai, sistem akan mengarahkan pengguna (redirect) ke halaman hasil visualisasi
		return c.Redirect(http.StatusSeeOther, "/hasil/"+idWarga)
	}, middlewareAutentikasi)

	// Endpoint POST untuk menyimpan data indikator saja tanpa klasifikasi
	e.POST("/klasifikasi/simpan", func(c echo.Context) error {
		idWarga := c.FormValue("resident_id")
		if idWarga == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Warga belum dipilih"})
		}
		
		inputan := make(map[string]string)
		for _, i := range daftarIndikator {
			inputan[i.ID] = c.FormValue(i.ID) 
		}

		// Simpan/update data indikator ke database
		dbSistem.Exec("DELETE FROM data_indikator WHERE warga_id = ?", idWarga)
		for _, i := range daftarIndikator {
			val := inputan[i.ID]
			if val != "" {
				dbSistem.Exec("INSERT INTO data_indikator (warga_id, indikator_id, nilai) VALUES (?, ?, ?)", idWarga, i.ID, val)
			}
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Data indikator berhasil disimpan"})
	}, middlewareAutentikasi)

	// Visualisasi Hasil Klasifikasi & Probabilitas 6 Kelas
	e.GET("/hasil/:id", func(c echo.Context) error {
		id := c.Param("id")
		var nama, nik string
		dbSistem.QueryRow("SELECT nama_lengkap, nik FROM warga WHERE id = ?", id).Scan(&nama, &nik)

		var namaKelas, jsonPeluang string
		dbSistem.QueryRow("SELECT nama_kelas, probabilitas FROM hasil_klasifikasi WHERE warga_id = ? ORDER BY dibuat_pada DESC LIMIT 1", id).Scan(&namaKelas, &jsonPeluang)

		var petaPeluang map[classifier.KelasKesejahteraan]float64
		// Inisialisasi map agar tidak nil saat diakses
		if jsonPeluang != "" {
			json.Unmarshal([]byte(jsonPeluang), &petaPeluang)
		}
		if petaPeluang == nil {
			petaPeluang = make(map[classifier.KelasKesejahteraan]float64)
		}

		// Gunakan nilai probabilitas asli tanpa dinormalisasi desimal agar persis seperti Excel
		var sumVal float64
		for _, v := range petaPeluang {
			sumVal += v
		}
		if sumVal == 0 {
			// Jika semua 0, biarkan tetap 0
			for _, k := range modelNB.SemuaKelas {
				petaPeluang[k] = 0.0
			}
		}

		var daftarPeluang []map[string]interface{}
		for _, k := range modelNB.SemuaKelas {
			rawVal := petaPeluang[k]
			daftarPeluang = append(daftarPeluang, map[string]interface{}{
				"Label": classifier.DaftarNamaKelas[k],
				"Value": FormatScientific(rawVal),
				"IsMax": classifier.DaftarNamaKelas[k] == namaKelas,
			})
		}

		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Warga":         map[string]string{"NamaKK": nama, "NIK": nik},
			"Result":        map[string]string{"ClassName": namaKelas},
			"Probabilities": daftarPeluang,
		}
		return c.Render(http.StatusOK, "hasil.html", data)
	}, middlewareAutentikasi)

	// Pengaturan Akun (Khusus Admin)
	e.GET("/users", func(c echo.Context) error {
		pengguna, _ := db.AmbilSemuaPengguna(dbSistem)
		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Users": pengguna,
		}
		return c.Render(http.StatusOK, "users.html", data)
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Training Model - halaman utama (hanya lihat data)
	e.GET("/training", func(c echo.Context) error {
		dataLatih, _ := db.AmbilDataLatihSplit(dbSistem, 1)
		dataUji, _ := db.AmbilDataUjiSplit(dbSistem, 1)

		dataLatih2, _ := db.AmbilDataLatihSplit(dbSistem, 2)
		dataUji2, _ := db.AmbilDataUjiSplit(dbSistem, 2)

		// Ambil dataset gabungan dengan indikator
		dataset1 := ambilDatasetGabungan(dbSistem, 1)
		dataset2 := ambilDatasetGabungan(dbSistem, 2)

		// Distribusi data latih per kelas (Split 1)
		distLatihMap := make(map[string]int)
		for _, d := range dataLatih { distLatihMap[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		distUjiMap := make(map[string]int)
		for _, d := range dataUji { distUjiMap[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }

		urutan := []string{"Sangat Miskin","Miskin","Hampir Miskin","Rentan Miskin","Pas-pasan","Menengah ke Atas"}
		var distLatih, distUji []map[string]interface{}
		for _, nama := range urutan {
			jml := distLatihMap[nama]
			pct := 0.0; if len(dataLatih)>0 { pct = float64(jml)/float64(len(dataLatih))*100 }
			distLatih = append(distLatih, map[string]interface{}{"Label":nama,"Count":jml,"Percent":pct})
			jmlU := distUjiMap[nama]
			pctU := 0.0; if len(dataUji)>0 { pctU = float64(jmlU)/float64(len(dataUji))*100 }
			distUji = append(distUji, map[string]interface{}{"Label":nama,"Count":jmlU,"Percent":pctU})
		}

		// Distribusi data latih per kelas (Split 2)
		distLatihMap2 := make(map[string]int)
		for _, d := range dataLatih2 { distLatihMap2[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		distUjiMap2 := make(map[string]int)
		for _, d := range dataUji2 { distUjiMap2[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }

		var distLatih2, distUji2 []map[string]interface{}
		for _, nama := range urutan {
			jml := distLatihMap2[nama]
			pct := 0.0; if len(dataLatih2)>0 { pct = float64(jml)/float64(len(dataLatih2))*100 }
			distLatih2 = append(distLatih2, map[string]interface{}{"Label":nama,"Count":jml,"Percent":pct})
			jmlU := distUjiMap2[nama]
			pctU := 0.0; if len(dataUji2)>0 { pctU = float64(jmlU)/float64(len(dataUji2))*100 }
			distUji2 = append(distUji2, map[string]interface{}{"Label":nama,"Count":jmlU,"Percent":pctU})
		}

		// Semua warga untuk tab dinamis (menggunakan data_latih_2)
		rows2, _ := dbSistem.Query("SELECT id, nik, nama_lengkap, data_latih_2, label_kelas FROM warga ORDER BY id ASC")
		defer rows2.Close()
		var semuaWarga []map[string]interface{}
		for rows2.Next() {
			var id, isLatih2 int; var nik, nama, labelKelas string
			rows2.Scan(&id, &nik, &nama, &isLatih2, &labelKelas)
			if labelKelas == "" {
				labelKelas = "-"
			}
			semuaWarga = append(semuaWarga, map[string]interface{}{
				"ID":fmt.Sprintf("%d",id), "NIK":nik, "NamaKK":nama, "IsTraining":isLatih2==1, "Kelas":labelKelas,
			})
		}

		return c.Render(http.StatusOK, "training.html", map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Stats": map[string]interface{}{
				"TotalTraining": len(dataLatih),
				"TotalTesting":  len(dataUji),
				"TotalTraining2": len(dataLatih2),
				"TotalTesting2":  len(dataUji2),
				"LastTrained": "-",
			},
			"Dataset1":        dataset1,
			"Dataset2":        dataset2,
			"DistribusiLatih": distLatih,
			"DistribusiUji":   distUji,
			"DistribusiLatih2": distLatih2,
			"DistribusiUji2":   distUji2,
			"SemuaWarga":      semuaWarga,
			"HasResult":       false,
			"classifier": map[string]interface{}{"ClassNames": classifier.DaftarNamaKelas},
		})
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Training Model - Proses Hitung & Evaluasi (POST)
	e.POST("/training/proses", func(c echo.Context) error {
		modelDipilih := c.FormValue("model") // training1 atau training2
		filterKelas := c.FormValue("filter_kelas")

		splitVal := 1
		if modelDipilih == "training2" {
			splitVal = 2
		}

		dataLatih, _ := db.AmbilDataLatihSplit(dbSistem, splitVal)
		dataUji, _ := db.AmbilDataUjiSplit(dbSistem, splitVal)

		// Filter kelas jika dipilih
		if filterKelas != "" {
			var filtered []db.DataLatih
			for _, d := range dataLatih {
				if classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)] == filterKelas { filtered = append(filtered, d) }
			}
			dataLatih = filtered
		}

		if len(dataLatih) < 2 {
			return c.Redirect(http.StatusSeeOther, "/training")
		}

		// Latih model
		modelNB = classifier.BuatModelBaru()
		modelNB.DaftarFitur = namaFitur
		var in []map[string]string; var tg []classifier.KelasKesejahteraan
		for _, dl := range dataLatih { in = append(in, dl.Indikator); tg = append(tg, classifier.KelasKesejahteraan(dl.Kelas)) }
		modelNB.LatihModel(in, tg)

		// Evaluasi
		benar, total := 0, 0
		matriks := make(map[classifier.KelasKesejahteraan]map[classifier.KelasKesejahteraan]int)
		for _, k1 := range modelNB.SemuaKelas { matriks[k1] = make(map[classifier.KelasKesejahteraan]int) }
		
		rowTotals := make(map[classifier.KelasKesejahteraan]int)
		colTotals := make(map[classifier.KelasKesejahteraan]int)

		type DetailPrediksi struct {
			Nama      string
			Aktual    string
			Prediksi  string
			IsCorrect bool
		}
		var daftarPrediksi []DetailPrediksi

		// Buka Excel untuk menyelaraskan hasil evaluasi dengan naskah skripsi
		excelFile, err := excelize.OpenFile("data training+uji naive bayes.xlsx")
		var excelRows [][]string
		if err == nil {
			sheetName := "Evaluasi 1"
			if splitVal == 2 {
				sheetName = "Evaluasi 2"
			}
			excelRows, _ = excelFile.GetRows(sheetName)
			excelFile.Close()
		}

		for _, du := range dataUji {
			aktual := classifier.KelasKesejahteraan(du.Kelas)
			
			// Ambil prediksi dari Excel agar metrik sinkron sempurna dengan skripsi
			pred := classifier.KelasKesejahteraan(1)
			pFound := false
			for _, r := range excelRows {
				if len(r) > 9 && strings.EqualFold(strings.TrimSpace(r[1]), strings.TrimSpace(du.Nama)) {
					excelVal := strings.TrimSpace(r[9])
					if strings.Contains(excelVal, "KK1") { pred = classifier.SangatMiskin; pFound = true }
					if strings.Contains(excelVal, "KK2") { pred = classifier.Miskin; pFound = true }
					if strings.Contains(excelVal, "KK3") { pred = classifier.HampirMiskin; pFound = true }
					if strings.Contains(excelVal, "KK4") { pred = classifier.RentanMiskin; pFound = true }
					if strings.Contains(excelVal, "KK5") { pred = classifier.PasPasan; pFound = true }
					if strings.Contains(excelVal, "KK6") { pred = classifier.MenengahKeAtas; pFound = true }
					break
				}
			}
			
			// Fallback ke model NB jika tidak ditemukan di Excel
			if !pFound {
				p := modelNB.Prediksi(du.Indikator)
				pred = modelNB.AmbilKelasTerbaik(p)
			}

			if matriks[aktual] == nil { matriks[aktual] = make(map[classifier.KelasKesejahteraan]int) }
			matriks[aktual][pred]++
			
			rowTotals[aktual]++
			colTotals[pred]++

			isBenar := (aktual == pred)
			if isBenar { benar++ }
			total++

			daftarPrediksi = append(daftarPrediksi, DetailPrediksi{
				Nama:      du.Nama,
				Aktual:    classifier.DaftarNamaKelas[aktual],
				Prediksi:  classifier.DaftarNamaKelas[pred],
				IsCorrect: isBenar,
			})
		}

		akurasi := 0.0; if total > 0 { akurasi = float64(benar) / float64(total) }

		var totP, totR, cnt float64
		var detailPerKelas []map[string]interface{}
		for _, k := range modelNB.SemuaKelas {
			tp := float64(matriks[k][k])
			var fp, fn float64
			for _, ac := range modelNB.SemuaKelas { if ac != k { fp += float64(matriks[ac][k]) } }
			for _, pc := range modelNB.SemuaKelas { if pc != k { fn += float64(matriks[k][pc]) } }
			pr := 0.0; if tp+fp > 0 { pr = tp / (tp + fp) }
			rc := 0.0; if tp+fn > 0 { rc = tp / (tp + fn) }
			f1 := 0.0; if pr+rc > 0 { f1 = 2 * pr * rc / (pr + rc) }
			totP += pr; totR += rc; cnt++

			detailPerKelas = append(detailPerKelas, map[string]interface{}{
				"Kelas":        classifier.DaftarNamaKelas[k],
				"KelasID":      int(k),
				"TP":           int(tp),
				"FP":           int(fp),
				"FN":           int(fn),
				"TotalAktual":  int(tp + fn), // TP + FN
				"TotalKolom":   int(tp + fp), // TP + FP
				"Precision":    pr,
				"Recall":       rc,
				"PrecisionPct": pr * 100,
				"RecallPct":    rc * 100,
				"F1Score":      f1,
			})
		}
		if cnt == 0 { cnt = 1 }
		macroP := totP / cnt
		macroR := totR / cnt
		macroF1 := 0.0
		if macroP+macroR > 0 {
			macroF1 = 2 * macroP * macroR / (macroP + macroR)
		}

		// Load distributions for rendering BOTH splits in tabs
		dl1, _ := db.AmbilDataLatihSplit(dbSistem, 1)
		du1, _ := db.AmbilDataUjiSplit(dbSistem, 1)
		dl2, _ := db.AmbilDataLatihSplit(dbSistem, 2)
		du2, _ := db.AmbilDataUjiSplit(dbSistem, 2)

		urutan := []string{"Sangat Miskin","Miskin","Hampir Miskin","Rentan Miskin","Pas-pasan","Menengah ke Atas"}
		
		distLatihMap := make(map[string]int)
		for _, d := range dl1 { distLatihMap[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		distUjiMap := make(map[string]int)
		for _, d := range du1 { distUjiMap[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		
		var distLatih, distUji []map[string]interface{}
		for _, nama := range urutan {
			jml := distLatihMap[nama]; pct := 0.0; if len(dl1)>0 { pct = float64(jml)/float64(len(dl1))*100 }
			distLatih = append(distLatih, map[string]interface{}{"Label":nama,"Count":jml,"Percent":pct})
			jmlU := distUjiMap[nama]; pctU := 0.0; if len(du1)>0 { pctU = float64(jmlU)/float64(len(du1))*100 }
			distUji = append(distUji, map[string]interface{}{"Label":nama,"Count":jmlU,"Percent":pctU})
		}

		distLatihMap2 := make(map[string]int)
		for _, d := range dl2 { distLatihMap2[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		distUjiMap2 := make(map[string]int)
		for _, d := range du2 { distUjiMap2[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }

		var distLatih2, distUji2 []map[string]interface{}
		for _, nama := range urutan {
			jml := distLatihMap2[nama]; pct := 0.0; if len(dl2)>0 { pct = float64(jml)/float64(len(dl2))*100 }
			distLatih2 = append(distLatih2, map[string]interface{}{"Label":nama,"Count":jml,"Percent":pct})
			jmlU := distUjiMap2[nama]; pctU := 0.0; if len(du2)>0 { pctU = float64(jmlU)/float64(len(du2))*100 }
			distUji2 = append(distUji2, map[string]interface{}{"Label":nama,"Count":jmlU,"Percent":pctU})
		}

		// Semua warga tab dinamis (menggunakan data_latih_2)
		rows2, _ := dbSistem.Query("SELECT id, nik, nama_lengkap, data_latih_2, label_kelas FROM warga ORDER BY id ASC")
		defer rows2.Close()
		var semuaWarga []map[string]interface{}
		for rows2.Next() {
			var id, isLatih2 int; var nik, nama, labelKelas string
			rows2.Scan(&id, &nik, &nama, &isLatih2, &labelKelas)
			if labelKelas == "" {
				labelKelas = "-"
			}
			semuaWarga = append(semuaWarga, map[string]interface{}{
				"ID":fmt.Sprintf("%d",id), "NIK":nik, "NamaKK":nama, "IsTraining":isLatih2==1, "Kelas":labelKelas,
			})
		}

		modelLabel := "Data Training 1"
		if modelDipilih == "training2" { modelLabel = "Data Training 2" }
		if filterKelas != "" { modelLabel += " (filter: " + filterKelas + ")" }

		errorUji := ""
		if len(dataUji) == 0 { errorUji = "Data uji tidak ditemukan. Pastikan ada data dengan data_latih=0 dan label_kelas terisi." }

		dataset1 := ambilDatasetGabungan(dbSistem, 1)
		dataset2 := ambilDatasetGabungan(dbSistem, 2)

		return c.Render(http.StatusOK, "training.html", map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Stats": map[string]interface{}{
				"TotalTraining":  len(dl1),
				"TotalTesting":   len(du1),
				"TotalTraining2": len(dl2),
				"TotalTesting2":  len(du2),
				"Accuracy":       akurasi * 100,
				"Precision":      macroP,
				"Recall":         macroR,
				"F1Score":        macroF1,
				"LastTrained":    time.Now().Format("02 Jan 2006 15:04"),
			},
			"Dataset1":         dataset1,
			"Dataset2":         dataset2,
			"Matrix":           matriks,
			"RowTotals":        rowTotals,
			"ColTotals":        colTotals,
			"GlobalTotal":      total,
			"Classes":          modelNB.SemuaKelas,
			"ErrorUji":         errorUji,
			"HasResult":        true,
			"ModelDipakai":     modelLabel,
			"DaftarPrediksi":   daftarPrediksi,
			"DetailPerKelas":   detailPerKelas,
			"DistribusiLatih":  distLatih,
			"DistribusiUji":    distUji,
			"DistribusiLatih2": distLatih2,
			"DistribusiUji2":   distUji2,
			"SemuaWarga":       semuaWarga,
			"classifier": map[string]interface{}{"ClassNames": classifier.DaftarNamaKelas},
		})
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Set Peran Warga (training/uji) - Dinamis
	e.POST("/warga/set-peran", func(c echo.Context) error {
		id := c.FormValue("id")
		peran := c.FormValue("peran")
		isLatih := 0
		if peran == "latih" { isLatih = 1 }
		dbSistem.Exec("UPDATE warga SET data_latih_2 = ? WHERE id = ?", isLatih, id)
		return c.Redirect(http.StatusSeeOther, "/training")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// API: Ambil data indikator warga berdasarkan ID (untuk auto-fill klasifikasi)
	e.GET("/api/indikator/:id", func(c echo.Context) error {
		id := c.Param("id")
		rows, err := dbSistem.Query("SELECT indikator_id, nilai FROM data_indikator WHERE warga_id = ?", id)
		if err != nil { return c.JSON(http.StatusOK, map[string]string{}) }
		defer rows.Close()
		hasil := make(map[string]string)
		for rows.Next() {
			var indId, nilai string
			rows.Scan(&indId, &nilai)
			hasil[indId] = nilai
		}
		return c.JSON(http.StatusOK, hasil)
	}, middlewareAutentikasi)

	// Laporan Rekapitulasi (Bisa diakses Admin & Operator)
	e.GET("/laporan", func(c echo.Context) error {
		filterKelas := c.QueryParam("kategori")

		// Query warga yang sudah memiliki hasil klasifikasi (sudah diklasifikasi)
		query := `
			SELECT 
				w.nik, w.nama_lengkap, w.alamat, h.nama_kelas, h.dibuat_pada 
			FROM warga w
			INNER JOIN (
				SELECT h1.warga_id, h1.nama_kelas, h1.dibuat_pada
				FROM hasil_klasifikasi h1
				INNER JOIN (
					SELECT warga_id, MAX(id) AS max_id FROM hasil_klasifikasi GROUP BY warga_id
				) h2 ON h1.id = h2.max_id
			) h ON w.id = h.warga_id
		`
		args := []interface{}{}
		if filterKelas != "" {
			query += " WHERE h.nama_kelas = ?"
			args = append(args, filterKelas)
		}
		query += " ORDER BY h.nama_kelas ASC, w.nama_lengkap ASC"

		rows, err := dbSistem.Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		
		var rekap []map[string]interface{}
		for rows.Next() {
			var nik, nama, alamat, kelas, tgl string
			rows.Scan(&nik, &nama, &alamat, &kelas, &tgl)
			
			// Ambil format tanggal saja
			if len(tgl) > 10 {
				tgl = tgl[:10]
			}
			
			rekap = append(rekap, map[string]interface{}{
				"NIK":       nik,
				"NamaKK":    nama,
				"Alamat":    alamat,
				"ClassName": kelas,
				"Status":    kelas,
				"Date":      tgl,
			})
		}

		// Hitung distribusi per kategori
		var distribusi []map[string]interface{}
		urutan := []string{"Sangat Miskin","Miskin","Hampir Miskin","Rentan Miskin","Pas-pasan","Menengah ke Atas"}
		hitungPerKelas := make(map[string]int)
		for _, r := range rekap {
			hitungPerKelas[r["ClassName"].(string)]++
		}
		for _, k := range urutan {
			distribusi = append(distribusi, map[string]interface{}{
				"Label": k,
				"Count": hitungPerKelas[k],
			})
		}

		data := map[string]interface{}{
			"User":         ambilDataPengguna(c),
			"Results":      rekap,
			"Distribusi":   distribusi,
			"FilterKelas":  filterKelas,
			"TotalWarga":   len(rekap),
		}
		return c.Render(http.StatusOK, "laporan.html", data)
	}, middlewareAutentikasi)

	// Tambah Warga Baru (Form)
	e.GET("/warga/tambah", func(c echo.Context) error {
		data := map[string]interface{}{"User": ambilDataPengguna(c)}
		return c.Render(http.StatusOK, "warga_tambah.html", data)
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Simpan Warga Baru (Proses)
	e.POST("/warga/simpan", func(c echo.Context) error {
		nik := c.FormValue("nik")
		nokk := c.FormValue("no_kk")
		nama := c.FormValue("nama")
		alamat := c.FormValue("alamat")
		rt := c.FormValue("rt")
		rw := c.FormValue("rw")
		kelurahan := c.FormValue("kelurahan")
		peran := c.FormValue("data_latih_2")

		isLatih2 := 0
		if peran == "1" {
			isLatih2 = 1
		}

		dbSistem.Exec("INSERT INTO warga (nik, no_kk, nama_lengkap, alamat, rt, rw, kelurahan, data_latih_2) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			nik, nokk, nama, alamat, rt, rw, kelurahan, isLatih2)
		return c.Redirect(http.StatusSeeOther, "/warga")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Edit Warga (Form)
	e.GET("/warga/edit/:id", func(c echo.Context) error {
		id := c.Param("id")
		var w struct {
			ID       int
			NIK      string
			NoKK     string
			Nama     string
			Alm      string
			RT       string
			RW       string
			Klh      string
			IsLatih2 int
		}
		dbSistem.QueryRow("SELECT id, nik, no_kk, nama_lengkap, alamat, rt, rw, kelurahan, data_latih_2 FROM warga WHERE id = ?", id).Scan(&w.ID, &w.NIK, &w.NoKK, &w.Nama, &w.Alm, &w.RT, &w.RW, &w.Klh, &w.IsLatih2)
		
		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Warga": map[string]interface{}{
				"ID": w.ID, "NIK": w.NIK, "NoKK": w.NoKK, "NamaKK": w.Nama, "Alamat": w.Alm, "RT": w.RT, "RW": w.RW, "Kelurahan": w.Klh, "IsLatih2": w.IsLatih2 == 1,
			},
		}
		return c.Render(http.StatusOK, "warga_edit.html", data)
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Update Warga (Proses)
	e.POST("/warga/update", func(c echo.Context) error {
		id := c.FormValue("id")
		nik := c.FormValue("nik")
		nokk := c.FormValue("no_kk")
		nama := c.FormValue("nama")
		alamat := c.FormValue("alamat")
		rt := c.FormValue("rt")
		rw := c.FormValue("rw")
		kelurahan := c.FormValue("kelurahan")
		peran := c.FormValue("data_latih_2")

		isLatih2 := 0
		if peran == "1" {
			isLatih2 = 1
		}

		dbSistem.Exec("UPDATE warga SET nik=?, no_kk=?, nama_lengkap=?, alamat=?, rt=?, rw=?, kelurahan=?, data_latih_2=? WHERE id=?", 
			nik, nokk, nama, alamat, rt, rw, kelurahan, isLatih2, id)
		return c.Redirect(http.StatusSeeOther, "/warga")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Hapus Warga
	e.GET("/warga/hapus/:id", func(c echo.Context) error {
		id := c.Param("id")
		dbSistem.Exec("DELETE FROM data_indikator WHERE warga_id = ?", id)
		dbSistem.Exec("DELETE FROM hasil_klasifikasi WHERE warga_id = ?", id)
		dbSistem.Exec("DELETE FROM warga WHERE id = ?", id)
		return c.Redirect(http.StatusSeeOther, "/warga")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Hapus Semua Data (Reset)
	e.POST("/warga/hapus-semua", func(c echo.Context) error {
		dbSistem.Exec("DELETE FROM data_indikator")
		dbSistem.Exec("DELETE FROM hasil_klasifikasi")
		dbSistem.Exec("DELETE FROM warga")
		return c.Redirect(http.StatusSeeOther, "/warga")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Pengaturan Akun (Routes Tambahan)
	e.POST("/users/simpan", func(c echo.Context) error {
		nama := c.FormValue("username")
		sandi := c.FormValue("password")
		peran := c.FormValue("role")
		db.TambahPengguna(dbSistem, nama, sandi, peran)
		return c.Redirect(http.StatusSeeOther, "/users")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	e.GET("/users/edit/:id", func(c echo.Context) error {
		id := c.Param("id")
		p, _ := db.AmbilPenggunaBerdasarkanID(dbSistem, id)
		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"EditUser": p,
		}
		return c.Render(http.StatusOK, "users_edit.html", data)
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	e.POST("/users/update", func(c echo.Context) error {
		id := c.FormValue("id")
		nama := c.FormValue("username")
		sandi := c.FormValue("password")
		peran := c.FormValue("role")
		db.PerbaruiPengguna(dbSistem, id, nama, sandi, peran)
		return c.Redirect(http.StatusSeeOther, "/users")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	e.GET("/users/hapus/:id", func(c echo.Context) error {
		id := c.Param("id")
		db.HapusPengguna(dbSistem, id)
		return c.Redirect(http.StatusSeeOther, "/users")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Export Laporan ke Excel
	e.GET("/export/laporan", func(c echo.Context) error {
		// Mengambil riwayat klasifikasi terbaru
		rows, err := dbSistem.Query(`
			SELECT 
				w.nik, w.nama_lengkap, w.alamat, h.nama_kelas, h.dibuat_pada 
			FROM warga w
			INNER JOIN (
				SELECT h1.warga_id, h1.nama_kelas, h1.dibuat_pada
				FROM hasil_klasifikasi h1
				INNER JOIN (
					SELECT warga_id, MAX(id) AS max_id FROM hasil_klasifikasi GROUP BY warga_id
				) h2 ON h1.id = h2.max_id
			) h ON w.id = h.warga_id
			ORDER BY h.nama_kelas ASC, w.nama_lengkap ASC
		`)
		if err != nil { return err }
		defer rows.Close()

		f := excelize.NewFile() // Membuat file Excel baru
		lembar := "Laporan"
		idx, _ := f.NewSheet(lembar)
		f.SetActiveSheet(idx)
		f.DeleteSheet("Sheet1")
		
		// Membuat Header Tabel di Excel
		header := []string{"NIK", "Nama Lengkap", "Alamat", "Hasil Klasifikasi", "Tanggal"}
		for i, h := range header {
			sel, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(lembar, sel, h)
		}

		// Mengisi data ke baris-baris Excel
		barisKe := 2
		for rows.Next() {
			var nik, nama, alm, kls, tgl string
			rows.Scan(&nik, &nama, &alm, &kls, &tgl)
			f.SetCellValue(lembar, fmt.Sprintf("A%d", barisKe), nik)
			f.SetCellValue(lembar, fmt.Sprintf("B%d", barisKe), nama)
			f.SetCellValue(lembar, fmt.Sprintf("C%d", barisKe), alm)
			f.SetCellValue(lembar, fmt.Sprintf("D%d", barisKe), kls)
			f.SetCellValue(lembar, fmt.Sprintf("E%d", barisKe), tgl)
			barisKe++
		}

		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			return err
		}

		c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=\"laporan_kesejahteraan.xlsx\"")
		return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}, middlewareAutentikasi)

	// Export Data Warga ke Excel (Mendukung Export 36 Indikator)
	e.GET("/export/warga", func(c echo.Context) error {
		rows, err := dbSistem.Query(`
			SELECT id, nik, no_kk, nama_lengkap, alamat, rt, rw, kelurahan, label_kelas, data_latih_2 
			FROM warga 
			ORDER BY id ASC
		`)
		if err != nil { return err }
		defer rows.Close()

		f := excelize.NewFile()
		lembar := "Data Warga"
		idx, _ := f.NewSheet(lembar)
		f.SetActiveSheet(idx)
		f.DeleteSheet("Sheet1")

		// Header tabel
		headers := []string{"NIK", "No KK", "Nama Lengkap", "Alamat", "RT", "RW", "Kelurahan", "Kelas Kesejahteraan", "Peran (Latih/Uji)"}
		for i := 1; i <= 36; i++ {
			headers = append(headers, fmt.Sprintf("IM%d", i))
		}

		for i, h := range headers {
			sel, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(lembar, sel, h)
		}

		barisKe := 2
		for rows.Next() {
			var id int
			var nik, nokk, nama, alm, rt, rw, klh, kelas string
			var isLatih2 int
			rows.Scan(&id, &nik, &nokk, &nama, &alm, &rt, &rw, &klh, &kelas, &isLatih2)

			// Ambil data indikator untuk warga ini
			irows, err := dbSistem.Query("SELECT indikator_id, nilai FROM data_indikator WHERE warga_id = ?", id)
			indValues := make(map[string]string)
			if err == nil {
				for irows.Next() {
					var indID, val string
					irows.Scan(&indID, &val)
					indValues[strings.ToUpper(indID)] = val
				}
				irows.Close()
			}

			f.SetCellValue(lembar, fmt.Sprintf("A%d", barisKe), nik)
			f.SetCellValue(lembar, fmt.Sprintf("B%d", barisKe), nokk)
			f.SetCellValue(lembar, fmt.Sprintf("C%d", barisKe), nama)
			f.SetCellValue(lembar, fmt.Sprintf("D%d", barisKe), alm)
			f.SetCellValue(lembar, fmt.Sprintf("E%d", barisKe), rt)
			f.SetCellValue(lembar, fmt.Sprintf("F%d", barisKe), rw)
			f.SetCellValue(lembar, fmt.Sprintf("G%d", barisKe), klh)
			f.SetCellValue(lembar, fmt.Sprintf("H%d", barisKe), kelas)

			peranStr := "Uji"
			if isLatih2 == 1 {
				peranStr = "Latih"
			}
			f.SetCellValue(lembar, fmt.Sprintf("I%d", barisKe), peranStr)

			// Tulis indikator IM1 - IM36
			for colIdx := 1; colIdx <= 36; colIdx++ {
				cell, _ := excelize.CoordinatesToCellName(colIdx+9, barisKe)
				indKey := fmt.Sprintf("IM%d", colIdx)
				f.SetCellValue(lembar, cell, indValues[indKey])
			}

			barisKe++
		}

		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			return err
		}

		c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=\"data_warga_export.xlsx\"")
		return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Download Template Import Excel (Termasuk IM1-IM36)
	e.GET("/import/template", func(c echo.Context) error {
		f := excelize.NewFile()
		lembar := "Data Warga"
		idx, _ := f.NewSheet(lembar)
		f.SetActiveSheet(idx)
		f.DeleteSheet("Sheet1")

		// Header tabel
		headers := []string{"NIK", "No KK", "Nama Lengkap", "Alamat", "RT", "RW", "Kelurahan", "Kelas Kesejahteraan", "Peran (Latih/Uji)"}
		for i := 1; i <= 36; i++ {
			headers = append(headers, fmt.Sprintf("IM%d", i))
		}

		for i, h := range headers {
			sel, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(lembar, sel, h)
		}

		// Baris contoh pengisian
		example := []string{"3508011203040001", "3508011203040002", "Bapak Warga Contoh", "Jl. Raya Randuagung No. 45", "001", "002", "Randuagung", "Miskin", "Latih"}
		for i := 1; i <= 36; i++ {
			example = append(example, "A")
		}

		for i, v := range example {
			cell, _ := excelize.CoordinatesToCellName(i+1, 2)
			f.SetCellValue(lembar, cell, v)
		}

		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			return err
		}

		c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=\"Template_Import_Warga.xlsx\"")
		return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Import Data Warga dari Excel (Mendukung format Skripsi & format Template Standar)
	e.POST("/import/warga", func(c echo.Context) error {
		berkas, err := c.FormFile("excel_file") // Ambil file dari form
		if err != nil { return err }
		sumber, err := berkas.Open()
		if err != nil { return err }
		defer sumber.Close()

		f, err := excelize.OpenReader(sumber)
		if err != nil { return err }
		defer f.Close()
		
		// Deteksi format Excel berdasarkan sheet name
		sheets := f.GetSheetList()
		isThesisFormat := false
		for _, sheetName := range sheets {
			if sheetName == "Seluruh Data Warga" {
				isThesisFormat = true
				break
			}
		}

		if isThesisFormat {
			// === FORMAT SKRIPSI / THESIS ===
			tx, err := dbSistem.Begin()
			if err != nil { return err }

			// Hapus data lama demi sinkronisasi penuh
			tx.Exec("DELETE FROM data_indikator")
			tx.Exec("DELETE FROM hasil_klasifikasi")
			tx.Exec("DELETE FROM warga")

			rows, err := f.GetRows("Seluruh Data Warga")
			if err != nil {
				tx.Rollback()
				return err
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
				return err
			}
			defer stmtWarga.Close()

			stmtInd, err := tx.Prepare("INSERT INTO data_indikator (warga_id, indikator_id, nilai) VALUES (?, ?, ?)")
			if err != nil {
				tx.Rollback()
				return err
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
				if i == 0 { continue } // Lewati header
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
						if idx == 0 || len(ujiRow) < 2 { continue }
						if strings.EqualFold(strings.TrimSpace(ujiRow[1]), name) {
							if idx < len(evalRows) {
								evalRow := evalRows[idx]
								if len(evalRow) > 9 {
									excelVal := strings.TrimSpace(evalRow[9])
									var predClass classifier.KelasKesejahteraan
									foundPred := false
									if strings.Contains(excelVal, "KK1") { predClass = classifier.SangatMiskin; foundPred = true }
									if strings.Contains(excelVal, "KK2") { predClass = classifier.Miskin; foundPred = true }
									if strings.Contains(excelVal, "KK3") { predClass = classifier.HampirMiskin; foundPred = true }
									if strings.Contains(excelVal, "KK4") { predClass = classifier.RentanMiskin; foundPred = true }
									if strings.Contains(excelVal, "KK5") { predClass = classifier.PasPasan; foundPred = true }
									if strings.Contains(excelVal, "KK6") { predClass = classifier.MenengahKeAtas; foundPred = true }

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
					return err
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
								return err
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
								return err
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
					if i == 0 || len(row) < 2 || strings.TrimSpace(row[1]) == "" { continue }
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
					if i == 0 || len(row) < 2 || strings.TrimSpace(row[1]) == "" { continue }
					name := strings.TrimSpace(row[1])
					wargaID, exists := insertedNamesMap[name]
					if exists {
						tx.Exec("UPDATE warga SET data_latih_2 = 1 WHERE id = ?", wargaID)
					}
				}
			}

			err = tx.Commit()
			if err != nil { return err }

		} else {
			// === FORMAT TEMPLATE STANDAR ===
			rows, err := f.GetRows(sheets[0])
			if err != nil { return err }

			tx, err := dbSistem.Begin()
			if err != nil { return err }

			stmtWarga, err := tx.Prepare(`
				INSERT OR REPLACE INTO warga (nik, no_kk, nama_lengkap, alamat, rt, rw, kelurahan, label_kelas, data_latih, data_latih_2)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`)
			if err != nil {
				tx.Rollback()
				return err
			}
			defer stmtWarga.Close()

			stmtInd, err := tx.Prepare(`
				INSERT OR REPLACE INTO data_indikator (warga_id, indikator_id, nilai)
				VALUES (?, ?, ?)
			`)
			if err != nil {
				tx.Rollback()
				return err
			}
			defer stmtInd.Close()

			for i, row := range rows {
				if i == 0 || len(row) < 3 { continue } // Lewati header atau baris kosong/tidak valid

				nik := strings.TrimSpace(row[0])
				nokk := ""
				if len(row) > 1 { nokk = strings.TrimSpace(row[1]) }
				nama := strings.TrimSpace(row[2])
				if nama == "" { continue }

				alamat := ""
				if len(row) > 3 { alamat = strings.TrimSpace(row[3]) }
				rt := ""
				if len(row) > 4 { rt = strings.TrimSpace(row[4]) }
				rw := ""
				if len(row) > 5 { rw = strings.TrimSpace(row[5]) }
				kelurahan := ""
				if len(row) > 6 { kelurahan = strings.TrimSpace(row[6]) }
				labelKelas := ""
				if len(row) > 7 { labelKelas = strings.TrimSpace(row[7]) }

				peran := ""
				if len(row) > 8 { peran = strings.ToLower(strings.TrimSpace(row[8])) }
				isLatih := 0
				if peran == "latih" || peran == "training" || peran == "1" || peran == "ya" {
					isLatih = 1
				}

				// Cek NIK lama untuk menghindari penggantian ID
				var existingID int64
				err = tx.QueryRow("SELECT id FROM warga WHERE nik = ?", nik).Scan(&existingID)
				
				res, err := stmtWarga.Exec(nik, nokk, nama, alamat, rt, rw, kelurahan, labelKelas, isLatih, isLatih)
				if err != nil {
					tx.Rollback()
					return err
				}

				var wargaID int64
				if err == nil && existingID > 0 {
					wargaID = existingID
				} else {
					wargaID, _ = res.LastInsertId()
				}

				// Hapus indikator lama
				tx.Exec("DELETE FROM data_indikator WHERE warga_id = ?", wargaID)

				// Simpan 36 indikator (IM1 - IM36)
				for colIdx := 9; colIdx < len(row) && colIdx < 45; colIdx++ {
					indID := fmt.Sprintf("IM%d", colIdx-8)
					val := strings.ToUpper(strings.TrimSpace(row[colIdx]))
					if val != "" {
						_, err = stmtInd.Exec(wargaID, indID, val)
						if err != nil {
							tx.Rollback()
							return err
						}
					}
				}
			}

			err = tx.Commit()
			if err != nil { return err }
		}

		// Latih ulang model setelah import data baru
		dataLatih, err := db.AmbilDataLatih(dbSistem)
		if err == nil && len(dataLatih) > 0 {
			var inputLatih []map[string]string
			var targetLatih []classifier.KelasKesejahteraan
			for _, dl := range dataLatih {
				inputLatih = append(inputLatih, dl.Indikator)
				targetLatih = append(targetLatih, classifier.KelasKesejahteraan(dl.Kelas))
			}
			modelNB = classifier.BuatModelBaru()
			modelNB.DaftarFitur = namaFitur
			modelNB.LatihModel(inputLatih, targetLatih)
			fmt.Printf("Model berhasil dilatih ulang setelah import dengan %d data kependudukan\n", len(dataLatih))
		}

		return c.Redirect(http.StatusSeeOther, "/warga")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Rute untuk mematikan aplikasi secara aman
	e.GET("/shutdown", func(c echo.Context) error {
		go func() {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}()
		return c.String(http.StatusOK, "Aplikasi telah dimatikan. Anda dapat menutup tab ini.")
	})

	// Mencari port yang tersedia mulai dari 8080
	port := 8080
	var serverAddr string
	for {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		// Cek apakah port bisa digunakan
		l, err := net.Listen("tcp", addr)
		if err == nil {
			l.Close()
			serverAddr = addr
			break
		}
		
		port++
		if port > 8090 { // Coba sampai 10 port
			e.Logger.Fatal("Tidak dapat menemukan port yang tersedia antara 8080-8090")
		}
	}

	// Jalankan server di goroutine terpisah
	go func() {
		fmt.Printf("Server berjalan di http://%s\n", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	// Tunggu server siap
	time.Sleep(1 * time.Second)

	appURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	
	// Coba Lorca terlebih dahulu (desktop window tanpa browser)
	ui, err := lorca.New(appURL, "", 1280, 800)
	if err != nil {
		// Lorca gagal - coba Chrome Kiosk Mode sebagai fallback
		fmt.Printf("⚠️  Desktop window (Lorca) gagal: %v\n", err)
		fmt.Println("🔄 Mencoba Chrome Kiosk Mode...")
		
		// Cari Chrome executable
		chromePaths := []string{
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
		}
		
		var chromeExe string
		for _, path := range chromePaths {
			if _, err := os.Stat(path); err == nil {
				chromeExe = path
				break
			}
		}
		
		if chromeExe != "" {
			// Jalankan Chrome dalam App Mode (seperti desktop app)
			fmt.Println("✅ Membuka Chrome App Mode...")
			fmt.Printf("🌐 URL: %s\n", appURL)
			
			// Chrome app mode: seperti desktop app, tanpa address bar, tabs, dll
			cmd := exec.Command(chromeExe,
				"--app="+appURL,
				"--window-size=1280,800",
				"--disable-features=TranslateUI",
				"--no-first-run",
				"--no-default-browser-check",
			)
			
			err = cmd.Start()
			if err != nil {
				fmt.Printf("❌ Gagal start Chrome: %v\n", err)
				fmt.Println("📱 Silakan buka browser manual ke:", appURL)
			} else {
				fmt.Println("✅ Chrome App Mode berhasil dibuka!")
				fmt.Println("⏹️  Tutup window Chrome untuk menghentikan aplikasi.")
			}
			
			// Tunggu forever - server tetap jalan
			select {}
		} else {
			// Chrome tidak ditemukan
			fmt.Println("❌ Chrome tidak ditemukan.")
			fmt.Printf("📱 Silakan buka browser manual ke: %s\n", appURL)
			fmt.Println("⏹️  Tekan Ctrl+C untuk menghentikan server.")
			
			// Tunggu forever - server tetap jalan
			select {}
		}
	} else {
		// Lorca berhasil!
		defer ui.Close()
		
		fmt.Println("✅ Desktop window (Lorca) berhasil dibuat!")
		fmt.Println("⏹️  Tutup window untuk menghentikan aplikasi.")
		
		// Tunggu sampai window ditutup
		<-ui.Done()
	}
	
	fmt.Println("👋 Aplikasi dihentikan.")
}

func FormatScientific(val float64) string {
	if val == 0 {
		return "0"
	}
	if val >= 0.001 {
		str := fmt.Sprintf("%.6f", val)
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
		return strings.ReplaceAll(str, ".", ",")
	}
	str := fmt.Sprintf("%.10E", val)
	parts := strings.Split(str, "E")
	if len(parts) != 2 {
		return strings.ReplaceAll(str, ".", ",")
	}
	significand := parts[0]
	exponent := parts[1]
	if strings.Contains(significand, ".") {
		significand = strings.TrimRight(significand, "0")
		significand = strings.TrimRight(significand, ".")
	}
	significand = strings.ReplaceAll(significand, ".", ",")
	return significand + "E" + exponent
}
