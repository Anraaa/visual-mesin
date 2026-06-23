# Tahapan Pengerjaan Visual Mesin — Migrasi Golang + React

---

## Fase 0: Project Setup & Foundation

| # | Task | Keterangan |
|---|------|-----------|
| 0.1 | Inisialisasi Go module | `go mod init`, struktur folder `cmd/`, `internal/`, `pkg/`, `migrations/` |
| 0.2 | Inisialisasi React + Vite | TypeScript + Ant Design (Antd) |
| 0.3 | Setup Docker Compose dev | MariaDB, Redis, Ollama — `docker-compose.yaml` + `Dockerfile` |
| 0.4 | Setup database migration | `golang-migrate/migrate` — 14 file SQL untuk tabel aplikasi |
| 0.5 | Setup CI/CD & linting | GitHub Actions, golangci-lint, ESLint, Prettier, lint-staged |

---

## Fase 1: Core Backend — Auth & RBAC

| # | Task | Keterangan |
|---|------|-----------|
| 1.1 | Struktur Go project | Gin router, middleware pattern, controller/service/repository |
| 1.2 | JWT authentication | Login, register, refresh token, logout, blacklist |
| 1.3 | RBAC middleware | Middleware `RequirePermission`, check role & permission |
| 1.4 | Swagger / OpenAPI | Setup `swaggo/swag`, dokumentasi otomatis |

---

## Fase 2: Core Backend — Dynamic Database Connection

| # | Task | Keterangan |
|---|------|-----------|
| 2.1 | DB Connection Manager | Pool map `map[string]*sql.DB`, health check, auto-reconnect |
| 2.2 | CRUD `db_connections` | API manage konfigurasi database eksternal |
| 2.3 | CRUD `resource_db_configs` | API mapping resource ke database |
| 2.4 | Generic Resource Query Service | Baca tabel resource secara dinamis via metadata |

---

## Fase 3: Core Backend — Resource Table APIs

| # | Task | Endpoints |
|---|------|-----------|
| 3.1 | Building module (RTBA) | `rtba1`, `rtba2`, `rtba3` |
| 3.2 | Building quality (RTBC) | `rtbc1`, `rtbc2`, `rtbc3`, `rtbc4` |
| 3.3 | Building evaluation (RTBE) | `rtbe1`, `rtbe2` |
| 3.4 | Extruder module | `rteex1`, `rteex2`, `rteex3head` |
| 3.5 | Curing & Green Tire | `curtire`, `item_measurement`, `gtentire` |
| 3.6 | Trimming & Monitoring | `trimming`, `rtc-tr1`, `monitoringtl1`, `rtl-tl1`, `rtltl1`, `alarm_history` |
| 3.7 | Recipe & Order | `recipe1`, `recipe1queue`, `recipe_history`, `order_report`, `batch_report` |
| 3.8 | Cyclic & Datalog | `recorddatacyclic`, `recorddatapcs`, `datalog` |
| 3.9 | Supporting | `mastermcn`, `bpbl`, `rsc_pc1`, `material` |

---

## Fase 4: Frontend Foundation

| # | Task | Keterangan |
|---|------|-----------|
| 4.1 | Struktur React project | Router (React Router), layout (sidebar + header), auth guard |
| 4.2 | Login page & auth flow | JWT storage (localStorage), Axios interceptor, redirect |
| 4.3 | Dashboard layout | Sidebar menu dinamis, breadcrumb, dark mode toggle, user dropdown |
| 4.4 | Theme provider | Zustand store untuk theme, persist ke localStorage |
| 4.5 | Error handling | Global error boundary, toast notifications, 404/403 pages |

---

## Fase 5: Frontend — Feature Pages

| # | Task | Keterangan |
|---|------|-----------|
| 5.1 | User management page | List, create, edit user; assign role & permission |
| 5.2 | Role & permission pages | CRUD roles, assign permissions |
| 5.3 | DB connection pages | CRUD database sources + test connection |
| 5.4 | Data produksi pages | Table view per resource table — filter, sort, pagination, search |
| 5.5 | Dashboard page | Summary cards, charts (Recharts), real-time updates |
| 5.6 | Detail page per entity | Modal atau drawer untuk lihat detail row |

---

## Fase 6: AI Chat Assistant (Phase 2)

| # | Task | Keterangan |
|---|------|-----------|
| 6.1 | AI Schema Map CRUD | Manage 17 intents, keyword matching, few-shot examples |
| 6.2 | Intent detection pipeline | NLP-based intent classification dari user query |
| 6.3 | SQL generation engine | Prompt engineering → Ollama → generate SQL |
| 6.4 | SQL firewall | Read-only validation, complexity limit, blacklist words |
| 6.5 | SQL execution & response | Execute aman, format hasil ke natural language |
| 6.6 | Chat UI | Chat interface, streaming response, history sidebar |
| 6.7 | Chat history | CRUD `ai_chat_history`, pagination, search |

---

## Fase 7: Export System (Phase 3)

| # | Task | Keterangan |
|---|------|-----------|
| 7.1 | Redis Stream queue | Setup queue + Go worker untuk export job |
| 7.2 | Chunked CSV generator | Stream data in chunks, avoid memory overflow |
| 7.3 | ZIP packaging | Gabung multiple CSV ke ZIP, progress tracking |
| 7.4 | Export API & download | Submit job, cek status, download file |
| 7.5 | Export UI | Tombol export, modal pilih kolom, progress bar, notifikasi |

---

## Fase 8: WebSocket & Real-time (Phase 3)

| # | Task | Keterangan |
|---|------|-----------|
| 8.1 | WebSocket Hub | `gorilla/websocket`, hub pattern, room-based broadcast |
| 8.2 | Auth over WebSocket | JWT verification on connect |
| 8.3 | Dashboard live updates | Push data produksi terbaru ke semua client |
| 8.4 | Notification broadcast | Notifikasi real-time (export selesai, error, dll) |
| 8.5 | Activity log | Audit trail user activity |

---

## Fase 9: Testing

| # | Task | Keterangan |
|---|------|-----------|
| 9.1 | Unit test backend | `testing` + `testify` untuk service & repository |
| 9.2 | Integration test API | Test endpoint dengan test container / mock DB |
| 9.3 | Unit test frontend | `vitest` + React Testing Library |
| 9.4 | E2E test | Playwright / Cypress untuk flow kritis |
| 9.5 | Load test | `k6` — 500+ concurrent users, <100ms p95 |

---

## Fase 10: Deployment

| # | Task | Keterangan |
|---|------|-----------|
| 10.1 | Dockerfile production | Multi-stage build Go (alpine) + React (nginx) |
| 10.2 | Docker Compose production | App, Nginx reverse proxy, MariaDB, Redis, Ollama |
| 10.3 | Environment config | `.env.example`, secrets management |
| 10.4 | Nginx config | Reverse proxy, static file serve, SSL |
| 10.5 | Deploy script | Manual deploy atau setup ke server production |
| 10.6 | Monitoring | Health check endpoint, logging (zerolog), metrics |

---

## Visual Dependency Graph

```
Fase 0 (Setup)
   ├──> Fase 1 (Auth Backend) ──> Fase 4 (Frontend Foundation)
   ├──> Fase 2 (DB Connection) ──> Fase 3 (Resource APIs)
   │                                   │
   └──> Fase 4 ◄───────────────────────┘
            │
            ├──> Fase 5 (Frontend Pages)
            ├──> Fase 6 (AI Chat) ── bisa jalan paralel
            ├──> Fase 7 (Export)
            └──> Fase 8 (WebSocket)
                     │
                     v
                Fase 9 (Testing)
                     │
                     v
                Fase 10 (Deploy)
```

---

## Catatan Penting

- **Docker**: Docker Compose untuk development (MariaDB + Redis + Ollama) harus ready di **Fase 0.3** agar semua developer punya environment seragam.
- **Multi-database**: Setiap endpoint di Fase 3 harus bisa membaca dari koneksi database yang berbeda-beda (dynamic DB connection dari Fase 2).
- **Paralel**: Fase 6 (AI), 7 (Export), dan 8 (WebSocket) bisa dikerjakan paralel setelah Fase 3 & 5 selesai.
- **Prioritas MVP**: Fase 0 → 1 → 2 → 3 (backend lengkap) + 4 → 5 (frontend lengkap) sudah cukup untuk MVP sebelum masuk fitur AI & advanced.
