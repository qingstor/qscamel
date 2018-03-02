# qscamel

[![Build Status](https://travis-ci.org/yunify/qscamel.svg?branch=master)](https://travis-ci.org/yunify/qscamel) [![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qscamel)](https://goreportcard.com/report/github.com/yunify/qscamel) [![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/yunify/qscamel/blob/master/LICENSE)

qscamel is a command line tool to migrate data between different endpoint efficiently.

## Features

- Easy task management
- Resume from where task stopped
- Automatic retry mechanism
- Concurrent migrating with goroutine pool
- Support both **copy** and **fetch** migrating methods
- Mutiple endpoint support

  - POSIX File System _(local fs, nfs, s3fs and so on)_
  - QingStor

## Quick start

Create a task file like following and save as `example-task.yaml`:

```yaml
name: example-task
type: copy

source:
  type: fs
  path: /path/to/source

destination:
  type: qingstor
  path: /path/to/destination
  options:
    bucket_name: example_bucket
    access_key_id: example_access_key_id
    secret_access_key: example_secret_access_key
```

Use qscamel:

```bash
qscamel run example-task.yaml
```

Have a cup of tea and you will see all file under `/path/to/source` will be migrated to qingstor bucket `example_bucket` with prefix `/path/to/destination`.

See the detailed usage with `qscamel -h` or `qscamel --help`, and read the following docs.

## Installation

Get the latest qscamel for Linux, macOS and Windows from [releases](https://github.com/yunify/qscamel/releases)

## Configuration

qscamel has following config options:

```yaml
# concurrency controls haw many goroutine run at the same time.
# if not set or set to 0, qscamel will use the numer of logic CPU * 100.
concurrency: 0
# log_level controls the log_level.
# Available value (from more to less): debug, info, warn, error, fatal, panic.
# Default value: info
log_level: info
# pid_file controls where the pid file will create.
# Default value: ~/.qscamel/qscamel.pid
pid_file: ~/.qscamel/qscamel.pid
# log_file controls where the log file will create.
# Default value: ~/.qscamel/qscamel.log
log_file: ~/.qscamel/qscamel.log
# database_file controls where the database file will create.
# Default value: ~/.qscamel/qscamel.db
database_file: ~/.qscamel/qscamel.db
```

The default config will read from `~/.qscamel/qscamel.yaml`, you can also specify the config path with flag `-c` or `--config`.

For example:

```bash
qscamel run example-task -c /path/to/config/file
```

## Task

Task file will define a task, and the task has following options:

```yaml
# name is the unique identifier for a task, qscamel will use it to distingush
# different tasks.
name: example-task
# type is the type for current task.
# Available value: copy, fetch, verify.
type: copy

# source is the source endpoint for current task.
source:
  # type is the type for endpoint.
  type: fs
  # path is the path for endpoint.
  path: "/home/xuanwo/Downloads/Telegram Desktop"

# destination is the destination endpoint for current task.
destination:
  # type is the type for endpoint.
  type: qingstor
  # path is the path for endpoint.
  path: /aaa
  # options is the options for differenty endpoint.
  options:
    bucket_name: example_bucket
    access_key_id: example_access_key_id
    secret_access_key: example_secret_access_key

# ignore_existing will control whether ignore existing object.
ignore_existing: false
# ignore_unmodified will control whether ignore unmodified object.
ignore_unmodified: false
```

### Endpoint fs

There is no more config for fs endpoint.

### Endpoint qingstor

qingstor endpoint has following options:

```yaml
# protocol controls protocol for qingstor.
# Available value: https, http
# Default value: https
protocol: https
# host controls host for qingstor.
# Default value: qingstor.com
host: qingstor.com
# port controls port for qingstor.
# Default value: 443
port: 443
# bucket_name is the bucket name for qingstor.
bucket_name: example_bucket
# access_key_id is the access_key_id for qingstor.
access_key_id: example_access_key_id
# secret_access_key is the secret_access_key for qingstor.
secret_access_key: example_secret_access_key
```

## Usage

### Run

Run is the main command for qscamel. We can use this command to create or resume a task.

In order to create a new task, we can use:

```bash
qscamel run /path/to/task/file
```

In order to resume a task, we can use:

```bash
qscamel run /path/to/task/file
```

or

```bash
qscamel run task-name
```

> When a new task created, we will calculate the sha256 checksum for it's content and save it
> to the database, and we will check if the content of the task file has been changed, if changed,
> qscamel will return an error. In other word, task can't be changed after created. If your need
> to update the task, please create a new one.

### Delete

Delete can delete a task.

```bash
qscamel delete task-name
```

### Status

Status will show the task status.

```bash
qscamel status
```

### Clean

Clean will delete all the finished tasks.

```bash
qscamel clean
```

### Version

Version will show the current version of qscamel.

```bash
qscamel version
```

## Contributing

Please see [_`Contributing Guidelines`_](./CONTRIBUTING.md) of this project before submitting patches.

## LICENSE

The Apache License (Version 2.0, January 2004).
