package zegrpc

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/amanbolat/zederr/zeerr"
	pbzederrv1 "github.com/amanbolat/zederr/zeproto/v1"
)

func Encode(err error) *status.Status {
	var zedErr zeerr.Error
	if !errors.As(err, &zedErr) {
		return status.New(codes.Unknown, "unknown error")
	}

	pbErr := encode(zedErr)

	sts := status.New(zedErr.GRPCCode(), zedErr.PublicMsg())
	sts, err = sts.WithDetails(pbErr)
	if err != nil {
		panic(fmt.Errorf("failed to attach details to status: %w", err))
	}

	return sts
}

func encode(zedErr zeerr.Error) *pbzederrv1.Error {
	pbArgs, err := structpb.NewStruct(zedErr.Args())
	if err != nil {
		panic(fmt.Errorf("failed to convert args to structpb: %w", err))
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
		return pbErr
	}

	causes := make([]*pbzederrv1.Error, 0, len(zedErr.Causes()))

	for _, cause := range zedErr.Causes() {
		encodedCause := encode(cause)
		causes = append(causes, encodedCause)
	}

	pbErr.Causes = causes

	return pbErr
}
