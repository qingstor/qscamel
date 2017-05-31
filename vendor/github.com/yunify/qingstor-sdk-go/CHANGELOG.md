# Change Log
All notable changes to QingStor SDK for Go will be documented in this file.

## [v2.2.5] - 2017-05-22

### Fixed

- Fix error in request URL query.
- Fix error in request header value.

## [v2.2.4] - 2017-03-28

### Fixed

- Fix type of Content-Type header.
- Add Content-Length to GetObjectOutput.
- Fix status code of DELETE CORS API.
- Fix type of object size for GET Bucket API.

### BREAKING CHANGES

- The type of content length and object size has been changed from `*int` to `*int64`.

## [v2.2.3] - 2017-03-10

### Added

- Allow user to append additional info to User-Agent

## [v2.2.2] - 2017-03-08

### Fixed

- Resource is not mandatory in bucket policy statement

## [v2.2.1] - 2017-03-05

### Changed

- Add "Encrypted" field to  "KeyType" struct

## [v2.2.0] - 2017-02-28

### Added

- Add ListMultipartUploads API.

### Fixed

- Fix request builder & signer.

## [v2.1.2] - 2017-01-16

### Fixed

- Fix request signer.

## [v2.1.1] - 2017-01-05

### Changed

- Fix logger output format, don't parse special characters.
- Rename package "errs" to "errors".

### Added

- Add type converters.

### BREAKING CHANGES

- Change value type in input and output to pointer.

## [v2.1.0] - 2016-12-23

### Changed

- Fix signer bug.
- Add more parameters to sign.

### Added

- Add request parameters for GET Object.
- Add IP address conditions for bucket policy.

## [v2.0.1] - 2016-12-15

### Changed

- Improve the implementation of deleting multiple objects.

## [v2.0.0] - 2016-12-14

### Added

- QingStor SDK for the Go programming language.

[v2.2.5]: https://github.com/yunify/qingstor-sdk-go/compare/v2.2.4...v2.2.5
[v2.2.4]: https://github.com/yunify/qingstor-sdk-go/compare/v2.2.3...v2.2.4
[v2.2.3]: https://github.com/yunify/qingstor-sdk-go/compare/v2.2.2...v2.2.3
[v2.2.2]: https://github.com/yunify/qingstor-sdk-go/compare/v2.2.1...v2.2.2
[v2.2.1]: https://github.com/yunify/qingstor-sdk-go/compare/v2.2.0...v2.2.1
[v2.2.0]: https://github.com/yunify/qingstor-sdk-go/compare/v2.1.2...v2.2.0
[v2.1.2]: https://github.com/yunify/qingstor-sdk-go/compare/v2.1.1...v2.1.2
[v2.1.1]: https://github.com/yunify/qingstor-sdk-go/compare/v2.1.0...v2.1.1
[v2.1.0]: https://github.com/yunify/qingstor-sdk-go/compare/v2.0.1...v2.1.0
[v2.0.1]: https://github.com/yunify/qingstor-sdk-go/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/yunify/qingstor-sdk-go/compare/v2.0.0...v2.0.0
