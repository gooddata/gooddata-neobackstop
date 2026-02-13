# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.21.1] - 2026-02-13

### Removed
- Unused `getBrowserName` function

## [0.21.0] - 2026-02-13

### Changed
- Refactored browser configuration: `browsers` is now a map where keys are user-defined aliases and values specify the browser `name` and `args`. Scenarios reference browsers by alias, allowing multiple configurations of the same browser engine
- Added `defaultBrowsers` field to specify which browser aliases to use when a scenario doesn't define its own

## [0.20.2] - 2026-02-12

### Fixed
- Avoid log.Panicf for ReloadAfterReady operation

## [0.20.1] - 2026-02-09

### Fixed
- Screenshot saving in wrong directory when reference is missing

## [0.20.0] - 2026-02-02

### Changed
- Support for retries
- Upgraded Go to 1.25.6

## [0.19.1] - 2026-02-02

### Changed
- Adjusted docker build stages

## [0.19.0] - 2025-12-16

### Changed
- Improved logging to allow for performance optimalisations

## [0.18.0] - 2025-12-16

### Added
- Support for RequireSameDimensions

## [0.17.0] - 2025-12-15

### Added
- Support for multiple ReadySelectors

## [0.16.0] - 2025-12-15

### Added
- State support for ReadySelector

## [0.15.0] - 2025-12-15

### Removed
- Extra delay

## [0.14.0] - 2025-12-15

### Removed
- Support for legacy postInteractionWait type

## [0.13.0] - 2025-12-15

### Removed
- Support for legacy hoverSelectors type

## [0.12.0] - 2025-12-15

### Removed
- Support for legacy clickSelectors type

## [0.11.1] - 2025-12-15

### Fixed
- Delay parsing issue

## [0.11.0] - 2025-12-15

### Changed
- Simplified input types

## [0.10.0] - 2025-12-15

### Removed
- Support for legacy delay type

## [0.9.0] - 2025-12-12

### Changed
- Improved logging

## [0.8.1] - 2025-12-12

### Fixed
- Removed unnecessary logging

## [0.8.0] - 2025-12-12

### Changed
- Image comparison is now performed in memory, without writing the test image to disk

## [0.7.1] - 2025-12-11

### Fixed
- Bug affecting multi-browser runs

## [0.7.0] - 2025-12-01

### Changed
- Install browser dependencies dynamically in Dockerfile

## [0.6.0] - 2025-12-01

### Changed
- Config structure to allow args per browser

## [0.5.2] - 2025-12-01

### Added
- Missing package (libx11-xcb1) to Dockerfile which is required for Firefox

## [0.5.1] - 2025-11-28

### Removed
- Headless boolean param for debug

## [0.5.0] - 2025-11-28

### Added
- Headless boolean param for debug

## [0.4.0] - 2025-11-28

### Changed
- Upgraded Go to 1.25.4

## [0.3.0] - 2025-11-28

### Added
- Support for browser-per-scenario configuration

## [0.2.0] - 2025-10-30

### Added
- Playwright browser installation added to docker build

## [0.1.0] - 2025-10-23

### Added
- Initial release

[0.21.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.21.0...v0.21.1
[0.21.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.20.2...v0.21.0
[0.20.2]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.20.1...v0.20.2
[0.20.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.20.0...v0.20.1
[0.20.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.19.1...v0.20.0
[0.19.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.19.0...v0.19.1
[0.19.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.18.0...v0.19.0
[0.18.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.17.0...v0.18.0
[0.17.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.16.0...v0.17.0
[0.16.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.15.0...v0.16.0
[0.15.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.11.1...v0.12.0
[0.11.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.11.0...v0.11.1
[0.11.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.5.2...v0.6.0
[0.5.2]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/gooddata/gooddata-neobackstop/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/gooddata/gooddata-neobackstop/releases/tag/v0.1.0
