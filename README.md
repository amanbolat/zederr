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
spec_version: "1"
domain: acme.com
namespace: auth
default_locale: en

errors:
  account_locked:
    description: Account is locked due to too many failed login attempts.
    deprecated: ""
    http_code: 401
    grpc_code: 1
    arguments:
      user_id:
        type: "string"
        description: "User ID"
        is_internal: "true"
      failed_attempts:
        type: "int"
        description: "Number of failed login attempts"
      unlock_time:
        type: "timestamp"
        description: "Time when the account will be unlocked"
    title: "Account is locked"
    internal_message: "User with id {{ .user_id }} failed to provide correct credentials {{ .failed_attempts }} times. The account is locked until {{ .unlock_time }}."
    public_message: "Your account is locked due to too many failed login attempts ({{ .failed_attempts }})."
    localization:
      arguments:
        user_id:
          description:
            zh: "用户ID"
      description:
        zh: "由于登录尝试失败次数过多，帐户已被锁定。"
      title:
        zh: "帐户已锁定"
      internal_message:
        zh: "用户ID为{{ .user_id }}的用户提供了错误的凭据{{ .failed_attempts }}次。帐户将被锁定直到{{ .unlock_time }}。"
      public_message:
        zh: "由于登录尝试失败次数过多，您的帐户已被锁定({{ .failed_attempts }})。"
      deprecated:
        zh: ""
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

