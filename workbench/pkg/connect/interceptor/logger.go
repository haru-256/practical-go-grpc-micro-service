package interceptor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
//   - *ReqRespLogger: 初期化されたReqRespLoggerインスタンス。connect.Interceptorインターフェースを実装します。
func NewReqRespLogger(logger *slog.Logger) *ReqRespLogger {
	return &ReqRespLogger{
		logger: logger,
	}
}

// logUnaryStart はUnaryリクエストの開始時にログを出力します。
//
// Parameters:
//   - ctx: リクエストのコンテキスト
//   - req: Connect RPCリクエスト
func (l *ReqRespLogger) logUnaryStart(ctx context.Context, req connect.AnyRequest) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"start processing request",
		slog.String("procedure", req.Spec().Procedure),
		slog.String("user-agent", req.Header().Get("User-Agent")),
		slog.String("http_method", req.HTTPMethod()),
		slog.String("addr", req.Peer().Addr),
	)
	// リクエストメッセージをDEBUGレベルで出力
	l.logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		"request message",
		slog.String("procedure", req.Spec().Procedure),
		slog.Any("requestMsg", formatLogValue(req.Any())),
	)
}

// logUnaryEnd はUnaryリクエストの終了時にログを出力します。
//
// Parameters:
//   - ctx: リクエストのコンテキスト
//   - req: Connect RPCリクエスト (procedureを取得するため)
//   - res: Connect RPCレスポンス
//   - startTime: リクエスト開始時刻
//   - err: 処理中に発生したエラー (正常終了時はnil)
func (l *ReqRespLogger) logUnaryEnd(ctx context.Context, req connect.AnyRequest, res connect.AnyResponse, startTime time.Time, err error) {
	duration := time.Since(startTime)
	logFields := []slog.Attr{
		slog.String("procedure", req.Spec().Procedure),
		slog.String("code", connect.CodeOf(err).String()),
		slog.Duration("duration", duration),
	}

	if err != nil {
		logFields = append(logFields, slog.String("error", err.Error()))
		l.logger.LogAttrs(ctx, slog.LevelError, "finished", logFields...)
	} else {
		l.logger.LogAttrs(ctx, slog.LevelInfo, "finished", logFields...)
	}

	// レスポンスメッセージをDEBUGレベルで出力
	if !isNilResponse(res) {
		if msg := res.Any(); msg != nil {
			l.logger.LogAttrs(
				ctx,
				slog.LevelDebug,
				"response message",
				slog.String("procedure", req.Spec().Procedure),
				slog.Any("responseMsg", formatLogValue(msg)),
			)
		}
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

// formatLogValue はメッセージをログ出力用に整形します。
// Protocol Buffers メッセージの場合は protojson.Marshal で JSON 文字列化し、
// json.RawMessage として返すことで、ログ内でJSON構造を保持します。
// API_OPAQUE モードでも適切にフィールドが出力されます。
func formatLogValue(msg any) any {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return msg
	}

	marshaler := protojson.MarshalOptions{
		Multiline:       false,
		EmitUnpopulated: false,
	}
	jsonBytes, err := marshaler.Marshal(protoMsg)
	if err == nil {
		return json.RawMessage(jsonBytes)
	}

	// Marshal失敗時はfmt.Sprintfにフォールバック
	return fmt.Sprintf("%v", protoMsg)
}

// NewUnaryInterceptor はリクエストのロギングを行うUnaryインターセプターを返します。
//
// Returns:
//   - connect.UnaryInterceptorFunc: ログ出力機能を持つインターセプター関数
func (l *ReqRespLogger) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startTime := time.Now()

			l.logUnaryStart(ctx, req)
			res, err := next(ctx, req)
			l.logUnaryEnd(ctx, req, res, startTime, err)

			return res, err
		}
	}
}

// WrapUnary はUnary RPCリクエストをラップしてログを記録します。
//
// このメソッドは connect.Interceptor インターフェースを実装するために使用されます。
// リクエストの開始時、終了時にログを出力し、処理時間とエラー情報を記録します。
//
// Parameters:
//   - next: 次に実行されるUnary RPC関数
//
// Returns:
//   - connect.UnaryFunc: ログ機能を持つラップされたUnary関数
func (l *ReqRespLogger) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		startTime := time.Now()

		l.logUnaryStart(ctx, req)
		res, err := next(ctx, req)
		l.logUnaryEnd(ctx, req, res, startTime, err)

		return res, err
	}
}

// WrapStreamingClient はストリーミングクライアントをラップしてメッセージ送受信をログに記録します。
//
// このメソッドは connect.Interceptor インターフェースを実装するために使用されます。
// 返されるコネクションは loggingClientConn でラップされ、以下のタイミングでログを出力します:
//   - 最初のメッセージ送信時: ストリーム開始のINFOログ
//   - 各メッセージ送信時: メッセージ内容のDEBUGログ
//   - 最初のメッセージ受信時: 受信開始のINFOログ
//   - 各メッセージ受信時: メッセージ内容のDEBUGログ
//   - ストリーム終了時: 全体の処理時間とステータスコードのINFOログ
//
// Parameters:
//   - next: 次に実行されるストリーミングクライアント関数
//
// Returns:
//   - connect.StreamingClientFunc: ログ機能を持つラップされたストリーミングクライアント関数
func (l *ReqRespLogger) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		conn := next(ctx, spec)
		return &loggingClientConn{
			StreamingClientConn: conn,
			logger:              l.logger,
			ctx:                 ctx,
			startTime:           time.Now(),
		}
	}
}

// WrapStreamingHandler はストリーミングハンドラーをラップしてメッセージ送受信をログに記録します。
//
// このメソッドは connect.Interceptor インターフェースを実装するために使用されます。
// コネクションは loggingHandlerConn でラップされ、以下のタイミングでログを出力します:
//   - ストリーム開始時: プロシージャ名、User-Agent、クライアントアドレスのINFOログ
//   - 各メッセージ送信時: メッセージ内容のDEBUGログ
//   - 各メッセージ受信時: メッセージ内容のDEBUGログ（エラー時はERRORログ）
//   - ストリーム終了時: 全体の処理時間、ステータスコード、エラー詳細のログ
//   - 正常終了: INFOレベル
//   - エラー終了: ERRORレベル
//
// Parameters:
//   - next: 次に実行されるストリーミングハンドラー関数
//
// Returns:
//   - connect.StreamingHandlerFunc: ログ機能を持つラップされたストリーミングハンドラー関数
func (l *ReqRespLogger) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		// ラップされたコネクションを作成してメッセージ送受信をインターセプト
		wrapped := &loggingHandlerConn{
			StreamingHandlerConn: conn,
			logger:               l.logger,
			ctx:                  ctx,
			startTime:            time.Now(),
		}

		// ストリーム開始をログ出力
		l.logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"start processing streaming request",
			slog.String("procedure", conn.Spec().Procedure),
			slog.String("user-agent", conn.RequestHeader().Get("User-Agent")),
			slog.String("addr", conn.Peer().Addr),
		)

		// ラップされたコネクションでハンドラーを実行
		err := next(ctx, wrapped)

		// エラーの有無に応じてログレベルを決定（EOFは正常終了）
		logLevel := slog.LevelInfo
		if err != nil && !errors.Is(err, io.EOF) {
			logLevel = slog.LevelError
		}
		// ストリーム終了をログ出力（処理時間とステータスコードを含む）
		l.logger.LogAttrs(
			ctx,
			logLevel,
			"finished streaming request",
			slog.String("procedure", conn.Spec().Procedure),
			slog.String("code", connect.CodeOf(err).String()),
			slog.Duration("duration", time.Since(wrapped.startTime)),
		)
		// EOF以外のエラー詳細を追加ログとして出力
		if err != nil && !errors.Is(err, io.EOF) {
			l.logger.LogAttrs(ctx, slog.LevelError, "error details", slog.String("error", err.Error()))
		}

		return err
	}
}

// loggingClientConn は StreamingClientConn をラップしてメッセージ送受信をログに記録します。
//
// このコネクションは WrapStreamingClient によって返され、クライアント側のストリーミングRPCにおいて
// 各メッセージの送受信を透過的にログに記録します。
//
// フィールド:
//   - StreamingClientConn: 元のストリーミングクライアントコネクション（埋め込み）
//   - logger: ログ出力用のslogロガー
//   - ctx: リクエストのコンテキスト
//   - startTime: ストリーム開始時刻（duration計測用）
//   - sentFirst: 最初のメッセージ送信済みフラグ（重複ログ防止用）
//   - recvFirst: 最初のメッセージ受信済みフラグ（重複ログ防止用）
type loggingClientConn struct {
	connect.StreamingClientConn
	logger    *slog.Logger
	ctx       context.Context
	startTime time.Time
	sentFirst bool
	recvFirst bool
}

// Send はメッセージをサーバーに送信し、その内容をログに記録します。
//
// 最初の送信時には、ストリーム開始のINFOログ（プロシージャ名、User-Agent、サーバーアドレス）を出力します。
// 各メッセージ送信時には、メッセージ内容をDEBUGレベルでログに記録します。
//
// DEBUGログは本番環境では無効化することで、パフォーマンスへの影響を最小限に抑えられます。
//
// Parameters:
//   - msg: 送信するメッセージ（Protocol Buffers メッセージなど）
//
// Returns:
//   - error: 送信中にエラーが発生した場合のエラー
func (c *loggingClientConn) Send(msg any) error {
	// 最初の送信時のみストリーム開始をログ出力
	if !c.sentFirst {
		c.logger.LogAttrs(
			c.ctx,
			slog.LevelInfo,
			"start streaming client request",
			slog.String("procedure", c.Spec().Procedure),
			slog.String("user-agent", c.RequestHeader().Get("User-Agent")),
			slog.String("addr", c.Peer().Addr),
		)
		c.sentFirst = true
	}

	// 各メッセージの内容をDEBUGレベルで記録
	c.logger.LogAttrs(
		c.ctx,
		slog.LevelDebug,
		"sending message",
		slog.String("procedure", c.Spec().Procedure),
		slog.Any("responseMsg", formatLogValue(msg)),
	)

	return c.StreamingClientConn.Send(msg)
}

// Receive はサーバーからメッセージを受信し、その内容をログに記録します。
//
// 最初の受信時には、レスポンス受信開始のINFOログを出力します。
// 各メッセージ受信時には、メッセージ内容をDEBUGレベルでログに記録します。
// ストリーム終了時（EOFまたはエラー）には、全体の処理時間とステータスコードをINFOレベルで記録します。
//
// EOFは正常なストリーム終了を示すため、ERRORログは出力されません。
// EOF以外のエラーの場合は、ERRORレベルでエラー詳細を記録します。
//
// Parameters:
//   - msg: 受信したメッセージを格納するポインタ
//
// Returns:
//   - error: 受信中にエラーが発生した場合のエラー、ストリーム終了時はio.EOF
func (c *loggingClientConn) Receive(msg any) error {
	err := c.StreamingClientConn.Receive(msg)

	// 最初の受信時のみレスポンス開始をログ出力
	if !c.recvFirst {
		c.logger.LogAttrs(
			c.ctx,
			slog.LevelInfo,
			"received first response",
			slog.String("procedure", c.Spec().Procedure),
		)
		c.recvFirst = true
	}

	// エラー（EOF含む）が発生した場合はストリーム終了をログ出力
	if err != nil {
		c.logger.LogAttrs(
			c.ctx,
			slog.LevelInfo,
			"finished streaming client",
			slog.String("procedure", c.Spec().Procedure),
			slog.String("code", connect.CodeOf(err).String()),
			slog.Duration("duration", time.Since(c.startTime)),
		)
		// EOFは正常終了なので、それ以外のエラーのみERRORログを出力
		if !errors.Is(err, io.EOF) {
			c.logger.LogAttrs(c.ctx, slog.LevelError, "receive error", slog.String("error", err.Error()))
		}
	} else {
		// 正常に受信できた場合はメッセージ内容をDEBUGレベルで記録
		c.logger.LogAttrs(
			c.ctx,
			slog.LevelDebug,
			"received message",
			slog.String("procedure", c.Spec().Procedure),
			slog.Any("responseMsg", formatLogValue(msg)),
		)
	}

	return err
}

// loggingHandlerConn は StreamingHandlerConn をラップしてメッセージ送受信をログに記録します。
//
// このコネクションは WrapStreamingHandler によって使用され、サーバー側のストリーミングRPCにおいて
// 各メッセージの送受信を透過的にログに記録します。
//
// フィールド:
//   - StreamingHandlerConn: 元のストリーミングハンドラーコネクション（埋め込み）
//   - logger: ログ出力用のslogロガー
//   - ctx: リクエストのコンテキスト
//   - startTime: ストリーム開始時刻（duration計測用）
type loggingHandlerConn struct {
	connect.StreamingHandlerConn
	logger    *slog.Logger
	ctx       context.Context
	startTime time.Time
}

// Send はメッセージをクライアントに送信し、その内容をログに記録します。
//
// 各メッセージ送信時に、メッセージ内容をDEBUGレベルでログに記録します。
// DEBUGログは本番環境では無効化することで、パフォーマンスへの影響を最小限に抑えられます。
//
// Note: Protocol Buffers メッセージは String() メソッドで文字列化して記録されます。
//
// Parameters:
//   - msg: 送信するメッセージ（Protocol Buffers メッセージなど）
//
// Returns:
//   - error: 送信中にエラーが発生した場合のエラー
func (c *loggingHandlerConn) Send(msg any) error {
	// メッセージ内容をDEBUGレベルで記録
	c.logger.LogAttrs(
		c.ctx,
		slog.LevelDebug,
		"sending message",
		slog.String("procedure", c.Spec().Procedure),
		slog.Any("responseMsg", formatLogValue(msg)),
	)

	return c.StreamingHandlerConn.Send(msg)
}

// Receive はクライアントからメッセージを受信し、その内容をログに記録します。
//
// 各メッセージ受信時に、メッセージ内容をDEBUGレベルでログに記録します。
// EOFは正常なストリーム終了を示すため、ログには記録されません。
// EOF以外のエラーの場合は、ERRORレベルでエラー詳細を記録します。
//
// Parameters:
//   - msg: 受信したメッセージを格納するポインタ
//
// Returns:
//   - error: 受信中にエラーが発生した場合のエラー、ストリーム終了時はio.EOF
func (c *loggingHandlerConn) Receive(msg any) error {
	err := c.StreamingHandlerConn.Receive(msg)

	// EOF以外のエラーが発生した場合はERRORログを出力
	if err != nil && !errors.Is(err, io.EOF) {
		c.logger.LogAttrs(
			c.ctx,
			slog.LevelError,
			"receive error",
			slog.String("procedure", c.Spec().Procedure),
			slog.String("error", err.Error()),
		)
	} else {
		// 正常に受信できた場合はメッセージ内容をDEBUGレベルで記録
		c.logger.LogAttrs(
			c.ctx,
			slog.LevelDebug,
			"received message",
			slog.String("procedure", c.Spec().Procedure),
			slog.Any("responseMsg", formatLogValue(msg)),
		)
	}

	return err
}

var _ connect.Interceptor = (*ReqRespLogger)(nil)
