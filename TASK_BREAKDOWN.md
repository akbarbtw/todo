# Task Breakdown for Personal Todo List App

Prioritas tugas-tugas ini dirancang agar mudah didelegasikan ke junior programmer atau model AI yang lebih murah. Setiap task bersifat sederhana, terstruktur, dan tidak terlalu kompleks. Outputnya berupa file checklist yang bisa langsung dikerjakan.

---

## Phase 1: Project Setup & Environment

- [ ] **1.1 Inisialisasi repository dan konfigurasi dasar**
  - Buat repository Git (jika belum ada).
  - Tambahkan file `.gitignore` untuk Go, Python (jika ada), dan database.
  - Pastikan editor/IDE siap digunakan (VS Code recommended).

- [ ] **1.2 Setup Docker Compose untuk backend dan database**
  - Buat file `docker-compose.yml` yang menyediakan:
    * Layanan `app` (Go backend).
    * Layanan `db` (PostgreSQL).
  - Konfigurasi environment variables untuk koneksi database (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB).
  - Test koneksi dengan `docker compose up -d`.

- [ ] **1.3 Buat struktur folder dasar backend Go**
  - Inisialisasi module Go: `go mod init todo-list`.
  - Buat folder `internal/` untuk kode aplikasi.
  - Buat file `main.go` minimal yang menampilkan pesan "Hello World" di port 8080.

---

## Phase 2: Database Schema

- [ ] **2.1 Buat SQL migration untuk tabel tasks**
  - Tulis file migrasi (misalnya `migrations/20240101_init.sql`).
  - Lakukan pembuatan tabel `tasks` sesuai schema dari PRD:
    * Kolom: `id` (UUID PRIMARY KEY), `title` (VARCHAR 255 NOT NULL), `description` (TEXT), `priority` (VARCHAR 20 DEFAULT 'medium'), `status` (VARCHAR 20 DEFAULT 'todo'), `due_date` (TIMESTAMP NULL), `created_at` (TIMESTAMP NOT NULL DEFAULT NOW()), `updated_at` (TIMESTAMP NOT NULL DEFAULT NOW()).
  - Tambahkan constraint nilai enumerated untuk `priority` (low, medium, high) dan `status` (todo, in_progress, completed).

- [ ] **2.2 Verifikasi schema di PostgreSQL**
  - Jalankan migrasi via `docker exec` atau tools klien (psql/pgAdmin).
  - Pastikan tabel dibuat dengan benar dan nilai enum tersimpan dengan benar.

---

## Phase 3: Backend API - CRUD Operations

- [ ] **3.1 Implementasikan endpoint GET /api/tasks**
  - Kembalikan daftar semua task beserta field-fieldnya.
  - Filter optional: `status`, `priority`, `due_date` (sekadar parameter query, logika sederhana).

- [ ] **3.2 Implementasikan endpoint GET /api/tasks/:id**
  - Ambil detail satu task berdasarkan ID.
- [ ] **3.3 Implementasikan endpoint POST /api/tasks**
  - Terima JSON: `title`, `description`, `priority`, `status`, `due_date`.
  - Lakukan validasi sederhana: `title` wajib ada, `priority` harus low/medium/high, `status` harus todo/in_progress/completed.
  - Simpan task ke database dan kembalikan response dengan status 201.

- [ ] **3.4 Implementasikan endpoint PUT /api/tasks/:id**
  - Perbarui field-task yang dipilih (hanya field yang dikirim, sisa tetap).
  - Validasi sama seperti POST.

- [ ] **3.5 Implementasikan endpoint DELETE /api/tasks/:id**
  - Hapus task berdasarkan ID.
  - (Opsional) Tambahkan konfirmasi logika di level backend dengan catatan di comment.

- [ ] **3.6 Implementasikan endpoint PATCH /api/tasks/:id/status**
  - Hanya memperbarui field `status` task.
  - Berguna untuk menandai task selesai atau mengubah progress.

---

## Phase 4: Frontend UI Halaman Dasar

- [ ] **4.1 Buat halaman Dashboard sederhana**
  - Gunakan React + TypeScript atau HTML templat Go (HTMX/SQLite alternatif).
  - Tampilkan ringkasan statistik:
    * Total tasks.
    * Completed tasks.
    * In Progress tasks.
    * Pending tasks.
    * Today's tasks (deadline hari ini).
  - Gunakan chart/simple list untuk menampilkan info.

- [ ] **4.2 Buat halaman All Tasks**
  - Tampilkan daftar semua task dalam tabel atau card.
  - Setiap card menampilkan: title, priority, status, due date.

- [ ] **4.3 Buat komponen Task Form**
  - Form dengan field: title (wajib), description (opsional), priority (dropdown low/medium/high), status (dropdown todo/in_progress/completed), due date (date picker).
  - Tombol Save dan Cancel.
  - Form bisa digunakan untuk Create (menambah task baru) maupun Update (mengedit task existing).

- [ ] **4.4 Integrasikan form dengan backend API**
  - Kirim request POST/PUT ke `/api/tasks` saat submit form.
  - Tangani response sukses/error dengan pesan alert sederhana.

---

## Phase 5: Filtering, Sorting & Usability

- [ ] **5.1 Tambahkan fitur filter status**
  - Tombol/tabs: All, Active (todo), Completed, In Progress.
  - Saat diklik, filter task yang ditampilkan di halaman All Tasks.

- [ ] **5.2 Tambahkan fitur filter priority**
  - Filter: Low, Medium, High.

- [ ] **5.3 Tambahkan fitur sorting**
  - Opsi sorting: By Due Date (terlama/terbaru), By Priority (tinggi/ Rendah), By Created Date (terlama/terbaru).

- [ ] **5.4 Pastikan responsiveness di mobile**
  - Lihat tampilan di lebar layar < 600px.
  - Pastikan tombol, teks, dan form masih terbaca dan dapat diklik dengan mudah.

---

## Phase 6: Data Persistence & Uji Coba

- [ ] **6.1 Uji coba menyimpan data setelah restart aplikasi**
  - Buat task, restart container, pastikan data masih ada di database.

- [ ] **6.2 Lakukan uji manual CRUD lengkap**
  - Buat task baru → lihat daftar → edit task → hapus task.
  - Verifikasi setiap langkah melalui browser (Postman atau curl juga boleh).

- [ ] **6.3 Tambahkan basic error handling di backend**
  - Jika request JSON tidak valid, kembalikan respons 400 dengan pesan "Invalid input".
  - Jika task tidak ditemukan (GET/PUT/DELETE by ID), kembalikan 404 dengan pesan "Task not found".

---

## Phase 7: Deployment & Selesai

- [ ] **7.1 Pastikan Docker Compose berjalan di lingkungan pengembangan**
  - `docker compose up --build` jalan tanpa error.
- [ ] **7.2 Buat README.md sederhana**
  - Isi perintah untuk menjalankan aplikasi: `docker compose up -d`.
  - Catat alamat URL frontend (misih http://localhost:3000) dan API (http://localhost:8080).
- [ ] **7.3 Review acceptance criteria dari PRD**
  - Pastikan semua checklist di bagian 9 PRD terpenuhi:
    * Buat, lihat, edit, hapus task.
    * Ubah status.
    * Prioritas dan deadline.
    * Filter dan sorting.
    * Dashboard statistik akurat.
    * Data bertahan setelah restart.
    * UI bisa digunakan desktop dan mobile.

---

*Catatan: Semua task di atas disengajarkan dengan spesifikasi yang jelas dan file-file konfigurasi (seperti docker-compose.yml, migrasi SQL, rute API) dibuat sederhana demi kemudahan pembelajaran dan pengembangan oleh junior programmer atau model AI yang lebih murah.*