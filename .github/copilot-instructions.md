# Project: Lashes - Go Proxy Rotator

A robust proxy rotation library for Go applications that supports multiple proxy types (HTTP, SOCKS4, SOCKS5) with configurable rotation strategies and persistence options.

## Core Features

- Multiple proxy types support (HTTP, SOCKS4, SOCKS5)
- Configurable rotation strategies (round-robin, random, weighted, least-used)
- Database persistence (SQLite, MySQL, PostgreSQL)
- Proxy health checking and validation
- Performance metrics tracking
- Repository pattern implementation

## Development Guidelines

### Development Dependencies

- Zero dependencies for core functionality
- Use standard library solutions where possible
- Database drivers only loaded when explicitly requested
- No external automation tools
- No third-party services

### Implementation Requirements

- Use native Go HTTP client
- Implement all features in pure Go
- Lazy-load optional features
- Use standard library rate limiting
- Keep database operations optional

### Storage Implementation

- In-memory storage as default
- Optional database support
- Lazy database initialization
- Support for multiple backends
- Clean separation of storage logic

## Code Style Guide

### Repository Pattern

- Use interfaces for repository definitions
- All repository methods should accept context.Context as first parameter
- Return concrete errors rather than wrapping in custom types
- Follow standard CRUD operation naming: Create, GetByID, Update, Delete, List

### Project Structure

- `/internal/domain` - Core domain models and interfaces
- `/internal/repository` - Data access layer and implementations
- `/internal/rotation` - Proxy rotation strategies
- `/internal/validation` - Proxy validation and health checks
- `/internal/storage` - Database configuration and migrations

### Code Dependencies

- GORM for database operations
- Multiple database drivers (SQLite, MySQL, PostgreSQL)
- Built with Go 1.24.0 features

### Testing Requirements

- Write table-driven tests
- Use meaningful test names
- Test both success and error cases
- Mock external dependencies

### Security Requirements

- Handle proxy credentials securely
- Implement rate limiting
- Support TLS configuration
- Validate proxy URLs
- Sanitize inputs

### Contributing Guidelines

- Follow semantic versioning
- Document all public APIs
- Write comprehensive tests
- Use pure Go when possible
- Keep dependencies minimal

### License Compliance

- MIT license for core functionality
- Compatible licenses for dependencies
- Clear attribution in docs
