package zegrpc

import (
	"google.golang.org/grpc/codes"

	"github.com/amanbolat/zederr/internal/transport"
	"github.com/amanbolat/zederr/zeerr"
	pbzederrv1 "github.com/amanbolat/zederr/zeproto/v1"
)

func Decode(pbErr *pbzederrv1.Error) (zeerr.Error, error) {
	return decode(pbErr)
}

func decode(pbErr *pbzederrv1.Error) (zeerr.Error, error) {
	args := make(zeerr.Arguments)

	if pbErr.Arguments != nil {
		args = pbErr.Arguments.AsMap()
	}

	zedErr := transport.NewError(
		pbErr.Code,
		pbErr.Domain,
		pbErr.Namespace,
		pbErr.Uid,
		int(pbErr.HttpCode),
		codes.Code(pbErr.GrpcCode),
		args,
		pbErr.InternalMessage,
		pbErr.PublicMessage,
		nil,
	)

	if len(pbErr.Causes) == 0 {
		return zedErr, nil
	}

	for _, c := range pbErr.Causes {
		cause, err := decode(c)
		if err != nil {
			return nil, err
		}

		zedErr = zedErr.WithCauses(cause)
	}

	return zedErr, nil
}
