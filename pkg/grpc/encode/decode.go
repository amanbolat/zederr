package encode

import (
	"github.com/amanbolat/zederr/pkg/core"
	pbzederrv1 "github.com/amanbolat/zederr/pkg/proto/v1"
	"github.com/amanbolat/zederr/pkg/stderr"
	"google.golang.org/grpc/codes"
)

func Decode(pbErr *pbzederrv1.Error) (core.Error, error) {
	return decode(pbErr)
}

func decode(pbErr *pbzederrv1.Error) (core.Error, error) {
	zedErr := stderr.NewError(
		pbErr.Code,
		pbErr.Domain,
		pbErr.Namespace,
		int(pbErr.HttpCode),
		codes.Code(pbErr.GrpcCode),
		pbErr.InternalMessage,
		pbErr.PublicMessage,
		pbErr.Arguments.AsMap(),
	)

	if len(pbErr.Causes) == 0 {
		return zedErr, nil
	}

	for _, c := range pbErr.Causes {
		cause, err := decode(c)
		if err != nil {
			return nil, err
		}

		zedErr = zedErr.WithCause(cause)
	}

	return zedErr, nil
}
