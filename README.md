# Dunhayat Coffee Roastery API

A clean and performant Go backend for Dunhayat (دان حیات) coffee roastery's
e-commerce system, built with Clean Architecture, DDD, and Hexagonal
Architecture principles using vertical slicing.

## Architecture

This project follows **vertical slicing** architecture where each domain
(products, orders, users) is self-contained with its own:
- Domain entities and business logic
- Use cases (application layer)
- HTTP handlers (delivery/HTTP transport layer)
- Repository implementations (data access layer)

**Dependency direction**: `domain` ← `usecase` ← adapters (`http`,
`repository`). Inner layers do not depend on outer layers.

**Interfaces placement**: Repository and use case interfaces live in
`internal/<slice>/domain`. Concrete implementations live in
`internal/<slice>/repository` and are injected inward.

**Cross-cutting infrastructure**: Shared concerns live in `pkg/`
(`config`, `db`, `redis`, `logger`, `router`) and are treated as
infrastructure, not domain.

**Migrations**: Database schema changes are managed with Atlas (no GORM
auto-migrate).

## Features

- **OTP-based Authentication**: Phone number verification using Kavenegar SMS
  service with Redis-based OTP storage
- **Product Management**: Coffee products with detailed attributes (roast level,
   bitterness, body, etc.)
- **Cart System**: 10-minute product reservation system
- **Order Management**: Complete purchase flow with status tracking
- **Payment Integration**: Secure payment processing implemented with Zibal API
- **Clean Architecture**: Separation of concerns with dependency injection
- **Database Migrations**: Atlas-based schema management with versioning
- **Hybrid Storage**: PostgreSQL for persistent data, Redis for OTPs and caching
- **Future-Ready**: Designed for easy extension (user profiles, additional
  payment methods, etc.)

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: net/http (standard library) with Go 1.22+ method-based routing
- **Configuration**: Viper
- **Database**: PostgreSQL with GORM (user data, orders, products)
- **Migrations**: Atlas (modern schema management)
- **SMS Service**: Kavenegar
- **Payment Service**: Zibal
- **Cache/OTP Storage**: Redis (OTP codes, rate limiting, caching)
- **Architecture**: Clean Architecture + DDD + Hexagonal

## Project Structure

```hier
┌── api/              # OpenAPI documents
├── cmd/              # Entrypoints
│   └── api/          # Main application
├── internal/         # Domain-specific code (vertical slices)
│   ├── auth/         # Authentication domain
│   ├── orders/       # Order domain
│   ├── products/     # Product domain
│   └── users/        # User domain
├── migrations/       # Atlas migration files
├── pkg/              # Shared utilities
│   ├── config/       # Configuration management (Viper)
│   ├── db/           # Database utilities (PostgreSQL)
│   ├── logger/       # Logging utilities
|   ├── middleware/   # HTTP middleware
│   ├── redis/        # Redis connection utilities
│   └── router/       # HTTP routing (Go 1.22+ method-based)
├── .air.toml         # Hot reloading configuration
├── .gitignore        # Git ignore rules
├── atlas.hcl         # Atlas configuration
├── config.yaml       # Application settings
├── go.mod, go.sum    # Go module dependencies
├── Makefile          # Development commands
└── README.md         # Comprehensive documentation
```

## Quick Start

### Prerequisites

- **Go 1.25+**
- **PostgreSQL 17+**
- **Redis 7+**
- **Atlas CLI** (installed with `go install`, `brew`, or package manager of
  your choice)

### Configuration Setup

The application uses Viper for configuration management, which supports:
- YAML configuration files
- Environment variables
- Default values
- Hot reloading (can be enabled)

To set up, copy the configuration template (`config.yaml.example`) to
`config.yaml`, and modify it as required.

### Database Setup

1. Create the database:
```sql
CREATE DATABASE dunhayat;
```

2. Apply database migrations using Atlas:
```sh
# Apply migrations
make migrate

# Check migration status
make migrate-status
```

### Migration Workflow

1. **Create Migration**: `make migrate-new name=description`
2. **Edit Migration**: Modify the generated SQL file in `migrations/`
3. **Apply Migration**: `make migrate`
4. **Verify Status**: `make migrate-status`

### Redis Setup

Install Redis with the package manager of your choice, and start the service.

### Build and Fly

Consult the `Makefile` and proceed to get airborne.

## Database Schema

Check out the
[initial database schema](migrations/20250828055134_initial_schema.sql).

## API Endpoints

Check out interactive Swagger API documents enabled in development mode.

## **OTP Flow**
```
1. User requests OTP
   ↓
2. Backend generates "123456"
   ↓
3. Backend saves to Redis: key="otp:+989123456789", value=OTP_JSON, TTL=600s
   ↓
4. Backend tells Kavenegar: "Send SMS with code 123456"
   ↓
5. Kavenegar sends SMS to user
   ↓
6. User types "123456" in frontend
   ↓
7. Frontend sends to backend: phone="+989123456789", code="123456"
   ↓
8. Backend fetches from Redis: "What's the OTP for +989123456789?"
   ↓
9. Backend compares: "123456" == "123456" ✓
   ↓
10. Backend creates user session, returns token
   ↓
11. Redis automatically removes OTP after 10 minutes
   ↓
12. User could proceed with their purchase etc.
```

## Future Enhancements

- User profile management
- Purchase history
- Email verification
- Admin dashboard
- Inventory management
- Advanced Redis caching for products and user data

_...and the sky's the limit..._

## 📝 License

This project is proprietary software for Dunhayat Coffee Roastery.
