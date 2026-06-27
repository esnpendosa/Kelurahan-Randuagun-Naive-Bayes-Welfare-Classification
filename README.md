# Sistem Klasifikasi Tingkat Kesejahteraan Masyarakat - Kelurahan Randuagung
## Berbasis Metode Naive Bayes dengan Arsitektur Go Server-Side Rendering (SSR) & Desktop App Mode

Aplikasi ini dikembangkan untuk mengklasifikasikan tingkat kesejahteraan keluarga di Kelurahan Randuagung ke dalam **6 kelas kesejahteraan** berdasarkan **36 indikator kemiskinan** terstandardisasi. Menggunakan algoritma **Naive Bayes** dengan perbaikan numerik (**Laplace Smoothing** dan **Log-Probability**) guna menjamin akurasi dan stabilitas perhitungan pada skala data kecil hingga menengah.

Aplikasi ini berjalan secara mandiri (self-contained desktop-shell) memanfaatkan runtime **Go (Golang)**, database **SQLite**, dan interface rendering **HTML5 Go Template (SSR)**.

---

## 📌 DAFTAR ISI
1. [Fitur Utama Sistem](#-fitur-utama-sistem)
2. [Arsitektur & Teknologi Stack](#-arsitektur--teknologi-stack)
3. [Rancangan Database (Skema SQLite)](#-rancangan-database-skema-sqlite)
4. [Algoritma Naive Bayes (Detail Perhitungan & Formula)](#-algoritma-naive-bayes-detail-perhitungan--formula)
5. [Daftar 36 Indikator Kesejahteraan (IM1 - IM36)](#-daftar-36-indikator-kesejahteraan-im1---im36)
6. [Struktur Direktori Proyek](#-struktur-direktori-proyek)
7. [Panduan Instalasi & Menjalankan Aplikasi](#-panduan-instalasi--menjalankan-aplikasi)
8. [Panduan Pembangunan Biner Mandiri (Build.exe)](#-panduan-pembangunan-biner-mandiri-buildexe)
9. [Kontributor & Lisensi](#-kontributor--lisensi)

---

## 🌟 Fitur Utama Sistem

Sistem ini didesain khusus untuk kebutuhan kantor Kelurahan Randuagung dengan mengutamakan kemudahan operasional (user-friendly) dan fleksibilitas:
*   **Dashboard Utama**: Statistik ringkas kependudukan, grafik komposisi tingkat kesejahteraan secara real-time menggunakan *Chart.js*, dan status performa model (akurasi terakhir).
*   **Manajemen Data Warga (CRUD)**:
    *   Pengelolaan profil warga (NIK, No KK, Nama, Alamat, RT/RW, Kelurahan).
    *   Pencarian instan berdasarkan NIK/Nama serta filter dinamis per-kategori kesejahteraan berbasis JavaScript client-side tanpa reload halaman.
    *   Opsi **Import dari File Excel (.xlsx)** dan **Export Master Data ke Excel**.
*   **Klasifikasi Real-Time (Formulir Cerdas)**:
    *   Formulir kuesioner interaktif untuk 36 indikator (IM1 - IM36).
    *   **Fitur Auto-Fill**: Mendeteksi data historis warga secara otomatis dari database dan mengisi formulir secara instan, menghemat waktu petugas kelurahan saat pembaharuan data.
*   **Training & Evaluasi Model (Excel-Style)**:
    *   Melatih model Naive Bayes secara dinamis menggunakan dua dataset (Dataset 1 dan Dataset 2).
    *   Menampilkan **Confusion Matrix 6x6** lengkap dengan baris/kolom Jumlah serta total data uji global.
    *   Metrik evaluasi detail (Akurasi, Precision per-kelas, Recall per-kelas, dan Macro F1-Score) yang identik dengan format perhitungan Microsoft Excel.
    *   Pilihan peran warga (data latih vs data uji) yang dapat ditukar secara dinamis.
*   **Laporan Hasil & Cetak Bersih**:
    *   Tabel laporan rekapitulasi kesejahteraan yang dilengkapi filter kategori dan export ke Excel.
    *   **Styling Ramah Cetak (@media print)**: Menghilangkan sidebar, topbar, dan tombol-tombol navigasi secara otomatis saat dicetak ke kertas (Ctrl+P) sehingga menghasilkan dokumen resmi kelurahan yang bersih dan rapi.

---

## 💻 Arsitektur & Teknologi Stack

Aplikasi ini menggunakan pendekatan monolitik modern dengan performa tinggi dan konsumsi memori rendah:
*   **Back-End Engine**: Go (Golang) versi 1.20+, memanfaatkan framework **Echo v4** yang sangat cepat dan efisien.
*   **Front-End Rendering**: **Go HTML Template (`html/template`)** dengan sistem modular (template inheritance) dan Vanilla CSS. Tidak ada ketergantungan JavaScript framework eksternal (seperti React/Vue) untuk menjaga portabilitas.
*   **Database**: **SQLite 3** menggunakan driver murni Go (`modernc.org/sqlite`) tanpa memerlukan instalasi CGO compiler (GCC) di Windows.
*   **Desktop Shell Mode**: Menggunakan **Chrome App Mode Kiosk Mode** sebagai fallback desktop shell lokal, membungkus web server lokal sehingga berjalan seperti aplikasi desktop mandiri tanpa address bar browser.

---

## 🗄️ Rancangan Database (Skema SQLite)

Berikut adalah skema tabel database yang didefinisikan dalam database lokal `data_skripsi.db`:

```sql
-- 1. Tabel Pengguna (Manajemen Akun Akses)
CREATE TABLE IF NOT EXISTS pengguna (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama_pengguna TEXT UNIQUE NOT NULL,
    kata_sandi TEXT NOT NULL, -- Enkripsi bcrypt hash
    peran TEXT NOT NULL -- 'Admin' atau 'Operator'
);

-- 2. Tabel Warga (Data Kependudukan)
CREATE TABLE IF NOT EXISTS warga (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nik TEXT UNIQUE NOT NULL,
    no_kk TEXT NOT NULL,
    nama_lengkap TEXT NOT NULL,
    alamat TEXT NOT NULL,
    rt TEXT NOT NULL,
    rw TEXT NOT NULL,
    kelurahan TEXT NOT NULL,
    data_latih INTEGER DEFAULT 0, -- 1 = Data Latih Dataset 1, 0 = Data Uji Dataset 1
    data_latih_2 INTEGER DEFAULT 0, -- 1 = Data Latih Dataset 2, 0 = Data Uji Dataset 2
    label_kelas TEXT DEFAULT '', -- Label kesejahteraan aktual terakhir
    idpengguna INTEGER REFERENCES pengguna(id)
);

-- 3. Tabel Data Indikator (Transaksi Pilihan Kuesioner Warga)
CREATE TABLE IF NOT EXISTS data_indikator (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    warga_id INTEGER NOT NULL,
    indikator_id TEXT NOT NULL, -- IM1 s.d. IM36
    nilai TEXT NOT NULL, -- Pilihan: 'A', 'B', 'C', 'D'
    FOREIGN KEY(warga_id) REFERENCES warga(id)
);

-- 4. Tabel Hasil Klasifikasi (Penyimpanan Prediksi Sistem)
CREATE TABLE IF NOT EXISTS hasil_klasifikasi (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    warga_id INTEGER NOT NULL,
    nama_kelas TEXT NOT NULL, -- Hasil kelas prediksi terbaik
    probabilitas TEXT NOT NULL, -- Nilai peluang 6 kelas dalam format JSON string
    dibuat_pada DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(warga_id) REFERENCES warga(id)
);
```

---

## 🧮 Algoritma Naive Bayes (Detail Perhitungan & Formula)

Sistem ini mengimplementasikan algoritma Naive Bayes Classifier untuk memprediksi tingkat kesejahteraan ke dalam 6 kelas:
1.  **Sangat Miskin (KK1)**
2.  **Miskin (KK2)**
3.  **Hampir Miskin (KK3)**
4.  **Rentan Miskin (KK4)**
5.  **Pas-pasan (KK5)**
6.  **Menengah ke Atas (KK6)**

### Rumus Perhitungan Langkah-Demi-Langkah:

#### 1. Prior Probability (Peluang Awal Kelas)
Peluang awal kemunculan dari setiap kelas kesejahteraan $C_k$:
$$P(C_k) = \frac{\text{Jumlah Warga Latih Kelas } C_k}{\text{Total Seluruh Warga Latih}}$$

#### 2. Likelihood Probability dengan Laplace Smoothing
Menghitung peluang bersyarat kemunculan suatu nilai indikator $X_i$ (misal: IM1 = 'A') di kelas $C_k$. Kita menerapkan **Laplace Smoothing** untuk menghindari probabilitas nol jika suatu nilai atribut tidak ada di data training kelas tertentu:
$$P(X_i = v | C_k) = \frac{Count(X_i = v \text{ di kelas } C_k) + 1}{Count(C_k) + |V|}$$
*Keterangan:*
*   $Count(X_i = v \text{ di kelas } C_k)$: Jumlah kemunculan nilai $v$ (misalnya 'A') untuk indikator $X_i$ pada data training kelas $C_k$.
*   $Count(C_k)$: Jumlah seluruh data latih di kelas $C_k$.
*   $|V|$: Jumlah kategori pilihan jawaban pada indikator tersebut (untuk IM1 - IM36 bernilai 4 pilihan: A, B, C, D).

#### 3. Perhitungan Menggunakan Log-Probability (Pencegahan Underflow)
Karena probabilitas berupa pecahan desimal kecil antara $0$ dan $1$, perkalian 36 indikator secara berturut-turut dapat menyebabkan nilai desimal mendekati nol hingga tidak terbaca oleh komputer (*floating-point underflow*).
Untuk itu, kita konversikan perkalian probabilitas menjadi **penjumlahan logaritma**:
$$\ln P(C_k | X) = \ln P(C_k) + \sum_{i=1}^{36} \ln P(X_i | C_k)$$

#### 4. Normalisasi Peluang Akhir (Softmax)
Setelah menjumlahkan nilai logaritma untuk masing-masing ke-6 kelas, nilai tersebut dikonversi kembali ke bentuk desimal normal melalui operasi eksponensial dan dibagi dengan jumlah total peluang seluruh kelas agar total akhir probabilitas bernilai $1.0$ ($100\%$):
$$P(C_k | X) = \frac{e^{\ln P(C_k | X)}}{\sum_{j=1}^{6} e^{\ln P(C_j | X)}}$$
Kelas dengan probabilitas tertinggi ($P(C_k | X)$ maksimal) dipilih sebagai kelas prediksi terbaik.

---

## 📋 Daftar 36 Indikator Kesejahteraan (IM1 - IM36)

Indikator ini dikelompokkan ke dalam tiga bagian besar di formulir input:

### A. Bagian 1: Kondisi Rumah
1.  **IM1**: Status Kepemilikan Bangunan Tempat Tinggal
2.  **IM2**: Luas Lantai Bangunan per Anggota Keluarga
3.  **IM3**: Jenis Lantai Terluas
4.  **IM4**: Jenis Dinding Rumah Terluas
5.  **IM5**: Kondisi Dinding Rumah
6.  **IM6**: Jenis Atap Terluas
7.  **IM7**: Kondisi Atap Rumah
8.  **IM8**: Sumber Air Minum Utama
9.  **IM9**: Jarak ke Penampungan Kotoran Terdekat
10. **IM10**: Sumber Penerangan Utama
11. **IM11**: Fasilitas Tempat Buang Air Besar
12. **IM12**: Jenis Kloset yang Digunakan

### B. Bagian 2: Ekonomi Keluarga
13. **IM13**: Bahan Bakar Utama untuk Memasak
14. **IM14**: Frekuensi Makan Anggota Keluarga per Hari
15. **IM15**: Frekuensi Membeli Daging/Susu per Minggu
16. **IM16**: Frekuensi Membeli Pakaian Baru per Tahun
17. **IM17**: Kemampuan Berobat ke Puskesmas/Klinik
18. **IM18**: Lapangan Pekerjaan Utama Kepala Keluarga
19. **IM19**: Pendidikan Tertinggi Kepala Keluarga
20. **IM20**: Status Kepemilikan Tabungan/Aset Cair
21. **IM21**: Pendapatan Rata-rata Keluarga per Bulan
22. **IM22**: Jumlah Tanggungan dalam Keluarga
23. **IM23**: Kepemilikan Jaminan Kesehatan (BPJS/KIS)
24. **IM24**: Status Penerimaan Bantuan Sosial Pemerintah

### C. Bagian 3: Aset & Fasilitas
25. **IM25**: Kepemilikan Sepeda Motor pribadi
26. **IM26**: Kepemilikan Mobil pribadi
27. **IM27**: Kepemilikan Pendingin Ruangan (AC)
28. **IM28**: Kepemilikan Alat Transportasi Non-Motor
29. **IM29**: Kepemilikan Lahan Pertanian/Perkebunan
30. **IM30**: Kepemilikan Hewan Ternak Besar (Sapi/Kambing)
31. **IM31**: Kepemilikan Emas/Perhiasan Bernilai Tinggi
32. **IM32**: Kepemilikan Handphone/Smartphone Aktif
33. **IM33**: Akses Jaringan Internet Rumah (WiFi/Paket Data)
34. **IM34**: Kepemilikan Komputer/Laptop Pribadi
35. **IM35**: Kepemilikan Tabung Gas Masak Bersubsidi
36. **IM36**: Kepemilikan Mesin Cuci Rumah Tangga

Setiap indikator memiliki 4 pilihan jawaban terstandardisasi (**A**, **B**, **C**, **D**) yang merepresentasikan tingkat kondisi (dari yang paling memprihatinkan hingga yang paling mapan).

---

## 📂 Struktur Direktori Proyek

```
Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification/
│
├── main.go                            ← Entry point aplikasi (web server, routing, init database)
├── data_skripsi.db                    ← Database SQLite utama tempat menyimpan semua data
├── go.mod / go.sum                    ← Dependensi runtime Golang
├── build.bat                          ← Script batch Windows untuk build exe dengan ikon
├── README.md                          ← Dokumentasi utama sistem
├── MANUAL_BOOK.md                     ← Buku panduan operasional pengguna kelurahan
│
├── internal/
│   ├── classifier/
│   │   ├── naive_bayes.go             ← Logika utama perhitungan algoritma Naive Bayes
│   │   └── indicators.go              ← Metadata 36 indikator & pilihan jawabannya
│   └── db/
│       └── db.go                      ← Penanganan database (inisialisasi, CRUD warga, seeding user)
│
├── templates/
│   ├── base.html                      ← Layout master (navigasi, topbar, sidebar, layout dasar)
│   ├── index.html                     ← Template Dashboard Utama
│   ├── warga.html                     ← Template Daftar Warga & Filter Instan
│   ├── warga_tambah.html              ← Template Form Tambah Warga
│   ├── warga_edit.html                ← Template Form Edit Warga
│   ├── klasifikasi.html               ← Template Form Pengisian 36 Indikator
│   ├── hasil.html                     ← Template Hasil Klasifikasi & Grafik Probabilitas
│   ├── training.html                  ← Template Evaluasi Model (Confusion Matrix Excel)
│   └── laporan.html                   ← Template Rekap Laporan & Export Excel
│
└── static/
    └── css/
        └── style.css                  ← CSS kustom (Design System Flat Slate/Blue Palette)
```

---

## 🚀 Panduan Instalasi & Menjalankan Aplikasi

### Persyaratan Sistem:
- **Sistem Operasi**: Windows 7 / 8 / 10 / 11 (32-bit atau 64-bit).
- **Runtime**: Golang v1.20 atau versi di atasnya (jika dijalankan dari source code).

### A. Menjalankan dari Source Code (Development)

1.  **Unduh & Instal Golang**: Download installer Go di [golang.org/dl](https://golang.org/dl/) dan selesaikan instalasi.
2.  **Clone Repositori**:
    ```bash
    git clone https://github.com/esnpendosa/Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification.git
    cd Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification
    ```
3.  **Unduh Dependensi**:
    ```bash
    go mod tidy
    ```
4.  **Jalankan Aplikasi**:
    ```bash
    go run main.go
    ```
    *Aplikasi akan mendeteksi port kosong secara otomatis (default: `8082`) dan membuka interface desktop secara mandiri.*

### B. Menjalankan Langsung via Biner Executable (Produksi)

1.  Temukan file **`Klasifikasi-Warga-Randuagung.exe`** di dalam folder root.
2.  **Klik dua kali** pada file tersebut.
3.  Browser bawaan Windows Anda akan terbuka secara otomatis pada halaman login.

---

## 🛠️ Panduan Pembangunan Biner Mandiri (Build .exe)

Untuk mengompilasi kode program menjadi satu file executable tunggal Windows yang mencakup seluruh aset statis dan template HTML (menggunakan fitur `//go:embed` bawaan Go):

1.  Buka terminal/Command Prompt di direktori root aplikasi.
2.  Jalankan perintah script batch yang disediakan:
    ```cmd
    build.bat
    ```
3.  Script akan secara otomatis memproses kompilasi dan menghasilkan file **`Klasifikasi-Warga-Randuagung.exe`** yang sudah tersemat icon desktop resmi.

---

## 📄 Lisensi
Hak Cipta © 2026 Muhammad As'ad Muhibbin Akbar.
*Skripsi Program Studi Teknik Informatika - Kelurahan Randuagung*
