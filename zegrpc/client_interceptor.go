package zegrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pbzederrv1 "github.com/amanbolat/zederr/zeproto/v1"
)

type clientInterceptorConfig struct {
	decoder Decoder
}

func defaultClientInterceptorConfig() *clientInterceptorConfig {
	return &clientInterceptorConfig{
		decoder: SimpleDecoder{},
	}
}

type ClientInterceptorOption func(config *clientInterceptorConfig)

func StreamClientInterceptor(opts ...ClientInterceptorOption) grpc.StreamClientInterceptor {
	cfg := defaultClientInterceptorConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		cltStream, err := streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			sts, ok := status.FromError(err)
			if !ok {
				return nil, err
			}

			for _, detail := range sts.Details() {
				if v, ok := detail.(*pbzederrv1.Error); ok {
					zedErr := cfg.decoder.Decode(v)

					return nil, zedErr
				}
			}
		}

		return cltStream, nil
	}
}

func UnaryClientInterceptor(opts ...ClientInterceptorOption) grpc.UnaryClientInterceptor {
	cfg := defaultClientInterceptorConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			sts, ok := status.FromError(err)
			if !ok {
				return err
			}

			for _, detail := range sts.Details() {
				if v, ok := detail.(*pbzederrv1.Error); ok {
					zedErr := cfg.decoder.Decode(v)

					return zedErr
				}
			}
		}

		return nil
	}
}
