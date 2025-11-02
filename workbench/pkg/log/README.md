# Logger パッケージ

このパッケージは、プロジェクト全体で使用される構造化ログ（structured logging）の実装を提供します。

## 概要

`log/slog` パッケージをベースにした、コンテキスト対応の構造化ログ機能を提供します。OpenTelemetryとの統合により、トレースIDとスパンIDを自動的にログに含めることができます。

## 機能

- **構造化ログ**: キーと値のペアでログを記録
- **コンテキスト対応**: `context.Context`からトレース情報を自動抽出
- **設定可能なレベル**: debug/info/warn/errorの4段階
- **複数のフォーマット**: text（開発用）とjson（本番用）
- **OpenTelemetry統合**: トレースIDとスパンIDの自動記録

## 使用方法

### 初期化

Viperの設定を使用してロガーを初期化します：

```go
import (
    "github.com/spf13/viper"
    "yourproject/pkg/logger"
)

// Viperの設定からロガーを作成
log, err := log.NewLogger(config)
if err != nil {
    // エラーハンドリング
}
```

### 設定

設定ファイル（config.toml）で以下の値を指定：

```toml
[log]
level = "info"    # debug, info, warn, error
format = "text"   # text, json
```

環境変数での上書きも可能：

```bash
export LOG_LEVEL=debug
export LOG_FORMAT=json
```

### ログ出力

#### 基本的な使い方

```go
// 情報ログ
logger.InfoContext(ctx, "商品を作成しました")

// エラーログ
logger.ErrorContext(ctx, "商品の作成に失敗しました")

// 警告ログ
logger.WarnContext(ctx, "商品名が重複しています")

// デバッグログ
logger.DebugContext(ctx, "デバッグ情報")
```

#### 構造化データの追加

`slog.Attr`を使用して、キーと値のペアを追加：

```go
logger.InfoContext(ctx, "商品を作成しました",
    slog.String("product_id", "123e4567-e89b-12d3-a456-426614174000"),
    slog.String("product_name", "テスト商品"),
    slog.Int("price", 1000))
```

#### エラー情報の記録

```go
if err != nil {
    logger.ErrorContext(ctx, "データベースエラー",
        slog.String("operation", "create_product"),
        slog.Any("error", err))
}
```

#### グループ化されたログ

関連する情報をグループ化：

```go
logger.InfoContext(ctx, "トランザクション完了",
    slog.Group("product",
        slog.String("id", productID),
        slog.String("name", productName)),
    slog.Group("category",
        slog.String("id", categoryID),
        slog.String("name", categoryName)))
```

## ログレベル

### debug

開発環境で詳細な動作を確認するために使用：

```go
logger.DebugContext(ctx, "リクエスト詳細",
    slog.Any("request", req))
```

### info

通常の操作フローを記録（本番環境の推奨デフォルト）：

```go
logger.InfoContext(ctx, "商品を作成しました",
    slog.String("product_id", id))
```

### warn

潜在的な問題や注意が必要な状況を記録：

```go
logger.WarnContext(ctx, "接続プールが枯渇しています",
    slog.Int("active_connections", count))
```

### error

エラーが発生した場合に記録：

```go
logger.ErrorContext(ctx, "商品の作成に失敗しました",
    slog.Any("error", err))
```

## ログフォーマット

### text形式

人間が読みやすいテキスト形式（開発環境推奨）：

```
2024-01-15 10:30:45 INFO 商品を作成しました product_id=123e4567-e89b-12d3-a456-426614174000 product_name=テスト商品
```

### json形式

ログ集約システムでの解析に適したJSON形式（本番環境推奨）：

```json
{
  "time": "2024-01-15T10:30:45.123456789+09:00",
  "level": "INFO",
  "msg": "商品を作成しました",
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "product_name": "テスト商品",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7"
}
```

## OpenTelemetry統合

OtelHandlerを使用することで、コンテキストからトレース情報を自動的に抽出してログに追加します：

```go
// トレース情報が自動的に含まれる
logger.InfoContext(ctx, "処理を開始しました")

// 出力例（JSON）:
// {
//   "msg": "処理を開始しました",
//   "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
//   "span_id": "00f067aa0ba902b7"
// }
```

## ベストプラクティス

### 1. 常にコンテキストを渡す

`*Context`メソッドを使用してトレース情報を含める：

```go
// Good
logger.InfoContext(ctx, "メッセージ")

// Avoid
logger.Info("メッセージ")  // トレース情報が含まれない
```

### 2. 構造化データを使用

文字列の連結ではなく、構造化データを使用：

```go
// Good
logger.InfoContext(ctx, "商品を作成",
    slog.String("product_id", id),
    slog.Int("price", price))

// Avoid
logger.InfoContext(ctx, fmt.Sprintf("商品を作成: id=%s, price=%d", id, price))
```

### 3. 適切なログレベルを選択

- **debug**: 開発時のみ必要な詳細情報
- **info**: 通常の動作フロー
- **warn**: 潜在的な問題
- **error**: 実際のエラー

### 4. 機密情報を含めない

パスワードやトークンなどの機密情報をログに出力しない：

```go
// Bad
logger.InfoContext(ctx, "ログイン",
    slog.String("password", password))  // 絶対にやらない！

// Good
logger.InfoContext(ctx, "ログイン",
    slog.String("user_id", userID))
```

### 5. エラーは`slog.Any`で記録

エラーオブジェクトは`slog.Any`を使用：

```go
logger.ErrorContext(ctx, "処理失敗",
    slog.Any("error", err))
```

## 依存性注入との統合

Uber Fxを使用した依存性注入の例：

```go
// モジュール定義
var Module = fx.Module("logger",
    fx.Provide(NewLogger),
)

// 使用例
type ProductService struct {
    logger *slog.Logger
    // ...
}

func NewProductService(logger *slog.Logger, ...) *ProductService {
    return &ProductService{
        logger: logger,
        // ...
    }
}
```

## トラブルシューティング

### ログが出力されない

1. ログレベルを確認：`LOG_LEVEL=debug`に設定
2. ログフォーマットを確認：`LOG_FORMAT=text`で読みやすく

### 無効なログレベルエラー

設定で無効なログレベルが指定された場合、自動的に`info`レベルにフォールバックします：

```go
### 無効なログレベルエラー

設定で無効なログレベルが指定された場合、自動的に`info`レベルにフォールバックします：

```go
// 設定: level = "invalid"
// 結果: Infoレベルでロガーが初期化される
```

```

### トレースIDが含まれない

1. コンテキストを正しく渡しているか確認
2. OpenTelemetryが正しく初期化されているか確認
3. `*Context`メソッドを使用しているか確認

## 参考リンク

- [log/slog公式ドキュメント](https://pkg.go.dev/log/slog)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Structured Logging Best Practices](https://www.sentinelone.com/blog/the-10-commandments-of-logging/)
