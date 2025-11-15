package interceptor

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	"connectrpc.com/connect"
)

// ReqRespLogger はConnect RPCリクエストのロギングを管理する構造体です。
//
// リクエストの開始・終了、エラー、レスポンス情報などを構造化ログとして出力します。
type ReqRespLogger struct {
	logger *slog.Logger
}

// NewReqRespLogger はReqRespLoggerの新しいインスタンスを作成します。
//
// Parameters:
//   - logger: 構造化ロギング用のslogロガー
//
// Returns:
//   - *ReqRespLogger: 初期化されたReqRespLoggerインスタンス
func NewReqRespLogger(logger *slog.Logger) *ReqRespLogger {
	return &ReqRespLogger{
		logger: logger,
	}
}

// logStart はリクエストの開始時にログを出力します。
//
// Parameters:
//   - ctx: リクエストのコンテキスト
//   - req: Connect RPCリクエスト
func (l *ReqRespLogger) logStart(ctx context.Context, req connect.AnyRequest) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"start processing request",
		slog.String("procedure", req.Spec().Procedure),
		slog.String("user-agent", req.Header().Get("User-Agent")),
		slog.String("http_method", req.HTTPMethod()),
		slog.String("peer", req.Peer().Addr),
		slog.Any("Msg", req.Any()),
	)
}

// logEnd はリクエスト/ストリームの終了時にログを出力します。
//
// Parameters:
//   - ctx: リクエストのコンテキスト
//   - res: Connect RPCレスポンス
//   - startTime: リクエスト開始時刻
//   - err: 処理中に発生したエラー (正常終了時はnil)
func (l *ReqRespLogger) logEnd(ctx context.Context, res connect.AnyResponse, startTime time.Time, err error) {
	duration := time.Since(startTime)
	logFields := []slog.Attr{
		slog.String("code", connect.CodeOf(err).String()),
		slog.Duration("duration", duration),
	}

	if !isNilResponse(res) {
		if msg := res.Any(); msg != nil {
			logFields = append(logFields, slog.Any("Msg", msg))
		}
	}

	if err != nil {
		logFields = append(logFields, slog.String("error", err.Error()))
		l.logger.LogAttrs(ctx, slog.LevelError, "finished with error", logFields...)
	} else {
		l.logger.LogAttrs(ctx, slog.LevelInfo, "finished", logFields...)
	}
}

// isNilResponse は型付きnilを含むレスポンスかどうかを判定します。
func isNilResponse(res connect.AnyResponse) bool {
	if res == nil {
		return true
	}

	v := reflect.ValueOf(res)
	return (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) && v.IsNil()
}

// NewUnaryInterceptor はリクエストのロギングを行うUnaryインターセプターを返します。
//
// Returns:
//   - connect.UnaryInterceptorFunc: ログ出力機能を持つインターセプター関数
func (l *ReqRespLogger) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startTime := time.Now()

			l.logStart(ctx, req)
			res, err := next(ctx, req)
			l.logEnd(ctx, res, startTime, err)

			return res, err
		}
	}
}
