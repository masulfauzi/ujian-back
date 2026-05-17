# Project Structure Overview

Complete modular backend project with clean architecture using Go, Fiber, PostgreSQL, and GORM.

## Directory Tree

```
backend/
├── cmd/
│   └── server/
│       └── main.go                          # Application entry point
│
├── configs/
│   ├── app.go                               # Application configuration (port, env, etc)
│   ├── database.go                          # Database connection & DSN setup
│   └── jwt.go                               # JWT configuration
│
├── internal/
│   │
│   ├── modules/                             # Business logic modules
│   │   │
│   │   ├── auth/                            # Authentication module
│   │   │   ├── controller/
│   │   │   │   └── auth_controller.go       # HTTP handlers
│   │   │   ├── service/
│   │   │   │   └── auth_service.go          # Business logic (register, login, token generation)
│   │   │   ├── repository/
│   │   │   │   └── auth_repository.go       # Data access (user queries)
│   │   │   ├── dto/
│   │   │   │   └── auth_dto.go              # Request/Response DTOs
│   │   │   ├── model/
│   │   │   │   └── auth.go                  # Domain models
│   │   │   ├── validator/
│   │   │   │   └── auth_validator.go        # Input validation
│   │   │   └── routes/
│   │   │       └── auth_routes.go           # Route definitions
│   │   │
│   │   └── user/                            # User management module
│   │       ├── controller/
│   │       │   └── user_controller.go       # HTTP handlers
│   │       ├── service/
│   │       │   └── user_service.go          # Business logic (CRUD operations)
│   │       ├── repository/
│   │       │   └── user_repository.go       # Data access (user CRUD)
│   │       ├── dto/
│   │       │   └── user_dto.go              # Request/Response DTOs
│   │       ├── model/
│   │       │   └── user.go                  # User database model
│   │       ├── validator/
│   │       │   └── user_validator.go        # Input validation
│   │       └── routes/
│   │           └── user_routes.go           # Route definitions
│   │
│   ├── middleware/                          # HTTP middleware
│   │   ├── jwt_auth.go                      # JWT authentication middleware
│   │   ├── logger.go                        # Request/response logging
│   │   ├── cors.go                          # CORS configuration
│   │   └── recovery.go                      # Panic recovery
│   │
│   ├── database/                            # Database management
│   │   ├── connection.go                    # Database connection singleton
│   │   └── migrate.go                       # Auto-migration setup
│   │
│   ├── helpers/                             # Helper functions
│   │   └── response.go                      # HTTP response formatting
│   │
│   ├── utils/                               # Utility functions
│   │   └── password.go                      # Password hashing & verification
│   │
│   └── constants/                           # Application constants
│       └── constants.go                     # Roles, statuses, error messages
│
├── migrations/                              # Database migration files (future)
├── pkg/                                     # Public packages (future)
├── docs/                                    # API documentation (future)
│
├── go.mod                                   # Go module definition
├── go.sum                                   # Go module checksums
├── .env                                     # Environment variables (local)
├── .env.example                             # Environment variables template
├── .gitignore                               # Git ignore rules
├── Dockerfile                               # Docker container image
├── docker-compose.yml                       # Docker services composition
├── Makefile                                 # Build automation
├── README.md                                # Full documentation
├── QUICKSTART.md                            # Quick start guide
└── PROJECT_STRUCTURE.md                     # This file
```

## File Count Summary

- **Total Go Files**: 27
- **Modules**: 2 (auth, user)
- **Middleware**: 4
- **Configuration Files**: 3
- **Support Files**: 6 (README, Makefile, Docker, etc)

## Layered Architecture

Each module follows the same pattern:

```
HTTP Request
    ↓
Controller (HTTP handler)
    ↓
Service (Business logic)
    ↓
Repository (Data access)
    ↓
Database
```

### Layer Responsibilities

1. **Controller**
   - Parse HTTP requests
   - Validate input format
   - Call services
   - Format responses

2. **Service**
   - Business logic implementation
   - Data validation
   - Service-to-service coordination
   - Error handling

3. **Repository**
   - Database queries
   - CRUD operations
   - Query optimization

4. **Model**
   - Database schema definition
   - GORM decorators
   - Data structures

5. **DTO**
   - Request/Response objects
   - Data transformation
   - API contracts

## Configuration Files

### `configs/app.go`
- Application metadata (name, port, env)
- Environment variable loading
- Configuration helpers

### `configs/database.go`
- PostgreSQL connection setup
- GORM initialization
- Connection pooling

### `configs/jwt.go`
- JWT secret management
- Token expiration settings
- Token validation configuration

## Middleware Stack

Applied in order:

1. **CORS** - Cross-origin resource sharing
2. **Logger** - Request/response logging
3. **Recovery** - Panic recovery and error handling
4. **JWT Auth** - (Applied selectively to protected routes)

## Database Layer

- **ORM**: GORM with PostgreSQL driver
- **Auto-migration**: Automatic table creation
- **Models**: Typed database entities with hooks
- **Repository Pattern**: Abstracted data access

## Authentication

- **JWT Token-based**
- **Bcrypt password hashing**
- **Protected endpoints** via JWT middleware
- **Token claims**: user_id, email, role

## Adding New Features

To add a new module (e.g., `products`):

1. Create folder: `internal/modules/products/`
2. Create subfolders for: controller, service, repository, model, dto, routes, validator
3. Implement interfaces following auth/user patterns
4. Register routes in `cmd/server/main.go`
5. Add database model to `internal/database/migrate.go`

## Key Design Patterns

1. **Repository Pattern** - Abstract data access
2. **Dependency Injection** - Services depend on repositories
3. **Interface-based design** - Easy testing and swapping
4. **Separation of concerns** - Each layer has single responsibility
5. **DTOs** - Decouple API contracts from internal models
6. **Middleware** - Cross-cutting concerns

## Error Handling

- Centralized error responses
- Consistent error format
- Meaningful error messages
- HTTP status codes

## Security

- Password hashing with bcrypt
- JWT token validation
- CORS enabled
- Input validation
- SQL injection prevention (GORM parameterized queries)

## Scalability Features

- Modular architecture (easy to split into microservices)
- Interface-based design (easy to swap implementations)
- Connection pooling (database efficiency)
- Middleware extensibility
- Clean code structure (maintainability)

## Future Enhancements

- [ ] Unit tests
- [ ] Integration tests
- [ ] API documentation (Swagger)
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Message queue integration
- [ ] Logging service (Zerolog)
- [ ] Monitoring/metrics
- [ ] CI/CD pipeline
- [ ] Load testing
- [ ] API versioning
- [ ] Soft delete support
- [ ] Audit logging
- [ ] Permission-based access control

## Tech Stack Reference

| Component | Package | Version |
|-----------|---------|---------|
| Framework | Fiber v2 | 2.50.0 |
| ORM | GORM | 1.25.4 |
| Database Driver | postgres | 1.5.4 |
| Authentication | JWT v5 | 5.0.0 |
| Validation | validator/v10 | 10.16.0 |
| Password Hashing | bcrypt | crypto/bcrypt |
| UUID Generation | google/uuid | 1.5.0 |
| Environment | godotenv | 1.5.1 |

---

**Created**: 2024
**Go Version**: 1.21+
**License**: MIT
