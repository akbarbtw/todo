# Dokumentasi Project: Personal Todo List API

Aplikasi backend REST API sederhana untuk mengelola daftar tugas (Todo List) pribadi menggunakan **Go (Golang)** dan **PostgreSQL**.

---

## 📋 Daftar Isi
1. [Arsitektur & Tech Stack](#-arsitektur--tech-stack)
2. [Struktur Project](#-struktur-project)
3. [Konfigurasi & Menjalankan Aplikasi](#-konfigurasi--menjalankan-aplikasi)
4. [API Specification](#-api-specification)
5. [Database Schema](#-database-schema)

---

## 🛠 Arsitektur & Tech Stack

- **Bahasa Pemrograman:** Go (Golang) 1.23
- **Framework Router:** `github.com/gorilla/mux`
- **Driver Database:** `github.com/lib/pq` (PostgreSQL)
- **Database:** PostgreSQL 16 (via Podman/Docker)
- **Utilities:** `github.com/google/uuid` (UUID v4 generation)

---

## 📂 Struktur Project

```text
todo-list/
├── cmd/
│   └── main.go                 # Entrypoint aplikasi & inisialisasi server
├── internal/
│   ├── handler/
│   │   └── task_handler.go     # HTTP handlers & routing
│   ├── model/
│   │   └── task.go             # Struct domain & validasi
│   ├── repository/
│   │   └── task_repository.go  # Database queries
│   └── service/
│       └── task_service.go     # Business logic & validasi
├── migrations/
│   └── 20240101_init.sql       # SQL migration untuk tabel tasks
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## 🚀 Konfigurasi & Menjalankan Aplikasi

### 1. Menjalankan Database PostgreSQL
Gunakan Podman/Docker untuk menjalankan container PostgreSQL:
```bash
podman run -d --name todo-db -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=todo \
  docker.io/library/postgres:16-alpine
```

### 2. Menjalankan Migrasi Database
```bash
podman exec -i todo-db psql -U postgres -d todo < migrations/20240101_init.sql
```

### 3. Menjalankan Server Backend
```bash
go run ./cmd/main.go
```
Server akan berjalan di `http://localhost:8080`.

---

## 🔌 API Specification

Base URL: `http://localhost:8080/api`

### 1. Health Check
- **Endpoint:** `GET /health`
- **Response:** `200 OK`

### 2. Get All Tasks
- **Endpoint:** `GET /api/tasks`
- **Query Parameters (Opsional):**
  - `status` (`todo`, `in_progress`, `completed`)
  - `priority` (`low`, `medium`, `high`)
  - `due_date` (`YYYY-MM-DD`)
- **Response (200 OK):**
```json
[
  {
    "id": "00bd25fb-e1a4-48ff-b404-b8f62f938bf6",
    "title": "Belajar Go dan PostgreSQL",
    "description": "Membuat aplikasi todo list backend",
    "priority": "high",
    "status": "in_progress",
    "created_at": "2026-10-31T12:51:27.847928Z",
    "updated_at": "2026-10-31T12:51:56.055343Z"
  }
]
```

### 3. Get Task by ID
- **Endpoint:** `GET /api/tasks/{id}`
- **Response (200 OK):** Single task object.
- **Response (404 Not Found):** `{"error": "task not found"}`

### 4. Create Task
- **Endpoint:** `POST /api/tasks`
- **Request Body (JSON):**
```json
{
  "title": "Belajar Go dan PostgreSQL",
  "description": "Membuat aplikasi todo list backend",
  "priority": "high",
  "status": "todo",
  "due_date": "2026-11-05T18:00:00Z"
}
```
- **Response (201 Created):** Created task object.

### 5. Update Task
- **Endpoint:** `PUT /api/tasks/{id}`
- **Request Body (JSON):** Mengirim field yang ingin diubah (partial/full).

### 6. Update Task Status
- **Endpoint:** `PATCH /api/tasks/{id}/status`
- **Request Body (JSON):**
```json
{
  "status": "completed"
}
```
- **Response (200 OK):** Updated task object.

### 7. Delete Task
- **Endpoint:** `DELETE /api/tasks/{id}`
- **Response (204 No Content):** Task berhasil dihapus.

---

## 🗄 Database Schema (`tasks`)

```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    status VARCHAR(20) NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'completed')),
    due_date TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```
