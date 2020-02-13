# qscamel

[![Build Status](https://travis-ci.org/yunify/qscamel.svg?branch=master)](https://travis-ci.org/yunify/qscamel) [![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qscamel)](https://goreportcard.com/report/github.com/yunify/qscamel) [![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/yunify/qscamel/blob/master/LICENSE)

qscamel is a command line tool to migrate data between different endpoint efficiently.

## Features

- Easy task management
- Resume from where task stopped
- Automatic retry mechanism
- Concurrent migrating with goroutine pool
- Support **copy**, **fetch** and **delete** migrating methods
- Support data verify
- Multiple endpoint support

  - POSIX File System _(local fs, nfs, s3fs and so on)_
  - Local file list
  - [QingStor](https://www.qingcloud.com/products/qingstor)
  - [Aliyun OSS](https://www.aliyun.com/product/oss)
  - [Google Cloud Storage](https://cloud.google.com/storage/)
  - [Qiniu](https://www.qiniu.com/)
  - [AWS S3](https://amazonaws-china.com/cn/s3)
  - [Upyun](https://www.upyun.com/)
  - [Tencent COS](https://cloud.tencent.com/product/cos)
  - [HDFS](http://hadoop.apache.org/)

## Quick start

Create a task file like following and save as `example-task.yaml`:

```yaml
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
qscamel run example-task -t example-task.yaml
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
# database_file controls where the database folder will create.
# Default value: ~/.qscamel/db
database_file: ~/.qscamel/db
# Proxy that qscamel used to connect endpoint.
proxy: ""
```

The default config will read from `~/.qscamel/qscamel.yaml`, you can also specify the config path with flag `-c` or `--config`.

For example:

```bash
qscamel run example-task -t example-task.yaml -c /path/to/config/file
```

## Task

Task file will define a task, and the task has following options:

```yaml
# type is the type for current task.
# Available value: copy, fetch, delete
type: copy

# source is the source endpoint for current task.
source:
  # type is the type for endpoint.
  # Available value: aliyun, cos, fs, filelist, gcs, qingstor, qiniu, s3, upyun.
  type: fs
  # path is the path for endpoint.
  path: "/path/to/source"

# destination is the destination endpoint for current task.
destination:
  # type is the type for endpoint.
  # Available value: fs, qingstor.
  type: qingstor
  # path is the path for endpoint.
  path: /aaa
  # options is the options for differenty endpoint.
  options:
    bucket_name: example_bucket
    access_key_id: example_access_key_id
    secret_access_key: example_secret_access_key

# ignore_existing controls whether and how to ignore exist file.
# If set to empty or not set, this config will be disabled.
# `last_modified` will check whether the dst object's last modified is greater than src's one.
# `md5sum` will calculate the whole object's md5.
# Available value: last_modified, md5sum.
ignore_existing: last_modified
# multipart boundary size controls when qscamel will use multipart
# unit is Byte ，when file size is bigger then this value, qscamel
# will use multipart API.
# Available value: 1 ~ 5368709120
# Default value: 2147483648
multipart_boundary_size: 2147483648
```

### Endpoint aliyun

Can be used as **source** endpoint.

Aliyun is the object storage service provided by [Alibaba](https://www.aliyun.com/product/oss).

aliyun endpoint has following options:

```yaml
endpoint: example_endpoint
bucket_name: example_bucket
access_key_id: example_access_key_id
access_key_secret: example_access_key_secret
```

### Endpoint cos

Can be used as **source** endpoint.

COS is the object storage service provided by [Tencent Cloud](https://cloud.tencent.com/product/cos).

cos endpoint has following options:

```yaml
bucket_url: https://example-123456789.cos.ap-beijing.myqcloud.com
secret_id: example_secret_id
secret_key: example_secret_key
```

### Endpoint fs

Can be used as **source** and **destination** endpoint.

There is no more config for fs endpoint.

### Endpoint filelist

Can be used as **source** endpoint.

qscamel will read filelist by line.

```yaml
list_path: /path/to/list
```

### Endpoint gcs

Can be used as **source** endpoint.

GCS(Google Cloud Storage) is the object storage service provided by [Google](https://cloud.google.com/storage/).

gcs endpoint has following options:

```yaml
api_key: example_api_key
bucket_name: exmaple_bukcet
```

### Endpoint hdfs

Can be used as **source** endpoint.

hdfs endpoint has following options:

```yaml
address: 127.0.0.1:8080
```

### Endpoint qingstor

Can be used as **source** and **destination** endpoint.

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
# zone controls zone for qingstor.
# This will auto detected, no need to set.
zone: pek3b
# bucket_name is the bucket name for qingstor.
bucket_name: example_bucket
# access_key_id is the access_key_id for qingstor.
access_key_id: example_access_key_id
# secret_access_key is the secret_access_key for qingstor.
secret_access_key: example_secret_access_key

# storage class is the storage class used for qingstor.
# Available value: STANDARD, STANDARD_IA
# Default value: STANDARD
storage_class: STANDARD
# disable_uri_cleaning will control whether or not the SDK will do
# cleaning on object key: `abc//def` -> `abc/def`
# Available value: true, false
# Default value: false
disable_uri_cleaning: false
```

### Endpoint qiniu

Can be used as **source** endpoint.

Qiniu is the object storage service provided by [Qiniu](https://www.qiniu.com/).

qiniu endpoint has following options:

```yaml
bucket_name: example_bucket
access_key: example_access_key
secret_key: example_secret_key
domain: example_domain
use_https: false
use_cdn_domains: false
```

### Endpoint s3

Can be used as **source** endpoint.

S3 is the object storage service provided by [AWS](https://amazonaws-china.com/cn/s3).

s3 endpoint has following options.

```yaml
bucket_name: example_bucket
endpoint: example_endpoint
region: example_region
access_key_id: example_access_key_id
secret_access_key: example_secret_access_key
disable_ssl: false
use_accelerate: false
path_style: false
enable_list_objects_v2: false
enable_signature_v2: false
disable_uri_cleaning: false
```

- `enable_signature_v2` is added for compatible usage in ceph and other S3-alike service.
- `disable_uri_cleaning` is added to control aws s3 sdk's url clean behavior.

### Endpoint upyun

Can be used as **source** endpoint.

upyun is the object storage service provided by [Upyun](https://www.upyun.com/).

upyun endpoint has following options.

```yaml
bucket_name: example_bucket
operator: example_operator
password: example_password
```

## Usage

### Run

Run is the main command for qscamel. We can use this command to create or resume a task.

In order to create a new task, we can use:

```bash
qscamel run task-name -t /path/to/task/file
```

In order to resume a task, we can use:

```bash
qscamel run task-name -t /path/to/task/file
```

or

```bash
qscamel run task-name
```

> When a new task created, we will calculate the sha256 checksum for it's content and save it to the database, and we will check if the content of the task file has been changed, if changed, qscamel will return an error. In other word, task can't be changed after created. If your need to update the task, please create a new one.

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
