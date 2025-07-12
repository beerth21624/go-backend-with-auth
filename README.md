# Beerdosan Go Backend API Template (In Progress)

Go backend API with authentication, session management, and hybrid architecture (clean architecture and domain driven design).

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
└── Makefile              # Development commands
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

## License

This project is licensed under the MIT License.

## Made by Beerdosan
