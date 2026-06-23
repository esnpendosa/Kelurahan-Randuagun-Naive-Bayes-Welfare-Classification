# Sistem Klasifikasi Tingkat Kesejahteraan Masyarakat - Kelurahan Randuagung

Aplikasi klasifikasi tingkat kesejahteraan masyarakat menggunakan algoritma **Naive Bayes** yang dibangun dengan bahasa pemrograman **Go (Golang)** dan database **SQLite**. Proyek ini dikembangkan untuk membantu Kelurahan Randuagung dalam menentukan klasifikasi ekonomi warga secara objektif berdasarkan 36 indikator kesejahteraan.

## Fitur Utama
- **Dashboard**: Visualisasi statistik kependudukan, akurasi model, dan distribusi kelas.
- **Data Warga**: Manajemen data kependudukan (CRUD) dan fitur import dari Excel.
- **Klasifikasi Baru**: Melakukan prediksi klasifikasi warga secara real-time.
- **Training & Evaluasi Model**: Evaluasi performa algoritma Naive Bayes menggunakan data uji dengan Confusion Matrix.
- **Laporan**: Rekapitulasi hasil klasifikasi yang dapat diekspor ke format Excel.
- **Pengaturan Akun**: Pengaturan akun akses admin.

---

## Struktur Folder Proyek

Berikut adalah struktur folder dan file dari sistem klasifikasi ini:

```
Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification/
│
├── main.go                          ← Otak utama aplikasi (server, routing, semua endpoint)
├── data_skripsi.db                  ← File database SQLite (semua data tersimpan di sini)
├── Klasifikasi-Warga-Randuagung.exe ← File yang bisa langsung dijalankan (hasil build)
├── go.mod / go.sum                  ← Daftar library yang digunakan
│
├── internal/
│   ├── classifier/
│   │   ├── naive_bayes.go           ← Logika perhitungan algoritma Naive Bayes
│   │   └── indicators.go            ← Daftar 36 indikator kemiskinan & pilihan jawabannya
│   └── db/
│       └── db.go                    ← Semua fungsi untuk baca/tulis database (CRUD)
│
├── templates/
│   ├── base.html                    ← Template induk (navbar, footer — dipakai semua halaman)
│   ├── index.html                   ← Halaman Dashboard
│   ├── warga.html                   ← Daftar seluruh warga
│   ├── warga_tambah.html            ← Form tambah warga baru
│   ├── warga_edit.html              ← Form edit data warga
│   ├── klasifikasi.html             ← Form input 36 indikator untuk klasifikasi
│   ├── hasil.html                   ← Visualisasi hasil & grafik probabilitas
│   ├── training.html                ← Halaman training model & confusion matrix
│   ├── laporan.html                 ← Riwayat hasil klasifikasi & export Excel
│   └── users.html / users_edit.html ← Manajemen akun Admin
│
├── data skripsi/
│   ├── erd.puml                     ← Diagram ERD (PlantUML)
│   ├── flowchart.puml               ← Flowchart alur sistem (PlantUML)
│   └── usecase.puml                 ← Use Case diagram (PlantUML)
│
└── MANUAL_BOOK.md                   ← Panduan lengkap penggunaan aplikasi
```

---

## Penjelasan Detail Setiap File & Komentar Kode

Setiap file kode dalam sistem ini telah dilengkapi dengan komentar bahasa Indonesia yang sangat detail agar alur logika pemrograman mudah dipahami. Berikut adalah penjelasan fungsi dari masing-masing file:

### 1. File Utama (Root)
*   **`main.go`**
    *   **Fungsi**: Berperan sebagai pusat kontrol (otak utama) aplikasi. Bertanggung jawab atas inisialisasi server web menggunakan framework Echo v4, pemetaan (routing) rute URL, manajemen middleware (logger, recover, session), verifikasi autentikasi user (Admin/Operator), dan pencarian port dinamis (mulai dari 8080-8090).
    *   **Logika Penting**:
        *   Melakukan pelatihan model Naive Bayes saat aplikasi pertama kali dijalankan dengan memuat data training dari SQLite.
        *   Menyediakan fungsi otomatis untuk membuka browser bawaan sistem saat server aktif.
        *   Memiliki mekanisme shutdown terintegrasi untuk menutup aplikasi secara aman (`os.Exit(0)`).

### 2. Modul `internal/classifier/` (Logika Algoritma)
*   **`naive_bayes.go`**
    *   **Fungsi**: Berisi implementasi penuh dari rumus algoritma Klasifikasi Naive Bayes dengan perbaikan numerik untuk menghindari bias.
    *   **Implementasi Formula**:
        1.  **Prior Probability** $P(C_k)$: Peluang awal kemunculan suatu kelas $C_k$ (dari 6 tingkat kesejahteraan):
            $$P(C_k) = \frac{\text{Jumlah Warga Kelas } C_k}{\text{Total Warga Latih}}$$
        2.  **Likelihood Probability** $P(X_i | C_k)$ dengan **Laplace Smoothing**: Untuk mencegah probabilitas bernilai nol jika suatu nilai atribut belum pernah muncul pada data latih:
            $$P(X_i = v | C_k) = \frac{count(X_i = v \text{ di kelas } C_k) + 1}{count(C_k) + |V|}$$
            *Dimana $|V|$ adalah total jumlah kategori pilihan jawaban pada indikator tersebut.*
        3.  **Log-Probability** (Pencegahan *Underflow*): Perkalian peluang desimal kecil dapat menghasilkan angka yang mendekati nol (underflow). Oleh karena itu, perkalian digantikan dengan penjumlahan logaritma:
            $$\ln P(C_k | X) = \ln P(C_k) + \sum_{i=1}^{n} \ln P(X_i | C_k)$$
        4.  **Normalisasi (Softmax)**: Mengubah kembali nilai eksponensial logaritma agar jumlah total peluang dari ke-6 kelas sama dengan $1.0$ ($100\%$).
*   **`indicators.go`**
    *   **Fungsi**: Mendefinisikan struktur data untuk 36 indikator kemiskinan (`IM1` sampai `IM36`) beserta pilihan jawabannya (A, B, C, D) yang dikelompokkan ke dalam 3 bagian besar:
        1.  **Kondisi Rumah** (Akses air minum, jenis dinding, luas lantai, toilet, dll.)
        2.  **Ekonomi Keluarga** (Penghasilan kepala keluarga, jumlah tanggungan, jaminan kesehatan, dll.)
        3.  **Aset & Fasilitas** (Kepemilikan AC, motor, mobil, HP, perhiasan, dll.)

### 3. Modul `internal/db/` (Manajemen Database)
*   **`db.go`**
    *   **Fungsi**: Mengelola semua komunikasi langsung dengan database SQLite (`data_skripsi.db`) menggunakan driver Go murni `modernc.org/sqlite`.
    *   **Daftar Skema Tabel**:
        *   `pengguna`: Menyimpan data username, hash password (menggunakan enkripsi bcrypt), dan role (Admin/Operator).
        *   `warga`: Menyimpan data profil kependudukan (NIK, KK, nama, alamat, kelurahan) serta status apakah data tersebut data latih atau klasifikasi real-time.
        *   `data_indikator`: Tabel transaksi yang menyimpan nilai pilihan jawaban dari 36 indikator warga.
        *   `hasil_klasifikasi`: Menyimpan hasil prediksi akhir model berupa nama kelas terbaik beserta struktur JSON data probabilitas lengkap.
    *   **Fitur CRUD**: Fungsi untuk menambah warga, seeding user admin pertama (`admin` / `admin123`), dan mengambil riwayat data latih.

### 4. Folder `templates/` (Antarmuka Web)
Seluruh file HTML menggunakan sistem modularisasi template bawaan Go (`html/template`) dengan format pewarisan layout:
*   **`base.html`**: Master layout (sidebar navigasi, topbar, profil user, tombol logout, dan tombol matikan aplikasi) yang menyelimuti halaman konten.
*   **`index.html`**: Dashboard visualisasi grafis interaktif dengan **Chart.js** yang mengambil data dari database via REST API internal.
*   **`warga.html`**: Tabel daftar warga terintegrasi dengan upload data format `.xlsx` (Excel) dan tombol hapus/edit.
*   **`warga_tambah.html` & `warga_edit.html`**: Form entry data profil kependudukan warga secara manual.
*   **`klasifikasi.html`**: Formulir kuesioner dinamis 36 indikator. Telah diperbaiki agar radio button berkelompok per ID indikator secara dinamis dengan parameter `name="{{$indicatorID}}"`.
*   **`hasil.html`**: Halaman pengumuman hasil kelas kesejahteraan beserta diagram batang probabilitas ke-6 kelas. Menyediakan opsi cetak print ramah kertas.
*   **`training.html`**: Dashboard evaluasi model yang melatih ulang model secara real-time dan menampilkan Confusion Matrix berukuran 6x6 untuk mengecek akurasi, presisi, recall, dan f1-score.
*   **`laporan.html`**: Halaman riwayat klasifikasi warga dan tombol download rekap dalam format Excel.
*   **`users.html` & `users_edit.html`**: Halaman untuk manajemen akun pengguna (tambah admin baru, edit password, hapus user).
*   **`login.html` & `register.html`**: Halaman mandiri untuk autentikasi sistem.

### 5. Folder `data skripsi/` (Diagram Sistem)
*   **`erd.puml`**: Rancangan relasi database SQLite (Entitas Pengguna, Warga, Data Indikator, dan Hasil Klasifikasi).
*   **`flowchart.puml`**: Alur proses jalannya program dari login, input warga, input indikator, klasifikasi Naive Bayes, hingga cetak laporan.
*   **`usecase.puml`**: Hak akses aktor Admin (bisa akses semua fitur termasuk training) vs Operator (hanya klasifikasi & melihat laporan).

---

## Cara Menjalankan Aplikasi

### Menjalankan dari Source Code (Development)

1.  **Persyaratan**:
    *   Pastikan Anda sudah menginstal **Go (Golang)** versi 1.20 ke atas.
    *   Jalankan perintah sinkronisasi library untuk pertama kali:
        ```bash
        go mod tidy
        ```

2.  **Langkah-Langkah**:
    *   **Inisialisasi Database (Jika file db kosong)**:
        ```powershell
        go run scripts/sync_db/sync_db.go
        ```
    *   **Menjalankan Server Web**:
        ```powershell
        go run main.go
        ```
    *   Buka browser secara otomatis atau akses secara manual melalui: `http://localhost:8080`

### Menjalankan Aplikasi Hasil Build (Produksi)

Jika Anda ingin langsung menjalankan tanpa menginstal compiler Go, Anda dapat mengklik ganda file binary executable yang sudah disediakan:
*   **`Klasifikasi-Warga-Randuagung.exe`**
    *   File ini merupakan kompilasi mandiri (self-contained binary) yang menyematkan (embed) seluruh folder `templates/`, `static/`, dan konfigurasi server di dalamnya.
    *   Cukup jalankan file tersebut, dan browser default sistem Anda akan langsung terbuka secara otomatis ke halaman login.

### Akun Default Awal
*   **Username**: `admin`
*   **Password**: `admin123`
*   **Role**: `Admin`

---

## Lisensi
Hak Cipta © 2026 Muhammad As'ad Muhibbin Akbar.
*Skripsi Program Studi Teknik Informatika - Kelurahan Randuagung*
