package encode

import (
	"fmt"

	"github.com/amanbolat/zederr/pkg/core"
	pbzederrv1 "github.com/amanbolat/zederr/pkg/proto/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func Encode(zedErr core.Error) (*status.Status, error) {
	pbErr, err := encode(zedErr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode zed error: %w", err)
	}

	sts := status.New(zedErr.GRPCCode(), zedErr.UID())
	sts, err = sts.WithDetails(pbErr)
	if err != nil {
		return nil, fmt.Errorf("failed to attach details to status: %w", err)
	}

	return sts, nil
}

func encode(zedErr core.Error) (*pbzederrv1.Error, error) {
	pbArgs, err := structpb.NewStruct(zedErr.Args())
	if zedErr != nil {
		return nil, fmt.Errorf("failed to convert args to structpb: %w", err)
	}

	pbErr := &pbzederrv1.Error{
		Uid:             zedErr.UID(),
		Domain:          zedErr.Domain(),
		Namespace:       zedErr.Namespace(),
		Code:            zedErr.Code(),
		HttpCode:        int64(zedErr.HTTPCode()),
		GrpcCode:        uint64(zedErr.GRPCCode()),
		PublicMessage:   zedErr.PublicMsg(),
		InternalMessage: zedErr.InternalMsg(),
		Arguments:       pbArgs,
		Causes:          nil,
	}

	if len(zedErr.Causes()) == 0 {
		return pbErr, nil
	}

	causes := make([]*pbzederrv1.Error, 0, len(zedErr.Causes()))
	for _, cause := range zedErr.Causes() {
		enc, err := encode(cause)
		if err != nil {
			return nil, fmt.Errorf("failed to encode cause: %w", err)
		}
		causes = append(causes, enc)
	}

	pbErr.Causes = causes

	return pbErr, nil
}
