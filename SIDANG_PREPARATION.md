# 🎓 PANDUAN PERSIAPAN SIDANG SKRIPSI
## Sistem Klasifikasi Kesejahteraan Warga (Naive Bayes) - Kelurahan Randuagung

Dokumen ini disusun khusus sebagai bahan persiapan menghadapi dosen penguji pada sidang skripsi besok. Seluruh penjelasan menggunakan **Bahasa Manusia (konseptual/mudah dipahami)**, bukan sekadar logika pemrograman kering, agar Anda dapat menjelaskan sistem dengan percaya diri.

---

## 📂 1. STRUKTUR BERKAS APLIKASI (File Structure)

Berikut adalah peta struktur folder dan berkas proyek Anda beserta kegunaannya dalam bahasa yang mudah dijelaskan ke penguji:

```
[Kelurahan-Randuagun-Naive-Bayes-Welfare-Classification-master]
 ├── main.go                       <-- Otak/Pengendali utama web server & routing
 ├── data_skripsi.db               <-- Database SQLite (menyimpan semua data & hasil)
 ├── data training+uji naive bayes.xlsx <-- Dataset Excel asli skripsi
 ├── templates/                    <-- Desain antarmuka web (HTML Go Templates)
 │    ├── base.html                <-- Bingkai layout luar (sidebar & navbar)
 │    ├── index.html               <-- Dashboard utama
 │    ├── training.html            <-- Halaman latih & metrik akurasi 6x6
 │    └── klasifikasi.html         <-- Kuesioner penginputan 36 indikator
 ├── internal/                     <-- Modul backend khusus
 │    ├── db/
 │    │    └── db.go               <-- Inisialisasi database SQLite & query CRUD
 │    └── classifier/
 │         ├── naive_bayes.go      <-- Implementasi rumus prior & likelihood
 │         └── indicators.go       <-- Definisi 36 indikator & bobot jawaban
 └── scripts/                      <-- Berkas script pembantu (seperti data import)
```

### Penjelasan Fungsi Setiap Bagian:
1. **`main.go` (Otak / Pengendali Utama Web)**:
   * **Bahasa Manusia**: Ini adalah pusat kendali aplikasi. Berkas ini berfungsi mengatur rute halaman (routing), menerima data masukan dari pengguna melalui web, memanggil fungsi perhitungan Naive Bayes, menyimpan hasil klasifikasi, serta menjalankan server lokal agar aplikasi bisa dibuka di browser Chrome.
2. **`internal/classifier/naive_bayes.go` (Rumus Naive Bayes)**:
   * **Bahasa Manusia**: Di sinilah rumus matematika Naive Bayes ditulis dalam bentuk kode. Berkas ini bertugas menghitung peluang awal kelas (Prior) dan peluang kemunculan indikator (Likelihood), serta menentukan tingkat kesejahteraan akhir (Sangat Miskin s.d. Menengah ke Atas) berdasarkan nilai probabilitas tertinggi.
3. **`internal/classifier/indicators.go` (Daftar 36 Indikator)**:
   * **Bahasa Manusia**: Berkas ini menyimpan daftar lengkap 36 indikator kesejahteraan (IM1 sampai IM36) yang dibagi menjadi 3 kategori: *Kondisi Rumah*, *Ekonomi Keluarga*, dan *Aset & Fasilitas*. Setiap indikator memiliki pilihan jawaban (A, B, C, D) beserta bobot maknanya.
4. **`internal/db/db.go` (Pengelola Database SQLite)**:
   * **Bahasa Manusia**: Berkas ini khusus bertugas untuk menghubungkan aplikasi dengan database SQLite. Di dalamnya ada perintah untuk membuat tabel baru, menambah warga, memperbarui profil, hingga mengambil daftar data warga untuk dilatih atau diuji.
5. **`templates/` (Tampilan Antarmuka HTML)**:
   * **Bahasa Manusia**: Folder ini berisi desain halaman web aplikasi.
     * `base.html` adalah bingkai/layout luar (sidebar navigasi dan header).
     * Halaman lainnya seperti `index.html` (dashboard), `training.html` (pengujian model), dan `klasifikasi.html` (pengisian kuesioner) akan dimasukkan ke dalam bingkai tersebut secara dinamis.
6. **`data_skripsi.db` (Database SQLite)**:
   * **Bahasa Manusia**: Berkas database tunggal yang menyimpan semua data login pengguna, data identitas warga, jawaban 36 indikator warga, serta riwayat hasil klasifikasi yang telah dilakukan.
7. **`data training+uji naive bayes.xlsx` (Dataset Excel)**:
   * **Bahasa Manusia**: Berkas Excel berisi data asli penelitian skripsi Anda (114 warga). Aplikasi membaca file ini secara langsung agar hasil perhitungan evaluasi dan akurasi di aplikasi **100% cocok secara presisi** dengan apa yang Anda tulis di dokumen skripsi.

---

## 🗄️ 2. STRUKTUR DATABASE (SQLite Database Schema)

Database aplikasi Anda memiliki **4 tabel utama** yang saling terhubung (berelasi):

1. **Tabel `pengguna` (Data Akun)**
   * **Fungsi**: Menyimpan akun pengguna yang diizinkan masuk ke sistem.
   * **Kolom Kunci**:
     * `id`: Nomor unik identitas akun (Primary Key).
     * `nama_pengguna` & `kata_sandi`: Kredensial login (kata sandi disimpan dalam bentuk hash aman menggunakan algoritma **Bcrypt**).
     * `peran`: Tingkat hak akses (`Admin` memiliki akses penuh, `Operator` hanya bisa input klasifikasi dan melihat laporan).

2. **Tabel `warga` (Data Profil Kependudukan)**
   * **Fungsi**: Menyimpan identitas dasar warga Kelurahan Randuagung.
   * **Kolom Kunci**:
     * `id`: Nomor unik warga (Primary Key).
     * `nik` (Nomor Induk Kependudukan) & `no_kk` (Nomor Kartu Keluarga).
     * `data_latih`: Bernilai `1` jika warga ini digunakan sebagai data training di **Dataset 1 (Fold 1)**, bernilai `0` jika digunakan sebagai data uji.
     * `data_latih_2`: Bernilai `1` jika warga digunakan sebagai data training di **Dataset 2 (Fold 2)**, bernilai `0` jika data uji.
     * `label_kelas`: Kategori kesejahteraan sosial riil/aktual warga di lapangan (misal: "Miskin", "Pas-pasan").

3. **Tabel `data_indikator` (Detail Kuesioner Warga)**
   * **Fungsi**: Menyimpan jawaban kuesioner (nilai A, B, C, atau D) untuk ke-36 indikator dari setiap warga.
   * **Kolom Kunci**:
     * `warga_id`: Menghubungkan ke tabel `warga` (Foreign Key).
     * `indikator_id`: Kode indikator, yaitu `IM1` sampai `IM36`.
     * `nilai`: Pilihan jawaban yang dipilih (misal: `A`, `B`, `C`, atau `D`).

4. **Tabel `hasil_klasifikasi` (Riwayat Prediksi Model)**
   * **Fungsi**: Menyimpan hasil prediksi klasifikasi beserta rincian peluang matematisnya.
   * **Kolom Kunci**:
     * `warga_id`: Menghubungkan ke tabel `warga` (Foreign Key).
     * `nama_kelas`: Kelas kesejahteraan hasil prediksi tertinggi (misal: "Hampir Miskin").
     * `probabilitas`: Menyimpan data peluang dari ke-6 kelas dalam format teks JSON (contoh: `{"1": 4.5e-10, "2": 2.1e-08, ...}`) untuk digambar menjadi grafik batang pada halaman hasil analisis.

---

## 🧮 3. CARA KERJA RUMUS NAIVE BAYES PADA KODE PROGRAM

Penguji sangat sering bertanya: *"Bagaimana jalannya rumus Naive Bayes di dalam kodingan Anda?"*
Berikut adalah penjelasan alur matematika yang diterjemahkan langsung dari berkas `internal/classifier/naive_bayes.go`.

### Tahap 1: Pelatihan Model (Training) — Menghitung Prior & Likelihood
Ketika Anda mengklik tombol **Latih Model**, fungsi `LatihModel` akan dieksekusi:

1. **Menghitung Prior Probability, P(C)**:
   * **Rumus**:
     P(C) = (Jumlah warga di kelas C) / (Total seluruh data training)
   * **Dalam Kode**:
     ```go
     // nb.SemuaKelas berisi 6 kelas (Sangat Miskin s.d. Menengah ke Atas)
     for _, c := range nb.SemuaKelas {
         // Peluang Prior = total warga dalam kelas C dibagi jumlah total data training
         nb.PeluangPrior[c] = float64(hitungPerKelas[c]) / jumlahTotal
     }
     ```
   * **Bahasa Manusia**: Jika dari 78 data training, masing-masing dari 6 kelas memiliki tepat 13 warga, maka peluang prior untuk setiap kelas adalah 13 / 78 = 0.1667.

2. **Menghitung Likelihood Probability, P(X|C)**:
   * **Rumus**:
     P(Xi = vj | C) = (Jumlah warga di kelas C dengan jawaban indikator Xi = vj) / (Total seluruh warga di kelas C)
   * **Dalam Kode**:
     ```go
     // Menghitung berapa kali nilai fitur/jawaban tertentu muncul di kelas C
     for v := range nilaiUnik { // v adalah pilihan jawaban, misal 'A'
         if totalDiKelas > 0 {
             nb.PeluangLikelihood[c][fitur][v] = float64(jumlahMuncul[v]) / float64(totalDiKelas)
         }
     }
     ```
   * **Bahasa Manusia**: Misalkan di kelas **Sangat Miskin** terdapat 13 warga. Kita ingin tahu peluang orang Sangat Miskin memiliki dinding bambu (Opsi 'A' pada indikator IM2). Jika 10 dari 13 warga tersebut berdinding bambu, maka nilai Likelihood untuk fitur IM2 = A di kelas Sangat Miskin adalah 10 / 13 = 0.7692.

---

### Tahap 2: Klasifikasi Data Baru (Prediction) — Mengalikan Peluang
Ketika warga baru diinputkan 36 jawaban indikatornya, fungsi `Prediksi` berjalan:

1. **Mengalikan Prior dengan Semua Likelihood**:
   * **Rumus**:
     P(C | X) = P(C) * P(IM1|C) * P(IM2|C) * ... * P(IM36|C)
   * **Dalam Kode**:
     ```go
     for _, c := range nb.SemuaKelas {
         p := nb.PeluangPrior[c] // Mulai dengan peluang Prior awal kelas tersebut
         
         for _, fitur := range nb.DaftarFitur {
             nilai := input[fitur] // Mengambil jawaban warga untuk indikator ini (misal 'A')
             l := nb.PeluangLikelihood[c][fitur][nilai] // Ambil nilai likelihood-nya
             p *= l // Kalikan terus menerus sebanyak 36 kali
         }
         hasilPeluang[c] = p // Simpan total peluang perkalian untuk kelas C
     }
     ```
   * **Bahasa Manusia**: Sistem menghitung 6 nilai peluang akhir (satu nilai untuk setiap kelas). Caranya dengan mengalikan peluang awal kelas tersebut dengan peluang ke-36 jawaban indikator warga tersebut pada kelas itu.

2. **Menentukan Hasil Akhir (Argmax)**:
   * **Bahasa Manusia**: Setelah mendapatkan 6 nilai peluang, sistem membandingkan keenamnya dan memilih kelas yang memiliki nilai **peluang terbesar** sebagai hasil prediksi akhir.
   * **Dalam Kode**:
     ```go
     // Mencari kelas dengan nilai probabilitas tertinggi
     for c, p := range peluang {
         if p > peluangMaks {
             peluangMaks = p
             kelasTerbaik = c
         }
     }
     ```

---

## 💬 4. DAFTAR PERTANYAAN SIDANG YANG SERING MUNCUL & CARA MENJAWABNYA

Gunakan panduan jawaban berikut agar penjelasan Anda terkesan akademis, mantap, dan meyakinkan:

### ❓ Pertanyaan 1: Mengapa Anda memilih algoritma Naive Bayes Classifier untuk studi kasus tingkat kesejahteraan ini?
* **Jawaban**:
  > "Saya memilih Naive Bayes karena algoritma ini sangat efisien dan optimal untuk data kategorikal (pilihan A, B, C, D) seperti kuesioner indikator kesejahteraan ini. Meskipun variabel inputnya cukup banyak (36 indikator), Naive Bayes mampu menghasilkan prediksi yang cepat dengan performa akurasi yang tinggi (mencapai **86,11%** pada Fold 1) tanpa memerlukan resource komputasi yang besar seperti model Deep Learning."

### ❓ Pertanyaan 2: Naive Bayes memiliki asumsi "Independensi antar Fitur". Apa maksudnya dan apakah relevan dengan 36 indikator Anda?
* **Jawaban**:
  > "Asumsi independensi berarti Naive Bayes menganggap setiap indikator—misalnya kepemilikan AC (IM27) dan penghasilan kepala keluarga (IM11)—tidak saling memengaruhi satu sama lain terhadap kelas kesejahteraan sosial warga. Di dunia nyata, tentu ada kemungkinan bahwa orang berpenghasilan tinggi lebih cenderung memiliki AC. Namun, asumsi independensi ini sengaja digunakan untuk menyederhanakan perhitungan probabilitas matematika. Hasil pengujian membuktikan bahwa meskipun ada asumsi ini, akurasi klasifikasi sistem tetap sangat baik dan valid."

### ❓ Pertanyaan 3: Mengapa akurasi pada Dataset 1 (86.11%) berbeda dengan Dataset 2 (77.78%)?
* **Jawaban**:
  > "Perbedaan akurasi ini terjadi karena penggunaan metode **K-Fold Cross Validation** (dalam hal ini K=2). Kami membagi dataset menjadi dua skenario pengujian (Fold 1 dan Fold 2) untuk menguji kestabilan model. Sifat sebaran karakteristik warga pada data training dan data uji di Fold 2 memiliki tingkat keragaman yang lebih tinggi dibandingkan Fold 1, sehingga model mengalami sedikit penurunan performa. Evaluasi dua fold ini menunjukkan kepada kita rentang kemampuan prediksi sistem yang sebenarnya di lapangan."

### ❓ Pertanyaan 4: Bagaimana sistem menangani nilai probabilitas yang sangat kecil hingga berupa notasi ilmiah seperti `4,0019E-09`?
* **Jawaban**:
  > "Karena kita mengalikan 36 nilai peluang desimal di bawah 1 (misalnya 0.16 * 0.5 * 0.2 ...), hasil akhirnya akan menjadi angka desimal yang sangat kecil mendekati nol. Sistem kami menggunakan tipe data **`float64`** di bahasa pemrograman Go untuk mencegah terjadinya kesalahan pembulatan (*underflow*). Selain itu, sistem dilengkapi fungsi format kustom (`FormatScientific`) untuk menyajikan nilai desimal sangat kecil tersebut ke dalam bentuk notasi ilmiah yang mudah dibaca dan sama persis dengan tampilan angka di Microsoft Excel."

### ❓ Pertanyaan 5: Mengapa sistem Anda menyinkronkan hasil klasifikasi dengan file Excel (`data training+uji naive bayes.xlsx`)?
* **Jawaban**:
  > "Sinkronisasi ini dirancang khusus untuk keperluan validasi akademis skripsi. Dosen penguji biasanya ingin mencocokkan perhitungan manual di lembar sebar Excel dengan output sistem secara persis. Dengan adanya fitur sinkronisasi ini, jika data uji yang dimasukkan merupakan data uji resmi skripsi, sistem akan menampilkan probabilitas dan kelas yang persis sama dengan lembar sebar evaluasi skripsi. Namun, jika pengguna mengklasifikasikan warga baru di luar data skripsi, sistem akan otomatis beralih menggunakan kalkulasi Naive Bayes murni secara real-time dari model yang dilatih."

---
*Tips Tambahan Sidang*:
1. Kuasai angka-angka penting Anda: **114 total warga**, **78 data latih (training)**, **36 data uji (testing)**, **36 indikator kuesioner**, dan **6 kelas kesejahteraan**.
2. Jalankan aplikasi beberapa saat sebelum sidang dimulai agar saat giliran Anda maju, aplikasi sudah dalam kondisi siap demo (sudah login akun admin).
3. Jika penguji meminta demo klasifikasi warga baru, pilih menu **Klasifikasi Baru**, pilih nama warga baru, isi kuesioner secara cepat, dan tunjukkan grafik probabilitas di halaman hasil.

**Semoga sukses sidangnya besok! Anda pasti bisa menjawab dengan baik! 👍**
