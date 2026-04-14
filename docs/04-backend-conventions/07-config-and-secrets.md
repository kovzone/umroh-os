# Config & Secrets

Configuration uses **Viper** with JSON files and environment variable overrides. Secrets are never committed.

## Config file layout

Each service has:
- `<svc>/config.json.sample` — committed, with placeholder values
- `<svc>/config.json` — gitignored, populated locally
- Production: env vars override JSON values

### Structure

Mirrors the baseline template's config struct:

```json
{
  "app": {
    "name": "iam-svc",
    "env": "dev"
  },
  "api": {
    "port": 4001,
    "cors_origins": ["http://localhost:3000"]
  },
  "store": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "postgres",
    "database": "iam_db",
    "max_open_conns": 25,
    "max_idle_conns": 5
  },
  "external": {
    "iam_grpc": "localhost:50051",
    "broker_grpc": "localhost:50099"
  },
  "token": {
    "kind": "paseto",
    "symmetric_key": "dev-only-32-byte-key-do-not-use!!",
    "access_ttl": "15m",
    "refresh_ttl": "168h"
  },
  "otel_tracer": {
    "enabled": true,
    "endpoint": "localhost:4317",
    "service_name": "iam-svc"
  },
  "log": {
    "level": "info"
  }
}
```

## Loading

`util/config.Load()` (from the template) reads `config.json`, applies env-var overrides via Viper, and returns a typed `Config` struct. Never read config inside business logic — pass what you need into constructors.

## Env var overrides

Viper convention: nested keys use `_`:

| Config path | Env var |
|---|---|
| `api.port` | `API_PORT` |
| `store.password` | `STORE_PASSWORD` |
| `external.iam_grpc` | `EXTERNAL_IAM_GRPC` |
| `token.symmetric_key` | `TOKEN_SYMMETRIC_KEY` |

## Secrets

- **Never commit secrets.** `config.json` is gitignored. Only `config.json.sample` is committed.
- **Production secrets** come from env vars (typically injected by the orchestrator: K8s Secret, Cloud Run env, etc.).
- **API keys for third parties** (Midtrans, Xendit, WhatsApp, GCP) live in env vars and are loaded into the `external` section of config.
- **Token signing keys** must be at least 32 bytes for PASETO. Generate with `openssl rand -base64 32`.
- **Encrypt at rest** for any third-party credential persisted in the database (e.g. tenant-installed integrations).
- **Audit access to secrets** — every read of an encrypted credential goes through a service method that logs the access.

## What you do NOT do

- ❌ Hardcode API keys in source code.
- ❌ Print secrets in logs (even at debug level).
- ❌ Share dev `config.json` files via Slack — use a password manager.
- ❌ Reuse production keys in dev.
- ❌ Read `os.Getenv` directly inside business logic — go through the typed config.
