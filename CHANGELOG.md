# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of HomeGuard WOL Service
- Device management with YAML configuration
- HTTP listener for wake-on-LAN requests
- MQTT listener for wake-on-LAN requests
- Support for waking devices by name or MAC address
- Graceful shutdown with signal handling
- Structured logging with slog
- Command-line configuration options
- Docker support with Dockerfile and docker-compose
- Comprehensive documentation (README, CONTRIBUTING)
- CI/CD with GitHub Actions
- GoReleaser configuration for multi-platform releases

### Features
- 🚀 Multi-protocol support (HTTP and MQTT)
- 📝 Device management via YAML configuration
- 🔄 Flexible wake-up methods (device name or direct MAC)
- 🛡️ Graceful shutdown
- 📊 Structured logging
- ⚡ Lightweight single binary

## [0.1.0] - 2025-01-XX

### Added
- Initial project structure
- Core WOL functionality
- Basic HTTP listener
- MQTT listener implementation
- Device configuration management
- Example configuration files
- Basic documentation

---

## Version History

### Version Format
- **Major**: Incompatible API changes
- **Minor**: Backward-compatible functionality additions
- **Patch**: Backward-compatible bug fixes

### Release Notes Template

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security updates
```

---

[Unreleased]: https://github.com/p3ddd/HomeGuard/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/p3ddd/HomeGuard/releases/tag/v0.1.0

