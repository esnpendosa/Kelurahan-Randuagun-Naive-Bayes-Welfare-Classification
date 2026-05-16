# 📘 Panduan Pengguna (Manual Book)
## Aplikasi Klasifikasi Kesejahteraan Ekonomi Keluarga - Kelurahan Randuagung

Aplikasi ini adalah sistem berbasis Desktop yang menggunakan algoritma **Naive Bayes** untuk menentukan klasifikasi kesejahteraan warga berdasarkan 36 indikator ekonomi.

---

## 🚀 1. Cara Menjalankan Aplikasi
Aplikasi ini sekarang tersedia dalam format desktop (.exe) yang portabel:

1.  Buka folder aplikasi.
2.  Cari file bernama **`Klasifikasi-Warga-Randuagung.exe`**.
3.  Klik dua kali (Double Click) pada file tersebut.
4.  **Browser Otomatis**: Tunggu beberapa detik, browser default Anda akan otomatis terbuka dan menampilkan halaman login.
5.  **Menutup Aplikasi**: Untuk mematikan aplikasi, klik tombol merah **"Keluar Aplikasi"** di bagian kanan atas layar.

---

## 🔐 2. Login & Hak Akses
Gunakan akun yang telah didaftarkan untuk masuk:
- **Username Default**: `admin`
- **Password Default**: `admin123`

### Peran (Roles):
- **Admin**: Akses penuh ke semua fitur (Manajemen Warga, Training Model, Manajemen User).
- **Operator**: Hanya bisa melakukan klasifikasi baru, melihat dashboard, dan mengunduh laporan.

---

## ✨ 3. Fitur Utama

### 📊 A. Dashboard
Halaman utama yang menampilkan ringkasan data kelurahan:
- **Statistik Cepat**: Total warga terdaftar dan total klasifikasi yang dilakukan.
- **Grafik Distribusi**: Visualisasi jumlah warga berdasarkan kelas kesejahteraan (Kelas 1 - 6).
- **Aktivitas Terakhir**: Log aktivitas sistem terbaru.

### 👥 B. Manajemen Data Warga (Admin Only)
Kelola database kependudukan sebelum dilakukan klasifikasi:
- **Tambah Warga**: Masukkan NIK, No KK, Nama, dan Alamat.
- **Import Excel**: Masukkan data warga secara massal menggunakan file Excel (.xlsx).
- **Edit/Hapus**: Memperbarui identitas warga jika ada perubahan.

### 🧠 C. Klasifikasi Baru
Fitur inti untuk memprediksi tingkat kesejahteraan:
1.  Pilih warga dari daftar dropdown.
2.  Isi **36 Indikator** ekonomi (Kondisi Rumah, Ekonomi, Aset, dll).
3.  Klik **"Proses Klasifikasi"**.
4.  Sistem akan menghitung probabilitas menggunakan rumus Naive Bayes.

### 📈 D. Visualisasi Hasil
Setelah klasifikasi, sistem menampilkan:
- **Kelas Terpilih**: Rekomendasi tingkat kesejahteraan (misal: "Miskin").
- **Grafik Probabilitas**: Menampilkan persentase kemiripan warga dengan 6 kategori kelas yang ada.

### 📋 E. Laporan & Export
- **Riwayat**: Melihat semua hasil klasifikasi yang pernah dilakukan.
- **Export Excel**: Unduh laporan rekapitulasi dalam format Excel untuk keperluan administrasi kantor.

### ⚙️ F. Training Model (Admin Only)
Menjaga kecerdasan sistem:
- **Data Latih**: Sistem memerlukan minimal 60 data latih (10 per kelas) agar akurat.
- **Confusion Matrix**: Visualisasi akurasi model untuk melihat seberapa tepat prediksi sistem dibandingkan data aktual.

---

## 📂 4. Penting Diketahui
- **Database**: File `data_skripsi.db` berisi semua data Anda. **Jangan menghapus file ini**, karena semua data warga dan hasil klasifikasi akan hilang.
- **Penyimpanan**: Aplikasi ini tidak memerlukan koneksi internet untuk berfungsi (Offline-ready).

---

## 🏷️ Kategori Kelas
1.  **Kelas 1**: Sangat Miskin
2.  **Kelas 2**: Miskin
3.  **Kelas 3**: Hampir Miskin
4.  **Kelas 4**: Rentan Miskin
5.  **Kelas 5**: Pas-pasan
6.  **Kelas 6**: Menengah ke Atas

---
*© 2026 Arina - Sistem Klasifikasi Kesejahteraan Kelurahan Randuagung*
