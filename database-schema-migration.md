# Database Schema — Visual Mesin (Target: Golang + React)

## 1. Arsitektur Database

Sistem menggunakan arsitektur **Multi-Database**:

```
┌─────────────────────────────────────────────────────┐
│              Application Database (Local)            │
│  Menyimpan data aplikasi: users, roles, chat, dll   │
├─────────────────────────────────────────────────────┤
│  Tabel: users, roles, permissions, ai_chat_history, │
│         ai_schema_map, db_connections, dll           │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│              Resource Database(s) (External)         │
│  Menyimpan data produksi mesin (bisa multi-server)  │
├─────────────────────────────────────────────────────┤
│  Tabel: rtba1..3, rtbc1..4, curtire, trimming, dll  │
│         rteex1, rteex2, rteex3head, recorddatacyclic│
│         alarm_history, material, order_report, dll   │
└─────────────────────────────────────────────────────┘
```

**Konsep:** Aplikasi punya DB sendiri untuk data operasional. Data produksi mesin bisa berada di DB terpisah (bahkan server berbeda) dan dikonfigurasi secara dinamis via tabel `resource_db_configs` atau `db_connections`.

---

## 2. Application Tables (Local DB)

### 2.1. `users`

Menyimpan data user aplikasi.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| nip | VARCHAR(50) UNIQUE | NIP/NIK karyawan |
| user_id | VARCHAR(100) UNIQUE | Login ID |
| user_name | VARCHAR(100) | Nama user |
| user_level | ENUM('admin','eng','tech','prod') | Level/Role user |
| email | VARCHAR(255) UNIQUE | Email (opsional) |
| email_verified_at | TIMESTAMP NULL | Verifikasi email |
| password | VARCHAR(255) | Password ter-hash (bcrypt) |
| avatar_url | VARCHAR(255) | URL avatar |
| remember_token | VARCHAR(100) | Token session |
| department | VARCHAR(100) | Departemen |
| jabatan | VARCHAR(100) | Jabatan |
| themes_settings | JSON | Preferensi tema |
| timestamp | TIMESTAMP(3) | Waktu record |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diupdate |

### 2.2. `roles`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| name | VARCHAR(255) | Nama role (admin, eng, tech, prod) |
| guard_name | VARCHAR(255) | Guard (default: web) |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |
| UNIQUE(name, guard_name) | | |

### 2.3. `permissions`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| name | VARCHAR(255) | Nama permission |
| guard_name | VARCHAR(255) | Guard |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |
| UNIQUE(name, guard_name) | | |

### 2.4. `model_has_roles`

Relasi many-to-many user ke role.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| role_id | BIGINT UNSIGNED | FK ke roles.id |
| model_type | VARCHAR(255) | Morph type (App\\Models\\User) |
| model_id | BIGINT UNSIGNED | FK ke users.id |
| PRIMARY KEY(role_id, model_id, model_type) | | |

### 2.5. `model_has_permissions`

Relasi many-to-many user ke permission.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| permission_id | BIGINT UNSIGNED | FK ke permissions.id |
| model_type | VARCHAR(255) | Morph type |
| model_id | BIGINT UNSIGNED | FK ke users.id |
| PRIMARY KEY(permission_id, model_id, model_type) | | |

### 2.6. `role_has_permissions`

Relasi many-to-many role ke permission.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| permission_id | BIGINT UNSIGNED | FK ke permissions.id |
| role_id | BIGINT UNSIGNED | FK ke roles.id |
| PRIMARY KEY(permission_id, role_id) | | |

### 2.7. `ai_schema_map`

Konfigurasi intent AI dan mapping tabel untuk AI Chat.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| intent_name | VARCHAR(100) | Nama intent (alarm, produksi, kualitas, dll) |
| keywords | JSON | Array kata kunci untuk matching intent |
| tables_involved | JSON | Array nama tabel relevan |
| schema_context | TEXT | Deskripsi kolom dalam bahasa natural untuk konteks AI |
| few_shot_examples | JSON NULL | Array [{question, sql}] contoh few-shot |
| description | VARCHAR(255) | Deskripsi untuk admin |
| is_active | BOOLEAN | Aktif/nonaktif |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### 2.8. `ai_chat_history`

Riwayat chat user dengan AI.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| user_id | BIGINT UNSIGNED | FK ke users.id |
| session_id | VARCHAR(255) | ID sesi chat |
| question | TEXT | Pertanyaan user |
| detected_intent | VARCHAR(100) | Intent terdeteksi |
| generated_sql | TEXT | SQL yang digenerate |
| sql_status | ENUM('pending','valid','invalid','error') | Status validasi SQL |
| ai_response | LONGTEXT | Response AI |
| status | ENUM('queued','processing','completed','failed','rejected') | Status proses |
| rejection_reason | VARCHAR(255) | Alasan ditolak |
| started_at | TIMESTAMP NULL | Waktu mulai proses |
| completed_at | TIMESTAMP NULL | Waktu selesai |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |
| INDEX(user_id) | | |
| INDEX(status) | | |
| INDEX(created_at) | | |

### 2.9. `db_connections`

Konfigurasi koneksi database dinamis.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| name | VARCHAR(255) UNIQUE | Nama koneksi |
| driver | ENUM('mysql','mariadb','postgresql','sqlite') | Driver |
| host | VARCHAR(255) | Host/IP |
| port | INT | Port (default 3306) |
| database | VARCHAR(255) | Nama database |
| username | VARCHAR(255) | Username |
| password | TEXT | Password (terenkripsi AES-256) |
| is_active | BOOLEAN | Status aktif |
| is_last_test_success | BOOLEAN NULL | Hasil test terakhir |
| last_tested_at | TIMESTAMP NULL | Waktu test terakhir |
| last_test_message | TEXT | Pesan test terakhir |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |
| INDEX(is_active) | | |
| INDEX(driver) | | |

### 2.10. `resource_db_configs`

Konfigurasi database per-resource/mesin.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| resource_name | VARCHAR(255) UNIQUE | Nama resource (rteex1, rtba1, dll) |
| label | VARCHAR(255) | Label human-readable |
| driver | VARCHAR(20) | Driver (default mariadb) |
| host | VARCHAR(255) | Host |
| port | INT | Port (default 3306) |
| database | VARCHAR(255) | Nama database |
| username | VARCHAR(255) | Username |
| password | TEXT | Password terenkripsi |
| is_active | BOOLEAN | Status aktif |
| is_last_test_success | BOOLEAN NULL | Hasil test |
| last_tested_at | TIMESTAMP NULL | |
| last_test_message | TEXT NULL | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |
| INDEX(is_active) | | |
| INDEX(resource_name) | | |

### 2.11. `notifications`

Notifikasi database (untuk filament/admin).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | CHAR(36) UUID | Primary Key |
| type | VARCHAR(255) | Tipe notifikasi |
| notifiable_type | VARCHAR(255) | Morph type |
| notifiable_id | BIGINT UNSIGNED | ID penerima |
| data | TEXT | Data JSON |
| read_at | TIMESTAMP NULL | Waktu dibaca |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### 2.12. `activity_log`

Log aktivitas pengguna.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| log_name | VARCHAR(255) NULL | Nama log |
| description | TEXT | Deskripsi aktivitas |
| subject_type | VARCHAR(255) NULL | |
| subject_id | BIGINT UNSIGNED NULL | |
| causer_type | VARCHAR(255) NULL | |
| causer_id | BIGINT UNSIGNED NULL | |
| properties | JSON NULL | Data tambahan |
| batch_uuid | CHAR(36) NULL | UUID batch |
| event | VARCHAR(255) NULL | Event |
| created_at | TIMESTAMP | |

### 2.13. `cache` & `cache_locks`

Cache store Redis fallback.

### 2.14. `jobs`, `job_batches`, `sessions`

Queue & session management.

---

## 3. Resource Tables (Production Data)

Tabel-tabel berikut menyimpan data produksi mesin. Bisa berada di database terpisah (dikonfigurasi via `resource_db_configs`).

### 3.1. Building / Assembly (RTBA)

#### `rtba1`, `rtba2`, `rtba3`

Tire building machine cycle time data. Struktur identik untuk 3 mesin.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | Primary Key |
| Timestamp | VARCHAR(50) | Waktu record |
| barcode | VARCHAR(50) | Barcode tire |
| specification | VARCHAR(50) | Spesifikasi |
| pattern | VARCHAR(100) | Pattern |
| GT_CT | VARCHAR(50) | Green tire cycle time |
| preAssy_pudCT | DECIMAL(10,1) | Pre-assembly PUD CT |
| chafer_pudCT | DECIMAL(10,1) | Chafer PUD CT |
| bodyplyApply_pudCT | DECIMAL(10,1) | Body ply apply PUD CT |
| bodyplyStitch_pudCT | DECIMAL(10,1) | Body ply stitch PUD CT |
| beadTrf&BECapply_pudCT | DECIMAL(10,1) | Bead transfer & BEC apply PUD CT |
| transfer_carcas_pudCT | DECIMAL(10,1) | Transfer carcas PUD CT |
| PUD_CT | DECIMAL(10,1) | Total PUD cycle time |
| belt1_btdCT | DECIMAL(10,1) | Belt 1 BTD CT |
| belt2_btdCT | DECIMAL(10,1) | Belt 2 BTD CT |
| belt3_btdCT | DECIMAL(10,1) | Belt 3 BTD CT |
| belt4_btdCT | DECIMAL(10,1) | Belt 4 BTD CT |
| tread_btdCT | DECIMAL(10,1) | Tread BTD CT |
| trdstitching_btdCT | DECIMAL(10,1) | Tread stitching BTD CT |
| BeltTreadtransfer_btdCT | DECIMAL(10,1) | Belt tread transfer BTD CT |
| BTD_CT | DECIMAL(10,1) | Total BTD cycle time |
| CarcassTRF_bddCT | DECIMAL(10,1) | Carcass transfer BDD CT |
| shapping_bddCT | DECIMAL(10,1) | Shapping BDD CT |
| treadStict_bddCT | DECIMAL(10,1) | Tread stitch BDD CT |
| turnUp_bddCT | DECIMAL(10,1) | Turn up BDD CT |
| SWStict_bddCT | DECIMAL(10,1) | Sidewall stitch BDD CT |
| gt_unload_bddCT | DECIMAL(10,1) | Green tire unload BDD CT |
| BDD_CT | DECIMAL(10,1) | Total BDD cycle time |
| treadlength | DECIMAL(10,1) | Panjang tread |
| shapping_press | DECIMAL(10,1) | Shapping pressure |
| beadlock_press | DECIMAL(10,1) | Beadlock pressure |
| BPstc_press | DECIMAL(10,1) | BP stitch pressure |
| btoc | DECIMAL(10,1) | BTOC |
| BTB_width | DECIMAL(10,1) | BTB width |
| CTR_width | DECIMAL(10,1) | CTR width |
| pud_oc_col | DECIMAL(10,1) | PUD OC col |
| pud_oc_exp | DECIMAL(10,1) | PUD OC exp |
| TreadHead_midP1 | DECIMAL(10,1) | Tread head mid P1 |
| TreadHead_edgeP6 | DECIMAL(10,1) | Tread head edge P6 |
| TreadHead_treadPos | DECIMAL(10,1) | Tread head position |
| ... (dan seterusnya untuk TreadMid, TreadTail) | | |
| INDEX(barcode) | | |

### 3.2. Building Quality (RTBC)

#### `rtbc1`, `rtbc2`, `rtbc3`, `rtbc4`

Building quality inspection. **Struktur identik dengan RTBA** (recid, Timestamp, barcode, specification, pattern, GT_CT, PUD/BTD/BDD CT, dimensi tread).

### 3.3. Building Evaluation (RTBE)

#### `rtbe1`, `rtbe2`

Building evaluation data. **Struktur identik dengan RTBA** (recid, Timestamp, barcode, specification, pattern, GT_CT, semua kolom PUD/BTD/BDD CT).

### 3.4. Extruder

#### `rteex1`

Data kualitas extruder (berat, cutting, line speed).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| trxtime | TIMESTAMP | Waktu transaksi |
| date_shift | VARCHAR(50) | Tanggal shift |
| ydate_shift | VARCHAR(50) | Yesterday date shift |
| machine | VARCHAR(25) | Kode mesin |
| size | VARCHAR(25) | Ukuran produk |
| recipe | VARCHAR(15) | Recipe |
| SpecWS | VARCHAR(15) | Spec weight scale |
| ActWS | VARCHAR(15) | Actual weight scale |
| USLWS / LSLWS | VARCHAR(15) | Upper/Lower spec limit weight scale |
| UCLWS / LCLWS | VARCHAR(15) | Upper/Lower control limit weight scale |
| SetCutting | VARCHAR(15) | Set cutting |
| SpecCutting | VARCHAR(15) | Spec cutting |
| LineSpeed | VARCHAR(15) | Line speed |
| Deviasi | VARCHAR(15) | Deviasi |
| OK / NG / Total | VARCHAR(15) | OK/NG/Total count |
| PersentaseOK / NG | VARCHAR(15) | Persentase OK/NG |
| UCLRS / LCLRS | VARCHAR(15) | Control limits rubber sheet |
| USLRS / LSLRS | VARCHAR(15) | Spec limits rubber sheet |
| SpecRS / ActRS | VARCHAR(15) | Spec/Actual rubber sheet |
| INDEX(trxtime, machine) | | |

#### `rteex2`

Statistik extruder per size per mesin.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| mesin | VARCHAR(20) | Nama mesin |
| trxtime | TIMESTAMP | Waktu |
| date_shift | VARCHAR(30) | Tanggal shift |
| ydate_shift | VARCHAR(30) | |
| ok / ng / total | VARCHAR(10) | Count |
| ngover / ngunder | VARCHAR(10) | Detail NG |
| okpersen / ngpersen | VARCHAR(10) | Persentase |
| size | VARCHAR(20) | Ukuran |
| spec | VARCHAR(20) | Spesifikasi |
| lsl / usl | VARCHAR(10) | Spec limits |
| actual | VARCHAR(10) | Actual value |
| speedline / speedbooking | VARCHAR(20) | Kecepatan |

#### `rteex3head`

Data 3-head extruder (current, temperature, RPM).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| trxtime | TIMESTAMP | Waktu |
| date_shift | VARCHAR(25) | |
| ydate_shift | VARCHAR(25) | |
| machine | VARCHAR(15) | Mesin |
| size | VARCHAR(15) | Ukuran |
| ExtUpCurrent | INT | Current extruder upper |
| ExtMidCurrent | INT | Current extruder middle |
| ExtLowCurrent | INT | Current extruder lower |
| TempHeadExtUp/Mid/Low | INT | Temperature head |
| SpeedFeedUp/Mid/Low | INT | Speed feed |
| SpecRPMScrewUp/Mid/Low | INT | Spec RPM screw |
| SpeedRPMUp/Mid/Low | INT | Actual RPM screw |

### 3.5. Curing

#### `curtire`

Data proses curing/vulkanisasi ban.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| bcd | VARCHAR(50) | Barcode |
| itemcode | VARCHAR(50) | Kode item |
| codespec | VARCHAR(50) | Kode spesifikasi |
| mcn | VARCHAR(50) | Mesin |
| mcntype | VARCHAR(50) | Tipe mesin |
| operator | VARCHAR(50) | Operator |
| nip | VARCHAR(50) | NIP operator |
| curein | VARCHAR(50) | Waktu cure in |
| cureout | VARCHAR(50) | Waktu cure out |
| dateshift | VARCHAR(50) | Tanggal shift |
| grpshift | VARCHAR(50) | Group shift |
| rsn | VARCHAR(50) | RSN |
| extemp / extemp_jdg | VARCHAR(50) | External temperature + judgment |
| intemp / intemp_jdg | VARCHAR(50) | Internal temperature + judgment |
| platen / platen_jdg | VARCHAR(50) | Platen temperature + judgment |
| jacket / jacket_jdg | VARCHAR(50) | Jacket temperature + judgment |
| inpress_st / _jdg | VARCHAR(50) | Steam pressure + judgment |
| inpress_n2 / _jdg2 | VARCHAR(50) | N2 pressure + judgment |
| curtime / curtime_jdg | VARCHAR(50) | Cure time + judgment |
| finaljdg | VARCHAR(50) | Final judgment (OK/NG) |
| finaldfc | VARCHAR(50) | Final defect code |
| desc | VARCHAR(50) | Deskripsi |
| eventdate | DATETIME | Tanggal event |
| created_at | TIMESTAMP | |
| INDEX(bcd) | | |

#### `item_measurement`

Data pengukuran curing (left/right cavity).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| tanggal | VARCHAR(50) | Tanggal |
| operator | VARCHAR(50) | |
| grup | VARCHAR(50) | |
| nip | VARCHAR(50) | |
| shift | VARCHAR(50) | |
| date_shift | VARCHAR(50) | |
| codespec | VARCHAR(50) | |
| mesinl / mesinr | VARCHAR(50) | Mesin left/right |
| bcdl / bcdr | VARCHAR(50) | Barcode left/right |
| rsnl / rsnr | VARCHAR(50) | |
| cureinl/r, cureoutl/r | VARCHAR(50) | |
| Parameter L/R: step, extemp, platen, jacket, intemp, inpressN2, inpressSt, curtime | | |
| Judgment per parameter (extemp_jdgL/R, dll) | | |
| jdg | VARCHAR(50) | Final judgment |
| eventdate | DATETIME | |
| created_at | TIMESTAMP | |

### 3.6. Green Tire

#### `gtentire`

Tracking green tire (building -> curing).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| bc_entried | VARCHAR(50) | Barcode entry |
| serialnumb | VARCHAR(50) | Serial number |
| bld_date | VARCHAR(50) | Tanggal building |
| bld_mcn01 | VARCHAR(50) | Mesin building |
| bld_shift | VARCHAR(50) | Shift building |
| cur_in / cur_out | VARCHAR(50) | Curing in/out |
| cur_shift | VARCHAR(50) | Shift curing |
| cur_mcn | VARCHAR(50) | Mesin curing |
| cur_opr01 | VARCHAR(50) | Operator curing |
| jdge_date | VARCHAR(50) | Tanggal judgment |
| jdge | VARCHAR(50) | Judgment |
| probcode | VARCHAR(50) | Problem code |
| status | VARCHAR(50) | Status |
| pic | VARCHAR(50) | PIC |
| whs_in / whs_out | VARCHAR(50) | Warehouse in/out |
| sarana | VARCHAR(50) | Sarana transport |
| INDEX(bc_entried) | | |

### 3.7. Trimming

#### `trimming`

Data proses trimming ban.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| Machine_number | VARCHAR(50) | Nomor mesin |
| Tirecode | VARCHAR(50) | Kode tire |
| Barcode_tire | VARCHAR(50) | Barcode |
| Start_Trimming | DATETIME | Waktu mulai |
| End_Trimming | DATETIME | Waktu selesai |
| Duration_Process | TIME | Durasi proses |
| Tires_Number | VARCHAR(50) | Jumlah ban |
| Machine_Speed | VARCHAR(50) | Kecepatan mesin |
| Pressure_Trimming | VARCHAR(50) | Pressure |
| Temperature | VARCHAR(50) | Temperatur |
| Operator_ID | VARCHAR(50) | Operator |

#### `rtc-tr1`

Real-time trimming quality.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| Trimming_MachineNumber | VARCHAR(100) | |
| Tire_Code | VARCHAR(100) | |
| Barcode_Tire | VARCHAR(100) | |
| Start_Triming | TIMESTAMP | |
| End_Triming | TIMESTAMP | |
| Duration_TrimingProcess | TIME | |
| Number_tiresTrimed | DECIMAL(20,6) | |
| Triming_MachineSpeed | DECIMAL(20,6) | |
| Pressure_Triming | DECIMAL(20,6) | |
| Temperature_Triming | DECIMAL(20,6) | |
| Triming_OperatorID | VARCHAR(50) | |

### 3.8. Material

#### `material`

Tracking material/compound.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| item | VARCHAR(50) | Nama material |
| bc_entried | VARCHAR(50) | Barcode |
| mcn | VARCHAR(50) | Mesin |
| opr | VARCHAR(50) | Operator |
| txndate | VARCHAR(50) | Tanggal transaksi |
| shift | VARCHAR(50) | Shift |
| qty | VARCHAR(50) | Quantity |
| stock | VARCHAR(50) | Stock |
| lokasi | VARCHAR(50) | Lokasi |
| sarana | VARCHAR(50) | Sarana |
| jdge | VARCHAR(50) | Judgment |

### 3.9. Monitoring & Yield

#### `monitoringtl1`

Monitoring line speed TL1.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| date_shift | VARCHAR(50) | |
| ydate_shift | VARCHAR(50) | |
| codesize | VARCHAR(20) | |
| SpecSpeed / SpeedSP / SpeedSrv / SpeedProd / SpeedActual | VARCHAR | Speed parameters |
| ToleransiMin / ToleransiPlus | VARCHAR(20) | |
| AverageSpeed | VARCHAR(20) | |
| TotalDataUnder | VARCHAR(20) | |
| TotalOutput | VARCHAR(20) | |

#### `rtl-tl1`

Line speed tracking with operators.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| tgl / jam / shift | VARCHAR | |
| codesize / SpecSpeed / SpeedSP / SpeedSrv / SpeedProd / SpeedA | VARCHAR | |
| tol-- / tol++ | VARCHAR(20) | |
| leader / operator1-5 | VARCHAR(20) | |
| average | VARCHAR(50) | |

#### `rtltl1`

Yield summary dengan breakdown defect.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| date_shift / ydate_shift | VARCHAR(50) | |
| txndate | VARCHAR(100) | |
| mesin | VARCHAR(20) | |
| size | VARCHAR(25) | |
| TotalProd / TotalOK / TotalReject | INT | |
| PersenTotalOK / PersenTotalReject | VARCHAR(50) | |
| LebarOverSq / LebarUnSq | INT | Width defects |
| TebalOverSq / TebalUnderSq | INT | Thickness defects |
| LebarOverAp / LebarUnderAp | INT | Width defects (AP) |
| TebalOverAp / TebalUnderAp | INT | Thickness defects (AP) |
| Metal | INT | Metal contamination |
| LebarTotalUnder / LebarTotalOver | INT | Total width defects |
| Persentase untuk semua kategori | VARCHAR | |

### 3.10. Alarm

#### `alarm_history`

Riwayat alarm mesin.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| timeOn | TIMESTAMP(3) | Waktu alarm ON |
| timeOff | TIMESTAMP(3) | Waktu alarm OFF |
| source | TEXT | Sumber alarm |
| message | TEXT | Pesan alarm |

### 3.11. Recipe

#### `recipe1`

Master recipe produksi.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| size | TEXT | Ukuran |
| speed_extupper / speed_extmiddle / speed_extlower | TEXT | Speed extruder |
| speed_line / speed_calender | TEXT | Kecepatan line/calender |
| run_scale / run_scale_up / run_scale_low | TEXT | Run scale |
| weight_scale / weight_up / weight_low | TEXT | Weight scale spec |
| width / width_up / width_low | TEXT | Width spec |
| tcu_upper/middle/lower_screw/barrel/head | TEXT | TCU temperatures |
| tcu_preformerup/down | TEXT | |
| tcu_calender parameters | TEXT | |
| gap_calender / cutter_calender | TEXT | |
| length_skiver | TEXT | |
| material_extupper/middle/lower/calender/rubbersheet/die | TEXT | |
| created_at / updated_at | TIMESTAMP | |

#### `recipe1queue`

Antrian produksi recipe.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| order_id | VARCHAR(20) | Nomor order |
| recipe_id | BIGINT UNSIGNED | FK ke recipe1.id |
| size | TEXT | |
| qty | TEXT | |
| shift | TEXT | |
| user_name | TEXT | |
| time_create / time_start / time_finish | TIMESTAMP(3) | |

#### `recipe_history`

Riwayat revisi recipe.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| recipe_id | BIGINT UNSIGNED | FK ke recipe1.id |
| created_at | DATETIME | |
| created_by | BIGINT UNSIGNED | FK ke users.id |
| modified_at | DATETIME NULL | |
| modified_by | BIGINT UNSIGNED NULL | FK ke users.id |
| revision | VARCHAR(20) | Nomor revisi |
| note | TEXT NULL | Catatan |

### 3.12. Order & Batch

#### `order_report`

Laporan produksi per order.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| order_id | VARCHAR(20) | Nomor order |
| recipe_id | INT | |
| recipe | TEXT | |
| set_qty / act_qty | TEXT | Target vs actual |
| shift | TEXT | |
| user_name | TEXT | |
| time_create / time_start / time_finish | TIMESTAMP(3) | |

#### `batch_report`

Laporan batch produksi.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| order_id | VARCHAR(20) | |
| batch_id | VARCHAR(100) | |
| user_name / shift / recipe | TEXT | |
| spec_weightscale / _up / _low | TEXT | |
| act_weightscale | TEXT | |
| spec_lengthskiver / act_lengthskiver | TEXT | |
| timestamp | TIMESTAMP(3) | |

### 3.13. Cyclic & PCS Data

#### `recorddatacyclic`

Data cyclic extruder (90+ kolom sensor).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| code_recipe | VARCHAR(50) | Kode recipe |
| timestamp_record | DATETIME | Waktu record |
| Material: comp_porkchop, comp_e2middle, comp_e3lower, comp_qsm60/120, preformer_name, finaldie_name | VARCHAR | |
| Speed: aspeed_e1/e2/e3/shrink3, sspeed_e1/e2/e3/shrink3 | FLOAT | Actual/Set speed extruder |
| Weight: sberat_control, aberat_control, tolerances | FLOAT | |
| Width: slebar_control, alebar_control, tolerances | FLOAT | |
| TCU: atcu/stcu screw/barrel/zone temperatures (PC, E1, E2, E3) | FLOAT | ~30 kolom TCU |
| Head temperatures: PC, E1, E2, E3, preformer, calender | FLOAT | |
| E120 temperatures | FLOAT | |
| Position incline/decline (d1-d7) | FLOAT | Actual & Set |
| Speed incline/decline (d1-d7) | FLOAT | |
| Pressure incline/decline (d1-d7) | FLOAT | |
| ampere | FLOAT | |
| prod_length_cooling | FLOAT | |

#### `recorddatapcs`

PCS recording data.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| recipe_code | VARCHAR(50) | |
| waktu | DATETIME | |
| sberat_finish / aberat_finish | FLOAT | Spec/Actual weight |
| scut_skiver / acut_skiver | FLOAT | Spec/Actual cut skiver |
| slength / alength | FLOAT | Spec/Actual length |
| switdhtol / awidthtol | FLOAT | Width tolerance |
| Tolerance columns (--) / (-) / (+) / (++) | FLOAT | |
| auto_rejectStat | FLOAT | |
| prod_OK / prod_NG / prod_Tot | FLOAT | |

### 3.14. Datalog

#### `datalog`

Datalog extruder (spec vs actual).

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| datetime | TIMESTAMP(3) | |
| recipe1 / recipe2 | TEXT | |
| spec_speedextupper/middle/lower, act_speed... | TEXT | |
| act_ampextupper/middle/lower | TEXT | |
| spec_speedline/calender, act_speed... | TEXT | |
| spec_runscale/up/low, act_runscale/up/low | TEXT | |
| spec_weightscale/up/low, act_weightscale | TEXT | |
| TCU data: spec/act untuk upper/middle/lower/preformer/calender/calenderext (masing-masing 4 channel) | TEXT | 40+ kolom TCU |
| spec_gapcalender, act_gapcalender | TEXT | |
| spec_cuttercalender, act_cuttercalender | TEXT | |
| act_compoundcalender | TEXT | |
| spec_lengthskiver, act_lengthskiver | TEXT | |
| spec_width/up/upline/low/lowline, act_width | TEXT | |

### 3.15. Supporting Tables

#### `mastermcn`

Master data mesin.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | INT AUTO_INCREMENT | |
| process | VARCHAR(50) | Proses (Extruder, Building, Curing, dll) |
| type | VARCHAR(50) | Tipe mesin |
| mcn | VARCHAR(50) | Kode mesin |
| mcnside | VARCHAR(50) | Sisi mesin (L/R) |

#### `bpbl`

BPBL inspection data.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| recid | BIGINT UNSIGNED AUTO_INCREMENT | |
| bc_entried | VARCHAR(50) | Barcode |
| mcncode | VARCHAR(50) | Mesin |
| probcode | VARCHAR(50) | Problem code |
| jdge | VARCHAR(50) | Judgment |
| opr / oprname | VARCHAR(50) | Operator |
| jdge_date / date_shift | VARCHAR(50) | |

#### `rsc_pc1`

Selection/sorting machine data.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | |
| machine | VARCHAR(20) | |
| position | VARCHAR(40) | |
| code | VARCHAR(10) | |
| barcode | VARCHAR(30) | |
| totalseleksi | VARCHAR(20) | |
| shift | VARCHAR(10) | |
| starttime / stoptime | VARCHAR(40) | |

---

## 4. Entity Relationship Summary

```
USERS ──1:N── AI_CHAT_HISTORY
USERS ──N:M── ROLES (via model_has_roles)
USERS ──N:M── PERMISSIONS (via model_has_permissions)
ROLES ──N:M── PERMISSIONS (via role_has_permissions)

RECIPE1 ──1:N── RECIPE1QUEUE
RECIPE1 ──1:N── RECIPE_HISTORY

DB_CONNECTIONS (konfigurasi DB dinamis)
RESOURCE_DB_CONFIGS (konfigurasi DB per resource)
AI_SCHEMA_MAP (konfigurasi intent AI)
```

### Grup Tabel Resource (setiap grup biasanya di DB terpisah):

**Building (RTBA/RTBC/RTBE):** rtba1, rtba2, rtba3, rtbc1, rtbc2, rtbc3, rtbc4, rtbe1, rtbe2
**Extruder:** rteex1, rteex2, rteex3head, recorddatacyclic, recorddatapcs, datalog
**Curing:** curtire, item_measurement
**Trimming:** trimming, rtc-tr1
**Monitoring:** monitoringtl1, rtl-tl1, rtltl1
**Order/Recipe:** order_report, batch_report, recipe1, recipe1queue, recipe_history
**Material:** material
**Alarm:** alarm_history
**Master:** mastermcn, rsc_pc1, bpbl, gtentire

---

## 5. Migration Strategy (Go)

Untuk Go, gunakan **golang-migrate/migrate** dengan file SQL:

```
backend/migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_permission_tables.up.sql
├── 000003_create_ai_schema_map_table.up.sql
├── 000004_create_ai_chat_history_table.up.sql
├── 000005_create_db_connections_table.up.sql
├── 000006_create_resource_db_configs_table.up.sql
├── 000007_create_notifications_table.up.sql
└── ...
```

Tabel resource (produksi) **tidak perlu migrasi** karena sudah ada di database masing-masing. Backend Go cukup membaca strukturnya secara dinamis.

---

## 6. Dynamic DB Connection Flow (Go)

```
1. Admin menambah konfigurasi DB via API → INSERT ke resource_db_configs
2. Backend Go menyimpan pool koneksi di memory map:
   map[string]*sql.DB {"rteex1": dbPool1, "rtba1": dbPool2, ...}
3. Saat request ke endpoint /api/v1/rteex1:
   - Backend cek resource_db_configs untuk resource "rteex1"
   - Ambil/simpan pool koneksi
   - Execute query ke database target
4. Setiap pool punya health check + auto-reconnect
```
