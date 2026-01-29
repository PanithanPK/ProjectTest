# Backend Golang Test

A production-ready REST API built with Go + Fiber v2 + MySQL, featuring JWT authentication,  and Swagger documentation.

## Tech Stack

- Go
- Fiber v2
- MySQL
- JWT (Bearer Token)
- bcrypt
- dotenv
- Swagger (swaggo)

1) Create database schema

- Create a MySQL database (example: `projecttest`)
- Run the migration:

```sql
source database/migration.sql;
```

2) Configure environment

- Copy `.env` and edit values.

3) Install dependencies

```bash
go mod tidy
```

4) Run the server

```bash
go run .
```

## Environment Variables

- `APP_PORT`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`


