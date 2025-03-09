# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Logo image in README.md
- Test coverage badge

## [0.1.8] - 2025-03-09

### Fixed

- Circuit breaker state transitions now properly handle half-open state
- Fixed state transitions from half-open to closed after success
- Fixed state transitions from half-open to open after failure
- Improved circuit breaker testing to better validate state transitions

### Security

- Added stronger security warnings for TLS certificate validation bypass
- Enhanced security documentation for InsecureSkipVerify usage
- Removed hardcoded credentials from example code
- Improved input validation throughout the codebase

### Changed

- Restructured README.md with clearer documentation and examples
- Improved code organization in client.go and breaker.go
- Reduced unused code to improve maintainability

### Removed

- Eliminated several unused functions flagged by linters
- Removed code duplication in the repository implementations

## [0.1.7]

### Added

- Circuit breaker pattern implementation for proxy failure handling
- New weighted rotation strategy using cryptographically secure randomization
- Support for in-memory metrics caching to improve performance

### Changed

- Improved proxy validation with better error handling
- Enhanced metrics collection with more detailed statistics
- Updated documentation with examples for new features

### Fixed

- Race condition in proxy rotation when using concurrent requests
- Memory leak in HTTP client creation
- Timeout handling in validation requests

## [0.1.6]

### Added

- Health check system for monitoring proxy status
- Support for PostgreSQL backend storage
- Rate limiting capabilities

### Changed

- Improved error handling with structured error types
- Better performance in rotation strategies
- Updated documentation with usage examples

## [0.1.5]

### Added

- Initial public release
- Support for HTTP, SOCKS4, and SOCKS5 proxies
- Multiple rotation strategies (round-robin, random)
- In-memory and SQLite storage options
- Basic proxy validation

## [0.1.1]

### Security

- Replaced `math/rand` with `crypto/rand` in the weighted strategy for secure random number generation
- Added TLS validation options for secure connections

### Fixed

- Fixed deterministic rotation in round-robin strategy by sorting proxies by URL
- Improved weighted strategy to properly handle zero-weight proxies (95/5 split)
- Fixed mock URL parsing for consistent behavior in tests
- Fixed test flakiness in rotation strategy tests
- Added error handling for HTTP response body closing

### Changed

- Enhanced `Next()` method in weighted strategy to use cryptographically secure randomization
- Restructured metrics collection for better performance tracking
- Added exponential backoff for retry logic
- Improved error handling and propagation

### Added

- GitHub Actions workflow for automated releases
- Additional test coverage for edge cases
- Added more examples demonstrating rotation strategies
- Comprehensive test suite for all rotation strategies

## [0.1.0]

### Added

- Initial release with core functionality
- Support for HTTP, SOCKS4, and SOCKS5 proxies
- Four rotation strategies: round-robin, random, weighted, least-used
- In-memory proxy storage
- Optional database backends (SQLite, MySQL, PostgreSQL)
- Basic proxy health checking and validation
- Simple metrics collection
- Context-aware API
