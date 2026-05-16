# Manual Penggunaan Aplikasi Klasifikasi Kesejahteraan Ekonomi Keluarga
## Kelurahan Randuagung

Aplikasi ini dirancang untuk membantu petugas Kelurahan Randuagung dalam menentukan tingkat kesejahteraan ekonomi keluarga menggunakan metode Naive Bayes berdasarkan 36 indikator ekonomi.

---

## 1. Persyaratan Sistem
- Windows 10/11
- Microsoft Edge WebView2 Runtime (biasanya sudah terinstal)
- Koneksi internet (hanya saat instalasi/update)

## 2. Cara Menjalankan Aplikasi
1. Buka terminal/CMD di folder aplikasi.
2. Jalankan perintah: `go run main.go`.
3. Buka browser dan akses: `http://localhost:8080`.
4. Login menggunakan akun yang diberikan oleh Admin.

## 3. Fitur Utama

### A. Dashboard
Menampilkan ringkasan data:
- Total warga yang terdaftar.
- Grafik distribusi kelas kesejahteraan (Kelas 1 - 6).
- Akurasi model Naive Bayes saat ini.

### B. Manajemen Data Warga
- **Cari Warga**: Gunakan kolom pencarian berdasarkan NIK atau Nama.
- **Tambah Warga**: Masukkan identitas dasar (NIK, No KK, Nama Kepala Keluarga, Alamat, RT/RW, Dusun).
- **Edit/Hapus**: Klik ikon pada tabel warga untuk mengubah atau menghapus data.

### C. Klasifikasi (Input 36 Indikator)
Ini adalah fitur inti untuk menentukan kelas kesejahteraan:
1. Pilih warga dari daftar atau klik tombol "Klasifikasi Baru".
2. Isi 36 indikator ekonomi yang dibagi menjadi 5 bagian:
   - **Kondisi Rumah**: Dinding, atap, lantai, air, dll.
   - **Data Kepala Keluarga**: Pendidikan, pekerjaan, penghasilan.
   - **Komposisi Keluarga**: Jumlah anggota, disabilitas, lansia.
   - **Akses & Pengeluaran**: Akses kesehatan, nilai aset.
   - **Kepemilikan Aset**: Motor, mobil, AC, dll.
3. Klik "Hitung Klasifikasi".
4. Hasil akan muncul berupa **Probabilitas** untuk setiap kelas dan **Kelas Akhir** yang direkomendasikan.

### D. Training Model (Admin)
Fitur ini digunakan untuk memperbarui kecerdasan aplikasi:
1. Klik menu "Training Model".
2. Klik tombol "Mulai Training".
3. Aplikasi akan mempelajari data warga yang sudah memiliki label kelas tetap.
4. Tinjau **Confusion Matrix** dan nilai **Accuracy, Precision, Recall, F1-Score** untuk memastikan model sudah akurat.

### E. Laporan
- Filter laporan berdasarkan tanggal, dusun, atau kelas kesejahteraan.
- Klik tombol **Export Excel** atau **PDF** untuk mengunduh laporan resmi kelurahan.

## 4. Kategori Kelas Kesejahteraan
- **Kelas 1**: Sangat Miskin
- **Kelas 2**: Miskin
- **Kelas 3**: Hampir Miskin
- **Kelas 4**: Rentan Miskin
- **Kelas 5**: Pas-pasan
- **Kelas 6**: Menengah ke Atas

---
*© 2026 Muhammad As'ad Muhibbin Akbar - Sistem Klasifikasi Naive Bayes*
