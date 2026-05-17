# Quick Start Guide

Get the Fiber Backend API running in minutes.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation Steps

### 1. Setup Database

Using PostgreSQL CLI:

```bash
createdb fiber_backend
```

Or using Docker:

```bash
docker compose up -d postgres
```

Wait for PostgreSQL to be ready:

```bash
docker compose exec postgres pg_isready -U postgres
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` if needed (default values should work with Docker Postgres).

### 3. Install Dependencies

```bash
go mod download
go mod tidy
```

### 4. Run the Application

```bash
# Development mode with hot reload
make dev

# Or using go directly
go run cmd/server/main.go

# Or build and run
make build
./bin/backend
```

The server will start on `http://localhost:3000`

## Test the API

### Health Check

```bash
curl http://localhost:3000/health
```

### Register

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

Save the returned `token` for the next requests.

### Login

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Current User

```bash
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer <token>"
```

Replace `<token>` with the token from register/login response.

### Get All Users

```bash
curl http://localhost:3000/api/users
```

### Create User

```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "password123"
  }'
```

### Get User by ID

```bash
curl http://localhost:3000/api/users/<user-id>
```

### Update User

```bash
curl -X PUT http://localhost:3000/api/users/<user-id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Jane Updated"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:3000/api/users/<user-id> \
  -H "Authorization: Bearer <token>"
```

## Docker Setup (Optional)

Run PostgreSQL with Docker Compose:

```bash
docker compose up -d postgres
```

Build and run backend with Docker:

```bash
docker build -t fiber-backend .
docker run -p 3000:3000 --env-file .env --network fiber-network fiber-backend
```

## Troubleshooting

### Port 3000 Already in Use

Change `APP_PORT` in `.env` to another port, e.g., `3001`.

### Database Connection Error

Check your `.env` database settings match your PostgreSQL configuration.

### JWT Token Errors

Make sure the JWT token format is correct: `Authorization: Bearer <token>`

## Next Steps

1. Read the full [README.md](README.md) for detailed API documentation
2. Check the project structure in [README.md](README.md#project-structure)
3. Add more modules following the existing pattern
4. Implement custom business logic in service layer
5. Add unit tests for your services
6. Deploy to production

## Support

See [README.md](README.md) for more information and troubleshooting.
