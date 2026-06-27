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

---

## 2. Langkah 2: Impor Dataset Latih & Uji dari Excel

Penelitian skripsi ini menggunakan data dari file Excel (`klasifikasi naive bayes tambahan data.xlsx`). Berikut cara memuatnya ke dalam sistem:

### 2.1 Impor Data Warga
1.  Pilih menu **Data Warga** di sidebar kiri.
2.  Pada panel **Import Data Warga (Excel)**, klik tombol **Choose File** / **Pilih File**.
3.  Pilih file Excel dataset skripsi Anda (`klasifikasi naive bayes tambahan data.xlsx`).
4.  Klik tombol **Import**.
5.  Daftar warga akan muncul di dalam tabel utama.

### 2.2 Penentuan Split Peran (Training vs Testing)
Sistem memisahkan warga ke dalam **Data Latih (Training)** dan **Data Uji (Testing)**.
1.  Untuk **Dataset 2** (komposisi kustom):
    *   Secara bawaan, semua warga hasil import Excel berstatus sebagai **Data Uji (Testing)**.
    *   Buka menu **Training Model** di sidebar kiri.
    *   Scroll ke bawah pada tab **"5. Manajemen Peran (Split 2)"**.
    *   Pilihlah 78 warga yang akan dijadikan data latih (masing-masing 13 warga untuk setiap ke-6 kelas kesejahteraan), lalu klik tombol **"Jadikan Data Latih"** di samping nama warga bersangkutan.
    *   Pilihlah 36 warga lainnya untuk dijadikan data uji (masing-masing 6 warga untuk setiap ke-6 kelas kesejahteraan) dengan membiarkan perannya tetap sebagai **"Data Uji"**.

---

## 3. Langkah 3: Pelatihan Model (Training) & Pembelajaran Probabilitas

Langkah ini melatih algoritma Naive Bayes untuk mempelajari probabilitas Prior dan Likelihood dari 36 indikator (IM1 - IM36) yang sudah dimasukkan:

1.  Akses menu **Training Model** di sidebar kiri.
2.  Pilih dataset yang ingin Anda latih (misalnya **Dataset 2**).
3.  Klik tombol **Mulai Proses Hitung / Latih Model**.
4.  Sistem akan mengeksekusi fungsi `modelNB.Latih()` di background:
    *   Menghitung jumlah kemunculan ($Count$) pilihan jawaban A, B, C, D untuk setiap indikator IM1 - IM36 pada masing-masing kelas.
    *   Menerapkan rumus **Laplace Smoothing** untuk mengantisipasi nilai nol pada data baru.
5.  Akan muncul notifikasi sukses berwarna hijau di atas layar: *"Model berhasil dilatih menggunakan 78 data kependudukan."*

---

## 4. Langkah 4: Pengujian Akurasi & Analisis Confusion Matrix 6x6

Setelah model berhasil dilatih, sistem akan menguji performa model tersebut terhadap data uji secara otomatis dan menyajikan matriks evaluasi untuk Bab IV Skripsi:

1.  Periksa tabel **Confusion Matrix 6x6** yang dirender di halaman Training.
2.  **Membaca Kebenaran Prediksi**:
    *   Lihat sel berwarna kuning (diagonal dari pojok kiri atas ke pojok kanan bawah). Angka di dalam sel kuning menunjukkan jumlah data uji yang **berhasil diprediksi secara tepat**.
    *   Lihat sel di luar warna kuning. Angka tersebut adalah data uji yang **salah diprediksi** (misal: aktualnya *Sangat Miskin* tapi diprediksi oleh sistem sebagai *Miskin*).
3.  **Mencatat Metrik Evaluasi**:
    *   **Akurasi**: Persentase total tebakan benar global. Catat angka akurasi ini (misal: $88.89\%$).
    *   **Recall per Kelas**: Catat kemampuan sistem mendeteksi kelas aktual tertentu (tabel sebelah kanan matrix).
    *   **Precision per Kelas**: Catat ketepatan tebakan sistem per kelas (tabel di bawah matrix).
    *   **F1-Score**: Rata-rata harmonis global dari presisi dan recall.
4.  Dokumentasikan screenshot tabel Confusion Matrix ini ke dalam laporan skripsi Anda sebagai bukti validasi keakuratan sistem.

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
