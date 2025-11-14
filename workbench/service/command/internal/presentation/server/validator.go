package server

import (
	"context"
	"errors"
	"log/slog"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

type Validator struct {
	logger    *slog.Logger
	validator protovalidate.Validator
}

func NewValidator(logger *slog.Logger) (*Validator, error) {
	v, err := protovalidate.New()
	return &Validator{
		logger:    logger,
		validator: v,
	}, err
}

func (v *Validator) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			msg, ok := req.Any().(proto.Message)
			if !ok {
				// リクエストが proto.Message ではない場合 (Connectでは通常あり得ない)
				// ここではエラーとして扱う (もしくは単に next に流す)
				v.logger.ErrorContext(ctx, "request type is not proto.Message")
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("request type is not proto.Message"),
				)
			}

			// validate the message
			if err := v.validator.Validate(msg); err != nil {
				// validationエラーはクライアント側の問題なのでInfoログとする
				v.logger.InfoContext(ctx, "request validation failed", "error", err)
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					err,
				)
			}

			return next(ctx, req)
		}
	}
}
