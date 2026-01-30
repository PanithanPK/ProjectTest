# Backend Golang Test

A production-ready REST API built with Go + Fiber v2 + MySQL, featuring JWT authentication, and Swagger documentation.

## Tech Stack

- Go
- Fiber v2
- MySQL
- JWT (Bearer Token)
- bcrypt
- dotenv
- Swagger (swaggo)

## Setup

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
- `JWT_SECRET`
- `JWT_EXPIRE` (example: `24h`)

## API Documentation (Swagger)

API documentation is available at:

```
http://localhost:5000/swagger/index.html
```

### Using Swagger UI

1. **Open Swagger UI** - Go to `http://localhost:5000/swagger/index.html`

2. **For endpoints that require authentication:**
   - Click the **Authorize** button (lock icon) at the top right
   - In the dialog, enter your JWT token in format: `Bearer <your_jwt_token>`
   - Example: `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - Click **Authorize** button

3. **Test endpoints:**
   - Try the request directly from Swagger UI
   - The Bearer token will be automatically added to the `Authorization` header

### Authentication Flow

1. **Register** - POST `/user/register` (no auth required)
2. **Login** - POST `/user/login` (no auth required)
   - Returns `access_token` in response
3. **Use Token** - For protected endpoints, use the token to authorize in Swagger UI
4. **Access Protected Routes** - GET `/user/me`, POST `/accounting/transfer`, etc.

### Example Bearer Token

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2NzY0NzgwMDB9.AbCdEfGhIjKlMnOpQrStUvWxYz
```

## API Endpoints

### User Management
- `POST /user/register` - Create new account
- `POST /user/login` - Login and get JWT token
- `GET /user/me` - Get current user info (requires auth)
- `PATCH /user/update` - Update user profile (requires auth)

### Money Transfer
- `POST /accounting/transfer` - Transfer money (requires auth)
- `GET /accounting/transfer-list` - View transfer history (requires auth)


