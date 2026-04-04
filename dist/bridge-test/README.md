# demo

Generated full-stack web application.

## Stack

| Layer    | Technology                      |
|----------|---------------------------------|
| Database | PostgreSQL              |
| Server   | Go + Gin + GORM                 |
| Client   | React + TypeScript + Vite       |

## Getting started

```bash
# Start the database, apply migrations, and run the server
./dev.sh

# Stop all containers
./shutdown.sh
```

The server listens on port **8080**.
The React dev server proxies API calls to the backend automatically.

```bash
# Start the frontend dev server (in a separate terminal)
cd client && npm install && npm run dev
```

## Environment variables

Configured in `.env`:

| Variable      | Description                        |
|---------------|------------------------------------|
| `DB_HOST`     | Database host                      |
| `DB_PORT`     | Database port                      |
| `DB_USER`     | Database user                      |
| `DB_PASSWORD` | Database password                  |
| `DB_NAME`     | Database name                      |
| `DB_SSLMODE`  | SSL mode (default: `disable`)      |

## API

### Models

#### `posts`

| Method   | Path                    | Description                                 |
|----------|-------------------------|---------------------------------------------|
| `GET`    | `/api/posts`            | List (pagination, search, filter, sort)     |
| `GET`    | `/api/posts/:id`        | Get by id                                   |
| `POST`   | `/api/posts`            | Create                                      |
| `PUT`    | `/api/posts/:id`        | Update                                      |
| `DELETE` | `/api/posts/:id`        | Delete                                      |
| `DELETE` | `/api/posts/batch`      | Batch delete — body: `{"ids": [1, 2, 3]}`   |

### Query parameters (list endpoints)

| Parameter     | Description                                                        |
|---------------|--------------------------------------------------------------------|
| `q`           | Full-text search across all text fields (case-insensitive)         |
| `<field>`     | Filter by exact value (numeric, boolean, enum, or foreign key)     |
| `sort_by`     | Field to sort by. Default: `id`                                    |
| `sort_dir`    | `asc` or `desc`. Default: `desc`                                   |
| `page`        | Page number (1-based). Default: `1`                                |
| `limit`       | Results per page. Default: `20`, max: `100`                        |
