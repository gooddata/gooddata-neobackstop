# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
