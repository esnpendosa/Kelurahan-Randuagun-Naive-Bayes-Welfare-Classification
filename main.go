package main // Paket utama sebagai titik masuk (entry point) aplikasi

import (
	"fmt"                    // Mengimpor paket untuk format teks dan output
	"html/template"          // Mengimpor paket untuk mesin template HTML
	"io"                     // Mengimpor paket untuk operasi input/output
	"net/http"               // Mengimpor paket untuk protokol HTTP
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
	"os/exec"  // Paket untuk menjalankan perintah eksternal
	"runtime"  // Paket untuk mendeteksi sistem operasi saat ini
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
		dbSistem.QueryRow("SELECT COUNT(*) FROM hasil_klasifikasi").Scan(&jumlahKlasifikasi)
		var jumlahLatih int
		dbSistem.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih = 1").Scan(&jumlahLatih)

		// Distribusi per kelas dari seluruh data (latih + uji yang punya label)
		rows, _ := dbSistem.Query(`
			SELECT label_kelas, COUNT(*) as jml
			FROM warga
			WHERE label_kelas != ''
			GROUP BY label_kelas
			ORDER BY jml DESC
		`)
		defer rows.Close()
		distribusi := []map[string]interface{}{}
		totalLabel := 0
		for rows.Next() {
			var label string; var jml int
			rows.Scan(&label, &jml)
			distribusi = append(distribusi, map[string]interface{}{"Label": label, "Count": jml})
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
		rows, err := dbSistem.Query("SELECT id, nik, no_kk, nama_lengkap, alamat, kelurahan, data_latih FROM warga ORDER BY data_latih DESC, id ASC")
		if err != nil {
			return err
		}
		defer rows.Close()

		var daftarWarga []map[string]interface{}
		for rows.Next() {
			var id, isLatih int
			var nik, nokk, nama, alamat, kelurahan string
			if err := rows.Scan(&id, &nik, &nokk, &nama, &alamat, &kelurahan, &isLatih); err != nil {
				continue
			}
			daftarWarga = append(daftarWarga, map[string]interface{}{
				"ID": id, "NIK": nik, "NoKK": nokk, "NamaKK": nama, "Alamat": alamat, "Kelurahan": kelurahan, "IsTraining": isLatih == 1,
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

		// Memanggil fungsi Prediksi pada model Naive Bayes menggunakan data inputan
		// Hasilnya adalah map probabilitas kemiripan warga tersebut ke setiap kelas (1-6)
		peluang := modelNB.Prediksi(inputan)
		
		// Memilih kelas kesejahteraan dengan probabilitas tertinggi dari hasil perhitungan
		kelasTerbaik := modelNB.AmbilKelasTerbaik(peluang)
		
		// Mengubah ID kelas (1-6) menjadi teks nama kelas (contoh: "Miskin", "Pas-pasan")
		namaKelas := classifier.DaftarNamaKelas[kelasTerbaik]

		// Mengubah hasil perhitungan peluang menjadi format string JSON agar bisa disimpan di database
		peluangJSON, _ := json.Marshal(peluang)
		
		// Menjalankan query SQL untuk menyimpan hasil akhir klasifikasi ke dalam tabel hasil_klasifikasi
		dbSistem.Exec("INSERT INTO hasil_klasifikasi (warga_id, nama_kelas, probabilitas) VALUES (?, ?, ?)", 
			idWarga, namaKelas, string(peluangJSON))

		// Setelah perhitungan selesai, sistem akan mengarahkan pengguna (redirect) ke halaman hasil visualisasi
		return c.Redirect(http.StatusSeeOther, "/hasil/"+idWarga)
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

		var daftarPeluang []map[string]interface{}
		for _, k := range modelNB.SemuaKelas {
			daftarPeluang = append(daftarPeluang, map[string]interface{}{
				"Label":   classifier.DaftarNamaKelas[k],
				"Percent": petaPeluang[k] * 100, // Jika k tidak ada, otomatis 0.0
				"IsMax":   classifier.DaftarNamaKelas[k] == namaKelas,
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
		dataLatih, _ := db.AmbilDataLatih(dbSistem)
		dataUji, _ := db.AmbilDataUji(dbSistem)

		// Distribusi data latih per kelas
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

		// Semua warga untuk tab dinamis
		rows2, _ := dbSistem.Query("SELECT id, nik, nama_lengkap, data_latih FROM warga ORDER BY id ASC")
		defer rows2.Close()
		var semuaWarga []map[string]interface{}
		for rows2.Next() {
			var id, isLatih int; var nik, nama string
			rows2.Scan(&id, &nik, &nama, &isLatih)
			semuaWarga = append(semuaWarga, map[string]interface{}{
				"ID":"" + fmt.Sprintf("%d",id), "NIK":nik, "NamaKK":nama, "IsTraining":isLatih==1,
			})
		}

		var jumlahLatih2 int
		dbSistem.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih = 1").Scan(&jumlahLatih2)

		return c.Render(http.StatusOK, "training.html", map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Stats": map[string]interface{}{
				"TotalTraining": len(dataLatih),
				"TotalTesting":  len(dataUji),
				"TotalTraining2": jumlahLatih2,
				"LastTrained": "-",
			},
			"DistribusiLatih": distLatih,
			"DistribusiUji":   distUji,
			"SemuaWarga":      semuaWarga,
			"HasResult":       false,
			"classifier": map[string]interface{}{"ClassNames": classifier.DaftarNamaKelas},
		})
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Training Model - Proses Hitung & Evaluasi (POST)
	e.POST("/training/proses", func(c echo.Context) error {
		modelDipilih := c.FormValue("model") // training1 atau training2
		filterKelas := c.FormValue("filter_kelas")
		_ = modelDipilih

		dataLatih, _ := db.AmbilDataLatih(dbSistem)
		dataUji, _ := db.AmbilDataUji(dbSistem)

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
		for _, du := range dataUji {
			p := modelNB.Prediksi(du.Indikator)
			pred := modelNB.AmbilKelasTerbaik(p)
			aktual := classifier.KelasKesejahteraan(du.Kelas)
			if matriks[aktual] == nil { matriks[aktual] = make(map[classifier.KelasKesejahteraan]int) }
			matriks[aktual][pred]++
			if aktual == pred { benar++ }
			total++
		}
		akurasi := 0.0; if total > 0 { akurasi = float64(benar) / float64(total) }

		var totP, totR, totF1, cnt float64
		for _, k := range modelNB.SemuaKelas {
			tp := float64(matriks[k][k])
			var fp, fn float64
			for _, ac := range modelNB.SemuaKelas { if ac != k { fp += float64(matriks[ac][k]) } }
			for _, pc := range modelNB.SemuaKelas { if pc != k { fn += float64(matriks[k][pc]) } }
			pr := 0.0; if tp+fp > 0 { pr = tp / (tp + fp) }
			rc := 0.0; if tp+fn > 0 { rc = tp / (tp + fn) }
			f1 := 0.0; if pr+rc > 0 { f1 = 2 * pr * rc / (pr + rc) }
			totP += pr; totR += rc; totF1 += f1; cnt++
		}
		if cnt == 0 { cnt = 1 }

		urutan := []string{"Sangat Miskin","Miskin","Hampir Miskin","Rentan Miskin","Pas-pasan","Menengah ke Atas"}
		distLatihMap := make(map[string]int)
		for _, d := range dataLatih { distLatihMap[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		distUjiMap := make(map[string]int)
		for _, d := range dataUji { distUjiMap[classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(d.Kelas)]]++ }
		var distLatih, distUji []map[string]interface{}
		for _, nama := range urutan {
			jml := distLatihMap[nama]; pct := 0.0; if len(dataLatih)>0 { pct = float64(jml)/float64(len(dataLatih))*100 }
			distLatih = append(distLatih, map[string]interface{}{"Label":nama,"Count":jml,"Percent":pct})
			jmlU := distUjiMap[nama]; pctU := 0.0; if len(dataUji)>0 { pctU = float64(jmlU)/float64(len(dataUji))*100 }
			distUji = append(distUji, map[string]interface{}{"Label":nama,"Count":jmlU,"Percent":pctU})
		}

		// Semua warga tab dinamis
		rows2, _ := dbSistem.Query("SELECT id, nik, nama_lengkap, data_latih FROM warga ORDER BY id ASC")
		defer rows2.Close()
		var semuaWarga []map[string]interface{}
		for rows2.Next() {
			var id, isLatih int; var nik, nama string
			rows2.Scan(&id, &nik, &nama, &isLatih)
			semuaWarga = append(semuaWarga, map[string]interface{}{
				"ID":fmt.Sprintf("%d",id), "NIK":nik, "NamaKK":nama, "IsTraining":isLatih==1,
			})
		}

		modelLabel := "Data Training 1"
		if modelDipilih == "training2" { modelLabel = "Data Training 2" }
		if filterKelas != "" { modelLabel += " (filter: " + filterKelas + ")" }

		errorUji := ""
		if len(dataUji) == 0 { errorUji = "Data uji tidak ditemukan. Pastikan ada data dengan data_latih=0 dan label_kelas terisi." }

		var jumlahLatih2 int
		dbSistem.QueryRow("SELECT COUNT(*) FROM warga WHERE data_latih = 1").Scan(&jumlahLatih2)

		return c.Render(http.StatusOK, "training.html", map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Stats": map[string]interface{}{
				"TotalTraining":  len(dataLatih),
				"TotalTesting":   len(dataUji),
				"TotalTraining2": jumlahLatih2,
				"Accuracy":       akurasi * 100,
				"Precision":      totP / cnt,
				"Recall":         totR / cnt,
				"F1Score":        totF1 / cnt,
				"LastTrained":    time.Now().Format("02 Jan 2006 15:04"),
			},
			"Matrix":          matriks,
			"Classes":         modelNB.SemuaKelas,
			"ErrorUji":        errorUji,
			"HasResult":       true,
			"ModelDipakai":    modelLabel,
			"DistribusiLatih": distLatih,
			"DistribusiUji":   distUji,
			"SemuaWarga":      semuaWarga,
			"classifier": map[string]interface{}{"ClassNames": classifier.DaftarNamaKelas},
		})
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Set Peran Warga (training/uji) - Dinamis
	e.POST("/warga/set-peran", func(c echo.Context) error {
		id := c.FormValue("id")
		peran := c.FormValue("peran")
		isLatih := 0
		if peran == "latih" { isLatih = 1 }
		dbSistem.Exec("UPDATE warga SET data_latih = ? WHERE id = ?", isLatih, id)
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
		// Mengambil semua data yang sudah terklasifikasi
		rows, _ := dbSistem.Query(`
			SELECT nik, nama_lengkap, alamat, label_kelas, dibuat_pada 
			FROM warga 
			WHERE data_latih = 1
		`)
		defer rows.Close()
		
		var rekap []map[string]interface{}
		for rows.Next() {
			var nik, nama, alamat, kelas, tgl string
			rows.Scan(&nik, &nama, &alamat, &kelas, &tgl)
			
			// Konversi kode kelas (1-6) ke nama kelas yang bisa dibaca
			var kodeKelas int
			fmt.Sscanf(kelas, "%d", &kodeKelas)
			namaKelas := classifier.DaftarNamaKelas[classifier.KelasKesejahteraan(kodeKelas)]
			if namaKelas == "" { namaKelas = "Terklasifikasi" }

			rekap = append(rekap, map[string]interface{}{
				"NIK":       nik,
				"NamaKK":    nama,
				"Alamat":    alamat,
				"ClassName": namaKelas,
				"Date":      tgl,
				"Status":    "Data Latih",
			})
		}

		data := map[string]interface{}{
			"User":    ambilDataPengguna(c),
			"Results": rekap,
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

		db.TambahWarga(dbSistem, nik, nokk, nama, alamat, rt, rw, kelurahan)
		return c.Redirect(http.StatusSeeOther, "/warga")
	}, middlewareAutentikasi, middlewarePeran("Admin"))

	// Edit Warga (Form)
	e.GET("/warga/edit/:id", func(c echo.Context) error {
		id := c.Param("id")
		var w struct {
			ID   int
			NIK  string
			NoKK string
			Nama string
			Alm  string
			RT   string
			RW   string
			Klh  string
		}
		dbSistem.QueryRow("SELECT id, nik, no_kk, nama_lengkap, alamat, rt, rw, kelurahan FROM warga WHERE id = ?", id).Scan(&w.ID, &w.NIK, &w.NoKK, &w.Nama, &w.Alm, &w.RT, &w.RW, &w.Klh)
		
		data := map[string]interface{}{
			"User": ambilDataPengguna(c),
			"Warga": map[string]interface{}{
				"ID": w.ID, "NIK": w.NIK, "NoKK": w.NoKK, "NamaKK": w.Nama, "Alamat": w.Alm, "RT": w.RT, "RW": w.RW, "Kelurahan": w.Klh,
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

		dbSistem.Exec("UPDATE warga SET nik=?, no_kk=?, nama_lengkap=?, alamat=?, rt=?, rw=?, kelurahan=? WHERE id=?", 
			nik, nokk, nama, alamat, rt, rw, kelurahan, id)
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
			SELECT r.nik, r.nama_lengkap, r.alamat, h.nama_kelas, h.dibuat_pada 
			FROM hasil_klasifikasi h
			JOIN warga r ON h.warga_id = r.id
			ORDER BY h.dibuat_pada DESC
		`)
		if err != nil { return err }
		defer rows.Close()

		f := excelize.NewFile() // Membuat file Excel baru
		lembar := "Laporan"
		f.SetSheetName("Sheet1", lembar)
		
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

		// Mengirim file Excel sebagai unduhan di browser
		c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=laporan_kesejahteraan.xlsx")
		return f.Write(c.Response().Writer)
	}, middlewareAutentikasi)

	// Import Data Warga dari Excel
	e.POST("/import/warga", func(c echo.Context) error {
		berkas, err := c.FormFile("excel_file") // Ambil file dari form
		if err != nil { return err }
		sumber, err := berkas.Open()
		if err != nil { return err }
		defer sumber.Close()

		f, err := excelize.OpenReader(sumber)
		if err != nil { return err }
		
		baris, _ := f.GetRows(f.GetSheetName(0))
		for i, b := range baris {
			if i == 0 || len(b) < 2 { continue } // Lewati header atau baris kosong
			
			nik := b[0]
			nama := b[1]
			alm := ""
			if len(b) > 2 { alm = b[2] }
			klh := ""
			if len(b) > 3 { klh = b[3] }
			
			// Simpan atau perbarui data warga
			dbSistem.Exec("INSERT OR REPLACE INTO warga (nik, nama_lengkap, alamat, kelurahan) VALUES (?, ?, ?, ?)", nik, nama, alm, klh)
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
	var serverErr error
	for {
		addr := fmt.Sprintf(":%d", port)
		// Cek apakah port bisa digunakan
		l, err := net.Listen("tcp", addr)
		if err == nil {
			l.Close()
			
			// Membuka browser secara otomatis ke port yang ditemukan
			go func(p int) {
				time.Sleep(500 * time.Millisecond)
				bukaBrowser(fmt.Sprintf("http://localhost:%d", p))
			}(port)

			fmt.Printf("Server berjalan di http://localhost:%d\n", port)
			serverErr = e.Start(addr)
			break
		}
		
		port++
		if port > 8090 { // Coba sampai 10 port
			serverErr = fmt.Errorf("Tidak dapat menemukan port yang tersedia antara 8080-8090")
			break
		}
	}

	if serverErr != nil {
		e.Logger.Fatal(serverErr)
	}
}

// bukaBrowser membuka URL di browser default berdasarkan sistem operasi
func bukaBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default: // linux
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		fmt.Printf("Gagal membuka browser: %v\n", err)
	}
}
