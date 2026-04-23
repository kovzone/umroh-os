# UmrohOS — Deployment Runbook

Log semua operasi manual di production server.

**Panduan singkat copy-paste (first deploy + GitHub secret):** [`docs/deployment/COPY-PASTE.md`](COPY-PASTE.md)

---

## 2026-04-23 — Initial Server Provisioning (Fase 1 + 2)

**Operator:** DevSecOps Agent  
**Server:** infra@216.176.238.161  
**OS:** Ubuntu 24.04 LTS (bare)

### Yang dikerjakan

Script provisioning dibuat di `docs/deployment/provision-server.sh`.  
Dijalankan manual oleh Lutfi dari Mac karena SSH key tidak accessible dari sandbox.

#### Fase 1 — System Provisioning

| # | Action | Status |
|---|--------|--------|
| 1 | `apt-get update && upgrade` | Dijalankan via script |
| 2 | Install prerequisites (ca-certificates, curl, gnupg) | Dijalankan via script |
| 3 | Docker Engine + Compose plugin dari docker.com APT repo | Dijalankan via script |
| 4 | `usermod -aG docker infra` | Dijalankan via script |
| 5 | UFW: deny all in, allow 22/80/4000 | Dijalankan via script |
| 6 | SSH hardening: PasswordAuthentication no | Dijalankan via script |
| 7 | fail2ban: maxretry=3, bantime=3600s | Dijalankan via script |
| 8 | Nginx install + enable service | Dijalankan via script |
| 9 | `mkdir -p /home/infra/umrohos` | Dijalankan via script |

#### Fase 2 — Nginx Config

| Config file | Listen | Upstream |
|-------------|--------|----------|
| `/etc/nginx/sites-available/umrohos` | `:80` | `localhost:3001` (core-web) |
| `/etc/nginx/sites-available/umrohos-api` | `:4000` | `localhost:4001` (gateway-svc) |

Default Nginx site dihapus dari sites-enabled.

### Cara menjalankan script

```bash
cd /path/to/umroh-os/docs/deployment
chmod +x provision-server.sh
./provision-server.sh
```

### Expected output setelah selesai

```
Docker: Docker version 27.x.x
UFW: active, ports 22/80/4000 open
Nginx: active
fail2ban: active
HTTP GET http://216.176.238.161/ → 502 (expected, upstream belum running)
```

502 dari HTTP test adalah **normal** — Nginx berjalan tapi container upstream (core-web) belum di-deploy.

---

## 2026-04-23 — Production CD Pipeline + Observability Stack

**Operator:** DevSecOps Agent

### Yang dikerjakan

| Artefak | Perubahan |
|---------|-----------|
| `docker-compose.prod.yml` | Tambah observability stack (prometheus, grafana, loki, tempo, fluent-bit); resource limits (`mem_limit`); log driver `json-file` max 10m/3 file; `restart: unless-stopped` semua service |
| `docs/deployment/env.prod.example` | Template lengkap semua env vars (postgres, IAM secrets, Grafana admin, OTEL) |
| `.github/workflows/deploy.yml` | CD pipeline: SSH → git pull → build --parallel → migrate → up -d → health check; secrets: `PROD_SSH_KEY` + `PROD_HOST` |
| `.gitignore` | Tambah `.env.prod` agar tidak masuk git |

---

## First Deploy — Prosedur Awal

Lakukan ini **sekali** sebelum CD pipeline bisa jalan otomatis.

### 1. Setup SSH key di GitHub

Di repo GitHub → Settings → Secrets and variables → Actions → New repository secret:

| Secret name | Value |
|-------------|-------|
| `PROD_SSH_KEY` | Isi private key SSH (`~/.ssh/lutfi_id_ed25519` atau deploy key baru) |
| `PROD_HOST` | `216.176.238.161` |

### 2. Clone repo ke server (sumber kanon: **kovzone/umroh-os**)

**Remote resmi** untuk produksi: `https://github.com/kovzone/umroh-os.git` (atau `git@github.com:kovzone/umroh-os.git` jika SSH).

```bash
ssh -i ~/.ssh/lutfi_id_ed25519 infra@216.176.238.161
# Di server:
cd /home/infra
git clone https://github.com/kovzone/umroh-os.git umrohos
cd /home/infra/umrohos
```

**Sudah pernah `git clone` dari org/repo lain?** Arahkan `origin` ke kovzone, lalu pull (cabang sesuai kebiasaan, mis. `dev` atau `main`):

```bash
cd /home/infra/umrohos
git remote -v
git remote set-url origin https://github.com/kovzone/umroh-os.git
# SSH (kalau deploy key / user key sudah ke kovzone):
#   git remote set-url origin git@github.com:kovzone/umroh-os.git
git fetch origin
git checkout dev   # atau: main
git pull origin dev
```

Deploy key / PAT wajib terdaftar di **repo** `kovzone/umroh-os` (bukan repo lama). Admin org/repo `kovzone` yang menambah deploy key bila perlu.

```bash
# Contoh: generate deploy key khusus server (opsional, repo private)
ssh-keygen -t ed25519 -C "prod@umroh-os" -f ~/.ssh/deploy_key -N ""
cat ~/.ssh/deploy_key.pub
# Tempel .pub ke: kovzone/umroh-os → Settings → Deploy keys (perlu peran Admin repo)
```

### 3. Buat file `.env.prod`

```bash
cp docs/deployment/env.prod.example /home/infra/umrohos/.env.prod
nano /home/infra/umrohos/.env.prod
# Ganti semua nilai CHANGE_ME dengan nilai production yang sesungguhnya
```

Nilai yang WAJIB diisi:
- `POSTGRES_PASSWORD` — gunakan `openssl rand -base64 24`
- `DATABASE_URL` — sesuaikan password dengan POSTGRES_PASSWORD
- `POSTGRES_CONNECTION_STRING` — idem
- `TOKEN_KEY` — gunakan `openssl rand -base64 32`
- `TOTP_ENCRYPTION_KEY` — gunakan `openssl rand -base64 32`
- `GF_ADMIN_PASSWORD` — ganti dari default

### 4. Jalankan first deploy

```bash
cd /home/infra/umrohos
chmod +x docs/deployment/deploy-prod.sh
./docs/deployment/deploy-prod.sh
```

Script ini akan:
1. Build semua Docker images
2. Start postgres + otel-collector
3. Tunggu postgres ready
4. Jalankan DB migrations
5. Start full stack
6. Print status semua containers

### 5. Verifikasi

```bash
# Gateway health
curl -f http://localhost:4001/system/live

# Core web (via Nginx)
curl -I http://216.176.238.161/

# Stack status
docker compose -f docker-compose.prod.yml --env-file .env.prod ps
```

---

## Rollback — Cara Kembali ke Versi Sebelumnya

### Rollback cepat (tanpa schema change)

```bash
ssh -i ~/.ssh/lutfi_id_ed25519 infra@216.176.238.161

cd /home/infra/umrohos

# Lihat commit history
git log --oneline -10

# Checkout ke commit sebelumnya
git checkout <previous-sha>

# Rebuild dan restart services yang terpengaruh
docker compose -f docker-compose.prod.yml --env-file .env.prod \
  build --parallel gateway-svc core-web iam-svc

docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

# Verifikasi
curl -f http://localhost:4001/system/live
```

### Rollback dengan schema change (migrasi DB)

**PENTING:** Jangan jalankan `migrate down` di production tanpa persetujuan Lutfi.

```bash
# 1. Backup database DULU
docker compose -f docker-compose.prod.yml --env-file .env.prod \
  exec -T postgres pg_dump -U postgres umrohos \
  > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. Konfirmasi ke Lutfi bahwa rollback migration aman
# 3. Jalankan migrate down hanya setelah persetujuan eksplisit
docker compose -f docker-compose.prod.yml --env-file .env.prod \
  --profile tooling run --rm migrate \
  -path /migrations -database "${DATABASE_URL}" down 1

# 4. Checkout kode lama dan restart
git checkout <previous-sha>
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

---

## Cara Check Logs

### Logs semua services (stream)

```bash
docker compose -f docker-compose.prod.yml logs -f
```

### Logs service tertentu

```bash
# gateway-svc
docker compose -f docker-compose.prod.yml logs -f gateway-svc

# iam-svc
docker compose -f docker-compose.prod.yml logs -f iam-svc

# core-web
docker compose -f docker-compose.prod.yml logs -f core-web

# postgres
docker compose -f docker-compose.prod.yml logs -f postgres
```

### Logs dengan jumlah baris terbatas

```bash
docker compose -f docker-compose.prod.yml logs --tail=100 gateway-svc
```

### Lihat log file JSON langsung

Log disimpan di Docker default path dengan driver `json-file`, max 10MB/3 file:
```bash
# Cari path log container
docker inspect gateway-svc | grep -i logpath
# Contoh: /var/lib/docker/containers/<id>/<id>-json.log
```

---

## Akses Grafana via SSH Tunnel

Grafana dan Prometheus hanya bisa diakses dari dalam server (tidak expose ke publik).
Gunakan SSH tunnel untuk akses dari Mac/laptop:

### Buka tunnel (jalankan di Mac/laptop)

```bash
ssh -N -L 3000:localhost:3000 \
       -L 9090:localhost:9090 \
       -L 3100:localhost:3100 \
       -i /Users/lutfiaf/.ssh/lutfi_id_ed25519 \
       infra@216.176.238.161
```

Flag `-N` artinya tidak membuka shell — hanya tunnel.

### Akses browser

| Service | URL |
|---------|-----|
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| Loki (health) | http://localhost:3100/ready |

### Login Grafana

- Username: `admin` (atau nilai `GF_ADMIN_USER` di `.env.prod`)
- Password: nilai `GF_ADMIN_PASSWORD` di `.env.prod`

**Wajib ganti password default** sebelum mengekspos tunnel ke jaringan lain.

### Import dashboard

Dashboard tersedia di folder `grafana/dashboards/` di repo. Grafana di-provision otomatis via `grafana/provisioning/` — dashboard akan tersedia langsung setelah container start.

---

## Fase 3–5 — Compose prod, secrets, CD (dokumentasi di repo)

| Fase | Artefak | Keterangan |
|------|---------|------------|
| 3 | `docker-compose.prod.yml` (root) | Stack: Postgres, OTel (debug), semua service Go, `core-web` (target `prod`), gateway `4001:4000`, web `3001:3000`. Config mount dari `config.json.sample` + override env. |
| 3 | `monitoring/otel-collector-config.prod.yaml` | Tanpa Jaeger/Tempo; alamat gRPC `4317` sama pattern dev. |
| 3 | `docs/deployment/deploy-prod.sh` | `build` → up postgres+otel → `pg_isready` → **migrate** (`--profile tooling run --rm migrate`) → `up -d` full. |
| 4 | `env.prod.sample` (root) | Template variabel. **Salin ke** `/home/infra/umrohos/.env.prod` di server, isi password + token, **jangan commit**. |
| 5 | `.github/workflows/deploy.yml` | `workflow_dispatch` + `push` ke `main` (path-filtered); SSH + `git pull` + `deploy-prod.sh`. Butuh GitHub secrets: `PROD_SSH_HOST`, `PROD_SSH_USER`, `PROD_SSH_KEY`. |

**Migrasi database:** tool `golang-migrate` via service `migrate` (bukan lewat `gateway-svc`).

**First time di server (ringkas):**

1. `git clone` **`kovzone/umroh-os`** ke `/home/infra/umrohos` (HTTPS/PAT atau SSH; deploy key di **repo** kovzone).
2. `cp env.prod.sample .env.prod` → edit secret.
3. `chmod +x docs/deployment/deploy-prod.sh && ./docs/deployment/deploy-prod.sh`

Atau: `docker compose -f docker-compose.prod.yml --env-file .env.prod up -d` **setelah** migrate (lihat `deploy-prod.sh` urutan pasti).

---

## Port Architecture

```
Internet
    │
    ▼
UFW Firewall
    ├── :22    SSH (key-only)
    ├── :80    Nginx → core-web (:3001 internal)
    └── :4000  Nginx → gateway-svc (:4001 internal)

Internal only (SSH tunnel untuk akses lokal):
    ├── :3000  Grafana
    ├── :9090  Prometheus
    ├── :3100  Loki
    └── :5432  PostgreSQL (TIDAK expose publik)
```

## Akses monitoring

```bash
ssh -L 3000:localhost:3000 -L 9090:localhost:9090 \
    -i /Users/lutfiaf/.ssh/lutfi_id_ed25519 infra@216.176.238.161
```

Kemudian buka browser: http://localhost:3000 untuk Grafana.
