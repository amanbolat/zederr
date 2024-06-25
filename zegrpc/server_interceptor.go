package zegrpc

import (
	"context"

	"google.golang.org/grpc"
)

type ErrorMapperFunc func(context.Context, error) error

func defaultErrMapper(_ context.Context, err error) error {
	return err
}

type serverInterceptorConfig struct {
	errMapperFunc ErrorMapperFunc
	encoder       Encoder
}

func defaultServerInterceptorConfig() *serverInterceptorConfig {
	return &serverInterceptorConfig{
		errMapperFunc: defaultErrMapper,
		encoder:       NewSimpleEncoder(defaultStatusCode, defaultStatusMessage),
	}
}

type ServerInterceptorOption func(*serverInterceptorConfig)

func WithErrorMapper(errMapperFunc ErrorMapperFunc) ServerInterceptorOption {
	return func(c *serverInterceptorConfig) {
		c.errMapperFunc = errMapperFunc
	}
}

func WithEncoder(encoder Encoder) ServerInterceptorOption {
	return func(c *serverInterceptorConfig) {
		c.encoder = encoder
	}
}

func StreamServerInterceptor(opts ...ServerInterceptorOption) grpc.StreamServerInterceptor {
	cfg := defaultServerInterceptorConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err == nil {
			return nil
		}

		zedErr := cfg.errMapperFunc(ss.Context(), err)
		sts := cfg.encoder.Encode(zedErr)

		return sts.Err()
	}
}

func UnaryServerInterceptor(opts ...ServerInterceptorOption) grpc.UnaryServerInterceptor {
	cfg := defaultServerInterceptorConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		zedErr := cfg.errMapperFunc(ctx, err)
		sts := cfg.encoder.Encode(zedErr)

		return resp, sts.Err()
	}
}
