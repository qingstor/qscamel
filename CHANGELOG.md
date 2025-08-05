# Change Log

All notable changes to QingStor qscamel will be documented in this file.

## [v2.0.26] - 2025-08-05

### Added

- feat(s3/qingstor): Supports migration of empty directory files under the root directory (#341)
- chore: Delete unnecessary log printouts (#343)

## [v2.0.25] - 2025-07-24

### Fixed

- fix(copy): Incorrect metadata validation when qingstor-to-qingstor (#338)

## [v2.0.24] - 2024-08-02

### Added

- feat: Added MD5 check after migration (#333)
- feat(fs/qingstor): Added gbk encoding support (#335)

### Fixed

- fix: Fixed the List method repeatedly placing objects into the channe… (#336)

### Changed

- refactor: Upgrade qingstor SDK to v4 (#334)
- Updated the official document links (#332)

## [v2.0.23] - 2023-08-16

### Added

- endpoint(qingstor/s3/fs): Modified the method of multipart upload in the copy task  (#330)

## [v2.0.22] - 2023-06-02

### Added

- endpoint/qingstor: Add content-type for qingstor migration (#326)
- endpoint/fs: Add last modified time to the source endpoint (#327)
- feat: Add ignore_before and speed limit and migrate dynamic printing (#328)

## [v2.0.21] - 2022-06-29

### Added

- qingstor: Add migration folder and user-defined metadata function (#320)

## [v2.0.20] - 2022-02-26

### Added

- endpoint/qingstor: Added user modification timeout function for qingstor (#316)

## [v2.0.18] - 2021-03-15

### Fixed

- endpoint: Fix complete multipart too early (#310)

## [v2.0.17] - 2020-09-15

### Added

- ep/fs: Add support to copy files by symlinks (#242)

### Fixed

- endpoint: Fix complete multipart excuated too early (#232)
- ep/dst: Fix upload failed when file concurrent write (#238)

## [v2.0.16] - 2020-04-29

### Fixed

- endpoint/qingstor: Fix disable uri cleaning struct tag incorrect
- endpoint/azblob: Fix context deadline exceeded while reading (#139)

## [v2.0.15] - 2020-04-19

### Added

- endpoint: Add dst support for s3 (#89)
- endpoint: Add azblob src support (#124)
- utils: Support migrate object key starts with / (#126)

### Fixed

- Fix panic on hdfs (#78)
- endpoint/s3: Fix HeadObject not found not handled (#99)

## [v2.0.14] - 2020-02-13

### Added

- endpoint: Add hdfs support (#59)
- endpoint/qingstor: Add disable uri cleaning support (#72)

## [v2.0.13] - 2019-11-07

### Added

- endpoint/cos: Add support for tencent cloud cos

## [v2.0.12] - 2019-07-29

### Fixed

- endpoint/aliyun: Fix objects in "/" not listed

## [v2.0.11] - 2019-07-15

### Changed

- endpoint/qingstor: Update sdk to v3.0.2 (#51)

## [v2.0.10] - 2019-04-19

### Added

- endpoint: s3: Add disable_uri_cleaning support

### Fixed

- utils: Fix path not join correctly

## [v2.0.9] - 2019-04-18

### Fixed

- endpoint: s3: Fix canonicalized resource not encoded correctly

## [v2.0.8] - 2019-04-16

### Changed

- endpoint: s3: Use header signer instead

## [v2.0.7] - 2019-04-16

### Fixed

- endpoint: s3: Fix v2 signer not swap correctly

## [v2.0.6] - 2019-04-16

### Added

- endpoint: s3: Add signature v2 support

## [v2.0.5] - 2018-08-20

### Added

- endpoint: s3: Support path style
- endpoint: s3: Enable ListObjects with v1 support

### Changed

- endpoint: s3: Only detect next marker while it not returned

### Fixed

- endpoint: s3: Fix invalid memory address while some key missing
- endpoint: s3: Fix wrong params used in list objects
- endpoint: s3: Fix typo in enable list objects v2 config
- endpoint: qingstor: utils: Fix part size too large
- utils: file: Fix create folder failed on windows
- Fix directory not created with correct perm

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

[v2.0.25]: https://github.com/yunify/qscamel/compare/v2.0.25...v2.0.26
[v2.0.25]: https://github.com/yunify/qscamel/compare/v2.0.24...v2.0.25
[v2.0.24]: https://github.com/yunify/qscamel/compare/v2.0.23...v2.0.24
[v2.0.23]: https://github.com/yunify/qscamel/compare/v2.0.22...v2.0.23
[v2.0.22]: https://github.com/yunify/qscamel/compare/v2.0.21...v2.0.22
[v2.0.21]: https://github.com/yunify/qscamel/compare/v2.0.20...v2.0.21
[v2.0.20]: https://github.com/yunify/qscamel/compare/v2.0.19...v2.0.20
[v2.0.18]: https://github.com/yunify/qscamel/compare/v2.0.17...v2.0.18
[v2.0.17]: https://github.com/yunify/qscamel/compare/v2.0.16...v2.0.17
[v2.0.16]: https://github.com/yunify/qscamel/compare/v2.0.15...v2.0.16
[v2.0.15]: https://github.com/yunify/qscamel/compare/v2.0.14...v2.0.15
[v2.0.14]: https://github.com/yunify/qscamel/compare/v2.0.13...v2.0.14
[v2.0.13]: https://github.com/yunify/qscamel/compare/v2.0.12...v2.0.13
[v2.0.12]: https://github.com/yunify/qscamel/compare/v2.0.11...v2.0.12
[v2.0.11]: https://github.com/yunify/qscamel/compare/v2.0.10...v2.0.11
[v2.0.10]: https://github.com/yunify/qscamel/compare/v2.0.10...v2.0.9
[v2.0.9]: https://github.com/yunify/qscamel/compare/v2.0.9...v2.0.8
[v2.0.8]: https://github.com/yunify/qscamel/compare/v2.0.8...v2.0.7
[v2.0.7]: https://github.com/yunify/qscamel/compare/v2.0.7...v2.0.6
[v2.0.6]: https://github.com/yunify/qscamel/compare/v2.0.6...v2.0.5
[v2.0.5]: https://github.com/yunify/qscamel/compare/v2.0.4...v2.0.5
[v2.0.4]: https://github.com/yunify/qscamel/compare/v2.0.3...v2.0.4
[v2.0.3]: https://github.com/yunify/qscamel/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/yunify/qscamel/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/yunify/qscamel/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/yunify/qscamel/compare/v1.1.0...v2.0.0
[v1.1.0]: https://github.com/yunify/qscamel/compare/v1.0.0...v1.1.0
