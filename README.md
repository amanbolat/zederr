# zederr - Standardized Errors for Go

[![GoDoc](https://godoc.org/github.com/amanbolat/zederr?status.svg)](https://godoc.org/github.com/amanbolat/zederr)
[![Go Report Card](https://goreportcard.com/badge/github.com/amanbolat/zederr)](https://goreportcard.com/report/github.com/amanbolat/zederr)
[![codecov](https://codecov.io/gh/amanbolat/zederr/branch/master/graph/badge.svg)](https://codecov.io/gh/amanbolat/zederr)

# About

`zederr` is a tool for error codes documentation and code generation. You can define all the errors in one YAML file and generate strictly typed error constructors. All the errors will be automatically localized depending on the user locale. Errors can be passed from one service to another or returned to the end user. 

# Features

- Unified way to define and document errors.
- Localized error messages.
- Strictly typed error constructors in Go.
- gRPC middleware for easy error handling between services.

# How to use

1. Install `zederr` locally.
2. Create a YAML file with error definitions.
3. Run `zederr` to generate code.
4. Use generated error constructors in your code.
5. Setup middlewares to handle errors between services.

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

