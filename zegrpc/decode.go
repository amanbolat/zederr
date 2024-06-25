package zegrpc

import (
	"google.golang.org/grpc/codes"

	"github.com/amanbolat/zederr/zeerr"
	pbzederrv1 "github.com/amanbolat/zederr/zeproto/v1"
)

type Decoder interface {
	Decode(pbErr *pbzederrv1.Error) *zeerr.Error
}

type SimpleDecoder struct{}

func (d SimpleDecoder) Decode(pbErr *pbzederrv1.Error) *zeerr.Error {
	return d.decode(pbErr)
}

func (d SimpleDecoder) decode(pbErr *pbzederrv1.Error) *zeerr.Error {
	args := make(map[string]any)

	if pbErr.Arguments != nil {
		args = pbErr.Arguments.AsMap()
	}

	zedErr := zeerr.RestoreError(
		pbErr.Id,
		int(pbErr.HttpCode),
		codes.Code(pbErr.GrpcCode),
		args,
		pbErr.Message,
		nil,
	)

	if len(pbErr.Causes) == 0 {
		return zedErr
	}

	for _, c := range pbErr.Causes {
		cause := d.decode(c)

		zedErr = zedErr.WithCauses(cause)
	}

	return zedErr
}
