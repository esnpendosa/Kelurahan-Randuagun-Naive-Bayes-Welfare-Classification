# Panduan Penggunaan Sistem (Manual Book)
## Sistem Klasifikasi Kesejahteraan Kelurahan Randuagung

### 1. Memulai Aplikasi
1. Jalankan aplikasi melalui terminal dengan perintah `go run main.go`.
2. Buka browser dan ketik alamat `http://localhost:8080`.
3. Masukkan **Username** dan **Password** Anda (Default: admin/admin123).

### 2. Dashboard
Halaman ini menampilkan ringkasan data:
- **Total Data Warga**: Jumlah seluruh warga di database.
- **Warga Terklasifikasi**: Jumlah warga yang sudah melalui proses prediksi.
- **Akurasi Model**: Persentase keakuratan algoritma berdasarkan data latih (1.277 data).
- **Grafik Distribusi**: Perbandingan jumlah warga di setiap tingkatan kelas (1-6).

### 3. Data Warga
Digunakan untuk mengelola profil penduduk:
- **Cari Warga**: Gunakan kolom pencarian untuk mencari berdasarkan NIK atau Nama.
- **Tambah Warga**: Klik tombol "+ Tambah Warga Baru" untuk memasukkan data manual.
- **Import Excel**: Unggah file Excel untuk memasukkan data warga secara massal.
- **Reset Data**: Menghapus seluruh data warga (Hati-hati: Tindakan ini permanen).

### 4. Klasifikasi Baru (Prediksi)
Langkah-langkah melakukan klasifikasi:
1. Pilih nama warga dari dropdown (Data warga harus sudah ada di menu Data Warga).
2. Isi 36 indikator (IM1 - IM36) berdasarkan kondisi riil warga.
3. Klik tombol **"Hitung Klasifikasi"**.
4. Sistem akan menampilkan hasil kelas dan grafik probabilitas untuk masing-masing tingkatan.

### 5. Training Model
Halaman ini ditujukan untuk evaluasi teknis:
- Klik **"Mulai Pelatihan Ulang"** untuk melatih model jika ada data latih baru.
- Lihat **Confusion Matrix** untuk mengetahui di kelas mana model sering terjadi kesalahan prediksi.

### 6. Laporan
Halaman untuk merekap hasil:
- Menampilkan daftar warga yang sudah diklasifikasi beserta tanggalnya.
- Klik **"Export Excel"** untuk mengunduh laporan dalam format file spreadsheet (.xlsx).

### 7. Manajemen User
(Khusus Admin)
- Digunakan untuk menambah atau mengedit akun petugas (Operator/Admin).
- Pastikan kata sandi aman dan peran (Role) sesuai dengan tanggung jawab petugas.

---
*Kontak Bantuan: Tim IT Kelurahan Randuagung*
