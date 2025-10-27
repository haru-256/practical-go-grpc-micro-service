package server

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

// ServerLogger はConnect RPCリクエストのロギングを管理する構造体です。
//
// リクエストの開始・終了、エラー、レスポンス情報などを構造化ログとして出力します。
type ServerLogger struct {
	logger *slog.Logger
}

// NewServerLogger はServerLoggerの新しいインスタンスを作成します。
//
// Parameters:
//   - logger: 構造化ロギング用のslogロガー
//
// Returns:
//   - *ServerLogger: 初期化されたServerLoggerインスタンス
func NewServerLogger(logger *slog.Logger) *ServerLogger {
	return &ServerLogger{
		logger: logger,
	}
}

// logStart はリクエストの開始時にログを出力します。
//
// Parameters:
//   - ctx: リクエストのコンテキスト
//   - req: Connect RPCリクエスト
func (l *ServerLogger) logStart(ctx context.Context, req connect.AnyRequest) {
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
func (l *ServerLogger) logEnd(ctx context.Context, res connect.AnyResponse, startTime time.Time, err error) {
	duration := time.Since(startTime)
	logFields := []slog.Attr{
		slog.String("code", connect.CodeOf(err).String()),
		slog.Duration("duration", duration),
	}

	if res != nil {
		// resが非nilでもr.Any()がnilの場合があるためチェック。おそらく「interface が非nilだが、内部の値がnil」という状況
		if r, ok := res.(*connect.Response[any]); ok && r != nil {
			logFields = append(logFields, slog.Any("Msg", r.Any()))
		}
	}

	if err != nil {
		logFields = append(logFields, slog.String("error", err.Error()))
		l.logger.LogAttrs(ctx, slog.LevelError, "finished with error", logFields...)
	} else {
		l.logger.LogAttrs(ctx, slog.LevelInfo, "finished", logFields...)
	}
}

// NewUnaryInterceptor はリクエストのロギングを行うUnaryインターセプターを返します。
//
// Returns:
//   - connect.UnaryInterceptorFunc: ログ出力機能を持つインターセプター関数
func (l *ServerLogger) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
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
