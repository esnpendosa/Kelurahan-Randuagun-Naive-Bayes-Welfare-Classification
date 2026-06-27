# 📘 BUKU PANDUAN OPERASIONAL PENGGUNA (MANUAL BOOK)
## Sistem Klasifikasi Kesejahteraan Ekonomi Keluarga Kelurahan Randuagung
### Berbasis Algoritma Naive Bayes Classifier & Go SSR Desktop Mode

---

Buku panduan ini disusun sebagai manual operasional dan panduan teknis bagi Petugas Kelurahan Randuagung untuk mengoperasikan Aplikasi Klasifikasi Kesejahteraan. Panduan ini menjelaskan alur operasional sistem secara menyeluruh, pemeliharaan database, serta disertai contoh perhitungan manual matematika algoritma Naive Bayes.

---

## 📌 DAFTAR ISI
1. [Bagian 1: Pengenalan Sistem & Arsitektur](#-bagian-1-pengenalan-sistem--arsitektur)
2. [Bagian 2: Instalasi & Menjalankan Aplikasi](#-bagian-2-instalasi--menjalankan-aplikasi)
   - [2.1 Persyaratan Perangkat Keras & Lunak](#21-persyaratan-perangkat-keras--lunak)
   - [2.2 Cara Membuka Aplikasi (Kiosk Mode)](#22-cara-membuka-aplikasi-kiosk-mode)
   - [2.3 Cara Menutup Aplikasi Secara Bersih](#23-cara-menutup-aplikasi-secara-bersih)
3. [Bagian 3: Penggunaan Fitur Dashboard](#-bagian-3-penggunaan-fitur-dashboard)
4. [Bagian 4: Pengelolaan Data Warga Kelurahan](#-bagian-4-pengelolaan-data-warga-kelurahan)
   - [4.1 Pencarian Real-Time & Penyaringan Instan](#41-pencarian-real-time--penyaringan-instan)
   - [4.2 Menambah & Mengedit Profil Kependudukan](#42-menambah--mengedit-profil-kependudukan)
   - [4.3 Manajemen Peran Model (Latih & Uji - Dataset 2)](#43-manajemen-peran-model-latih--uji---dataset-2)
   - [4.4 Alur Import & Export Excel](#44-alur-import--export-excel)
5. [Bagian 5: Simulasi Klasifikasi Warga (Form Cerdas)](#-bagian-5-simulasi-klasifikasi-warga-form-cerdas)
   - [5.1 Alur Kerja Fitur Auto-Fill Indikator](#51-alur-kerja-fitur-auto-fill-indikator)
   - [5.2 Kuesioner 36 Indikator Kesejahteraan](#52-kuesioner-36-indikator-kesejahteraan)
   - [5.3 Hasil Prediksi & Cetak Dokumen Resmi Bersih](#53-hasil-prediksi--cetak-dokumen-resmi-bersih)
6. [Bagian 6: Modul Evaluasi Model & Confusion Matrix](#-bagian-6-modul-evaluasi-model--confusion-matrix)
   - [6.1 Membaca Struktur Confusion Matrix 6x6](#61-membaca-struktur-confusion-matrix-6x6)
   - [6.2 Penjelasan Metrik Evaluasi Matematika](#62-penjelasan-metrik-evaluasi-matematika)
7. [Bagian 7: Contoh Simulasi Perhitungan Naive Bayes Secara Manual](#-bagian-7-contoh-simulasi-perhitungan-naive-bayes-secara-manual)
   - [7.1 Pembentukan Dataset Contoh](#71-pembentukan-dataset-contoh)
   - [7.2 Langkah Perhitungan Probabilitas Prior](#72-langkah-perhitungan-probabilitas-prior)
   - [7.3 Langkah Perhitungan Probabilitas Likelihood (Laplace Smoothing)](#73-langkah-perhitungan-probabilitas-likelihood-laplace-smoothing)
   - [7.4 Penggabungan Log-Probability & Normalisasi Softmax](#74-penggabungan-log-probability--normalisasi-softmax)
8. [Bagian 8: Pemeliharaan Database & Backup Data](#-bagian-8-pemeliharaan-database--backup-data)
9. [Bagian 9: FAQ & Panduan Solusi Masalah (Troubleshooting)](#-bagian-9-faq--panduan-solusi-masalah-troubleshooting)

---

## 💻 Bagian 1: Pengenalan Sistem & Arsitektur

Sistem Klasifikasi Kesejahteraan Kelurahan Randuagung dirancang untuk mengelompokkan keluarga ke dalam 6 tingkat kemiskinan/kesejahteraan:
1. **Sangat Miskin** (Tingkat kesulitan ekonomi tertinggi)
2. **Miskin**
3. **Hampir Miskin**
4. **Rentan Miskin**
5. **Pas-pasan**
6. **Menengah ke Atas** (Keluarga mapan/mandiri)

Aplikasi dibangun menggunakan **Golang** sebagai server lokal, **SQLite 3** sebagai wadah data offline-first, dan dijalankan langsung di browser lokal dalam **App Mode** tanpa menggunakan internet.

---

## 🛠️ Bagian 2: Instalasi & Menjalankan Aplikasi

### 2.1 Persyaratan Perangkat Keras & Lunak
*   **Sistem Operasi**: Windows 10 / 11 (direkomendasikan).
*   **RAM**: Minimal 2 GB.
*   **Browser**: Google Chrome terpasang di komputer (wajib untuk desktop shell kiosk).

### 2.2 Cara Membuka Aplikasi (Kiosk Mode)
1.  Buka folder program aplikasi Anda.
2.  Cari file bernama **`Klasifikasi-Warga-Randuagung.exe`**.
3.  **Klik dua kali** pada file tersebut.
4.  Jendela command prompt hitam akan muncul sebagai inisialisasi server lokal, disusul langsung oleh terbukanya aplikasi dalam tampilan window desktop bersih.
5.  Masuk menggunakan akun default:
    *   **Username**: `admin`
    *   **Password**: `admin123`

### 2.3 Cara Menutup Aplikasi Secara Bersih
1.  Klik menu **"Keluar"** di pojok kanan atas untuk mengamankan akun.
2.  Tutup jendela aplikasi. Server di latar belakang akan otomatis berhenti beroperasi.

---

## 📊 Bagian 3: Penggunaan Fitur Dashboard

Dashboard menyajikan statistik ringkas tingkat kabupaten/kelurahan:
*   **Total Warga**: Jumlah kepala keluarga terdaftar.
*   **Total Hasil Klasifikasi**: Jumlah warga yang telah diuji/diklasifikasi secara real-time.
*   **Visualisasi Chart.js**: Menampilkan diagram persentase sebaran tingkat kesejahteraan warga. Batang yang tinggi menandakan mayoritas status ekonomi di kelurahan tersebut.

---

## 👥 Bagian 4: Pengelolaan Data Warga Kelurahan

Menu ini digunakan untuk mencatat identitas kependudukan dasar warga sebelum melakukan klasifikasi.

### 4.1 Pencarian Real-Time & Penyaringan Instan
*   **Pencarian**: Ketik nama warga atau NIK pada kolom pencarian. Sistem akan menyaring isi tabel seketika tanpa perlu menekan tombol Cari atau memuat ulang halaman.
*   **Filter Dropdown**: Klik dropdown **"Filter Status Kesejahteraan"**, lalu pilih kategori (misal: *Sangat Miskin*). Tabel hanya akan menampilkan warga yang tergolong dalam kelas tersebut.

### 4.2 Menambah & Mengedit Profil Kependudukan
*   Klik **+ Tambah Warga Baru** untuk membuka form input manual.
*   Isi NIK (16 digit), No KK, Nama Lengkap (Kepala Keluarga), Alamat, RT/RW, dan Kelurahan.
*   Pilih opsi **Peran dalam Evaluasi Model (Dataset 2)**:
    *   **Sebagai Data Uji**: Default warga baru untuk dievaluasi prediksinya.
    *   **Sebagai Data Latih**: Menjadikan warga tersebut dasar model training pembelajaran.
*   Klik **Simpan Data Warga**.

### 4.3 Manajemen Peran Model (Latih & Uji - Dataset 2)
Jika Anda ingin mengganti peran warga di kemudian hari:
1.  Buka menu **Training Model** di sidebar kiri.
2.  Scroll ke bawah menuju **5. Manajemen Peran (Split 2)**.
3.  Cari nama warga, lalu klik tombol **Jadikan Data Latih** atau **Jadikan Data Uji** sesuai keperluan Anda secara dinamis.

### 4.4 Alur Import & Export Excel
*   **Import Excel**:
    *   Format kolom file Excel: `Kolom A: NIK`, `Kolom B: Nama`, `Kolom C: Alamat`, `Kolom D: Kelurahan`.
    *   Klik **Pilih File**, pilih file `.xlsx` Anda, lalu klik **Import**.
*   **Export Excel**: Klik **Export Excel** untuk mengunduh master data warga saat ini ke komputer dalam bentuk tabel `.xlsx`.

---

## 🧠 Bagian 5: Simulasi Klasifikasi Warga (Form Cerdas)

### 5.1 Alur Kerja Fitur Auto-Fill Indikator
Sistem telah dibekali dengan kecerdasan penyimpanan dinamis.
1.  Pilih menu **Klasifikasi Baru**.
2.  Pada kolom **Pilih Data Warga**, ketik atau pilih NIK/Nama Warga.
3.  Jika warga tersebut **sudah pernah diinput indikatornya di masa lalu**, sistem akan otomatis mengisi ke-36 radio button pilihan jawaban di bawahnya secara instan.
4.  Petugas hanya perlu meninjau dan mengubah pilihan yang berubah, lalu memproses prediksi.

### 5.2 Kuesioner 36 Indikator Kesejahteraan
Formulir dikelompokkan ke dalam 3 kartu besar:
*   **Kondisi Rumah** (IM1 - IM12): Memotret status tempat tinggal fisik (lantai terluas, jenis dinding, kepemilikan rumah, atap, air bersih, dll).
*   **Ekonomi Keluarga** (IM13 - IM24): Mengukur pendapatan bulanan, jumlah tanggungan, tingkat pendidikan, bahan bakar memasak, frekuensi konsumsi protein, serta fasilitas BPJS.
*   **Aset & Fasilitas** (IM25 - IM36): Mendata kepemilikan barang sekunder (motor, mobil, AC, emas, internet, laptop, mesin cuci, dll).

Pilihan jawaban terstandardisasi dari **A** (paling memprihatinkan/pra-sejahtera) hingga **D** (paling sejahtera/mandiri).

### 5.3 Hasil Prediksi & Cetak Dokumen Resmi Bersih
Sistem akan memproses Naive Bayes dan mengarahkan ke halaman hasil:
*   Menampilkan nama kelas akhir.
*   Grafik probabilitas persentase untuk masing-masing 6 kelas.
*   **Cara Cetak Rapi**:
    1. Klik tombol **Cetak Hasil Prediksi** (atau tekan `Ctrl + P`).
    2. Pada jendela opsi cetak browser, pastikan mencentang **Hilangkan Header & Footer** (untuk menyembunyikan URL situs dan tanggal otomatis di kertas).
    3. Simpan sebagai PDF atau cetak ke printer A4. Sidebar menu dan tombol sistem akan disembunyikan otomatis oleh media print CSS.

---

## ⚙️ Bagian 6: Modul Evaluasi Model & Confusion Matrix

Halaman **Training Model** digunakan untuk menguji kualitas kecerdasan buatan Naive Bayes.

### 6.1 Membaca Struktur Confusion Matrix 6x6
Tabel Confusion Matrix memetakan data aktual dengan data prediksi:
*   **Baris**: Kelas Aktual (Kondisi Lapangan Sebenarnya).
*   **Kolom**: Kelas Prediksi (Tebakan Algoritma).
*   **Sel Diagonal Kuning**: Jumlah data uji yang berhasil ditebak dengan benar oleh sistem (True Positive).
*   **Sel Luar Diagonal**: Jumlah data yang tebakannya meleset.

### 6.2 Penjelasan Metrik Evaluasi Matematika
*   **Akurasi**: Rasio prediksi benar secara keseluruhan terhadap total data uji.
*   **Recall per Kelas**: Seberapa andal sistem menebak kategori tertentu berdasarkan seluruh data aktual di kelas tersebut.
*   **Precision per Kelas**: Seberapa akurat tebakan sistem ketika menunjuk kategori tertentu berdasarkan total tebakan sistem.
*   **F1-Score**: Nilai rata-rata harmonis untuk melihat performa global tanpa dipengaruhi bias data.

---

## 🧮 Bagian 7: Contoh Simulasi Perhitungan Naive Bayes Secara Manual

Untuk memudahkan penulisan Bab IV Skripsi, berikut disajikan simulasi perhitungan matematika Naive Bayes secara lengkap.

### 7.1 Pembentukan Dataset Contoh
Misalkan kita memiliki data latih sebanyak **10 kepala keluarga** yang terbagi ke dalam **2 Kelas Kesejahteraan**:
*   $C_1$: **Miskin** (5 data)
*   $C_2$: **Menengah** (5 data)

Kita hanya akan menguji **2 Indikator** untuk kesederhanaan contoh:
*   $X_1$: Jenis Dinding Rumah (Pilihan: A = Bambu, B = Kayu, C = Tembok) $\rightarrow$ Jumlah kategori $|V| = 3$.
*   $X_2$: Sumber Air Minum (Pilihan: A = Sumur Sungai, B = Sumur Pompa, C = PAM) $\rightarrow$ Jumlah kategori $|V| = 3$.

#### Tabel Data Latih (Training):
| No | Nama | Jenis Dinding ($X_1$) | Sumber Air ($X_2$) | Kelas Aktual ($C$) |
|:---:|---|:---:|:---:|:---:|
| 1 | Warga 1 | A | A | Miskin |
| 2 | Warga 2 | A | B | Miskin |
| 3 | Warga 3 | B | A | Miskin |
| 4 | Warga 4 | B | B | Miskin |
| 5 | Warga 5 | C | A | Miskin |
| 6 | Warga 6 | B | C | Menengah |
| 7 | Warga 7 | C | B | Menengah |
| 8 | Warga 8 | C | C | Menengah |
| 9 | Warga 9 | C | C | Menengah |
| 10| Warga 10| B | C | Menengah |

---

### 7.2 Langkah Perhitungan Probabilitas Prior
Probabilitas awal masing-masing kelas:
$$P(C_1 = \text{Miskin}) = \frac{5}{10} = 0.5$$
$$P(C_2 = \text{Menengah}) = \frac{5}{10} = 0.5$$

---

### 7.3 Langkah Perhitungan Probabilitas Likelihood (Laplace Smoothing)
Kita akan mengklasifikasikan warga baru (**Warga Uji**) dengan karakteristik:
*   $X_1$ (Jenis Dinding) = **A (Bambu)**
*   $X_2$ (Sumber Air) = **C (PAM)**

Kita hitung nilai **Likelihood** dengan rumus **Laplace Smoothing** ($|V| = 3$):
$$P(X_i | C_k) = \frac{Count(X_i \text{ di kelas } C_k) + 1}{Count(C_k) + 3}$$

#### A. Untuk Kelas $C_1$ (Miskin):
*   Untuk $X_1 = A$:
    *   Jumlah dinding 'A' di kelas Miskin = 2 (Warga 1 dan Warga 2).
    *   $$P(X_1 = A | C_1 = \text{Miskin}) = \frac{2 + 1}{5 + 3} = \frac{3}{8} = 0.375$$
*   Untuk $X_2 = C$:
    *   Jumlah sumber air 'C' di kelas Miskin = 0.
    *   $$P(X_2 = C | C_1 = \text{Miskin}) = \frac{0 + 1}{5 + 3} = \frac{1}{8} = 0.125$$

#### B. Untuk Kelas $C_2$ (Menengah):
*   Untuk $X_1 = A$:
    *   Jumlah dinding 'A' di kelas Menengah = 0.
    *   $$P(X_1 = A | C_2 = \text{Menengah}) = \frac{0 + 1}{5 + 3} = \frac{1}{8} = 0.125$$
*   Untuk $X_2 = C$:
    *   Jumlah sumber air 'C' di kelas Menengah = 4 (Warga 6, 8, 9, 10).
    *   $$P(X_2 = C | C_2 = \text{Menengah}) = \frac{4 + 1}{5 + 3} = \frac{5}{8} = 0.625$$

---

### 7.4 Penggabungan Log-Probability & Normalisasi Softmax

#### Langkah 1: Penggabungan Nilai Logaritma Alami ($\ln$)
Untuk mencegah underflow nilai pecahan kecil, kita lakukan penjumlahan logaritma:
$$\text{Nilai Log}(C_k) = \ln P(C_k) + \ln P(X_1 | C_k) + \ln P(X_2 | C_k)$$

*   **Untuk Kelas $C_1$ (Miskin)**:
    $$\text{Nilai Log}(C_1) = \ln(0.5) + \ln(0.375) + \ln(0.125)$$
    $$\text{Nilai Log}(C_1) = (-0.6931) + (-0.9808) + (-2.0794) = -3.7533$$

*   **Untuk Kelas $C_2$ (Menengah)**:
    $$\text{Nilai Log}(C_2) = \ln(0.5) + \ln(0.125) + \ln(0.625)$$
    $$\text{Nilai Log}(C_2) = (-0.6931) + (-2.0794) + (-0.4700) = -3.2425$$

#### Langkah 2: Mengembalikan ke Nilai Eksponensial ($e^x$)
$$E_1 = e^{-3.7533} = 0.02344$$
$$E_2 = e^{-3.2425} = 0.03907$$

#### Langkah 3: Normalisasi Softmax (Menghasilkan Persentase)
$$\text{Total Eksponensial} = E_1 + E_2 = 0.02344 + 0.03907 = 0.06251$$

*   **Probabilitas Akhir Warga Uji Tergolong Miskin ($C_1$)**:
    $$P(C_1 | X) = \frac{0.02344}{0.06251} = 0.3750 \rightarrow \mathbf{37.5\%}$$

*   **Probabilitas Akhir Warga Uji Tergolong Menengah ($C_2$)**:
    $$P(C_2 | X) = \frac{0.03907}{0.06251} = 0.6250 \rightarrow \mathbf{62.5\%}$$

#### **Kesimpulan**:
Sistem akan memprediksi warga baru tersebut masuk ke dalam kelas **Menengah** dengan probabilitas sebesar **$62.5\%$** karena nilai $P(Menengah | X) > P(Miskin | X)$.

---

## 📂 Bagian 8: Pemeliharaan Database & Backup Data

*   Seluruh data disimpan pada file database **`data_skripsi.db`** di folder root aplikasi.
*   **Cara melakukan backup**:
    1.  Tutup aplikasi terlebih dahulu.
    2.  Salin (copy) file `data_skripsi.db` dan simpan ke flashdisk atau penyimpanan awan (Google Drive) secara berkala.
    3.  Jika terjadi kerusakan komputer, cukup letakkan file backup tersebut kembali ke folder program dengan nama yang sama untuk mengembalikan semua data Anda.

---

## ❓ Bagian 9: FAQ & Panduan Solusi Masalah (Troubleshooting)

*   **Q: Aplikasi tidak terbuka saat diklik double.**
    *   *A*: Pastikan browser Google Chrome telah terinstal di Windows Anda. Atau, buka browser Anda sendiri lalu ketik alamat `http://127.0.0.1:8082`.
*   **Q: Mengapa hasil akurasi model bernilai 0%?**
    *   *A*: Hal ini terjadi karena tidak ada data uji berlabel aktual di database. Daftarkan minimal beberapa warga dengan peran "Data Uji" dan isi kelas aktualnya untuk memicu evaluasi.
*   **Q: Bagaimana cara mengganti akun password admin bawaan?**
    *   *A*: Pilih menu **Pengaturan Akun**, lalu klik **Edit** pada akun admin, ketik sandi baru dan klik simpan.

---
*Manual Book v2.5 - Kelurahan Randuagung*
*Sistem Klasifikasi Kesejahteraan Naive Bayes*
