# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Pagination support for TODO list (`ListPaged` method)
- Database connection pool configuration
- Health check endpoint with database status
- Graceful shutdown with signal handling
- Structured logging with `log/slog`
- Unified error codes (`internal/todo/errors.go`)
- Configuration management with viper (`internal/todo/config.go`)
- Request validation with go-playground/validator
- Frontend utilities (Toast, Loading, FormValidator, i18n)
- Production deployment guide
- CONTRIBUTING.md with development guidelines
- ARIA accessibility labels for frontend

### Changed
- API routes now use `/v1/` prefix for versioning
- JWT secret must be set via `JWT_SECRET` environment variable for non-memory storage

### Fixed
- Random ID collision in memory store
- Removed deprecated `rand.Seed()` call (Go 1.20+)
- XSS vulnerability in mock-api.js (innerHTML to textContent)
- Memory leak in refresh token and login failure stores
- RateLimiter goroutine leak (added Stop method)

### Security
- Password strength validation (min 8 chars, letters + numbers)
- Email format validation
- Login failure rate limiting with account lockout

## [1.0.0] - 2025-12-05

### Added
- Initial release
- TODO CRUD API with JWT authentication
- Memory and MySQL/SQLite storage backends
- User registration and login
- Role-based access control (RBAC)
- Rate limiting middleware
- Frontend portal with dark mode support
- Admin dashboard for user/role management
