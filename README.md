# VentureX Backend

A robust Go backend API with authentication, session management, and clean architecture.

## Features

- 🔐 **Authentication System** - JWT-based authentication with refresh tokens
- 👥 **User Management** - User registration, login, and profile management
- 🛡️ **Security** - Rate limiting, password hashing, session management
- 📊 **Database** - PostgreSQL with migrations using Goose
- 🏗️ **Clean Architecture** - Domain-driven design with proper separation of concerns
- 🔄 **Session Management** - Device fingerprinting and session tracking
- 📝 **Audit Trail** - Login attempt tracking and monitoring

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT tokens
- **Migrations**: Goose
- **Password Hashing**: bcrypt/Argon2

## Project Structure

```
backend/
├── cmd/api/                 # Application entry point
├── config/                  # Configuration files
├── internal/
│   ├── app/
│   │   ├── api/            # HTTP handlers and middleware
│   │   ├── domain/         # Business entities and rules
│   │   ├── repositories/   # Data access layer
│   │   ├── service/        # Business logic services
│   │   └── usecase/        # Application use cases
│   └── pkg/                # Shared packages
├── migrations/             # Database migrations
├── scripts/               # Setup scripts
└── Makefile              # Development commands
```

## Quick Start

### Prerequisites

- Go 1.24 or higher
- PostgreSQL database
- Make (for using Makefile commands)

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Configure database**

   - Copy `config/app.example.yaml` to `config/app.yaml`
   - Update database configuration:

   ```yaml
   database:
     host: your-db-host
     port: 5432
     user: your-db-user
     password: your-db-password
     dbname: your-db-name
     sslmode: disable
     timezone: Asia/Bangkok
   ```

4. **Run database migrations**

   ```bash
   make migrate-up
   ```

5. **Start the application**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

## Database Migration Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback one migration
make migrate-down

# Check migration status
make migrate-status

# Create new migration
make migrate-create name=migration_name

# Reset database (WARNING: drops all data)
make migrate-reset

# Show current migration version
make migrate-version
```

## API Endpoints

### Authentication

| Method | Endpoint                           | Description                |
| ------ | ---------------------------------- | -------------------------- |
| POST   | `/api/v1/auth/login`               | User login                 |
| POST   | `/api/v1/auth/logout`              | User logout                |
| POST   | `/api/v1/auth/refresh`             | Refresh access token       |
| GET    | `/api/v1/auth/me`                  | Get user profile           |
| GET    | `/api/v1/auth/sessions`            | Get user sessions          |
| DELETE | `/api/v1/auth/sessions/:sessionId` | Terminate specific session |
| DELETE | `/api/v1/auth/sessions`            | Terminate all sessions     |
| PUT    | `/api/v1/auth/password`            | Change password            |

### Health Check

| Method | Endpoint  | Description  |
| ------ | --------- | ------------ |
| GET    | `/health` | Health check |

## API Usage Examples

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Response:

```json
{
  "status": {
    "code": 200,
    "header": "Success",
    "description": "Request completed successfully"
  },
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
    "expires_at": "2025-07-11T01:08:00Z",
    "user": {
      "id": "1",
      "username": "admin",
      "email": "admin@venturex.com",
      "role": "user",
      "status": "active"
    }
  }
}
```

### Get Profile

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Development Commands

```bash
# Run the application
make run

# Build the application
make build

# Run tests
make test

# Clean build artifacts
make clean

# Show all available commands
make help
```

## Default Users

After running migrations, a default admin user is created:

- **Username**: `admin`
- **Password**: `admin123`
- **Email**: `admin@venturex.com`

⚠️ **Important**: Change the default password in production!

## Environment Variables

You can override configuration using environment variables:

```bash
# Database configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_DBNAME=venturex
DATABASE_SSLMODE=disable
DATABASE_TIMEZONE=Asia/Bangkok

# Server configuration
SERVER_PORT=8080

# JWT configuration
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h
JWT_ISSUER=venturex-backend
JWT_AUDIENCE=venturex-app
```

## Docker Support

```bash
# Start services
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs
```

## Architecture

This project follows Clean Architecture principles:

1. **Domain Layer** - Business entities and rules
2. **Use Case Layer** - Application business logic
3. **Service Layer** - Domain services and shared logic
4. **Repository Layer** - Data access abstraction
5. **Handler Layer** - HTTP request/response handling

## Security Features

- **Password Hashing**: bcrypt with configurable cost
- **JWT Tokens**: RSA-signed access and refresh tokens
- **Session Management**: Device fingerprinting and tracking
- **Rate Limiting**: Login attempt monitoring
- **Audit Trail**: Comprehensive login attempt logging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.
