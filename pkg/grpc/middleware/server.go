package middleware

import (
	"context"

	"github.com/amanbolat/zederr/pkg/core"
	"github.com/amanbolat/zederr/pkg/grpc/encode"
	"google.golang.org/grpc"
)

type ErrorMapper interface {
	MapError(err error) core.Error
}

func StreamServerInterceptor(errMapper ErrorMapper) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err == nil {
			return err
		}

		zedErr := errMapper.MapError(err)
		sts, err := encode.Encode(zedErr)
		if err != nil {
			panic(err)
		}

		return sts.Err()
	}
}

func UnaryServerInterceptor(errMapper ErrorMapper) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, err
		}

		zedErr := errMapper.MapError(err)
		sts, err := encode.Encode(zedErr)
		if err != nil {
			panic(err)
		}

		return resp, sts.Err()
	}
}
