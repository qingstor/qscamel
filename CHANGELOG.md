# Change Log

All notable changes to QingStor qscamel will be documented in this file.

## [Unreleased]

### Added

- endpoint: s3: Support path style
- endpoint: s3: Enable ListObjects with v1 support

### Fixed

- endpoint: s3: Fix invalid memory address while some key missing

## [v2.0.4] - 2018-08-13

### Added

- endpoint: qingstor: Support multi-thread resumable multipart upload
- endpoint: qingstor: Support storage class in multipart upload
- endpoint: qingstor: Add user agent support

### Changed

- DB format changed, old task can't run by this version

## [v2.0.3] - 2018-08-06

### Changed

- migrate: Do not retry the same object too many times

### Fixed

- endpoint: fs: Fix can't create folder on windows

## [v2.0.2] - 2018-06-14

### Changed

- endpoint: qingstor: Do not use delimiter while listing

### Fixed

- Fix json file not handled correctly

## [v2.0.1] - 2018-06-06

### Changed

- Handle the db closed error

### Fixed

- Fix connections not reused on Windows

## [v2.0.0] - 2018-06-06

### Added

- Support migrate from local fs to QingStor
- Support migrate from QingStor to local fs

### Changed

- Use task file instead of command line argument
- Split task config from qscamel's config

## [v1.1.0] - 2017-09-14

### Added

- Support migration from UPYUN Storage Service to QingStor Object Storage
- Support migration from one Bucket to another in QingStor Object Storage
- Add a simple cmd progress bar during migration
- Update [`qingstor-sdk-go`](https://github.com/yunify/qingstor-sdk-go) version to v2.5.5
- Use glide to manage dependencies

### Fixed

- Fix missing DefaultThreadNum
- Fix an infinite loop problem in migration

## v1.0.0 - 2016-12-22

### Added

- QingStor qscamel.

[v2.0.4]: https://github.com/yunify/qscamel/compare/v2.0.3...v2.0.4
[v2.0.3]: https://github.com/yunify/qscamel/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/yunify/qscamel/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/yunify/qscamel/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/yunify/qscamel/compare/v1.1.0...v2.0.0
[v1.1.0]: https://github.com/yunify/qscamel/compare/v1.0.0...v1.1.0
