# 🎓 PANDUAN LENGKAP PERSIAPAN SIDANG SKRIPSI
## Sistem Klasifikasi Kesejahteraan Warga (Naive Bayes) - Kelurahan Randuagung

Dokumen ini disusun sebagai panduan utama Anda untuk menghadapi dosen penguji pada sidang skripsi. Penjelasan di bawah ini dirancang menggunakan **Bahasa Manusia (konseptual & mudah dipahami)** tanpa menghilangkan esensi akademis dan teknisnya.

---

## 📂 1. PETA STRUKTUR BERKAS APLIKASI (File Structure)

Penguji sering kali meminta Anda menunjukkan di mana letak kode program tertentu. Gunakan peta ini untuk memandu mereka dengan percaya diri:

```
[Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification-master]
 ├── main.go                       <-- Otak Utama (Web Server, Routing, & Controller)
 ├── data_skripsi.db               <-- Database SQLite (Menyimpan data login, warga, & hasil)
 ├── data training+uji naive bayes.xlsx <-- Dataset Excel Asli Skripsi (Sumbu Validasi)
 ├── templates/                    <-- Antarmuka Pengguna / User Interface (HTML)
 │    ├── base.html                <-- Kerangka Layout (Sidebar navigasi & tema visual)
 │    ├── index.html               <-- Dashboard (Statistik warga, kelas, & log aktivitas)
 │    ├── training.html            <-- Simulasi Model (Pengujian akurasi, confusion matrix)
 │    ├── klasifikasi.html         <-- Kuesioner Input (Formulir pengisian 36 indikator)
 │    └── ...
 ├── internal/                     <-- Modul Logika Sistem (Backend)
 │    ├── db/
 │    │    └── db.go               <-- Koneksi SQLite & Query Database (CRUD)
 │    └── classifier/
 │         ├── naive_bayes.go      <-- Mesin Perhitungan Naive Bayes (Prior & Likelihood)
 │         └── indicators.go       <-- Definisi 36 Indikator Kesejahteraan & Pilihan Jawaban
 └── scripts/                      <-- Script Pendukung (Import data Excel ke SQLite)
```

### Penjelasan Fungsi Komponen Utama (Gaya Sidang):
1. **`main.go` (Pusat Kendali)**:
   * *Fungsi*: Menjalankan web server lokal, mengatur rute URL (Routing), dan menjadi jembatan antara tampilan HTML dengan database serta rumus klasifikasi.
   * *Penjelasan ke Penguji*: *"Ini adalah file utama yang pertama kali dijalankan. Di sini saya mengatur alur data dari browser pengguna, memprosesnya di backend, lalu mengembalikan hasilnya ke layar."*
2. **`naive_bayes.go` (Mesin Naive Bayes)**:
   * *Fungsi*: Tempat di mana rumus matematika Naive Bayes ditulis dalam bentuk kode pemrograman (perhitungan peluang kelas dan perkalian likelihood).
   * *Penjelasan ke Penguji*: *"Di file inilah seluruh logika rumus Naive Bayes dihitung secara otomatis, mulai dari menghitung peluang awal kelas (Prior) hingga peluang kecocokan indikator warga (Likelihood)."*
3. **`db.go` (Jalur Data)**:
   * *Fungsi*: Berkomunikasi dengan database SQLite untuk membaca, menyimpan, dan mengubah data warga serta akun pengguna.
   * *Penjelasan ke Penguji*: *"File ini bertugas mengelola koneksi database secara aman, serta menangani proses penyimpanan hasil klasifikasi agar tidak hilang ketika aplikasi dimatikan."*

---

## 🗄️ 2. SKEMA & RELASI DATABASE (SQLite Database Schema)

Aplikasi ini menggunakan database relasional ringan **SQLite** (`data_skripsi.db`). Berikut adalah struktur tabel dan bagaimana mereka saling terhubung:

### 📊 Gambar Relasi Tabel (Entity-Relationship Diagram / ERD)

Berikut adalah visualisasi ERD database Anda. Anda dapat merujuk gambar ini saat ditanya tentang struktur database:

![Diagram Relasi Database (ERD)](data%20skripsi/database_erd.png)

### Penjelasan Struktur Tabel:

1. **Tabel `pengguna` (Data Akun & Hak Akses)**
   * *Fungsi*: Menyimpan data akun untuk masuk ke aplikasi.
   * *Kolom Utama*: `id` (Primary Key), `nama_pengguna`, `kata_sandi` (di-hash aman dengan **Bcrypt**), `peran` (`Admin` atau `Operator`), dan `avatar` (foto profil).

2. **Tabel `warga` (Data Profil Kependudukan)**
   * *Fungsi*: Menyimpan biodata warga Kelurahan Randuagung.
   * *Kolom Utama*: 
     * `id` (Primary Key).
     * `nik` (Nomor Induk Kependudukan) & `no_kk` (Nomor Kartu Keluarga).
     * `nama_lengkap`, `alamat`, `rt`, `rw`, `kelurahan`.
     * `data_latih` (Split 1): Bernilai `1` jika menjadi data training, `0` jika data uji.
     * `data_latih_2` (Split 2): Bernilai `1` jika data training, `0` jika data uji.
     * `label_kelas`: Kelas kesejahteraan aktual di lapangan (misal: "Sangat Miskin", "Pas-pasan").
     * `idpengguna` (Foreign Key $\rightarrow$ `pengguna.id`): Mencatat petugas/admin mana yang mendaftarkan warga tersebut.

3. **Tabel `data_indikator` (Detail Jawaban Kuesioner Warga)**
   * *Fungsi*: Menyimpan pilihan jawaban warga untuk ke-36 indikator.
   * *Kolom Utama*: 
     * `id` (Primary Key).
     * `warga_id` (Foreign Key $\rightarrow$ `warga.id`): Menghubungkan jawaban dengan warga yang bersangkutan.
     * `indikator_id`: Kode pertanyaan kuesioner (`IM1` sampai `IM36`).
     * `nilai`: Nilai jawaban kategorikal (`A`, `B`, `C`, atau `D`).

4. **Tabel `hasil_klasifikasi` (Riwayat Prediksi Model)**
   * *Fungsi*: Menyimpan hasil tebakan/prediksi model Naive Bayes beserta kalkulasi probabilitasnya.
   * *Kolom Utama*:
     * `id` (Primary Key).
     * `warga_id` (Foreign Key $\rightarrow$ `warga.id`): Hasil prediksi untuk warga tersebut.
     * `nama_kelas`: Nama kelas hasil prediksi dengan nilai probabilitas tertinggi (misal: "Rentan Miskin").
     * `probabilitas`: Format JSON berisi rincian nilai peluang dari ke-6 kelas kesejahteraan untuk digambar sebagai grafik batang pada halaman laporan.

---

## 🧮 3. MATEMATIKA NAIVE BAYES (Konsep & Contoh Hitung Manual)

Dosen penguji hampir pasti akan meminta Anda menjelaskan core matematika dari Naive Bayes. Berikut adalah cara menjelaskannya:

### A. Rumus Utama Naive Bayes

$$P(C_k | X) = \frac{P(C_k) \times P(X | C_k)}{P(X)}$$

Karena nilai penyebut $P(X)$ selalu sama untuk setiap kelas, maka kita cukup menghitung pembilangnya saja:

$$P(C_k | X) \propto P(C_k) \times \prod_{i=1}^{36} P(X_i | C_k)$$

*Di mana:*
* $C_k$ = Kelas kesejahteraan ke-$k$ (ada 6 kelas, dari Sangat Miskin hingga Menengah ke Atas).
* $X$ = Kumpulan jawaban kuesioner warga (sebanyak 36 jawaban, $Xi$).
* $P(C_k)$ = **Peluang Prior** (peluang awal kelas $C_k$ sebelum melihat jawaban kuesioner).
* $P(X_i | C_k)$ = **Peluang Likelihood** (peluang seseorang memiliki jawaban $X_i$ jika dia berada di kelas $C_k$).

---

### B. Langkah Perhitungan Sistem (Contoh Kasus Sederhana)

Bayangkan kita hanya menggunakan **2 indikator** (misal: IM2 = Bahan Dinding dan IM11 = Penghasilan) dengan **2 Kelas** (Miskin dan Pas-pasan) dari total **10 data training** untuk memudahkan simulasi coretan di papan tulis sidang.

#### **Data Training (10 Orang):**
* Kelas **Miskin** (6 orang):
  * 4 orang dindingnya "Bambu" (Jawaban A), 2 orang dindingnya "Tembok" (Jawaban B).
  * 5 orang penghasilannya "< Rp 1 Juta" (Jawaban A), 1 orang penghasilannya "> Rp 1 Juta" (Jawaban B).
* Kelas **Pas-pasan** (4 orang):
  * 1 orang dindingnya "Bambu" (Jawaban A), 3 orang dindingnya "Tembok" (Jawaban B).
  * 1 orang penghasilannya "< Rp 1 Juta" (Jawaban A), 3 orang penghasilannya "> Rp 1 Juta" (Jawaban B).

---

#### **Tahap 1: Menghitung Peluang Prior $P(C_k)$**
Peluang awal masing-masing kelas dari total 10 data latih:
* $P(\text{Miskin}) = \frac{6}{10} = 0.6$
* $P(\text{Pas-pasan}) = \frac{4}{10} = 0.4$

---

#### **Tahap 2: Menghitung Peluang Likelihood $P(X_i|C_k)$**
Kita hitung peluang masing-masing jawaban di dalam tiap kelas:
* **Pada Kelas Miskin:**
  * Peluang Dinding Bambu: $P(\text{Bambu} | \text{Miskin}) = \frac{4}{6} = 0.67$
  * Peluang Penghasilan < Rp 1 Juta: $P(\text{< Rp 1 Jt} | \text{Miskin}) = \frac{5}{6} = 0.83$
* **On Kelas Pas-pasan:**
  * Peluang Dinding Bambu: $P(\text{Bambu} | \text{Pas-pasan}) = \frac{1}{4} = 0.25$
  * Peluang Penghasilan < Rp 1 Juta: $P(\text{< Rp 1 Jt} | \text{Pas-pasan}) = \frac{1}{4} = 0.25$

---

#### **Tahap 3: Prediksi Data Warga Baru**
Ada warga baru bernama **Budi** dengan karakteristik:
* **Dinding: Bambu**
* **Penghasilan: < Rp 1 Juta**

Sistem akan menghitung peluang Budi masuk ke masing-masing kelas:

1. **Peluang Budi Masuk Kelas "Miskin":**
   $$\begin{aligned}
   P(\text{Miskin} | \text{Budi}) &= P(\text{Miskin}) \times P(\text{Bambu} | \text{Miskin}) \times P(\text{< Rp 1 Jt} | \text{Miskin}) \\
   &= 0.6 \times 0.67 \times 0.83 \\
   &= \mathbf{0.333}
   \end{aligned}$$

2. **Peluang Budi Masuk Kelas "Pas-pasan":**
   $$\begin{aligned}
   P(\text{Pas-pasan} | \text{Budi}) &= P(\text{Pas-pasan}) \times P(\text{Bambu} | \text{Pas-pasan}) \times P(\text{< Rp 1 Jt} | \text{Pas-pasan}) \\
   &= 0.4 \times 0.25 \times 0.25 \\
   &= \mathbf{0.025}
   \end{aligned}$$

#### **Tahap 4: Penentuan Hasil Akhir (Argmax)**
Sistem membandingkan nilai peluang akhir:
* Peluang Miskin ($0.333$) > Peluang Pas-pasan ($0.025$).
* Kesimpulan: **Budi diklasifikasikan sebagai warga "Miskin"**.

---

## 📈 4. EVALUASI MODEL & AKURASI (K-Fold Cross Validation)

Penelitian skripsi Anda menggunakan metode **K-Fold Cross Validation (K=2)** untuk membuktikan keandalan model.

### Mengapa menggunakan 2-Fold Cross Validation?
Data penelitian berjumlah **114 warga**. Data ini dibagi menjadi **2 bagian (Fold / Split) yang sama besar** (masing-masing 57 data) untuk menghindari bias pengujian.
* **Pengujian Skenario 1 (Fold 1 / Dataset 1)**:
  * **Data Latih (Training)**: 78 data warga.
  * **Data Uji (Testing)**: 36 data warga.
  * **Hasil Akurasi**: **86,11%** (31 dari 36 data uji berhasil ditebak dengan benar oleh sistem).
* **Pengujian Skenario 2 (Fold 2 / Dataset 2)**:
  * **Data Latih (Training)**: 78 data warga.
  * **Data Uji (Testing)**: 36 data warga (posisi data uji dan latih ditukar dari Fold 1).
  * **Hasil Akurasi**: **77,78%** (28 dari 36 data uji berhasil ditebak dengan benar).
* **Akurasi Rata-rata**: **81,94%** (membuktikan bahwa model klasifikasi ini sangat layak dan akurat untuk digunakan di Kelurahan Randuagung).

---

## 💡 5. TIPS DEMO APLIKASI SAAT SIDANG (Langkah Uji Coba)

Jika dosen penguji meminta Anda mendemokan program, ikuti langkah tenang berikut:

1. **Jalankan Web Server**: Klik dua kali file `build.bat` atau jalankan `main.exe`.
2. **Login ke Sistem**: Buka Chrome pada alamat `http://localhost:8080`, lalu masuk menggunakan akun default:
   * **Username**: `admin`
   * **Password**: `admin123`
3. **Menu Dashboard**: Tunjukkan statistik data warga yang terbagi berdasarkan 6 kelas kesejahteraan.
4. **Uji Akurasi Model (Menu Model/Training)**:
   * Tunjukkan tabel evaluasi ke-2 Fold (Dataset 1 & Dataset 2).
   * Klik tombol **"Latih Model"**. Jelaskan bahwa saat tombol diklik, sistem menghitung ulang ribuan nilai likelihood berdasarkan data latih yang ada di database.
   * Tunjukkan **Confusion Matrix** (tabel 6x6 actual vs predicted) untuk memperlihatkan di mana letak kesalahan prediksi model.
5. **Klasifikasi Baru (Menu Klasifikasi)**:
   * Pilih nama warga baru yang ingin diuji.
   * Isi kuesioner 36 indikator secara acak/cepat.
   * Tekan tombol **"Klasifikasikan"**.
   * Sistem akan mengarahkan ke halaman hasil klasifikasi yang menampilkan **Grafik Batang Probabilitas** dari ke-6 kelas serta status kesejahteraan warga tersebut.

---

## 💬 6. BANK PERTANYAAN DOSEN PENGUJI & CARA MENJAWABNYA

Gunakan panduan jawaban ini agar Anda terdengar menguasai materi, percaya diri, dan profesional:

### ❓ P1: Kenapa Anda tidak menggunakan "Laplace Smoothing" dalam perhitungan Naive Bayes Anda?
* **Jawaban Konseptual (Bahasa Manusia)**:
  > *"Pada skripsi ini, saya tidak menggunakan Laplace Smoothing karena dataset training yang digunakan (78 data) sudah mencakup semua variasi jawaban indikator (A, B, C, D) untuk seluruh kelas kesejahteraan. Tidak ada indikator penting yang bernilai 0 kemunculannya dalam data training resmi kami. Selain itu, peniadaan Laplace Smoothing bertujuan agar hasil perhitungan sistem di program Go ini **100% cocok secara presisi** dengan perhitungan manual di lembar sebar (Excel) skripsi yang diajukan."*
* **Penjelasan Teknis**:
  > *"Sistem kami sudah menangani kasus pembagian nol di fungsi `Prediksi` dengan melakukan pengecekan `p == 0`. Jika nilai peluang mencapai 0 pada suatu kelas, sistem akan langsung menghentikan perkalian (*early break*) dan menetapkan peluang kelas tersebut sebagai 0 secara aman tanpa menyebabkan program crash."*

---

### ❓ P2: Apa maksud dari asumsi "Naive" (independensi) pada algoritma ini?
* **Jawaban Konseptual**:
  > *"Kata 'Naive' (naif/polos) berarti algoritma ini menganggap setiap indikator dari ke-36 pertanyaan kuesioner bersifat **mandiri dan tidak saling mempengaruhi satu sama lain** dalam menentukan kelas kesejahteraan warga. Contohnya, kepemilikan AC (IM27) dianggap tidak ada hubungannya dengan tingkat penghasilan kepala keluarga (IM11) menurut perhitungan rumus."*
* **Kenapa tetap dipakai?**:
  > *"Meskipun di dunia nyata indikator-indikator tersebut mungkin saling berkaitan, asumsi independensi ini sengaja digunakan untuk menyederhanakan perhitungan matematika yang sangat rumit menjadi perkalian peluang yang sederhana. Terbukti, meskipun berasumsi 'naive', hasil akurasi model kami masih sangat tinggi yaitu mencapai **86,11%**."*

---

### ❓ P3: Bagaimana cara sistem Anda menangani angka peluang yang sangat kecil seperti `2.34e-12` (notasi ilmiah)?
* **Jawaban Konseptual**:
  > *"Karena kita mengalikan 36 nilai peluang desimal (semuanya di bawah angka 1), hasil perkalian akhirnya pasti akan menghasilkan pecahan desimal yang sangat panjang mendekati nol. Jika menggunakan tipe data desimal biasa, komputer akan mengalami kesalahan pembulatan (*underflow*) dan menganggap nilainya 0 mutlak."*
* **Solusi Program**:
  > *"Untuk mengatasi hal ini, program kami menggunakan tipe data **`float64`** (presisi ganda) yang mampu menampung angka sangat kecil hingga puluhan digit di belakang koma. Selain itu, kami membuat fungsi pemformatan angka bernama `FormatScientific` agar nilai peluang sangat kecil tersebut ditampilkan dalam format notasi ilmiah yang rapi dan mudah dibaca oleh dosen penguji, mirip seperti tampilan di Excel skripsi."*

---

### ❓ P4: Kenapa Anda memilih menggunakan bahasa pemrograman Go (Golang) dan database SQLite?
* **Jawaban Konseptual (Golang)**:
  > *"Saya menggunakan Golang karena bahasa ini sangat cepat dalam hal eksekusi program dan hemat memori. Golang juga menghasilkan satu file executable mandiri (`main.exe` / `Klasifikasi-Warga-Randuagung.exe`) sehingga aplikasi dapat langsung dijalankan di komputer kelurahan tanpa perlu menginstal runtime tambahan seperti Python atau Node.js."*
* **Jawaban Konseptual (SQLite)**:
  > *"Saya menggunakan SQLite karena ia berbentuk berkas database tunggal (`data_skripsi.db`) yang tertanam langsung di dalam aplikasi. Pihak kelurahan tidak perlu menyewa server database MySQL atau PostgreSQL yang rumit. Proses backup data juga sangat mudah, cukup dengan menyalin file `.db` tersebut."*

---

### ❓ P5: Dari mana asal usul 36 indikator kuesioner yang Anda gunakan dalam skripsi ini?
* **Jawaban Konseptual**:
  > *"Ke-36 indikator tersebut merupakan parameter resmi yang diadopsi dari standar penilaian kesejahteraan sosial milik Kelurahan Randuagung, yang juga merujuk pada variabel DTKS (Data Terpadu Kesejahteraan Sosial) Kementerian Sosial RI dan kriteria BPS. Indikator ini mencakup 3 aspek utama: Kondisi Rumah Tinggal (IM1-IM9), Keadaan Ekonomi Keluarga (IM10-IM21), serta Kepemilikan Aset & Fasilitas (IM22-IM36)."*

---

### ❓ P6: Apa perbedaan antara Data Latih (Training) dan Data Uji (Testing) dalam aplikasi Anda?
* **Jawaban Konseptual**:
  > *"**Data Latih (Training)** adalah data warga yang sudah diketahui kelas kesejahteraannya dari lapangan, digunakan oleh sistem untuk mempelajari pola peluang (menghitung prior dan likelihood). Sedangkan **Data Uji (Testing)** adalah data warga yang digunakan untuk mengetes kepintaran model; sistem akan menebak kelas warga tersebut berdasarkan 36 indikatornya tanpa melihat label aslinya terlebih dahulu, lalu hasil tebakan dicocokkan dengan label asli untuk menghitung persentase akurasi."*

---

**Semoga Sukses Sidang Skripsinya! Tetap Tenang, Kuasai Angka-angka Utama Anda, dan Jawab dengan Yakin! 🎓👍**
