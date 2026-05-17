# Installation Guide

Complete step-by-step guide to set up and run the Fiber Backend API.

## System Requirements

- **Go**: 1.21 or higher
- **PostgreSQL**: 12 or higher (or Docker)
- **Git**: For version control
- **Make**: (Optional) For using Makefile shortcuts

## Quick Installation (5 minutes)

### 1. Install PostgreSQL

**Option A: Using Docker (Recommended)**

```bash
docker compose up -d postgres
```

Wait for PostgreSQL to be ready:

```bash
docker compose exec postgres pg_isready -U postgres
```

**Option B: Using Homebrew (macOS)**

```bash
brew install postgresql
brew services start postgresql
createdb fiber_backend
```

**Option C: Using apt (Ubuntu/Debian)**

```bash
sudo apt-get install postgresql postgresql-contrib
sudo -u postgres createdb fiber_backend
```

### 2. Clone/Download Project

If not already done:

```bash
cd /Applications/XAMPP/xamppfiles/htdocs/UJIAN-NEW/ujian-back
```

### 3. Install Go Dependencies

```bash
go mod download
```

If go.sum is missing or corrupted:

```bash
rm -f go.sum
go mod tidy
```

### 4. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` if using non-default PostgreSQL settings:

```bash
nano .env
# or
code .env
```

### 5. Run the Application

```bash
go run cmd/server/main.go
```

Or using Makefile:

```bash
make run
```

Or with the built binary:

```bash
go build -o bin/backend cmd/server/main.go
./bin/backend
```

## Detailed Installation Steps

### Step 1: Verify Go Installation

```bash
go version
# Expected output: go version go1.21 (or higher) darwin/amd64
```

If not installed, download from: https://golang.org/dl/

### Step 2: Verify PostgreSQL Installation

**Using Docker:**

```bash
docker ps
# Look for postgres container
```

**Using Local Installation:**

```bash
psql --version
# Expected output: psql (PostgreSQL) 12.0 (or higher)

# Test connection
psql -U postgres -c "SELECT version();"
```

### Step 3: Clone Repository

If using Git:

```bash
git clone <repository-url>
cd backend
```

### Step 4: Create Database

**Using Docker:**

```bash
docker compose up -d postgres
docker compose exec postgres psql -U postgres -c "CREATE DATABASE fiber_backend;"
```

**Using Local PostgreSQL:**

```bash
createdb -U postgres fiber_backend
```

Or using SQL:

```bash
psql -U postgres
# Then in psql prompt:
CREATE DATABASE fiber_backend;
\q
```

### Step 5: Setup Environment File

```bash
# Copy example to .env
cp .env.example .env

# Edit if needed (using your favorite editor)
nano .env
```

**Default .env Content:**

```env
APP_NAME=Fiber Backend API
APP_PORT=3000
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_backend

JWT_SECRET=supersecretkey
JWT_EXPIRED=72h

LOG_LEVEL=info
```

### Step 6: Download Dependencies

```bash
go mod download
```

This downloads all required packages without modifying go.mod or go.sum.

### Step 7: Build Application

**Option A: Build Binary**

```bash
mkdir -p bin
go build -o bin/backend cmd/server/main.go
```

**Option B: Run Directly**

```bash
go run cmd/server/main.go
```

### Step 8: Verify It Works

In another terminal:

```bash
curl http://localhost:3000/health
```

Expected response:

```json
{
  "status": "ok",
  "service": "Fiber Backend API"
}
```

## Using Docker Compose

### Start Services

```bash
docker compose up -d
```

This starts PostgreSQL with persistent volume.

### View Logs

```bash
docker compose logs -f postgres
```

### Stop Services

```bash
docker compose down
```

### Clean Up (Delete Data)

```bash
docker compose down -v
```

## Using Makefile

### Available Commands

```bash
make help              # Show all available commands
make install-deps      # Download dependencies
make build            # Build binary
make run              # Run the application
make test             # Run tests
make clean            # Clean build artifacts
make fmt              # Format code
make vet              # Run go vet
make lint             # Run linter (if installed)
make docker-build     # Build Docker image
make docker-run       # Run Docker container
```

### Example Workflow

```bash
# Install dependencies
make install-deps

# Format code
make fmt

# Build
make build

# Run
make run

# In another terminal, test
curl http://localhost:3000/health
```

## Troubleshooting

### Error: "port 5432 already in use"

```bash
# Find process using port 5432
lsof -i :5432

# Kill the process (if needed)
kill -9 <PID>

# Or change port in .env
# DB_PORT=5433
```

### Error: "database fiber_backend does not exist"

```bash
# Create the database
psql -U postgres -c "CREATE DATABASE fiber_backend;"

# Or with Docker
docker compose exec postgres psql -U postgres -c "CREATE DATABASE fiber_backend;"
```

### Error: "connection refused"

1. Verify PostgreSQL is running:
   ```bash
   docker compose ps  # for Docker
   pg_isready -h localhost  # for local
   ```

2. Verify database credentials in `.env`

3. Check firewall settings

### Error: "go: missing go.sum"

```bash
rm -f go.sum go.mod
go mod init backend
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10
go get github.com/google/uuid
go mod tidy
```

### Error: "port 3000 already in use"

Change port in `.env`:

```env
APP_PORT=3001
```

Then run again.

### Error: "invalid go version"

This error occurs if go.mod has `go 1.25.0` but you have `go 1.21`:

```bash
# Edit go.mod and change first line to:
go 1.21

# Then run:
go mod tidy
```

## Development Workflow

### 1. Start PostgreSQL

```bash
docker compose up -d postgres
# or local PostgreSQL if installed
```

### 2. Run in Development Mode

```bash
go run cmd/server/main.go
```

The server will restart if you modify source files (with hot reload tools like `air`):

```bash
go install github.com/cosmtrek/air@latest
air
```

### 3. Test API Endpoints

Use curl, Postman, or VS Code REST Client:

```bash
# Register user
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@test.com","password":"pass123"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@test.com","password":"pass123"}'

# Get user (with token)
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Production Setup

### Environment Variables

Create `.env.production`:

```env
APP_NAME=Fiber Backend API
APP_PORT=3000
APP_ENV=production

DB_HOST=prod-db-host
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=strong_password_here
DB_NAME=fiber_backend_prod

JWT_SECRET=very_long_random_secret_string_here
JWT_EXPIRED=24h

LOG_LEVEL=warn
```

### Build for Production

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/backend cmd/server/main.go
```

### Docker Build

```bash
docker build -t fiber-backend:latest .
docker run -p 3000:3000 --env-file .env.production fiber-backend:latest
```

### With Kubernetes (Optional)

Create `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fiber-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fiber-backend
  template:
    metadata:
      labels:
        app: fiber-backend
    spec:
      containers:
      - name: fiber-backend
        image: fiber-backend:latest
        ports:
        - containerPort: 3000
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: db-config
              key: host
```

Deploy:

```bash
kubectl apply -f deployment.yaml
```

## Verification Checklist

After installation, verify:

- [ ] Go version is 1.21+: `go version`
- [ ] PostgreSQL is running: `docker compose ps` or `pg_isready`
- [ ] Database created: `psql -l | grep fiber_backend`
- [ ] Dependencies installed: `go mod verify`
- [ ] .env file exists and configured
- [ ] Application starts: `go run cmd/server/main.go`
- [ ] Health endpoint works: `curl http://localhost:3000/health`
- [ ] Can register user: `curl -X POST http://localhost:3000/api/auth/register ...`
- [ ] Can login: `curl -X POST http://localhost:3000/api/auth/login ...`
- [ ] JWT token is received

## Next Steps

1. Read [QUICKSTART.md](QUICKSTART.md) for quick testing
2. Read [README.md](README.md) for API documentation
3. Check [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for code organization
4. Start adding features following the module pattern
5. Write tests for business logic
6. Setup CI/CD pipeline

## Getting Help

If you encounter issues:

1. Check error messages carefully
2. Review relevant section above
3. Verify all prerequisites are installed
4. Check file permissions: `ls -la`
5. Review logs: Check application output
6. Check PostgreSQL logs: `docker compose logs postgres`

## Additional Resources

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Documentation](https://jwt.io/introduction)
- [Go Best Practices](https://golang.org/doc/effective_go)

---

**Installation Complete!** 🎉

You now have a production-ready backend API ready for development.
