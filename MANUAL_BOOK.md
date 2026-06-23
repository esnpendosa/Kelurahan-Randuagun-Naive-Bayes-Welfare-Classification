# 📘 Panduan Pengguna (Manual Book)
## Sistem Klasifikasi Kesejahteraan Ekonomi Keluarga
### Kelurahan Randuagung — Berbasis Metode Naive Bayes

---

> **Untuk siapa panduan ini?**
> Panduan ini dibuat untuk petugas kelurahan yang akan menggunakan aplikasi ini sehari-hari. Tidak perlu paham teknologi — cukup ikuti langkah-langkah yang tertulis di sini.

---

## 🚀 Bagian 1: Cara Membuka Aplikasi

Aplikasi ini tidak perlu diinstal terlebih dahulu. Cukup ikuti langkah berikut:

1. Buka folder tempat file aplikasi tersimpan.
2. Cari file bernama **`Klasifikasi-Warga-Randuagung.exe`** (ikonnya seperti gambar jendela komputer).
3. **Klik dua kali** pada file tersebut.
4. Tunggu beberapa detik — browser (Chrome/Edge/Firefox) akan otomatis terbuka dan menampilkan halaman login.
5. Jika browser tidak langsung terbuka, tunggu sekitar 5 detik, lalu buka browser Anda sendiri dan ketik `http://localhost:8080` di kolom alamat.

> ⚠️ **Penting:** Jangan tutup jendela hitam (Command Prompt) yang muncul saat aplikasi berjalan. Jendela tersebut adalah "mesin" aplikasinya. Jika ditutup, aplikasi akan berhenti.

### Cara Menutup Aplikasi
Klik tombol merah **"Keluar Aplikasi"** di pojok kanan atas layar. Jangan langsung tutup browsernya saja.

---

## 🔐 Bagian 2: Login (Masuk ke Sistem)

Setelah browser terbuka, Anda akan melihat halaman login.

- **Username (Nama Pengguna):** `admin`
- **Password (Kata Sandi):** `admin123`

> 💡 **Tips:** Setelah pertama kali masuk, sebaiknya ganti password agar lebih aman. Caranya ada di **Bagian 7** panduan ini.

Klik tombol **"Masuk"** untuk melanjutkan. Anda akan dibawa ke halaman utama (Dashboard).

---

## 📊 Bagian 3: Dashboard (Halaman Utama)

Setelah login, Anda akan melihat tampilan seperti papan informasi yang berisi:

- **Total Warga Terdaftar** — Jumlah seluruh kepala keluarga yang sudah tercatat di sistem.
- **Total Klasifikasi** — Jumlah berapa kali proses penentuan kesejahteraan sudah dilakukan.
- **Grafik Distribusi Kelas** — Diagram batang yang menunjukkan berapa warga yang termasuk dalam setiap kelas kesejahteraan (Kelas 1 hingga 6).

Menu navigasi ada di sebelah kiri layar.

---

## 👥 Bagian 4: Manajemen Data Warga

Menu ini digunakan untuk **mencatat, mengubah, atau menghapus data kepala keluarga**.

### 4.1 Melihat Daftar Warga

1. Klik menu **"Data Warga"** di sebelah kiri.
2. Akan muncul tabel berisi daftar semua warga yang sudah terdaftar.
3. Warga yang bertanda **"Data Latih"** artinya data tersebut digunakan sebagai bahan pembelajaran sistem.

### 4.2 Menambah Warga Baru (Satu per Satu)

1. Klik tombol **"+ Tambah Warga Baru"** di pojok kanan atas tabel.
2. Isi formulir yang muncul:
   - **NIK** → Nomor Induk Kependudukan (16 digit, ada di KTP)
   - **Nomor KK** → Nomor Kartu Keluarga
   - **Nama Lengkap** → Nama Kepala Keluarga
   - **Alamat** → Nama jalan atau gang tempat tinggal
   - **RT / RW** → Nomor RT dan RW
   - **Kelurahan** → Nama kelurahan tempat tinggal
3. Klik tombol **"Simpan Data Warga"**.

### 4.3 Menambah Banyak Warga Sekaligus (Import Excel)

Jika Anda punya data banyak warga dalam file Excel, bisa diunggah sekaligus:

1. Siapkan file Excel (.xlsx) dengan format kolom: `NIK | Nama | Alamat | Kelurahan`
2. Di halaman Data Warga, klik tombol **"Pilih File"** di bagian Import.
3. Pilih file Excel Anda, lalu klik **"Import"**.
4. Data akan langsung masuk ke sistem secara otomatis.

> 📎 Tersedia file contoh bernama `Template_Import_Warga_filled.xlsx` di folder aplikasi sebagai panduan format yang benar.

### 4.4 Mengubah Data Warga

1. Temukan nama warga di tabel.
2. Klik tombol **"Edit"** (berwarna kuning) di sebelah namanya.
3. Ubah data yang perlu diperbaiki.
4. Klik **"Simpan Perubahan"**.

### 4.5 Menghapus Data Warga

1. Temukan nama warga di tabel.
2. Klik tombol **"Hapus"** (berwarna merah).
3. Akan muncul konfirmasi — klik **"OK"** untuk menghapus.

> ⚠️ **Perhatian:** Data yang dihapus tidak bisa dikembalikan. Pastikan Anda tidak salah memilih.

---

## 🧠 Bagian 5: Klasifikasi Baru (Menentukan Kelas Kesejahteraan)

Ini adalah fitur utama aplikasi — untuk **menentukan tingkat kesejahteraan sebuah keluarga** secara otomatis.

### Langkah-langkahnya:

**Langkah 1 — Pilih Warga:**
1. Klik menu **"Klasifikasi Baru"** di sebelah kiri.
2. Pada bagian **"Pilih Warga"**, klik kotak dropdown dan pilih nama warga yang ingin diklasifikasi.
   - Pastikan nama warga sudah terdaftar di Data Warga terlebih dahulu.

**Langkah 2 — Isi 36 Indikator:**
Formulir dibagi menjadi 3 bagian:
- 🏠 **Kondisi Rumah** (12 pertanyaan) — tentang kondisi fisik tempat tinggal
- 💰 **Ekonomi Keluarga** (12 pertanyaan) — tentang pendapatan, aset, dan beban keluarga
- 🛋️ **Aset & Fasilitas** (12 pertanyaan) — tentang barang-barang yang dimiliki keluarga

Untuk setiap pertanyaan, pilih **satu jawaban** yang paling sesuai dengan kondisi warga tersebut.

**Langkah 3 — Proses Klasifikasi:**
Setelah semua 36 pertanyaan dijawab, klik tombol **"Mulai Proses Klasifikasi"**.

Sistem akan langsung menghitung dan menampilkan hasilnya.

### Memahami Hasil Klasifikasi

Setelah diproses, Anda akan melihat:
- **Nama Kelas** yang ditentukan sistem (misalnya: "Miskin" atau "Hampir Miskin")
- **Grafik Batang** berisi persentase kemiripan warga tersebut dengan setiap kelas (Kelas 1–6). Batang yang paling panjang adalah kelas yang dipilih sistem.

### Kelas Kesejahteraan

| Kelas | Nama | Artinya |
|:---:|---|---|
| **Kelas 1** | Sangat Miskin | Keluarga dalam kondisi sangat sulit, perlu bantuan mendesak |
| **Kelas 2** | Miskin | Keluarga kekurangan, perlu bantuan pemerintah |
| **Kelas 3** | Hampir Miskin | Kondisi rentan, perlu dipantau agar tidak memburuk |
| **Kelas 4** | Rentan Miskin | Cukup stabil, tapi mudah terdampak jika ada masalah |
| **Kelas 5** | Pas-pasan | Kebutuhan sehari-hari terpenuhi, tidak lebih tidak kurang |
| **Kelas 6** | Menengah ke Atas | Keluarga sudah sejahtera dan mapan |

---

## 📋 Bagian 6: Laporan dan Export Excel

### Melihat Laporan

1. Klik menu **"Laporan"** di sebelah kiri.
2. Akan muncul tabel berisi seluruh riwayat hasil klasifikasi yang pernah dilakukan.

### Mengunduh Laporan ke Excel

1. Di halaman Laporan, klik tombol **"Export Excel"**.
2. File `.xlsx` akan otomatis terunduh ke komputer Anda.
3. File tersebut bisa dibuka dengan Microsoft Excel atau Google Sheets.
4. Isinya berupa rekap: NIK, Nama, Alamat, Hasil Klasifikasi, dan Tanggal.

---

## ⚙️ Bagian 7: Training & Evaluasi Model (Melatih Ulang Sistem)

Fitur ini digunakan untuk **melatih ulang sistem** menggunakan data latih dan **mengevaluasi performanya** menggunakan data uji.

### Cara Training dan Evaluasi Model:

1. Klik menu **"Training Model"** di sebelah kiri.
2. Sistem akan melatih model menggunakan data bertanda **"Data Latih"** (data_latih = 1).
3. Setelah model dilatih, sistem akan menguji kemampuannya dengan mengklasifikasikan data bertanda **"Data Uji"** (data_latih = 0) yang memiliki label kelas aktual.
4. Hasil pengujian dan evaluasi akan ditampilkan secara real-time, yaitu:
   - **Akurasi, Presisi, Recall, dan F1-Score** — Metrik untuk menilai seberapa baik model memprediksi kelas kesejahteraan yang sebenarnya (diambil dari rata-rata performa model terhadap data uji).
   - **Confusion Matrix** — Tabel perbandingan 6x6 antara kelas aktual data uji dengan kelas yang diprediksi oleh model.

> ⚠️ **Syarat:** Sistem membutuhkan minimal **60 data latih** (minimal 10 data untuk setiap kelas) agar proses training dapat dilakukan. Pastikan juga Anda telah mengunggah data uji dengan label kelas aktual untuk melakukan evaluasi performa model.

---

## 🔑 Bagian 8: Pengaturan Akun

Fitur ini digunakan oleh Admin untuk **mengubah atau menghapus kredensial akun** yang digunakan untuk mengakses sistem. Karena sistem hanya menggunakan satu akun Admin, Anda tidak perlu menambah akun baru.

### Melihat Informasi Akun

1. Klik menu **"Pengaturan Akun"** di sebelah kiri.
2. Anda akan melihat nama pengguna (username) akun Admin yang aktif.

### Mengubah Informasi Akun

1. Klik tombol **"Edit"** di sebelah nama akun Admin.
2. Ubah nama pengguna atau kata sandi jika diperlukan.
3. Klik **"Simpan Perubahan"**.

### Menghapus Akun

- Klik **"Hapus"** jika Anda ingin menghapus akun tersebut. Namun, harap berhati-hati karena jika akun dihapus, Anda harus membuat akun baru atau mereset database untuk dapat mengakses sistem kembali.

> 💡 **Tips Keamanan:** Ganti kata sandi bawaan (`admin123`) setelah pertama kali login. Gunakan kombinasi huruf dan angka yang kuat.

---

## 📂 Bagian 9: Hal-hal Penting yang Perlu Diketahui

### File Database
- Semua data warga dan hasil klasifikasi tersimpan di file bernama **`data_skripsi.db`**.
- **Jangan hapus atau pindahkan file ini** — semua data akan hilang jika file ini rusak atau terhapus.
- Disarankan untuk **menyalin (backup) file ini** secara rutin ke flashdisk atau tempat lain.

### Koneksi Internet
- Aplikasi ini **tidak memerlukan internet** untuk berjalan. Bisa digunakan meskipun tidak ada koneksi internet (offline).

### Jika Aplikasi Bermasalah
- Tutup aplikasi, tunggu beberapa detik, lalu buka kembali file `.exe`.
- Pastikan tidak ada aplikasi lain yang menggunakan port 8080 di komputer yang sama.

---

## ❓ Bagian 10: Pertanyaan yang Sering Ditanyakan (FAQ)

**Q: Apakah data warga aman?**
A: Ya. Semua data tersimpan di komputer lokal (tidak dikirim ke internet) dan hanya bisa diakses oleh orang yang punya username dan password.

**Q: Bagaimana jika saya lupa password?**
A: Jika Anda lupa kata sandi akun Admin tunggal Anda, silakan hubungi pengembang aplikasi untuk mereset kata sandi melalui modifikasi langsung pada file database `data_skripsi.db`.

**Q: Apakah hasil klasifikasi bisa salah?**
A: Sistem menggunakan metode statistik (Naive Bayes) sehingga hasilnya berupa rekomendasi. Akurasi sistem bergantung pada kualitas data latih yang dimasukkan. Semakin banyak dan akurat data latihnya, semakin tepat rekomendasinya.

**Q: Bolehkah satu warga diklasifikasi lebih dari sekali?**
A: Boleh. Setiap hasil klasifikasi akan tersimpan di riwayat dengan tanggal berbeda. Ini berguna jika kondisi warga berubah dari waktu ke waktu.

**Q: Apa perbedaan "Data Latih" dan "Data Uji"?**
A: **Data Latih** adalah data kependudukan (dengan flag data_latih = 1) yang digunakan oleh algoritma Naive Bayes untuk menghitung probabilitas Prior dan Likelihood. **Data Uji** adalah dataset terpisah (dengan flag data_latih = 0) yang digunakan khusus untuk mengevaluasi performa klasifikasi model (mengukur akurasi, presisi, recall, dan f1-score via Confusion Matrix).

---

*© 2026 Arina — Sistem Klasifikasi Kesejahteraan Kelurahan Randuagung*
*Dikembangkan sebagai bagian dari penelitian skripsi dengan metode Naive Bayes*
