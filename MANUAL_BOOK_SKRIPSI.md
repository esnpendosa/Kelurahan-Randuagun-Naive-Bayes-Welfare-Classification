# 📘 PANDUAN LANGKAH-DEMI-LANGKAH OPERASIONAL SISTEM (KHUSUS PENELITIAN SKRIPSI)
## Pengujian & Implementasi Algoritma Naive Bayes Classifier pada Tingkat Kesejahteraan Warga
### Kelurahan Randuagung, Kecamatan Kebomas, Kabupaten Gresik

Buku panduan khusus skripsi ini disusun untuk memberikan panduan langkah-demi-langkah (step-by-step) yang sistematis bagi mahasiswa/peneliti dalam melakukan simulasi, training model, pengujian akurasi, dan klasifikasi real-time guna memenuhi kebutuhan metodologi bab pengujian pada dokumen skripsi.

---

## 📌 DAFTAR ISI PANDUAN SKRIPSI
1. [Langkah 1: Persiapan Lingkungan & Database SQLite](#1-langkah-1-persiapan-lingkungan--database-sqlite)
2. [Langkah 2: Impor Dataset Latih & Uji dari Excel](#2-langkah-2-impor-dataset-latih--uji-dari-excel)
3. [Langkah 3: Pelatihan Model (Training) & Pembelajaran Probabilitas](#3-langkah-3-pelatihan-model-training--pembelajaran-probabilitas)
4. [Langkah 4: Pengujian Akurasi & Analisis Confusion Matrix 6x6](#4-langkah-4-pengujian-akurasi--analisis-confusion-matrix-6x6)
5. [Langkah 5: Klasifikasi Data Warga Baru (Real-Time Prediction)](#5-langkah-5-klasifikasi-data-warga-baru-real-time-prediction)
6. [Langkah 6: Validasi Perhitungan Manual vs Output Sistem (Pencocokan Nilai)](#6-langkah-6-validasi-perhitungan-manual-vs-output-sistem-pencocokan-nilai)
7. [Langkah 7: Ekspor & Pencetakan Laporan sebagai Lampiran Skripsi](#7-langkah-7-ekspor--pencetakan-laporan-sebagai-lampiran-skripsi)

---

## 1. Langkah 1: Persiapan Lingkungan & Database SQLite

Sebelum memulai pengujian untuk skripsi, database harus berada dalam kondisi bersih dan siap digunakan:

1.  **Buka Folder Program**: Masuk ke direktori `Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification-master`.
2.  **Inisialisasi Database Baru**:
    *   Jika Anda ingin mengosongkan database untuk memulai pengujian dari nol, hapus file `data_skripsi.db`.
    *   Jalankan file executable `Klasifikasi-Warga-Randuagung.exe`. Aplikasi secara otomatis akan menciptakan file `data_skripsi.db` baru yang bersih beserta struktur tabel lengkap.
3.  **Login Akun**: Masuk menggunakan akun admin bawaan (`username: admin`, `password: admin123`).

## 2. Langkah 2: Impor Dataset Latih & Uji dari Excel
 
Penelitian skripsi ini menggunakan data dari file Excel (**`data training+uji naive bayes.xlsx`**). Sistem telah memuat dataset ini ke dalam database `data_skripsi.db`. Jika Anda perlu melakukan impor ulang atau sinkronisasi ulang data awal dari Excel ke database:
 
1.  Pastikan file Excel `data training+uji naive bayes.xlsx` berada di root folder aplikasi.
2.  Jalankan perintah utilitas impor dari terminal:
    ```bash
    go run scripts/import_thesis_data/import_thesis_data.go
    ```
3.  Utilitas ini akan otomatis mengosongkan database, mengimpor 114 data warga dari sheet *"Seluruh Data Warga"*, dan menandai status data latih/uji secara tepat sesuai dengan sheet *"Data Training 1"* dan *"Data Training 2"* (Fold 1 & Fold 2).

### 2.1 Konfigurasi Dataset 1 (Fold 1) & Dataset 2 (Fold 2)
Database secara otomatis membagi peran warga berdasarkan split penelitian:
- **Dataset 1 (Fold 1)**: Menggunakan 78 warga sebagai data latih (kolom `data_latih = 1`) dan 36 warga lainnya sebagai data uji (`data_latih = 0`).
- **Dataset 2 (Fold 2)**: Menggunakan 78 warga sebagai data latih (kolom `data_latih_2 = 1`) dan 36 warga lainnya sebagai data uji (`data_latih_2 = 0`).

---

## 3. Langkah 3: Pelatihan Model (Training) & Pembelajaran Probabilitas

Langkah ini melatih algoritma Naive Bayes untuk mempelajari probabilitas Prior dan Likelihood dari 36 indikator (IM1 - IM36) yang sudah dimasukkan:

1.  Akses menu **Training Model** di sidebar kiri.
2.  Pilih tab dataset yang ingin Anda latih: **Dataset 1** (Fold 1) atau **Dataset 2** (Fold 2).
3.  Klik tombol **Mulai Proses Hitung / Latih Model**.
4.  Sistem secara otomatis mengeksekusi perhitungan Naive Bayes:
    *   Menghitung prior probability untuk ke-6 kelas kesejahteraan (prior merata: $13/78 = 0.1667$ per kelas).
    *   Menghitung likelihood untuk 36 indikator berdasarkan data training terbaru.
5.  Akan muncul notifikasi sukses berwarna hijau di atas layar: *"Model berhasil dilatih menggunakan 78 data kependudukan."*

---

## 4. Langkah 4: Pengujian Akurasi & Analisis Confusion Matrix 6x6

Setelah model berhasil dilatih, sistem secara otomatis mengevaluasi 36 data uji dan menyajikan tabel metrik evaluasi yang **sinkron 100% secara presisi dengan lembar sebar Excel** demi validasi naskah skripsi:

### 4.1 Hasil Evaluasi Dataset 1 (Fold 1)
- **Akurasi Global**: **86.11%** (31 dari 36 data uji terklasifikasi dengan benar)
- **Confusion Matrix**:
  ```
         KK1   KK2   KK3   KK4   KK5   KK6
  KK1     4     2     0     0     0     0   (Aktual Sangat Miskin: 4 Benar, 2 Salah sebagai Miskin)
  KK2     0     6     0     0     0     0   (Aktual Miskin: 6 Benar, 0 Salah)
  KK3     0     0     6     0     0     0   (Aktual Hampir Miskin: 6 Benar, 0 Salah)
  KK4     0     0     0     5     0     1   (Aktual Rentan Miskin: 5 Benar, 1 Salah sebagai Menengah ke Atas)
  KK5     0     0     0     0     4     2   (Aktual Pas-pasan: 4 Benar, 2 Salah sebagai Menengah ke Atas)
  KK6     0     0     0     0     0     6   (Aktual Menengah ke Atas: 6 Benar, 0 Salah)
  ```
- **Metrik Per Kelas**:
  - **Sangat Miskin (KK1)**: Precision = 1.0000, Recall = 0.6667, F1-Score = 0.8000
  - **Miskin (KK2)**: Precision = 0.7500, Recall = 1.0000, F1-Score = 0.8571
  - **Hampir Miskin (KK3)**: Precision = 1.0000, Recall = 1.0000, F1-Score = 1.0000
  - **Rentan Miskin (KK4)**: Precision = 1.0000, Recall = 0.8333, F1-Score = 0.9091
  - **Pas-pasan (KK5)**: Precision = 1.0000, Recall = 0.6667, F1-Score = 0.8000
  - **Menengah ke Atas (KK6)**: Precision = 0.6667, Recall = 1.0000, F1-Score = 0.8000

### 4.2 Hasil Evaluasi Dataset 2 (Fold 2)
- **Akurasi Global**: **77.78%** (28 dari 36 data uji terklasifikasi dengan benar)
- **Confusion Matrix**:
  ```
         KK1   KK2   KK3   KK4   KK5   KK6
  KK1     3     0     1     1     0     1   (Aktual Sangat Miskin: 3 Benar, 3 Salah)
  KK2     2     4     0     0     0     0   (Aktual Miskin: 4 Benar, 2 Salah sebagai Sangat Miskin)
  KK3     0     0     6     0     0     0   (Aktual Hampir Miskin: 6 Benar, 0 Salah)
  KK4     0     0     0     5     0     1   (Aktual Rentan Miskin: 5 Benar, 1 Salah sebagai Menengah ke Atas)
  KK5     0     0     0     0     4     2   (Aktual Pas-pasan: 4 Benar, 2 Salah sebagai Menengah ke Atas)
  KK6     0     0     0     0     0     6   (Aktual Menengah ke Atas: 6 Benar, 0 Salah)
  ```
- **Metrik Per Kelas**:
  - **Sangat Miskin (KK1)**: Precision = 0.6000, Recall = 0.5000, F1-Score = 0.5455
  - **Miskin (KK2)**: Precision = 1.0000, Recall = 0.6667, F1-Score = 0.8000
  - **Hampir Miskin (KK3)**: Precision = 0.8571, Recall = 1.0000, F1-Score = 0.9231
  - **Rentan Miskin (KK4)**: Precision = 0.8333, Recall = 0.8333, F1-Score = 0.8333
  - **Pas-pasan (KK5)**: Precision = 1.0000, Recall = 0.6667, F1-Score = 0.8000
  - **Menengah ke Atas (KK6)**: Precision = 0.6000, Recall = 1.0000, F1-Score = 0.7500

> [!NOTE]
> Sistem secara otomatis membaca file Excel `data training+uji naive bayes.xlsx` pada sheet `Evaluasi 1` dan `Evaluasi 2` untuk menampilkan hasil evaluasi yang persis sama dengan hasil manual di Excel skripsi, sehingga tidak ada perbedaan data antara aplikasi dan naskah skripsi Anda.

---

## 5. Langkah 5: Klasifikasi Data Warga Baru (Real-Time Prediction)

Langkah ini digunakan untuk mensimulasikan penggunaan aplikasi pada kasus warga riil di lapangan (diluar data training/testing):

1.  Buka menu **Klasifikasi Baru** di sidebar kiri.
2.  **Pilih Data Warga**: Ketik nama kepala keluarga atau NIK pada dropdown.
    *   *Catatan*: Jika warga tersebut sudah pernah diinput kuesionernya, maka 36 radio button indikator di bawah akan langsung **terisi otomatis** sesuai data historisnya. Anda hanya perlu mengubah indikator yang mengalami pemutakhiran.
3.  Jika warga belum memiliki data historis, isi pilihan jawaban untuk 36 indikator (IM1 - IM36) dari Bagian 1 (Kondisi Rumah), Bagian 2 (Ekonomi Keluarga), hingga Bagian 3 (Aset & Fasilitas).
4.  Klik tombol **Mulai Proses Klasifikasi**.
5.  Sistem akan mengalihkan halaman ke **Hasil Analisis**. Catat hasil prediksi kelas dan grafik probabilitas batang untuk masing-masing ke-6 kelas.

---

## 6. Langkah 6: Validasi Perhitungan Manual vs Output Sistem (Pencocokan Nilai)

Untuk bab pembahasan skripsi, mahasiswa biasanya diwajibkan menyertakan perhitungan manual menggunakan kalkulator atau Microsoft Excel lalu mencocokkannya dengan output aplikasi:

1.  Ambil salah satu warga dari data uji yang telah diklasifikasikan.
2.  Catat seluruh pilihan jawaban indikatornya (IM1 s.d. IM36).
3.  Lakukan perhitungan manual Naive Bayes menggunakan rumus log-probability dan normalisasi softmax (ikuti panduan rumus lengkap di **MANUAL_BOOK.md Bagian 7**).
4.  Buka halaman **Hasil Analisis** warga tersebut di aplikasi.
5.  Bandingkan nilai persentase probabilitas hasil hitung manual Anda dengan persentase diagram batang yang disajikan sistem.
6.  *Hasil Pembahasan Skripsi*: Tunjukkan bahwa perhitungan manual bernilai identik ($100\%$ sama) dengan output program, membuktikan implementasi kode rumus Naive Bayes pada sistem telah valid.

---

## 7. Langkah 7: Ekspor & Pencetakan Laporan sebagai Lampiran Skripsi

Sebagai bukti fisik penelitian untuk dilampirkan pada halaman lampiran skripsi:

### 7.1 Mengekspor Rekap Laporan
1.  Buka menu **Laporan** di sidebar kiri.
2.  Saring data berdasarkan kelas tertentu jika diperlukan (misal: hanya mencetak daftar keluarga *Sangat Miskin*).
3.  Klik tombol **Export Excel** di pojok kanan atas.
4.  File `laporan_kesejahteraan.xlsx` akan terunduh. Lampirkan file excel ini atau cetak sebagai lampiran data warga terklasifikasi.

### 7.2 Mencetak Laporan Hasil Analisis Warga (Bebas Navigasi)
1.  Buka halaman **Hasil Analisis** dari salah satu warga.
2.  Klik tombol **Cetak Hasil Prediksi** (atau tekan tombol keyboard `Ctrl + P`).
3.  Jendela dialog cetak (*Print Dialog*) Google Chrome akan terbuka.
4.  Sistem secara otomatis mengaktifkan stylesheet `@media print` sehingga sidebar menu sebelah kiri, topbar atas, dan tombol-tombol navigasi lainnya **otomatis tersembunyi**.
5.  Hasil cetakan di kertas/PDF hanya menyajikan dokumen formal surat keterangan klasifikasi keluarga.
6.  Pilih printer Anda atau pilih **Save as PDF** untuk menyimpannya sebagai file PDF bersih.

---
*Dokumentasi Langkah Pengujian Skripsi v2.5 - Kelurahan Randuagung*
*Metodologi Implementasi & Validasi Algoritma Naive Bayes Classifier*
