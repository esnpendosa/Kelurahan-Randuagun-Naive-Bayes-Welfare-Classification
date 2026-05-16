# Sistem Klasifikasi Tingkat Kesejahteraan Masyarakat - Kelurahan Randuagung

Aplikasi klasifikasi tingkat kesejahteraan masyarakat menggunakan algoritma **Naive Bayes** yang dibangun dengan bahasa pemrograman **Go (Golang)** dan database **SQLite**. Proyek ini dikembangkan untuk membantu Kelurahan Randuagung dalam menentukan klasifikasi ekonomi warga secara objektif berdasarkan 36 indikator kesejahteraan.

## Fitur Utama
- **Dashboard**: Visualisasi statistik kependudukan, akurasi model, dan distribusi kelas.
- **Data Warga**: Manajemen data kependudukan (CRUD) dan fitur import dari Excel.
- **Klasifikasi Baru**: Melakukan prediksi klasifikasi warga secara real-time.
- **Training Model**: Evaluasi performa algoritma Naive Bayes dengan Confusion Matrix.
- **Laporan**: Rekapitulasi hasil klasifikasi yang dapat diekspor ke format Excel.
- **Manajemen User**: Pengaturan akun akses untuk Admin dan Operator.

## Teknologi yang Digunakan
- **Backend**: Go (Golang)
- **Framework Web**: Echo v4
- **Database**: SQLite
- **Logika**: Naive Bayes (dengan Laplace Smoothing)
- **Frontend**: HTML5, Vanilla CSS, JavaScript

## Cara Menjalankan Aplikasi

1.  **Persyaratan**:
    - Pastikan Go sudah terinstal di sistem Anda.
    - Pastikan file `data_skripsi.db` ada di root folder (atau jalankan skrip sync untuk membuatnya).

2.  **Langkah-langkah**:
    ```powershell
    # 1. Jalankan sinkronisasi database (jika data kosong)
    go run scripts/sync_db/sync_db.go

    # 2. Jalankan aplikasi utama
    go run main.go
    ```
    Buka browser dan akses: `http://localhost:8080`

3.  **Akun Default**:
    - **Username**: `admin`
    - **Password**: `admin123`

## Struktur Folder
- `internal/classifier`: Logika inti Naive Bayes.
- `internal/db`: Fungsi interaksi database.
- `templates/`: Koleksi file tampilan HTML.
- `static/`: File CSS, JS, dan gambar.
- `scripts/`: Skrip utilitas untuk sinkronisasi dan seeding data.

## Lisensi
Hak Cipta © 2026 Muhammad As'ad Muhibbin Akbar.
v1.0.0 Desktop Web
