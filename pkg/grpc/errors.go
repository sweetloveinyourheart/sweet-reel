package grpc

import (
	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgumentErrorWithField(violations ...*errdetails.BadRequest_FieldViolation) error {
	statusInvalid := connect.NewError(connect.CodeInvalidArgument, errors.New("invalid parameters"))
	for _, violation := range violations {
		if detail, detailErr := connect.NewErrorDetail(violation); detailErr == nil {
			statusInvalid.AddDetail(detail)
		}
	}

	return statusInvalid
}

func InvalidArgumentError(err error) error {
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		return err
	}
	return connect.NewError(connect.CodeInvalidArgument, err)
}

func InternalError(err error) error {
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		return err
	}
	return connect.NewError(connect.CodeInternal, err)
}

func AlreadyExistsError(err error) error {
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		return err
	}
	return connect.NewError(connect.CodeAlreadyExists, err)
}

func NotFoundError(err error) error {
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		return err
	}
	return connect.NewError(connect.CodeNotFound, errors.Newf("not found: %s", err))
}

func UnauthenticatedError(err error) error {
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		return err
	}
	return connect.NewError(connect.CodeUnauthenticated, errors.Newf("unauthenticated: %s", err))
}

func PreconditionFailure(typ string, subject string, desc string) *errdetails.PreconditionFailure_Violation {
	return &errdetails.PreconditionFailure_Violation{
		Type:        typ,
		Subject:     subject,
		Description: desc,
	}
}

func PreconditionError(violations ...*errdetails.PreconditionFailure_Violation) error {
	statusFailed := connect.NewError(connect.CodeFailedPrecondition, errors.New("invalid request"))

	for _, violation := range violations {
		if detail, detailErr := connect.NewErrorDetail(violation); detailErr == nil {
			statusFailed.AddDetail(detail)
		}
	}

	return statusFailed
}
