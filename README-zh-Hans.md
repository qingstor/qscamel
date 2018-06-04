# qscamel

[![Build Status](https://travis-ci.org/yunify/qscamel.svg?branch=master)](https://travis-ci.org/yunify/qscamel) [![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qscamel)](https://goreportcard.com/report/github.com/yunify/qscamel) [![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/yunify/qscamel/blob/master/LICENSE)

qscamel 是一个用于在不同的端点 (Endpoint) 中高效迁移数据的工具。

## 功能

- 简单，便于使用的任务管理
- 从任务中断处续传，节省宝贵的时间
- 完全自动化的重试机制
- 基于 Goroutine 池实现的并发机制
- 支持 **copy**, **fetch**, **delete** 等迁移机制
- 支持数据校验
- 多端点支持

  - 符合 POSIX 标准的文件系统 _(local fs, nfs, s3fs 等)_
  - 本地文件列表
  - [QingStor](https://www.qingcloud.com/products/qingstor)
  - [Aliyun OSS](https://www.aliyun.com/product/oss)
  - [Google Cloud Storage](https://cloud.google.com/storage/)
  - [Qiniu](https://www.qiniu.com/)
  - [AWS S3](https://amazonaws-china.com/cn/s3)
  - [Upyun](https://www.upyun.com/)

## 快速入门

像下面一样创建一个任务文件，并保存为 `example-task.yaml`：

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

使用 `qscamel`：

```bash
qscamel run example-task -t example-task.yaml
```

坐下来，喝杯茶。你将会看到所有在 `/path/to/source` 目录下的文件都会被迁移到 QingStor 的 Bucket `example_bucket` 的 `/path/to/destination` 前缀下。

可以通过 `qscamel -h` 或者 `qscamel --help` 来了解更多使用上的细节，并且阅读下面的文档。

## 安装

可以在 [releases](https://github.com/yunify/qscamel/releases) 页面获取适用于 Linux, macOS 和 Windows 的最新版 qscamel。

## 配置

qscamel 有如下配置项：

```yaml
# concurrency 会控制同时启用的并发数量。
# 如果没有设置，或者设置为 0， qscamle 将会使用逻辑 CPU 数量 * 100 作为该项的值。
concurrency: 0
# log_level 控制日志的级别。
# 可选值（从更多到更少）： debug, info, warn, error, fatal, panic.
# 默认值： info
log_level: info
# pid_file 将会控制在何处创建 PID 文件。
# 默认值: ~/.qscamel/qscamel.pid
pid_file: ~/.qscamel/qscamel.pid
# log_file 将会控制在何处创建日志文件。
# 默认值: ~/.qscamel/qscamel.log
log_file: ~/.qscamel/qscamel.log
# database_file 将会控制在何处创建数据库。
# 默认值: ~/.qscamel/db
database_file: ~/.qscamel/db
```

qscamel 默认从 `~/.qscamel/qscamel.yaml` 读取配置文件，你也可以通过 `-c` 或者 `--config` 来指定配置文件的位置。
通过指定不同的配置文件，你可以同时运行多个 qscamel 实例。

比如:

```bash
qscamel run example-task -t example-task.yaml -c /path/to/config/file
```

## 任务

任务文件将会定义一个任务，每个任务都有如下配置：

```yaml
# type 是任务的类型。
# 可选值: copy, fetch, delete
# copy 将会从 source 处读取文件，并写入到 destination。
# fetch 将会从 source 处获取文件的下载链接，并使用 destination 的 fetch 功能进行拉取。
# delete 将会从 source 处获取文件的信息，并在 destination 处删除。
type: copy

# source 是任务的 source 端点。
source:
  # type 是当前端点的类型。
  # 可选值: aliyun, fs, filelist, gcs, qingstor, qiniu, s3, upyun.
  type: fs
  # path 是当前端点的路径。
  path: "/path/to/source"

# destination 是任务的 destination 端点。
destination:
  # type 是当前端点的类型。
  # 可选值: fs, qingstor.
  type: qingstor
  # path 是当前端点的路径。
  path: /aaa
  # options 是不同端点的配置，详情请参考下面的文档。
  options:
    bucket_name: example_bucket
    access_key_id: example_access_key_id
    secret_access_key: example_secret_access_key

# ignore_existing 控制是否跳过已经存在的文件。
# `disable` 将会禁用该配置，即总是覆盖
# `size` 当文件的 size 相同时会跳过
# `quick_md5sum` 将会对文件做一次快速 md5 计算，当 md5 相同时会跳过
# `full_md5sum` 将会对文件做完整的 md5 计算，当 md5 相同时会跳过
# 可选值: disable, size, quick_md5sum, full_md5sum.
# 默认值: disable
ignore_existing: disable
```

### Endpoint aliyun

能够用做 **source** 端点。

Aliyun 是 [阿里云](https://www.aliyun.com/product/oss) 提供的对象存储服务。

aliyun 端点有如下配置内容:

```yaml
endpoint: example_endpoint
bucket_name: example_bucket
access_key_id: example_access_key_id
access_key_secret: example_access_key_secret
```

### Endpoint fs

能够用做 **source** 和 **destination** 端点。

fs 端点没有更多的配置内容。

### Endpoint filelist

能够用做 **source** 端点。

qscamel 将会按照行来读取该列表。

```yaml
list_path: /path/to/list
```

### Endpoint gcs

能够用做 **source** 端点。

GCS(Google Cloud Storage) 是 [Google](https://cloud.google.com/storage/) 提供的对象存储服务。

gcs 端点有如下配置内容:

```yaml
api_key: example_api_key
bucket_name: exmaple_bukcet
```

### Endpoint qingstor

能够用做 **source** 和 **destination** 端点。

qingstor 端点有如下配置内容:

```yaml
# protocol 控制访问 QingStor 的协议类型。
# 可选值: https, http
# 默认值: https
protocol: https
# host 控制访问 QingStor 的主机名。
# 默认值: qingstor.com
host: qingstor.com
# port 控制访问 QingStor 的端口号。
# 默认值: 443
port: 443
# zone 控制访问 QingStor 的区域.
# 默认值：自动检测，不需要手动配置
# This will auto detected, no need to set.
zone: pek3b
# bucket_name 是 QingStor 的 bucket 名称。
bucket_name: example_bucket
# access_key_id 是 QingStor 的 access_key_id。
access_key_id: example_access_key_id
# secret_access_key 是 QingStor 的 secret_access_key。
secret_access_key: example_secret_access_key

# storage class 是 QingStor 所使用的存储级别
# 可选值: STANDARD, STANDARD_IA
# 默认值: STANDARD
storage_class: STANDARD

# multipart boundary size 控制 QingStor 何时使用分段上传
# 单位为 Byte ，当文件大于该数值时，将会使用分段上传
# 可选值: 1 ~ 5368709120
# 默认值: 2147483648
multipart_boundary_size: 2147483648
```

### Endpoint qiniu

能够用做 **source** 端点。

Qiniu 是 [Qiniu](https://www.qiniu.com/) 提供的对象存储服务。

qiniu 端点有如下配置内容:

```yaml
# bucket_name 是 qiniu 的 bucket 名称
bucket_name: example_bucket
# access_key 是 qiniu 的 access_key
access_key: example_access_key
# secret_key 是 qiniu 的 secret_key
secret_key: example_secret_key
# domain 是用于访问 qiniu bucket 的域名
domain: example_domain
# use_https 控制是否使用 https 来访问 qiniu
# 默认值： false
use_https: false
# use_cdn_domains 控制是否使用 CDN 加速域名来访问 qiniu
# 默认值： false
use_cdn_domains: false
```

### Endpoint s3

能够用做 **source** 端点。

S3 是 [AWS](https://amazonaws-china.com/cn/s3) 提供的对象存储服务。

s3 端点有如下配置内容:

```yaml
bucket_name: example_bucket
endpoint: example_endpoint
region: example_region
access_key_id: example_access_key_id
secret_access_key: example_secret_access_key
disable_ssl: false
use_accelerate: false
```

### Endpoint upyun

能够用做 **source** 端点。

upyun 是 [Upyun](https://www.upyun.com/) 提供的对象存储服务。

upyun 端点有如下配置内容:

```yaml
bucket_name: example_bucket
operator: example_operator
password: example_password
```

## 用法

### Run

Run 是 qscamel 最主要的命令。我们使用这个命令来创建或者恢复一个任务。

如果要创建一个任务，我们可以使用：

```bash
qscamel run task-name -t /path/to/task/file
```

如果要恢复一个任务，我们可以使用：

```bash
qscamel run task-name -t /path/to/task/file
```

or

```bash
qscamel run task-name
```

> 当一个新任务创建的时候就，我们将会计算任务内容的 sha256 校验和并且保存在数据库当中，同时我们还会检查任务文件的内容是否发生了修改。如果改变了，qscamel 将会返回一个错误并退出。换句话说，任务在创建完毕后就不能修改。如果你需要修改一个任务的内容，请创建一个新任务。

### Delete

Delete 能够删除一个任务。

```bash
qscamel delete task-name
```

### Status

Status 将会展示所有任务的状态。

```bash
qscamel status
```

### Clean

Clean 将会删除所有已经完成的任务。

```bash
qscamel clean
```

### Version

Version 将会显示当前 qscamel 的版本。

```bash
qscamel version
```

## 贡献

请在提交 Patch 前阅读 [_`Contributing Guidelines`_](./CONTRIBUTING.md)。

## 协议

The Apache License (Version 2.0, January 2004).
