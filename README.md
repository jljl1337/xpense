# xpense

[![Source](https://img.shields.io/badge/Source-GitHub-blue?logo=github)](https://github.com/jljl1337/xpense)
[![Docker](https://img.shields.io/docker/pulls/jljl1337/xpense?logo=docker&label=jljl1337%2Fxpense)](https://hub.docker.com/r/jljl1337/xpense)
![GitHub License](https://img.shields.io/github/license/jljl1337/xpense?label=License)

## Features

- Simple: Everything in one image, no separate database instance required
- Disaster-ready: Built-in automated database backups
- Efficient: Lightweight and fast, built with Go and SQLite
- Configurable: Easily customizable via environment variables

## Install

Sample `compose.yml`:

```yml
services:

  xpense:
    image: jljl1337/xpense:latest # Use specific tag for production
    container_name: xpense
    restart: unless-stopped
    volumes:
      - ./live:/xpense/data/live
      - ./backup:/xpense/data/backup
    environment:
      - BACKUP_CRON_SCHEDULE=0 * * * * # Backup every hour
```

## Configuration

| Environment Variable | Type | Default Value | Description |
|---------------------|------|---------------|-------------|
| `DB_PATH` | string | `data/live/db/live.db` | Path to the SQLite database file |
| `DB_BUSY_TIMEOUT` | string | `30000` | Database busy timeout in milliseconds |
| `BACKUP_DB_PATH` | string | `data/backup/db/backup.db` | Path to the backup database file |
| `BACKUP_CRON_SCHEDULE` | string | `0 0 * * *` | Cron schedule for database backups (daily at midnight) |
| `LOG_LEVEL` | int | `0` | Logging level for the application |
| `LOG_HEALTH_CHECK` | bool | `false` | Whether to log health check requests |
| `PORT` | string | `8080` | Port number for the HTTP server |
| `CORS_ORIGINS` | string | `*` | Allowed CORS origins (comma-separated) |
| `PASSWORD_BCRYPT_COST` | int | `12` | Bcrypt cost factor for password hashing |
| `SESSION_COOKIE_NAME` | string | `xpense_session_token` | Name of the session cookie |
| `SESSION_COOKIE_HTTP_ONLY` | bool | `true` | Whether the session cookie is HTTP-only |
| `SESSION_COOKIE_SECURE` | bool | `false` | Whether the session cookie requires HTTPS |
| `SESSION_TOKEN_LENGTH` | int | `32` | Length of generated session tokens |
| `SESSION_TOKEN_CHARSET` | string | alphanumeric characters (case-sensitive) | Character set for session token generation |
| `SESSION_LIFETIME_MIN` | int | `10080` (7 days) | Session lifetime in minutes |
| `PRE_SESSION_LIFETIME_MIN` | int | `15` | Pre-session lifetime in minutes |
| `CSRF_TOKEN_LENGTH` | int | `32` | Length of generated CSRF tokens |
| `CSRF_TOKEN_CHARSET` | string | alphanumeric characters (case-sensitive) | Character set for CSRF token generation |
| `PAGE_SIZE_MAX` | int64 | `100` | Maximum page size for paginated results |
| `PAGE_SIZE_DEFAULT` | int64 | `10` | Default page size for paginated results |
| `SESSION_COOKIE_SAME_SITE_MODE` | string | `lax` | SameSite mode for session cookie (`lax`, `strict`, or `none`), other values are treated as `none` |

## Development

```bash
go run cmd/xpense/main.go
```
