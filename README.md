# Visual Mesin

Sistem monitoring dan analisis produksi ban (tire manufacturing) berbasis **Golang + React**, migrasi dari Laravel + Filament.
 
note: fase 9-10 blm

## Arsitektur

```
┌─────────────────────────────────────────────┐
│               React SPA (Antd)               │
│  Dashboard · CRUD · AI Chat · Export         │
└────────────┬────────────────────┬────────────┘
             │ HTTP/REST          │ WebSocket
             ▼                    ▼
┌─────────────────────────────────────────────┐
│           Go Backend API (Gin)               │
│  Auth · CRUD · AI Pipeline · Export Engine   │
└───────────┬─────────────────────────────────┘
            │
     ┌──────┴──────┐
     │   Redis      │
     │ (Cache/Queue)│
     └──────┬──────┘
            │
     ┌──────┴──────┐
     │   MariaDB    │
     │ (Local +     │
     │  Resource)   │
     └─────────────┘
```

## Stack

| Layer | Teknologi |
|-------|-----------|
| Backend API | Go + Gin + GORM |
| Frontend | React 19 + TypeScript + Vite |
| UI | Ant Design 6 + Tailwind CSS |
| State | TanStack Query + Zustand |
| Database | MariaDB/MySQL (multi-database) |
| Cache/Queue | Redis |
| Auth | JWT + BCrypt + RBAC |
| AI/LLM | Ollama (via HTTP) |
| Dokumentasi API | Swagger/OpenAPI (swaggo) |

## Struktur Proyek

```
backend/                        # Go backend
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # App configuration
│   ├── db/                     # DB connection manager + migrations
│   ├── middleware/             # Auth, CORS, RBAC, Response helpers
│   ├── models/                 # GORM models
│   ├── handlers/               # HTTP handlers
│   ├── services/               # Business logic
│   ├── repository/             # Data access layer
│   ├── ai/                     # AI chat pipeline
│   ├── ws/                     # WebSocket hub
│   └── routes/                 # Route definitions
├── migrations/                 # SQL migration files
└── pkg/utils/                  # Shared utilities (crypto, etc.)

frontend/                       # React frontend
├── src/
│   ├── components/             # Shared components
│   ├── pages/                  # Page components
│   ├── layouts/                # Layout components
│   ├── stores/                 # Zustand stores
│   ├── services/               # API client (Axios)
│   ├── hooks/                  # Custom hooks
│   ├── types/                  # TypeScript types
│   └── utils/                  # Utilities
└── vite.config.ts
```

## Memulai

### Prasyarat

- Go 1.25+
- Node.js 22+
- Docker & Docker Compose (MariaDB, Redis, Ollama)

### Akun Default (Seeder)

| Level | Email | Password |
|-------|-------|----------|
| Admin | admin@admin.com | password |
| Produksi | user@visualmesin.com | user123 |

Seeder jalan otomatis via migration (`000010_seed_default_users.up.sql`).

### Development

```bash
# Clone & masuk
git clone https://github.com/Anraaa/visual-mesin.git
cd visual-mesin

# Setup environment
cp backend/.env.example backend/.env
# Edit .env sesuai kebutuhan

# Jalankan infrastructure (MariaDB, Redis, Ollama)
docker compose up -d

# Backend
cd backend
go run ./cmd/server/

# Frontend (terminal terpisah)
cd frontend
npm install
npm run dev
```

Backend berjalan di `http://localhost:8080`, Frontend di `http://localhost:5173`.

### API Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/auth/register` | Register user |
| GET | `/api/v1/auth/me` | Profile user |
| GET | `/api/v1/db-connections` | List DB connections |
| POST | `/api/v1/db-connections` | Tambah DB connection |
| PUT | `/api/v1/db-connections/:id` | Update DB connection |
| DELETE | `/api/v1/db-connections/:id` | Hapus DB connection |
| POST | `/api/v1/db-connections/test` | Test koneksi DB |
| GET | `/api/v1/resource-db-configs` | List resource DB configs |
| POST | `/api/v1/resource-db-configs` | Tambah resource DB config |
| GET | `/api/v1/resources/:resource` | List data produksi |
| GET | `/api/v1/resources/:resource/:id` | Detail data produksi |
| POST | `/api/v1/resources/:resource` | Tambah data produksi |
| PUT | `/api/v1/resources/:resource/:id` | Update data produksi |
| DELETE | `/api/v1/resources/:resource/:id` | Hapus data produksi |
| GET | `/api/v1/resources/:resource/columns` | Kolom tabel produksi |
| GET | `/api/v1/building/:table` | Building data (rtba1-3, rtbc1-4, rtbe1-2) |
| GET | `/api/v1/extruder/:table` | Extruder data (rteex1-3, cyclic, pcs, datalog) |
| GET | `/api/v1/curing/:table` | Curing data (curtire, item_measurement, gtentire) |
| GET | `/api/v1/trimming/:table` | Trimming data (trimming, rtc-tr1) |
| GET | `/api/v1/monitoring/:table` | Monitoring data (monitoringtl1, alarm, yield) |
| GET | `/api/v1/recipe/:table` | Recipe & order data |
| GET | `/api/v1/master/:table` | Master data (mesin, bpbl, material) |
| GET | `/swagger/*any` | Dokumentasi Swagger |

## Database

Arsitektur multi-database:
- **Local DB**: users, roles, permissions, ai_chat_history, db_connections, dll
- **Resource DBs**: Data produksi mesin (rtba, rtbc, rteex, curtire, trimming, dll) — bisa di server berbeda, dikonfigurasi dinamis

## Tahapan Migrasi

| Fase | Status | Deskripsi |
|------|--------|-----------|
| 0 | ✅ Selesai | Project setup, Docker, migrations |
| 1 | ✅ Selesai | Auth & RBAC (JWT, login, register) |
| 2 | ✅ Selesai | Dynamic DB connection manager + CRUD configs |
| 3 | ✅ Selesai | Resource table APIs (Building, Extruder, Curing, dll) |
| 4 | ⬜ Belum | Frontend foundation (layout, login, theme) |
| 5 | ⬜ Belum | Frontend feature pages |
| 6 | ✅ Selesai | AI Chat Assistant (Ollama, intent detection, SQL gen, firewall, UI) |
| 7 | ⬜ Belum | Export system |
| 8 | ⬜ Belum | WebSocket real-time |
| 9 | ⬜ Belum | Testing |
| 10 | ⬜ Belum | Deployment |

## Lisensi

Proprietary — Internal use only.
