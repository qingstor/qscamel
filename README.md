# qscamel


<span style="display: inline-block">
[![Build Status](https://travis-ci.org/yunify/qscamel.svg?branch=master)]
(https://travis-ci.org/yunify/qscamel)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qscamel)]
(https://goreportcard.com/report/github.com/yunify/qscamel)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)]
(https://github.com/yunify/qscamel/blob/master/LICENSE)
</span>

qscamel is a command line tool facilitating data migration to QingStor.
It can read from local file or other object storage platform to get source object list, then does batch migration from source site to QingStor.

## Getting Started

### Binary Download

Get the latest qscamel for Linux, macOS and Windows from [releases]

### Preparation

To use qscamel, there must be a configuration file to configure your own
`access_key_id` and `secret_access_key`, for example:

``` bash
access_key_id: 'QINGCLOUDACCESSKEYID'
secret_access_key: 'QINGCLOUDSECRETACCESSKEYEXAMPLE'
```

The configuration file is `~/.qingstor/config.yaml` by default, it also
can be specified by the option `-c /path/to/config`.

## Usage

### Command-line Flags

Command-line flags supported by qscamel are listed below.

##### General flags:

| short | full | type | required | usage |
| ----- |------|:------:|:----------:|------ |
| -t | --src-type    | string | Y | Specify source type, support "file" and other object storage like "s3", "qiniu" and "aliyun".
| -s | --src         | string | Y | Specify migration source. If --src-type is "file", --src specifies source list file. Otherwise, --src specifies source bucket name.
| -b | --bucket      | string | Y | Specify QingStor bucket
| -d | --description | string | Y | Describe current migration task. This description will be used as record filename for task resuming.
| -c | --config      | string | N | Specify QingStor YAML configuration file
| -T | --threads     | int    | N | Specify the number of objects being migrated concurrently (default 10, max 100)
| -l | --log-file    | string | N | Specify the path of log file
| -v | --version     | bool   | N | Print the version number of qscamel and exit
| -h | --help        | bool   | N | Print the usage of qscamel and exit

##### Overwrite flags:

By default, qscamel compares last modified time of source file with last modified time of object in QingStor.
If the object doesn't exist in QingStor or object in QingStor isn't the latest, then do migration.

| short | full | type | required | usage |
| ----- |------|:------:|:----------:|------ |
| -i | --ignore-existing | bool   | N | Ignore existing object
| -o | --overwrite       | bool   | N | Overwrite existing object
| -n | --dry-run         | bool   | N | Perform a trial run with no actual migration

##### Object storage source flags:

| short | full | type | required | usage |
| ----- |------|:------:|:----------:|------ |
| -z | --src-zone        | string | N | Specify source zone for object storage type source
| -a | --src-access-key  | string | N | Specify source access_key_id for object storage type source
| -S | --src-secret-key  | string | N | Specify source secret_access_key for object storage type source

### Source List File Format

Use `--src-type=file` or `-t file` to enable reading from file. Then use `--src` or `-s` to specify source list file.
Source list file defines http/https source site object for migration.
Each line of source site ends with `\n`. There are two format of one line:

1.Just source site link. Object name will be parsed from url, for example:

``` bash
# Http url with path, object name: public/cat.png
http://image.example.com/public/cat.png

# Http url with no path, object name: image.example.com
http://image.example.com
```

2.Specify object name. Format: source site`[spacing]`object name, for example:

``` bash
# Specify object name: archive/cat.png
http://image.example.com/public/cat.png archive/cat.png
```

### Supported Object Storage Source

Use `--src-type=<platform>` or `-t <platform>` to enable migrating from other platform, for example `--src-type=s3`. Then use `--src` or `-s` to specify source bucket name.

| platform | require --src-zone | require --src-access | require --src-secret |
| -------- |:------------------:|:--------------------:|:--------------------:|
| s3       | Y                  | Y                    | Y                    | 
| qiniu    | N                  | Y                    | Y                    |
| aliyun   | Y                  | Y                    | Y                    |

### Examples

``` bash
# Read from source list file
$ qscamel -t file -s ~/source-list -b QingStor-bucket-name -d "migrate 01"

# Overwrite existing object
$ qscamel -t file -s ~/source-list -b QingStor-bucket-name -d "migrate 02" -o

# Ignore existing object and dry-run
$ qscamel -t file -s ~/source-list -b QingStor-bucket-name -d "migrate 03" -i -n

# Specify threats and log-file
$ qscamel -t file -s ~/source-list -b QingStor-bucket-name -d "migrate 04" -T 5 -l ~/logfile

# Migrate from aws s3
$ qscamel -t s3 -s s3-bucket-name -z us-east-1 -a "S3ACCESSKEYID" -S "S3SECRETACCESSKEY" -b QingStor-bucket-name -d "migrate 05"

# Migrate from qiniu
$ qscamel -t qiniu -s qiniu-bucket-name -a "QINIUACCESSKEYID" -S "QINIUSECRETACCESSKEY" -b QingStor-bucket-name -d "migrate 06"

# Migrate from aliyun
$ qscamel -t aliyun -s aliyun-bucket-name -z oss-cn-shanghai -a "ALIYUNACCESSKEYID" -S "ALIYUNSECRETACCESSKEY" -b QingStor-bucket-name -d "migrate 07"
```

See the detailed usage with `qscamel -h` or `qscamel --help`.

## Contributing

Please see [_`Contributing Guidelines`_](./CONTRIBUTING.md) of this project before submitting patches.

## LICENSE

The Apache License (Version 2.0, January 2004).

  [releases]: https://github.com/yunify/qscamel/releases
