# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release
- Core `Notifier` interface for extensible notification providers
- `Manager` for handling multiple notification providers
- Telegram notification provider
  - Support for simple text messages
  - Markdown/HTML formatting support
  - Photo messages
  - Silent notifications for low priority messages
- Slack notification provider
  - Support for simple text messages
  - Rich messages with blocks
  - Attachments with fields
  - File uploads
  - Custom username and icon
- Async broadcasting to multiple providers
- Comprehensive test suite
- Examples for simple usage, manager, and custom providers
- Full documentation in README

### Features
- Synchronous and asynchronous message broadcasting
- Rich message support with titles, attachments, and fields
- Priority-based notifications (high, normal, low)
- Type-safe error handling with `NotificationError`
- Thread-safe manager implementation
- Context support for cancellation and timeouts

## [0.1.0] - 2025-10-06

### Added
- Initial project setup
- Go module configuration
- Basic project structure
- MIT License

---

## Template for future releases

## [Version] - YYYY-MM-DD

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

