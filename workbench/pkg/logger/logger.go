package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace" // OpenTelemetryのライブラリを今のうちからimportしておく
)

// NewLogger は設定に基づいてslog.Loggerを初期化します。
func NewLogger(config *viper.Viper) (*slog.Logger, error) {
	// ログレベルを設定
	level := slog.LevelInfo
	if err := level.UnmarshalText([]byte(config.GetString("log.level"))); err != nil {
		return nil, err
	}

	isDevelopment := config.GetString("env") == "development" // e.g. from config
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: isDevelopment,
	}

	// ハンドラの設定
	var handler slog.Handler
	logFormat := config.GetString("log.format")
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// OtelHandlerでラップしてトレース情報を自動追加
	otelHandler := &OtelHandler{Next: handler}

	// ロガーを作成して返す
	logger := slog.New(otelHandler)
	return logger, nil
}

// OtelHandler はログレコードにトレースIDとスパンIDを自動で追加するslog.Handlerです。
type OtelHandler struct {
	Next slog.Handler
}

func (h *OtelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Next.Enabled(ctx, level)
}

func (h *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &OtelHandler{Next: h.Next.WithAttrs(attrs)}
}

func (h *OtelHandler) WithGroup(name string) slog.Handler {
	return &OtelHandler{Next: h.Next.WithGroup(name)}
}

// Handleメソッドで、contextからトレース情報を取得しログ属性に追加します。
func (h *OtelHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := trace.SpanContextFromContext(ctx)

	// SpanContextが有効な場合のみ、trace_idとspan_idをログに追加
	if spanCtx.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return h.Next.Handle(ctx, r)
}
