package zegrpc

import (
	"context"

	"google.golang.org/grpc"
)

type ErrorMapperFunc func(context.Context, error) error

func defaultErrMapper(_ context.Context, err error) error {
	return err
}

type serverConfig struct {
	errMapperFunc ErrorMapperFunc
}

func defaultServerConfig() *serverConfig {
	return &serverConfig{
		errMapperFunc: defaultErrMapper,
	}
}

type Option func(*serverConfig)

func WithErrorMapper(errMapperFunc ErrorMapperFunc) Option {
	return func(c *serverConfig) {
		c.errMapperFunc = errMapperFunc
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	cfg := defaultServerConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err == nil {
			return nil
		}

		zedErr := cfg.errMapperFunc(ss.Context(), err)
		sts := Encode(zedErr)

		return sts.Err()
	}
}

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	cfg := defaultServerConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		zedErr := cfg.errMapperFunc(ctx, err)
		sts := Encode(zedErr)

		return resp, sts.Err()
	}
}
