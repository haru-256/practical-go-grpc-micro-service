package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

// Validator はProtovalidateを使用したConnect RPCバリデーションを提供します。
type Validator struct {
	logger    *slog.Logger
	validator protovalidate.Validator
}

// NewValidator はValidatorの新しいインスタンスを生成します。
func NewValidator(logger *slog.Logger) (*Validator, error) {
	v, err := protovalidate.New()
	return &Validator{
		logger:    logger,
		validator: v,
	}, err
}

// NewUnaryInterceptor はリクエストのバリデーションを行うUnaryインターセプターを返します。
func (v *Validator) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			msg, ok := req.Any().(proto.Message)
			if !ok {
				v.logger.ErrorContext(ctx, "request type is not proto.Message")
				return nil, connect.NewError(connect.CodeInternal, errors.New("request type is not proto.Message"))
			}

			if err := v.validator.Validate(msg); err != nil {
				v.logger.InfoContext(ctx, "request validation failed", "error", err)
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}

			return next(ctx, req)
		}
	}
}
