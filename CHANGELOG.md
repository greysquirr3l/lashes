# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Logo image in README.md
- Test coverage badge

## [0.1.1] - 2023-11-10

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

## [0.1.0] - 2023-11-01

### Added

- Initial release with core functionality
- Support for HTTP, SOCKS4, and SOCKS5 proxies
- Four rotation strategies: round-robin, random, weighted, least-used
- In-memory proxy storage
- Optional database backends (SQLite, MySQL, PostgreSQL)
- Basic proxy health checking and validation
- Simple metrics collection
- Context-aware API
