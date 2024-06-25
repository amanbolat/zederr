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

const (
	defaultStatusCode    = codes.Unknown
	defaultStatusMessage = "unknown error"
)

type Encoder interface {
	Encode(err error) *status.Status
}

type SimpleEncoder struct {
	statusCode    codes.Code
	statusMessage string
}

func NewSimpleEncoder(statusCode codes.Code, statusMessage string) SimpleEncoder {
	return SimpleEncoder{
		statusCode:    statusCode,
		statusMessage: statusMessage,
	}
}

func (e SimpleEncoder) Encode(err error) *status.Status {
	statusErr, ok := status.FromError(err)
	if ok {
		return statusErr
	}

	var zedErr *zeerr.Error
	if !errors.As(err, &zedErr) {
		return status.New(e.statusCode, e.statusMessage)
	}

	sts := status.New(zedErr.GRPCCode(), zedErr.Message())

	return sts
}

type FullEncoder struct {
	statusCode    codes.Code
	statusMessage string
}

func NewFullEncoder(statusCode codes.Code, statusMessage string) SimpleEncoder {
	return SimpleEncoder{
		statusCode:    statusCode,
		statusMessage: statusMessage,
	}
}

func (e FullEncoder) Encode(err error) *status.Status {
	statusErr, ok := status.FromError(err)
	if ok {
		return statusErr
	}

	var zedErr *zeerr.Error
	if !errors.As(err, &zedErr) {
		return status.New(e.statusCode, e.statusMessage)
	}

	pbErr := e.encode(zedErr)

	sts := status.New(zedErr.GRPCCode(), zedErr.Message())
	sts, err = sts.WithDetails(pbErr)
	if err != nil {
		panic(fmt.Errorf("failed to attach details to status: %w", err))
	}

	return sts
}

func (e FullEncoder) encode(zedErr *zeerr.Error) *pbzederrv1.Error {
	pbArgs, err := structpb.NewStruct(zedErr.Arguments())
	if err != nil {
		panic(fmt.Errorf("failed to convert args to structpb: %w", err))
	}

	pbErr := &pbzederrv1.Error{
		Id:        zedErr.ID(),
		GrpcCode:  int32(zedErr.GRPCCode()),
		HttpCode:  int32(zedErr.HTTPCode()),
		Message:   zedErr.Message(),
		Arguments: pbArgs,
		Causes:    nil,
	}

	if len(zedErr.Causes()) == 0 {
		return pbErr
	}

	causes := make([]*pbzederrv1.Error, 0, len(zedErr.Causes()))

	for _, cause := range zedErr.Causes() {
		encodedCause := e.encode(cause)
		causes = append(causes, encodedCause)
	}

	pbErr.Causes = causes

	return pbErr
}
