# zederr - Standardized Errors for Go

[![GoDoc](https://godoc.org/github.com/amanbolat/zederr?status.svg)](https://godoc.org/github.com/amanbolat/zederr)
[![Go Report Card](https://goreportcard.com/badge/github.com/amanbolat/zederr)](https://goreportcard.com/report/github.com/amanbolat/zederr)
[![codecov](https://codecov.io/gh/amanbolat/zederr/branch/master/graph/badge.svg)](https://codecov.io/gh/amanbolat/zederr)

## About

`zederr` is a tool for error codes documentation and code generation. You can define all the errors in one YAML file and generate strictly typed error constructors. Error public messages are automatically localized on initialization based on the user locale. 

## Why

If your application returns error,
it implies that the consumers rely on them and might have corresponding business logic.
The question is, how to identify those errors and write a reliable logic.
Comparing an error message can easily fail, as content may change.
Using gRPC or HTTP code is a nice solution, when you have a few different possible errors per endpoint.
One option left – custom error codes.

Introducing custom error codes solves one problem –
they become part of the API contract and consumer logic can rely on them,
but there are a few other challenges that you might face after.

- 📘**Documentation**. You have to document all the possible errors, provide a concise and comprehensive description.
- 🌎**Localization**. Yes, it's possible to maintain a map of error codes and corresponding translations on client side, in case you have SPA or mobile application. However, there are a few cases when you have to localize error messages on the backend:
  - The consumer of your API requires localized messages, for instance, if you are serving customers in different regions.
  - Error messages change frequently, and your web clients cannot be updated that frequently.
  - Errors are very complex objects, that should be pre-localized and easily be rendered on the client side.
  - You need to aggregate error messages from different microservices in API gateway or BFF.
- ⚠️️**Deprecation**. You cannot just remove errors from your API, as they are part of the contract between the service and the consumers.
- 🌀**Complex errors** with arguments and nested causes. Image you have a huge multistep form on your UI, that should be submitted with a single API call, and you need to return tens or hundreds of different errors constructed as a deeply nested tree. All those errors should be rendered on the client side in a way the user can easily understand each message and fix them at once or one-by-one. Or maybe you want to link that list of returned errors with fields in your form. 

`zederr` tries to address all those issues. It requires you to write an error specification the same way you write OpenAPI specification to describe your service API. All the errors are declared in Error Specification file. The latter will be used to generate Go code to easily construct errors in runtime. 

## Features

- Specification. Create a specification for all errors that your application might return.
- Documentation. Auto-generate a nice documentation to share with your service consumers.
- Localization. Localize all the error messages.
- Go code generation. Generate strictly typed error constructors.

## Error types

- Detailed error
- Minimal error


## Example

The example below is a good illustration of `zederr` functionality.

Specification file:

```yaml
# Specification version tells the parser which version of the specification to use.
# Required.
spec_version: "1"
# Default locale for error content localization. Example: `en`.
# Required.
default_locale: en

# The collection of error codes and their metadata.
errors:
  # Unique error code.
  # It will be used to generate a human-readable error code constructor.
  "account_locked":
    # Description of the error code.
    # Required.
    description: Account is locked due to too many failed login attempts.
    # Deprecation reason. If not empty, the error code is considered deprecated.
    # Optional.
    is_deprecated: false
    # HTTP status code for the error code.
    # Default: 500
    # Optional.
    http_code: 401
    # A gRPC status code for the error code.
    grpc_code: 1
    # Arguments that can be used in the error message templates.
    # Optional.
    arguments:
      user_id:
        # Argument type.
        # Required.
        type: "string"
        # Argument description.
        # Optional.
        description: "User ID"
      failed_attempts:
        type: "int"
        description: "Number of failed login attempts"
      unlock_time:
        type: "timestamp"
        description: "Time when the account will be unlocked"
    # Error message template.
    # Required.
    message: "Your account is locked due to too many failed login attempts ({{ .failed_attempts }}). It will be unlocked at {{ .unlock_time }}."
    # Localization for the error code.
    # Optional.
    localization:
      # Localization for the arguments.
      # Required for all the arguments listed in the `arguments` section.
      arguments:
        user_id:
          description:
            zh: "用户ID"
      description:
        zh: "由于登录尝试失败次数过多，帐户已被锁定。"
      message:
        zh: "由于登录尝试失败次数过多，您的帐户已被锁定({{ .failed_attempts }})。"
```

Running the command bellow will generate Go constructors for each error in the specification file:

```shell
zederr gen --go-out ./out --spec zederr_spec.yaml
```

As you can see `zederr` generates a constructor that requires you to provide all arguments,
and each argument has the correct type:

```go
func NewAccountLocked(ctx context.Context, user_id string, failed_attempts int, unlock_time time.Time) *zeerr.Error {
	return zeerr.NewError(
		ctx,
		localizer,
		"account_locked",
		401,
		pkgcodes.Canceled,
		map[string]any{
			"user_id":         user_id,
			"failed_attempts": failed_attempts,
			"unlock_time":     unlock_time,
		},
	)
}
```

If you construct that error and encode to JSON:

```go
err := zederr.NewAccountLocked(zeerr.ContextWithLocale(context.Background(), language.Chinese), "user_1", 1, time.Now())
```

The error message of `err.Error()` will be:

```text
由于登录尝试失败次数过多，您的帐户已被锁定(1)。其将在2024-06-26 00:36:06.33748 +0200 CEST m=+0.002228543自动解冻
```

## License

Apache License Version 2.0

