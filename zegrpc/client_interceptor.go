package zegrpc

import (
	"context"

	pbzederrv1 "github.com/amanbolat/zederr/zeproto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		cltStream, err := streamer(ctx, desc, cc, method, opts...)
		if err == nil {
			return cltStream, nil
		}

		sts, ok := status.FromError(err)
		if !ok {
			return cltStream, err
		}

		for _, detail := range sts.Details() {
			if v, ok := detail.(*pbzederrv1.Error); ok {
				zedErr, decodeErr := Decode(v)
				if decodeErr == nil {
					return cltStream, zedErr
				}
			}
		}

		return cltStream, err
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}

		sts, ok := status.FromError(err)
		if !ok {
			return err
		}

		for _, detail := range sts.Details() {
			if v, ok := detail.(*pbzederrv1.Error); ok {
				zedErr, err := Decode(v)
				if err == nil {
					return zedErr
				}
			}
		}

		return err
	}
}
