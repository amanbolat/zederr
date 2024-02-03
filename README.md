# zederr - Standardized Errors for Go

[![GoDoc](https://godoc.org/github.com/amanbolat/zederr?status.svg)](https://godoc.org/github.com/amanbolat/zederr)
[![Go Report Card](https://goreportcard.com/badge/github.com/amanbolat/zederr)](https://goreportcard.com/report/github.com/amanbolat/zederr)
[![codecov](https://codecov.io/gh/amanbolat/zederr/branch/master/graph/badge.svg)](https://codecov.io/gh/amanbolat/zederr)

## About

`zederr` is a tool for error standardization. It provides a way to define errors in a single place and generate code from YAML files. 

## What is generated?

**Go code:**

- Constructor functions to create strictly typed errors with custom parameters.
- Error method to convert its parameters to protobuf message.
- Error method to convert its parameters to `map[string]any`.

**Localization:**

- YAML file with error codes and translations.


## YAML file format

```yaml
common.file_too_large:
  http_code: 400
  grpc_code: 3
  description: If the file received by server is larger than applied limit this error should be returned
  translations:
    en: "File is too large, max size is {{ .MaxSize | int }}"
    zh_cn: "上传的文件不能大于{{ .MaxSize | int }}"

auth.unauthorized:
  http_code: 404
  grpc_code: 1
  description: User is not authorized to perform this action
  translations:
    en: "Please login to perform this action"
    zh_cn: "请登录再进行操作"

```

### Constraints

- Message parameter name can only contain ASCII characters.

| ✅Good            | ❌Wrong          |
|------------------|-----------------|
| `{{ .MaxSize }}` | `{{ .1Param }}` |


- The first character of message parameter name must be a letter.

| ✅Good            | ❌Wrong            |
|------------------|-------------------|
| `{{ .MaxSize }}` | `{{ ._MaxSize }}` |

# License

Apache License Version 2.0

