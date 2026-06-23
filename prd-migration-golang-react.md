# PRD: Migrasi Visual Mesin ke Golang + React

## 1. Ringkasan

Migrasi aplikasi Visual Mesin dari **Laravel (PHP) + Filament Admin Panel** ke **Golang (Backend) + React (Frontend)** untuk meningkatkan performa, skalabilitas, dan pengalaman pengguna.

### Target Stack

| Layer | Teknologi |
|-------|-----------|
| Backend API | Go (Gin/Echo/Fiber) |
| Frontend | React + TypeScript + Vite |
| UI Library | Ant Design (Antd) |
| State Management | React Query (TanStack Query) + Zustand |
| Charting | Recharts / Apache ECharts |
| Real-time | WebSocket (Gorilla/neffos) |
| Database | MariaDB/MySQL (sama seperti existing) |
| Cache/Queue | Redis |
| Auth | JWT + RBAC |
| AI/LLM | Ollama (sama, via HTTP client) |
| Dokumentasi API | Swagger/OpenAPI (swaggo) |

---

## 2. Tujuan

1. **Performa lebih tinggi** — Go memberikan latensi lebih rendah dan throughput lebih tinggi dibanding PHP-FPM
2. **Skalabilitas lebih baik** — Go ringan (goroutines) untuk ribuan koneksi WebSocket simultan
3. **Frontend modern** — React memberikan UX yang lebih responsif, SPA tanpa reload
4. **Perawatan lebih mudah** — Codebase terpisah backend/frontend dengan boundary yang jelas
5. **Real-time lebih powerful** — WebSocket native tanpa perantara seperti Pusher/Reverb
6. **Penghematan resource** — Go binary ringan tanpa runtime interpreter

---

## 3. Fitur yang Akan Dimigrasi

### 3.1. Fase 1 — Core Platform (MVP)

| Fitur | Prioritas | Keterangan |
|-------|-----------|------------|
| Auth & RBAC | P0 | Login JWT, role-based access (admin, eng, tech, prod) |
| REST API v1 | P0 | 20+ endpoint existing untuk akses data produksi |
| Dashboard | P0 | Widget overview, chart real-time |
| CRUD Data Produksi | P0 | Manajemen 30+ tabel produksi |
| Manajemen Users | P0 | CRUD user, role, permission |
| Manajemen DB Dinamis | P0 | Tambah/ubah koneksi database runtime |

### 3.2. Fase 2 — AI & Analytics

| Fitur | Prioritas | Keterangan |
|-------|-----------|------------|
| AI Chat Assistant | P0 | Chat dengan LLM lokal (Ollama) untuk tanya data produksi |
| Intent Detection Pipeline | P0 | 17 intent, SQL generation, SQL firewall |
| Predictive Analytics | P1 | Prediksi material, order, maintenance, yield |
| Barcode Ambiguity Handler | P1 | Deteksi barcode tanpa konteks |

### 3.3. Fase 3 — Advanced Features

| Fitur | Prioritas | Keterangan |
|-------|-----------|------------|
| Export System | P0 | CSV streaming dengan chunking + ZIP + Redis Stream |
| WebSocket Real-time | P1 | Update dashboard live, status export |
| Dark Mode | P1 | Toggle tema gelap/terang |
| Mobile Responsive | P1 | Dashboard bisa diakses dari tablet/hp |

---

## 4. Arsitektur Sistem

### 4.1. High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                   React SPA                         │
│  (Dashboard, CRUD, AI Chat, Export, Settings)       │
└────────────┬────────────────────────────┬───────────┘
             │ HTTP/REST                   │ WebSocket
             ▼                            ▼
┌─────────────────────────────────────────────────────┐
│              Go Backend API (Gin)                   │
│  ┌─────────┐ ┌──────────┐ ┌──────────────────┐     │
│  │ Auth    │ │ CRUD API │ │ AI Chat Pipeline │     │
│  │ (JWT)   │ │ (REST)   │ │ (Intent→SQL→AI)  │     │
│  ├─────────┤ ├──────────┤ ├──────────────────┤     │
│  │ Export  │ │ WebSocket│ │ Predictive       │     │
│  │ Service │ │ Hub      │ │ Analytics Engine │     │
│  └─────────┘ └──────────┘ └──────────────────┘     │
└───────────┬─────────────────────────────────────────┘
            │
     ┌──────┴──────┐
     │   Redis      │
     │ (Cache/Q)    │
     └──────┬──────┘
            │
     ┌──────┴──────┐
     │  MariaDB     │
     │ (Local +     │
     │  Resource)   │
     └─────────────┘
```

### 4.2. Alur Data AI Chat

```
User Input → Intent Detection → SQL Generation → SQL Firewall
    → Execute Query → Format Response → Return ke React UI
```

### 4.3. Alur Export

```
User Request → Redis Stream Queue → Go Worker (goroutine)
    → Chunked DB Query (100K rows) → CSV → ZIP → S3/Local
    → WebSocket notify selesai
```

---

## 5. Directory Structure (Target)

```
/visual-mesin/
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/            # App configuration
│   │   ├── db/                # Database connections
│   │   ├── middleware/        # Auth, CORS, logging
│   │   ├── models/            # Data models
│   │   ├── handlers/          # HTTP handlers
│   │   ├── services/          # Business logic
│   │   ├── repository/        # Data access layer
│   │   ├── ai/                # AI chat pipeline
│   │   ├── export/            # Export engine
│   │   ├── ws/                # WebSocket hub
│   │   └── routes/            # Route definitions
│   ├── pkg/
│   │   └── utils/             # Shared utilities
│   ├── migrations/            # SQL migrations
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # React app
│   ├── src/
│   │   ├── components/        # Shared components
│   │   ├── pages/             # Page components
│   │   ├── hooks/             # Custom hooks
│   │   ├── services/          # API client
│   │   ├── stores/            # Zustand stores
│   │   ├── types/             # TypeScript types
│   │   ├── utils/             # Utilities
│   │   ├── layouts/           # Layout components
│   │   └── features/          # Feature modules
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
│
├── docker-compose.yaml
└── docs/
    ├── prd.md
    ├── api-spec.yaml
    └── database-schema.md
```

---

## 6. Tech Stack Detail

### Backend (Go)

| Komponen | Pilihan | Alasan |
|----------|---------|--------|
| HTTP Framework | Gin | Mature, fast, banyak middleware |
| ORM | GORM | Popular, fitur lengkap, migration support |
| Auth | JWT (golang-jwt) + BCrypt | Standard, aman |
| Validasi | go-playground/validator | Deklaratif via struct tags |
| WebSocket | gorilla/websocket | Standard library untuk WS |
| Queue | go-redis + redis streams | Sama dengan existing Redis |
| AI Client | net/http (Ollama API) | Simple, tanpa dependency tambahan |
| Config | viper | Environment/file config |
| Logging | zerolog | Performa tinggi, structured log |
| Swagger | swaggo/swag | Auto-generate OpenAPI spec |
| Migration | golang-migrate | SQL-first migration |

### Frontend (React)

| Komponen | Pilihan | Alasan |
|----------|---------|--------|
| Framework | React 18 + TypeScript | Standard industri |
| Bundler | Vite | Cepat, modern |
| Routing | React Router v6 | Standard |
| UI Library | Ant Design (Antd) | Enterprise-grade, komponen lengkap untuk admin panel |
| State | TanStack Query + Zustand | Server state + client state |
| Chart | Recharts | React-native, mudah dikustom |
| Table | TanStack Table | Headless, powerful |
| Form | React Hook Form + Zod | Performa tinggi, validasi |
| HTTP | Axios | Interceptor, cancel token |
| WebSocket | Custom hook (gorilla/ws) | Langsung ke WS Go |
| CSS | Tailwind CSS | Utility-first |

---

## 7. Database Strategy

### Multi-Database Architecture (sama seperti existing)

```
┌─────────────────┐     ┌─────────────────────┐
│ Application DB  │     │ Resource DB(s)       │
│ (Local)         │     │ (Production Data)    │
│                 │     │                      │
│ - users         │     │ - rtba1/2/3          │
│ - roles         │     │ - rtbc1/2/3/4        │
│ - permissions   │     │ - curtire             │
│ - ai_schema_map │     │ - trimming           │
│ - chat_history  │     │ - rteex1/2/3head     │
│ - db_connections│     │ - material           │
│ - etc           │     │ - dan 25+ tabel lain │
└─────────────────┘     └─────────────────────┘
```

### Koneksi Dinamis
- Backend Go menyimpan konfigurasi DB di tabel `db_connections`
- Saat request masuk, backend resolve koneksi target berdasarkan resource
- Pool koneksi dikelola per database (sql.DB connection pooling)

---

## 8. REST API Design (Sama seperti existing)

Semua endpoint existing dipertahankan dengan response format yang sama:

### Format Response Standard

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": { ... },
  "meta": {
    "current_page": 1,
    "per_page": 25,
    "total": 100,
    "last_page": 4
  }
}
```

### Authentication
- Header: `X-API-Key` untuk akses machine-to-machine
- Header: `Authorization: Bearer <JWT>` untuk user session

---

## 9. Frontend Page Structure

```
/login                    → Login page
/dashboard                → Main dashboard (widgets + charts)
/analytics/ai-chat        → AI Chat Assistant
/data/{resource}          → CRUD table untuk setiap resource
/data/{resource}/{id}     → Detail view
/exports                  → Export files management
/settings                 → User settings
/admin/users              → User management
/admin/db-connections     → Database connection management
/admin/roles              → Role & permission management
```

### Komponen Per Halaman

**Dashboard:**
- Machine Status Overview (card grid)
- Alarm History Chart
- Curing Production Chart
- Datalog Monitoring
- Material Stock Chart
- Yield Defect Chart
- Side Comparison Chart
- Speed Monitoring

**CRUD Table (generik):**
- Search bar + filter panel
- Sortable columns
- Pagination
- Bulk actions
- Inline edit
- Export button

**AI Chat:**
- Sidebar: session history
- Main: chat bubbles (user + AI)
- Composer: text input + suggestion chips
- Status indicator (processing/complete/error)

---

## 10. Tahapan Migrasi

### Tahap 1 — Foundation (2-3 minggu)
- [x] Setup Go project structure + routing
- [x] Setup React project + Vite + routing
- [x] Auth system (JWT login, register, RBAC)
- [x] Database connection pool (local + dynamic)
- [x] Basic CRUD generator pattern
- [x] Docker compose untuk development

### Tahap 2 — Core Features (3-4 minggu)
- [ ] Migrasi 32 CRUD resources dari Filament ke React
- [ ] Dashboard dengan semua widget/chart
- [ ] REST API untuk semua endpoint
- [ ] Export system dengan streaming

### Tahap 3 — AI & Analytics (2-3 minggu)
- [ ] AI Chat Pipeline (intent detection, SQL gen, firewall)
- [ ] Predictive Analytics (material, order, maintenance, yield)
- [ ] WebSocket untuk real-time chat status

### Tahap 4 — Polish (1-2 minggu)
- [ ] Dark mode
- [ ] Responsive design
- [ ] Performance optimization
- [ ] Testing
- [ ] Dokumentasi

---

## 11. Non-Functional Requirements

| Kategori | Target |
|----------|--------|
| Response Time API | < 100ms (p95) tanpa AI |
| Response Time AI | < 5 detik (p95) |
| Concurrent Users | 500+ simultan |
| Concurrency | 10.000+ WebSocket koneksi |
| Uptime | 99.9% |
| Security | JWT, RBAC, SQL Firewall, Encrypted DB password |
| Data Consistency | ACID untuk transaksi lokal |

---

## 12. Risiko & Mitigasi

| Risiko | Dampak | Mitigasi |
|--------|--------|----------|
| Migrasi 50+ migrations | Tinggi | Auto-run migration di startup Go |
| Dynamic DB connections | Sedang | Pool management dengan timeout + retry |
| AI pipeline kompleks | Tinggi | Modular service, test per stage |
| Learning curve Go team | Sedang | Code review, pair programming |
| Data loss saat migrasi | Tinggi | Parallel run dulu, cutover terencana |

---

## 13. Glossary

| Istilah | Definisi |
|---------|----------|
| RTBA | Radial Tire Building Assembly |
| RTBC | Radial Tire Building Calendar |
| RTBE | Radial Tire Building Evaluation |
| RTEEX | Radial Tire Extruder |
| RTC | Real-Time Control |
| NG | Not Good (cacat) |
| OK | Good (lolos QC) |
| PUD | Pass Up Drum (tahap building) |
| BTD | Belt Tread Drum (tahap building) |
| BDD | Bead Drum (tahap building) |
| ETC | Estimated Time to Completion |
| LLM | Large Language Model (Ollama) |
| SPA | Single Page Application |
| RBAC | Role-Based Access Control |
