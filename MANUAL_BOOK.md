# 📘 PANDUAN PENGGUNA (MANUAL BOOK)
## Sistem Klasifikasi Kesejahteraan Ekonomi Keluarga
### Kelurahan Randuagung — Berbasis Metode Naive Bayes Go SSR

---

> **Untuk siapa panduan ini?**
> Panduan ini disusun untuk memudahkan Petugas Kelurahan Randuagung dalam mengoperasikan aplikasi Klasifikasi Kesejahteraan. Di dalamnya memuat panduan lengkap dari instalasi, pengelolaan data, hingga evaluasi akurasi model.

---

## 📌 DAFTAR ISI
1. [Bagian 1: Cara Membuka & Menutup Aplikasi](#-bagian-1-cara-membuka--menutup-aplikasi)
2. [Bagian 2: Autentikasi Login](#-bagian-2-autentikasi-login)
3. [Bagian 3: Dashboard & Ringkasan Statistik](#-bagian-3-dashboard--ringkasan-statistik)
4. [Bagian 4: Manajemen Data Warga (CRUD & Filter Instan)](#-bagian-4-manajemen-data-warga-crud--filter-instan)
   - [4.1 Pencarian & Filter Kesejahteraan Instan](#41-pencarian--filter-kesejahteraan-instan)
   - [4.2 Tambah Warga & Penentuan Peran Evaluasi (Dataset 2)](#42-tambah-warga--penentuan-peran-evaluasi-dataset-2)
   - [4.3 Import & Export Excel](#43-import--export-excel)
5. [Bagian 5: Klasifikasi Baru (Prediksi Otomatis)](#-bagian-5-klasifikasi-baru-prediksi-otomatis)
   - [5.1 Fitur Auto-Fill Indikator](#51-fitur-auto-fill-indikator)
   - [5.2 Melakukan Klasifikasi](#52-melakukan-klasifikasi)
   - [5.3 Hasil Analisis & Cetak Dokumen Rapi](#53-hasil-analisis--cetak-dokumen-rapi)
6. [Bagian 6: Training & Hasil Evaluasi Model (Model Excel-Style)](#-bagian-6-training--hasil-evaluasi-model-model-excel-style)
   - [6.1 Memahami Confusion Matrix 6x6](#61-memahami-confusion-matrix-6x6)
   - [6.2 Memahami Metrik Performa (Recall, Precision, F1-Score)](#62-memahami-metrik-performa-recall-precision-f1-score)
   - [6.3 Manajemen Peran Secara Dinamis](#63-manajemen-peran-secara-dinamis)
7. [Bagian 7: Laporan Rekapitulasi & Cetak Dokumen](#-bagian-7-laporan-rekapitulasi--cetak-dokumen)
8. [Bagian 8: Pengaturan Pengguna (Manajemen Akun)](#-bagian-8-pengaturan-pengguna-manajemen-akun)
9. [Bagian 9: Pemeliharaan Database & Penyelesaian Masalah (FAQ)](#-bagian-9-pemeliharaan-database--penyelesaian-masalah-faq)

---

## 🚀 Bagian 1: Cara Membuka & Menutup Aplikasi

Aplikasi ini bersifat portabel dan dirancang berjalan secara lokal sebagai program desktop native.

### 1.1 Cara Membuka Aplikasi
1. Buka folder tempat program aplikasi diletakkan.
2. Cari file biner **`Klasifikasi-Warga-Randuagung.exe`** (ikon jendela biru komputer).
3. **Klik ganda (double-click)** pada file tersebut.
4. Jendela Command Prompt (layar hitam) akan muncul sejenak, diikuti dengan terbukanya jendela aplikasi mandiri (Chrome App Mode) pada halaman Login.
5. Jika Chrome App Mode tidak terbuka secara otomatis, buka Google Chrome secara manual dan ketik alamat: `http://127.0.0.1:8082`

> ⚠️ **Penting:** Jangan tutup jendela Command Prompt hitam yang menyala di background komputer. Jendela tersebut merupakan server aplikasi lokal yang memproses seluruh algoritma dan database. Jika ditutup, aplikasi akan langsung mati.

### 1.2 Cara Menutup Aplikasi
Untuk menghentikan server secara bersih:
1. Klik menu **"Keluar"** di pojok kanan atas untuk keluar dari sesi akun.
2. Tutup jendela Chrome App Mode. Server lokal akan secara otomatis mati dalam beberapa saat atau Anda dapat menutup Command Prompt hitam yang menyala.

---

## 🔐 Bagian 2: Autentikasi Login

Sebelum masuk to sistem, pengguna wajib melakukan autentikasi demi keamanan data warga.

- **Username Bawaan**: `admin`
- **Password Bawaan**: `admin123`

Petugas kelurahan dapat mengganti nama pengguna dan password bawaan ini di menu **Pengaturan Akun** sesaat setelah berhasil masuk demi keamanan data kelurahan.

---

## 📊 Bagian 3: Dashboard & Ringkasan Statistik

Setelah login, halaman utama (**Dashboard**) menyajikan ringkasan data analitis kelurahan:
- **Total Warga**: Jumlah kepala keluarga yang saat ini tersimpan di database.
- **Total Klasifikasi**: Jumlah riwayat warga yang telah diproses klasifikasinya secara real-time.
- **Total Data Latih (Dataset 1)**: Jumlah basis data training awal.
- **Grafik Distribusi Kesejahteraan**: Chart interaktif dinamis yang menyajikan komposisi tingkat kesejahteraan warga di Kelurahan Randuagung.

---

## 👥 Bagian 4: Manajemen Data Warga (CRUD & Filter Instan)

Menu **Data Warga** adalah pusat pengelolaan profil kependudukan kelurahan.

### 4.1 Pencarian & Filter Kesejahteraan Instan
Sistem kini dilengkapi dengan fitur filter multi-kriteria tanpa memuat ulang halaman (*zero reload*):
* **Pencarian Nama/NIK**: Ketik nama kepala keluarga atau NIK di kotak pencarian, daftar warga akan menyaring secara otomatis saat Anda mengetik (*as you type*).
* **Dropdown Filter Kelas**: Pilih salah satu kelas kesejahteraan (misal: **Sangat Miskin**, **Miskin**, dll.). Tabel akan menyaring baris warga yang sesuai dengan label tersebut secara real-time.

### 4.2 Tambah Warga & Penentuan Peran Evaluasi (Dataset 2)
Ketika Anda menambahkan warga secara manual:
1. Klik **+ Tambah Warga Baru** di kanan atas.
2. Isi NIK (16 digit), No KK, Nama Lengkap, Alamat, RT/RW, dan Kelurahan.
3. Di bagian bawah, terdapat pilihan **"Peran dalam Evaluasi Model (Dataset 2)"**:
   - **Sebagai Data Uji (Testing)**: Menjadikan data warga ini sebagai bahan pengujian akurasi model di Dataset 2.
   - **Sebagai Data Latih (Training)**: Menjadikan data warga ini sebagai basis pembelajaran Naive Bayes di Dataset 2.
4. Klik **Simpan Data Warga**.

### 4.3 Import & Export Excel
* **Import Data**: Masukkan banyak data warga sekaligus dengan mengunggah file Excel (.xlsx). Pilih file pada bagian form import lalu klik **Import**.
* **Export Data**: Klik **Export Excel** untuk mengunduh seluruh master data warga ke dalam format tabel Excel.

---

## 🧠 Bagian 5: Klasifikasi Baru (Prediksi Otomatis)

Ini adalah fitur inti yang digunakan untuk menentukan kelas tingkat kesejahteraan warga.

### 5.1 Fitur Auto-Fill Indikator
Sistem telah dilengkapi dengan pendeteksi data historis.
* Saat Anda memilih nama warga dari dropdown kuesioner, sistem akan memanggil API database di background.
* Jika warga tersebut **sudah pernah memiliki data indikator**, semua radio button 36 indikator di formulir bawah akan **terisi otomatis** sesuai data terakhir yang tersimpan.
* Anda tidak perlu mengisi dari nol. Cukup periksa kembali, ubah indikator yang mengalami perubahan, lalu proses klasifikasi.

### 5.2 Melakukan Klasifikasi
1. Pilih nama warga.
2. Lengkapi 36 indikator kemiskinan yang terbagi atas:
   - **Kondisi Rumah** (IM1 - IM12)
   - **Ekonomi Keluarga** (IM13 - IM24)
   - **Aset & Fasilitas** (IM25 - IM36)
3. Klik tombol **Mulai Proses Klasifikasi**.

### 5.3 Hasil Analisis & Cetak Dokumen Rapi
Setelah diklik, Anda akan diarahkan ke halaman **Hasil Analisis** yang menampilkan:
- Hasil kelas kesejahteraan (misal: **Miskin**, **Pas-pasan**, dsb).
- Grafik batang probabilitas kecocokan terhadap 6 kelas.
- **Cetak Laporan Bersih**: Klik tombol **Cetak Hasil Prediksi** (atau tekan `Ctrl + P`). Sistem secara otomatis menyembunyikan sidebar kiri, tombol navigasi, dan header atas sehingga hasil cetakan/PDF berupa dokumen resmi hasil analisis keluarga yang bersih, formal, dan rapi.

---

## ⚙️ Bagian 6: Training & Hasil Evaluasi Model (Model Excel-Style)

Halaman ini digunakan oleh Admin untuk menguji performa model matematika Naive Bayes. Terdapat 2 dataset yang disediakan (Dataset 1 dan Dataset 2).

### 6.1 Memahami Confusion Matrix 6x6
Tabel ini membandingkan data aktual (baris) dengan hasil prediksi model (kolom) untuk data uji:
- **Diagonal Kuning**: Menunjukkan jumlah tebakan benar (True Positive - TP).
- **Kolom Jumlah (Kanan)**: Total aktual data uji per kelas.
- **Baris Jumlah (Bawah)**: Total prediksi model per kelas.
- **Pojok Kanan Bawah (Hijau)**: Total data uji global (misal: 36 data).

### 6.2 Memahami Metrik Performa (Recall, Precision, F1-Score)
Sistem menghitung performa dengan format desimal presisi tinggi ala Microsoft Excel:
* **Akurasi**: Persentase kebenaran prediksi global ($\frac{\text{Jumlah Benar}}{\text{Total Data Uji}}$).
* **Tabel Recall**: Mengukur kemampuan model dalam menemukan kembali data aktual per kelas ($\text{Recall} = \frac{\text{TP}}{\text{TP} + \text{FN}}$) dengan baris **Rata-Rata Recall** di bagian bawah.
* **Tabel Precision**: Mengukur ketepatan prediksi model terhadap data yang ditebak per kelas ($\text{Precision} = \frac{\text{TP}}{\text{TP} + \text{FP}}$) dengan baris **Rata-Rata Presisi** di bagian bawah.
* **F1-Score**: Rata-rata harmonis global ($2 \times \frac{\text{Presisi} \times \text{Recall}}{\text{Presisi} + \text{Recall}}$).

### 6.3 Manajemen Peran Secara Dinamis
Pada tab **5. Manajemen Peran (Split 2)**, Anda dapat meninjau peran setiap warga pada Dataset 2. Anda dapat mengubah peran warga dari **Uji** ke **Latih** secara dinamis hanya dengan sekali klik tombol di tabel.

---

## 📋 Bagian 7: Laporan Rekapitulasi & Cetak Dokumen

Menu **Laporan** menyajikan seluruh rekapitulasi data hasil klasifikasi warga.
* **Filter Kategori**: Menyaring tabel laporan berdasarkan kategori kesejahteraan.
* **Export Excel**: Mengunduh rekap laporan dalam format `.xlsx`.
* **Print Ramah Kertas**: Tekan `Ctrl + P` pada halaman ini. Halaman akan dicetak sebagai dokumen bersih tanpa menyertakan sidebar, topbar, dan tombol-tombol filter.

---

## 🔑 Bagian 8: Pengaturan Pengguna (Manajemen Akun)

Menu **Pengaturan Akun** digunakan untuk mengelola kredensial admin dan operator:
1. Klik menu **Pengaturan Akun**.
2. Anda dapat melihat daftar pengguna terdaftar.
3. Klik **Edit** pada user admin untuk mengganti username atau mengganti password default `admin123`.

---

## 📂 Bagian 9: Pemeliharaan Database & Penyelesaian Masalah (FAQ)

### 9.1 Pemeliharaan Database
- Database aplikasi ini berbentuk file tunggal **`data_skripsi.db`** di root folder.
- **Backup Rutin**: Disarankan untuk menyalin file `data_skripsi.db` ini ke flashdisk/cloud drive secara berkala. Jika komputer Anda mengalami kendala hardware, data Anda tetap aman.

### 9.2 Pertanyaan yang Sering Ditanyakan (FAQ)
* **Q: Mengapa hasil cetak di browser masih menampilkan sidebar menu?**
  * *A*: Pastikan Anda menekan tombol **"Cetak Hasil Prediksi"** di halaman hasil, atau refresh halaman terlebih dahulu sebelum menekan `Ctrl + P` agar stylesheet cetak termuat dengan benar.
* **Q: Mengapa metrik evaluasi bernilai 0% setelah saya klik Mulai Proses Hitung?**
  * *A*: Periksa kembali apakah data uji yang didefinisikan di manajemen peran sudah memiliki label kelas aktual di database. Jika tidak ada data uji berlabel aktual, model tidak memiliki pembanding untuk menghitung akurasi.
* **Q: Apakah aplikasi ini membutuhkan internet?**
  * *A*: Tidak. Aplikasi berjalan $100\%$ offline secara lokal di komputer kelurahan.

---
*Manual Book v2.0 - Sistem Klasifikasi Kesejahteraan Kelurahan Randuagung*
*Dikembangkan oleh Tim Developer Skripsi Teknik Informatika - Kelurahan Randuagung*
