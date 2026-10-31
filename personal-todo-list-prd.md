# Product Requirements Document (PRD)

# Personal Todo List App

**Version:** 1.0\
**Status:** Draft --- MVP\
**Product type:** Personal productivity application\
**Target user:** Single user\
**Release scope:** MVP

------------------------------------------------------------------------

## 1. Product Overview

### 1.1 Background

Personal Todo List adalah aplikasi sederhana untuk membantu pengguna
mengelola aktivitas dan pekerjaan sehari-hari.

Pengguna dapat mencatat tugas, menentukan prioritas, mengatur deadline,
serta memantau tugas yang sudah atau belum diselesaikan.

Aplikasi dibuat khusus untuk penggunaan pribadi sehingga tidak
membutuhkan fitur kolaborasi, multi-user, atau manajemen tim.

### 1.2 Product Goals

Tujuan utama aplikasi:

1.  Memudahkan pencatatan tugas sehari-hari.
2.  Membantu pengguna mengatur prioritas pekerjaan.
3.  Memantau deadline setiap tugas.
4.  Menampilkan progres penyelesaian tugas.
5.  Menyediakan antarmuka yang minimalis dan mudah digunakan.

### 1.3 Target User

-   Pengguna individu.
-   Membutuhkan task management sederhana.
-   Ingin mengelola pekerjaan, aktivitas pribadi, dan ide.
-   Tidak membutuhkan fitur kolaborasi.

## 2. Functional Requirements

### 2.1 Task Management

Pengguna dapat melakukan operasi CRUD (Create, Read, Update, Delete):

-   Membuat task baru.
-   Melihat seluruh task.
-   Mengedit task.
-   Menghapus task.
-   Menandai task sebagai completed.

Setiap task memiliki atribut berikut:

  Field         Type          Required Description
  ------------- ----------- ---------- ----------------------------
  ID            UUID               Yes Identitas unik task
  Title         String             Yes Nama atau judul task
  Description   Text                No Detail task
  Priority      Enum               Yes Tingkat prioritas
  Status        Enum               Yes Status pengerjaan
  Due Date      Date/Time           No Tenggat waktu
  Created At    Timestamp          Yes Waktu task dibuat
  Updated At    Timestamp          Yes Waktu terakhir task diubah

### 2.2 Task Priority

Setiap task memiliki tiga tingkat prioritas:

-   **Low** --- pekerjaan tidak mendesak.
-   **Medium** --- pekerjaan dengan prioritas normal.
-   **High** --- pekerjaan penting dan perlu segera diselesaikan.

### 2.3 Task Status

Status task:

-   **Todo** --- belum dikerjakan.
-   **In Progress** --- sedang dikerjakan.
-   **Completed** --- sudah selesai.

### 2.4 Task Filtering & Sorting

Pengguna dapat memfilter task berdasarkan:

-   All Tasks.
-   Today.
-   Upcoming.
-   Completed.

Pengguna dapat mengurutkan task berdasarkan:

-   Due date.
-   Priority.
-   Created date.

### 2.5 Dashboard

Dashboard menampilkan ringkasan aktivitas:

  Component       Description
  --------------- ------------------------------------
  Total Tasks     Jumlah seluruh task
  Completed       Jumlah task yang selesai
  In Progress     Jumlah task yang sedang dikerjakan
  Pending         Jumlah task yang belum dikerjakan
  Today's Tasks   Task dengan deadline hari ini

Dashboard juga menampilkan daftar task terbaru atau task yang perlu
segera diselesaikan.

## 3. User Interface & Pages

Aplikasi MVP memiliki tiga halaman utama.

### 3.1 Dashboard

Halaman utama untuk melihat ringkasan aktivitas.

Komponen:

-   Summary cards.
-   Today's tasks.
-   Upcoming tasks.
-   Quick add task.

### 3.2 All Tasks

Halaman untuk mengelola semua task.

Komponen:

-   Task list.
-   Search task.
-   Filter status.
-   Filter priority.
-   Sorting.
-   Edit dan delete action.

### 3.3 Task Form

Form untuk membuat atau mengedit task.

Komponen:

-   Title.
-   Description.
-   Priority.
-   Status.
-   Due date.
-   Save button.
-   Cancel button.

## 4. Non-Functional Requirements

  -----------------------------------------------------------------------
  Requirement                         Description
  ----------------------------------- -----------------------------------
  Performance                         Halaman harus responsif dan cepat.

  Usability                           UI sederhana dan intuitif.

  Security                            Data hanya dapat diakses oleh
                                      pengguna.

  Reliability                         Data tidak hilang setelah aplikasi
                                      ditutup atau direstart.

  Responsive                          Mendukung desktop dan mobile.

  Maintainability                     Struktur kode modular dan mudah
                                      dikembangkan.
  -----------------------------------------------------------------------

## 5. Database Design

Untuk MVP, aplikasi cukup menggunakan satu tabel utama: `tasks`.

### 5.1 SQL Schema

``` sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    status VARCHAR(20) NOT NULL DEFAULT 'todo',
    due_date TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 5.2 Field Constraints

Nilai yang diperbolehkan:

-   `priority`: `low`, `medium`, `high`.
-   `status`: `todo`, `in_progress`, `completed`.

Validasi nilai tersebut perlu diterapkan pada backend. Database
constraint juga dapat ditambahkan untuk menjaga integritas data.

## 6. API Requirements

Jika menggunakan backend REST API, endpoint yang diperlukan:

  Method   Endpoint                  Description
  -------- ------------------------- ----------------------------
  GET      `/api/tasks`              Mengambil seluruh task
  GET      `/api/tasks/:id`          Mengambil detail satu task
  POST     `/api/tasks`              Membuat task
  PUT      `/api/tasks/:id`          Memperbarui task
  DELETE   `/api/tasks/:id`          Menghapus task
  PATCH    `/api/tasks/:id/status`   Memperbarui status task

### 6.1 Example Request --- Create Task

``` json
{
  "title": "Learn PostgreSQL",
  "description": "Study database fundamentals",
  "priority": "high",
  "status": "todo",
  "due_date": "2026-11-02T18:00:00Z"
}
```

### 6.2 API Behaviour

-   Endpoint harus mengembalikan HTTP status code yang sesuai.
-   Request dengan data tidak valid harus ditolak dengan pesan error
    yang jelas.
-   Request untuk task yang tidak ditemukan harus mengembalikan
    `404 Not Found`.
-   Penghapusan task harus meminta konfirmasi dari UI sebelum request
    dikirim.
-   Daftar task dapat menerima parameter filter dan sorting bila
    diperlukan.

## 7. Recommended Tech Stack

Rekomendasi stack awal:

  Layer        Technology           Purpose
  ------------ -------------------- -----------------------------------
  Backend      Go (Golang)          REST API dan business logic
  Database     PostgreSQL           Persistent data storage
  Frontend     React + TypeScript   User interface
  Styling      Tailwind CSS         Responsive UI
  Deployment   Docker Compose       Menjalankan aplikasi dan database

**Alternatif sederhana:** Go + HTML templates + HTMX + SQLite.
Pendekatan ini mengurangi kompleksitas frontend dan cocok untuk aplikasi
personal dengan satu pengguna.

Pemilihan stack final dapat disesuaikan dengan tujuan project: belajar
full-stack development atau membangun aplikasi yang sesederhana mungkin.

## 8. Out of Scope

Fitur berikut tidak termasuk dalam MVP:

-   Multi-user authentication.
-   Team collaboration.
-   Task assignment.
-   File attachment.
-   Push notification.
-   Calendar integration.
-   AI task assistant.
-   Recurring tasks.
-   Third-party integrations.

Fitur tersebut dapat dipertimbangkan untuk versi berikutnya.

## 9. Acceptance Criteria

MVP dianggap selesai jika:

-   [ ] Pengguna bisa membuat task.
-   [ ] Pengguna bisa melihat seluruh task.
-   [ ] Pengguna bisa mengedit task.
-   [ ] Pengguna bisa menghapus task.
-   [ ] Pengguna bisa mengubah status task.
-   [ ] Pengguna bisa menentukan prioritas.
-   [ ] Pengguna bisa menentukan deadline.
-   [ ] Pengguna bisa melakukan filtering dan sorting.
-   [ ] Dashboard menampilkan statistik task secara akurat.
-   [ ] Data tetap tersimpan setelah aplikasi direstart.
-   [ ] UI dapat digunakan melalui desktop dan mobile.

## 10. Development Milestones

Estimasi berikut merupakan target awal untuk satu developer.

  Phase       Scope                          Estimate
  ----------- -------------------------- ------------
  1           Project setup & database         1 hari
  2           Backend CRUD API                 2 hari
  3           Frontend UI                      2 hari
  4           Integration                      1 hari
  5           Testing & deployment             1 hari
  **Total**                                **7 hari**

Estimasi dapat berubah sesuai stack, pengalaman developer, dan kebutuhan
deployment.

## 11. Future Improvements

### V1.1

-   Recurring tasks.
-   Task categories.
-   Dark mode.
-   Improved search.

### V1.2

-   Calendar view.
-   Task reminders.
-   Export/import JSON.
-   Backup dan restore.

### V2.0

-   PWA support.
-   Offline mode.
-   Synchronization antar perangkat.

## 12. Product Decisions --- To Be Confirmed

Keputusan berikut perlu ditetapkan sebelum implementasi:

  -----------------------------------------------------------------------
  Decision                Options                 Initial Recommendation
  ----------------------- ----------------------- -----------------------
  Platform                Web, PWA, Desktop       Web

  Storage                 PostgreSQL, SQLite,     PostgreSQL untuk
                          LocalStorage            latihan backend; SQLite
                                                  untuk deployment
                                                  sederhana

  Backend                 Go, Go + HTMX, No       Go
                          backend                 

  Authentication          None, Single-user login None jika hanya
                                                  berjalan lokal dan
                                                  tidak terekspos ke
                                                  internet

  UI style                Minimalist, Modern      Minimalist
                          dashboard, Catppuccin   
  -----------------------------------------------------------------------

### Recommended MVP Direction

Untuk project latihan full-stack, gunakan **Go + PostgreSQL + React +
TypeScript** dengan Docker Compose.

Jika prioritasnya adalah aplikasi personal yang ringan dan cepat
selesai, gunakan **Go + HTMX + SQLite**.

------------------------------------------------------------------------

## Appendix: MVP Scope Summary

Fokus MVP adalah pengelolaan task dasar: CRUD, prioritas, status,
deadline, filter, sorting, dan dashboard ringkas.

Prinsip pengembangan:

1.  Implementasikan kebutuhan inti terlebih dahulu.
2.  Hindari fitur yang belum dibutuhkan pengguna.
3.  Pastikan data tersimpan dengan andal.
4.  Jaga struktur kode tetap sederhana dan mudah dirawat.
5.  Tambahkan fitur lanjutan berdasarkan kebutuhan nyata.
