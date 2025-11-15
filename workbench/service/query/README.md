# Query Service

CQRS（Command Query Responsibility Segregation）パターンにおけるQuery Serviceの実装です。
このサービスは、データの検索・参照といった読み込み操作を担当します。

## 目次

- [概要](#概要)
- [ディレクトリ構造](#ディレクトリ構造)
- [各層の責務](#各層の責務)
- [アーキテクチャパターン](#アーキテクチャパターン)
- [実装ガイドライン](#実装ガイドライン)
- [開発時の注意点](#開発時の注意点)
- [ビルドとテスト](#ビルドとテスト)
- [参考リンク](#参考リンク)

## 概要

Query Serviceは、CQRSアーキテクチャの一部として、以下の責務を持ちます：

- データの検索、参照操作の処理
- 読み込み専用のデータモデル（Read Model）の提供
- パフォーマンスを重視したデータ取得
- 複雑な検索条件への対応

## ディレクトリ構造

```text
query/
├── cmd/                    # アプリケーションのエントリポイント
│   └── server/            # サーバー起動処理
│       └── main.go        # メイン関数
├── internal/              # 内部パッケージ（外部から直接importできない）
│   ├── domain/           # ドメイン層
│   │   ├── models/       # ドメインモデル定義
│   │   └── repository/   # リポジトリインターフェース
│   ├── infrastructure/   # インフラストラクチャ層
│   │   ├── config/       # 設定管理
│   │   └── db/          # データベースアクセス
│   ├── presentation/     # プレゼンテーション層
│   │   └── server/       # gRPCサーバー実装
│   ├── testhelpers/      # テストヘルパー関数
│   └── mock/             # モック（gomockで生成）
└── README.md             # このファイル
```

## 各層の責務

### cmd/

アプリケーションのエントリポイントを含みます。

- **server/main.go**: アプリケーションの起動処理、Uber Fxによる依存関係の注入、サーバーの設定を行います

### internal/domain/

ドメイン層（Domain Layer）を実装します。Query側は読み込み専用のモデルを定義します。

- **責務**:
    - 読み込み専用のドメインモデル
    - リポジトリインターフェース（検索メソッド）
    - ドメインロジック（必要に応じて）

- **models/**: ドメインモデル定義
    - **Product**: 商品エンティティ（ID、名前、価格、カテゴリ）
    - **Category**: カテゴリエンティティ（ID、名前）

- **repository/**: リポジトリインターフェース
    - **ProductRepository**: 商品検索メソッド（List/FindById/FindByNameLike）
    - **CategoryRepository**: カテゴリ検索メソッド（List/FindById）

#### ドメインモデルの特徴

Query側のドメインモデルは、Command側と以下の点で異なります：

1. **読み込み専用**: 更新メソッドを持たない
2. **シンプルな構造**: 値オブジェクトではなく基本型を使用する場合が多い
3. **パフォーマンス重視**: 必要なデータのみを持つ

```go
// Product モデル例
type Product struct {
 id       string
 name     string
 price    uint32
 category *Category
}

// アクセサメソッドのみ提供（変更メソッドなし）
func (p *Product) Id() string { return p.id }
func (p *Product) Name() string { return p.name }
func (p *Product) Price() uint32 { return p.price }
func (p *Product) Category() *Category { return p.category }
```

### internal/infrastructure/

インフラストラクチャ層（Infrastructure Layer）を実装します。

- **責務**:
    - データベースへのアクセス（リポジトリ実装）
    - 設定管理
    - ロギング

- **config/**: 設定管理
    - **config.go**: Viperを使用した設定ファイルの読み込み

- **db/**: データベースアクセス
    - **database.go**: GORM接続の初期化
    - **repository.go**: ProductRepositoryImpl、CategoryRepositoryImplの実装
    - **module.go**: Uber Fxモジュール定義（インフラ層の依存関係を構成）

### internal/presentation/

プレゼンテーション層（Presentation Layer）を実装します。

- **責務**:
    - gRPC/Connect RPCハンドラーの実装
    - リクエスト/レスポンスの変換
    - 入力値の検証（Protovalidateを使用）
    - サーバーのライフサイクル管理

- **server/**: gRPC/Connect RPCサーバー実装
    - **handler.go**: CategoryServiceとProductServiceのハンドラー実装
        - `CategoryServiceHandlerImpl`: カテゴリ一覧・詳細取得のエンドポイント
        - `ProductServiceHandlerImpl`: 商品一覧・詳細取得・検索のエンドポイント
    - **server.go**: HTTPサーバーとルーティングの設定
    - **共通インターセプター**: `pkg/connect/interceptor/logger.go`（リクエスト/レスポンスロギング）、`pkg/connect/interceptor/validate.go`（Protovalidate検証）
    - **handler_test.go**: ハンドラーのユニットテスト（mockを使用）
    - **handler_integration_test.go**: ハンドラーの統合テスト（実際のDBを使用）

- **module.go**: Uber Fxモジュール定義（プレゼンテーション層の依存関係とライフサイクル管理）

### internal/testhelpers/

テストで使用するヘルパー関数を提供します。

- **責務**:
    - データベースのセットアップとクリーンアップ
    - テスト用のロガーとバリデータ
    - 統合テスト用の共通セットアップ

- **testhelpers.go**: テストヘルパー関数
    - `SetupDB`: テスト用データベース接続の初期化
    - `TeardownDB`: データベース接続のクローズ
    - `TestLogger`: テスト用ロガー（出力破棄）
    - `IntegrationTestSetup`: 統合テスト用セットアップ構造体

## アーキテクチャパターン

このQuery Serviceは以下のアーキテクチャパターンを採用しています：

### クリーンアーキテクチャ

- 依存関係の向きがドメイン層に向かっている
- 外部の詳細（データベース、Webフレームワーク等）に依存しないドメイン層

### CQRS（Command Query Responsibility Segregation）

- コマンド（書き込み）とクエリ（読み込み）の責務を分離
- このサービスはクエリ側の実装を担当
- Command側とは異なるデータモデルを使用可能

### レイヤードアーキテクチャ

Query側は比較的シンプルな構造を持ち、以下の3層で構成されています：

1. **Presentation Layer**: gRPCハンドラー、リクエスト/レスポンス変換
2. **Domain Layer**: 読み込み専用モデル、リポジトリインターフェース
3. **Infrastructure Layer**: データベースアクセス、リポジトリ実装

## 実装ガイドライン

### 1. 依存関係の管理と設計パターン

#### 「インターフェースを受け入れ、具象（struct）を返す」原則

このプロジェクトでは、Goの重要な設計思想に従っています：

**コンストラクタの設計原則:**

```go
// ✅ 推奨: インターフェースを受け取り、具象型を返す
func NewCategoryServiceHandlerImpl(
    logger *slog.Logger,
    repo repository.CategoryRepository,  // インターフェースを受け入れる
) (*CategoryServiceHandlerImpl, error) {  // 具象型を返す
    return &CategoryServiceHandlerImpl{
        logger: logger,
        repo:   repo,
    }, nil
}
```

#### Uber Fxによる依存性注入

依存性注入には**Uber Fx**を使用しています：

```go
// internal/presentation/module.go
var Module = fx.Module(
    "presentation",
    infrastructure.Module,  // インフラ層を含める
    fx.Provide(
        // 具象型を返すコンストラクタをインターフェースに変換
        fx.Annotate(
            server.NewCategoryServiceHandlerImpl,
            fx.As(new(queryconnect.CategoryServiceHandler)),
        ),
        fx.Annotate(
            server.NewProductServiceHandlerImpl,
            fx.As(new(queryconnect.ProductServiceHandler)),
        ),
        server.NewQueryServer,
    ),
    fx.Invoke(registerLifecycleHooks),
)
```

### 2. エラーハンドリング

Query側では主に以下のエラーを扱います：

#### CRUDエラー (`errs.CRUDError`)

リポジトリ層で発生する、データアクセスに関するエラーです。

**エラーコード:**

- **NOT_FOUND**: リソースが見つからない

**使用例:**

```go
// リソースが見つからない場合
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, errs.NewCRUDError("NOT_FOUND", 
        fmt.Sprintf("商品番号: %s は存在しません", id))
}
```

#### 内部エラー (`errs.InternalError`)

予期しないエラーや内部処理のエラーです。

**使用例:**

```go
if err != nil {
    return nil, errs.NewInternalError("商品一覧の取得に失敗しました", err)
}
```

#### gRPCエラーへの変換

プレゼンテーション層では、ドメインエラーを適切なgRPCエラーコードに変換します：

```go
func handleError(err error, operation string) error {
    var crudErr *errs.CRUDError
    var internalErr *errs.InternalError

    if errors.As(err, &crudErr) {
        switch crudErr.Code {
        case "NOT_FOUND":
            return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s: %w", operation, err))
        default:
            return connect.NewError(connect.CodeInternal, fmt.Errorf("%s: %w", operation, err))
        }
    }

    if errors.As(err, &internalErr) {
        return connect.NewError(connect.CodeInternal, fmt.Errorf("%s: %w", operation, err))
    }

    return connect.NewError(connect.CodeInternal, fmt.Errorf("%s: %w", operation, err))
}
```

### 3. トランザクション管理

Query側は読み込み専用のため、トランザクション管理は不要です。
データベースへの接続は読み込み専用モードで行います。

### 4. テスト戦略

各層ごとにユニットテストと統合テストを実装し、包括的なテストカバレッジを確保しています。

#### テストツール

- **testify**: アサーションライブラリ
- **gomock**: モック生成ツール

#### ドメイン層のテスト

ドメイン層（リポジトリ実装）のテストは統合テストとして実装されています：

```bash
# リポジトリの統合テスト（実際のDBを使用）
go test -tags=integration ./internal/infrastructure/db/...
```

#### プレゼンテーション層のテスト

プレゼンテーション層（ハンドラー）には、ユニットテストと統合テストがあります。

**ユニットテスト（リポジトリをモック化）:**

```bash
# ハンドラーのユニットテスト（リポジトリをモック化）
go test ./internal/presentation/server/... -run TestCategoryServiceHandler
go test ./internal/presentation/server/... -run TestProductServiceHandler
```

リポジトリをモック化し、ハンドラーのロジック（バリデーション、エラーハンドリング、レスポンス生成）を検証します。

**統合テスト（実際のDB使用）:**

```bash
# ハンドラーの統合テスト（実際のDBを使用）
go test -tags=integration ./internal/presentation/server/...
```

実際のMySQLデータベースを使用して、gRPCハンドラーの動作を検証します。

#### moduleテスト

FXモジュールの依存関係と初期化を検証：

```bash
# モジュールの統合テスト
go test -tags=integration ./internal/presentation/... -run TestModule
```

#### すべてのテストを実行

```bash
# Makefileを使用（推奨）
make test

# 統合テストを含むすべてのテスト
go test -tags=integration ./...

# ユニットテストのみ（統合テストを除外）
go test ./...
```

#### テストヘルパーの活用

`internal/testhelpers` パッケージには、テストで使用できる便利な関数が用意されています：

```go
// データベースのセットアップ
dbConn, err := testhelpers.SetupDB("../../", "config")

// データベースのクリーンアップ
defer testhelpers.TeardownDB(dbConn)

// テスト用ロガー（出力破棄）
logger := testhelpers.TestLogger

// 統合テスト用セットアップ
setup := &testhelpers.IntegrationTestSetup{
    DBConn: dbConn,
    Logger: logger,
}
```

### 5. API仕様

#### CategoryService

| メソッド | リクエスト | レスポンス | 説明 |
|---------|----------|----------|------|
| ListCategories | ListCategoriesRequest | ListCategoriesResponse | カテゴリ一覧を取得 |
| GetCategoryById | GetCategoryByIdRequest | GetCategoryByIdResponse | カテゴリIDでカテゴリを取得 |

#### ProductService

| メソッド | リクエスト | レスポンス | 説明 |
|---------|----------|----------|------|
| ListProducts | ListProductsRequest | ListProductsResponse | 商品一覧を取得 |
| GetProductById | GetProductByIdRequest | GetProductByIdResponse | 商品IDで商品を取得 |
| SearchProductsByKeyword | SearchProductsByKeywordRequest | SearchProductsByKeywordResponse | 商品名のキーワードで商品を検索 |

## ロギング

このサービスは構造化ログ（structured logging）として`log/slog`を使用しています。

### ロガーの使用方法

すべての層で、依存性注入されたロガーを使用してログを出力します：

```go
// エラーログ
logger.ErrorContext(ctx, "Failed to list categories", "error", err)

// 情報ログ
logger.InfoContext(ctx, "Request validation failed", "error", err)
```

### ログレベル

設定ファイルまたは環境変数で変更可能：

- `debug`: デバッグ情報（開発環境）
- `info`: 一般的な情報（本番環境推奨）
- `warn`: 警告
- `error`: エラー

### ログフォーマット

- `text`: 人間が読みやすいテキスト形式（開発環境）
- `json`: JSON形式（本番環境推奨、ログ集約システムに適している）

## 設定管理

このサービスは、設定管理に[Viper](https://github.com/spf13/viper)を使用し、TOMLファイルと環境変数の両方をサポートしています。

### 設定ファイル

デフォルトの設定ファイルは`config.toml`です：

```toml
[log]
level = "info"
format = "text"

[db]
dbname = "query_db"
host = "localhost"
port = 3307
user = "query_user"
password = "query_pass"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = "1h"
conn_max_idle_time = "10m"
```

### 環境変数による上書き

環境変数を使用して設定を上書きできます：

**ログ設定:**

- `LOG_LEVEL`: ログレベル（debug/info/warn/error）
- `LOG_FORMAT`: ログフォーマット（text/json）

**データベース設定:**

- `DB_DBNAME`: データベース名
- `DB_HOST`: ホスト名
- `DB_PORT`: ポート番号
- `DB_USER`: ユーザー名
- `DB_PASSWORD`: パスワード
- `DB_MAX_IDLE_CONNS`: 最大アイドル接続数
- `DB_MAX_OPEN_CONNS`: 最大オープン接続数
- `DB_CONN_MAX_LIFETIME`: 接続最大ライフタイム
- `DB_CONN_MAX_IDLE_TIME`: アイドル接続最大ライフタイム

**使用例:**

```bash
# 本番環境でJSON形式のログを使用
export LOG_LEVEL=info
export LOG_FORMAT=json

# データベース接続設定
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secure_password
```

## 開発時の注意点

1. **Query側の特徴を理解する**
   - 読み込み専用の操作のみを担当
   - パフォーマンスを重視（キャッシュ、非正規化など）
   - Command側とは異なるデータモデルを使用可能

2. **シンプルな設計を心がける**
   - Command側に比べてシンプルな構造
   - 複雑なビジネスロジックは不要
   - データ取得に集中

3. **パフォーマンスを重視**
   - 必要なデータのみを取得
   - 適切なインデックスの使用
   - N+1問題に注意（JOINやEager Loadingの活用）

4. **読み込み専用を徹底**
   - データの更新操作は一切行わない
   - トランザクションは不要
   - 読み込み専用接続の使用

5. **依存関係の向きを守る**
   - 依存関係は常にドメイン層に向かう（依存性逆転の原則）
   - リポジトリはインターフェースのみドメイン層で定義、実装はインフラ層

6. **「インターフェースを受け入れ、具象を返す」原則の遵守**
   - コンストラクタは必ず具象型（`*XxxImpl`）を返す
   - コンストラクタ名に`Impl`サフィックスを付ける
   - Fxモジュールで`fx.Annotate`と`fx.As`を使用してインターフェースに変換

## ビルドとテスト

### ビルド

```bash
# プロジェクトルートから
make build

# または直接ビルド
go build -o bin/query-service ./cmd/server
```

### テスト実行

```bash
# すべてのテストを実行
make test

# 統合テストを含むすべてのテスト
go test -tags=integration ./...

# ユニットテストのみ
go test ./...

# 特定の層のテスト
go test ./internal/domain/...           # ドメイン層
go test -tags=integration ./internal/infrastructure/...  # インフラ層（統合テスト）
go test ./internal/presentation/...     # プレゼンテーション層

# カバレッジレポート付きで実行
go test -cover ./...
```

### データベースセットアップ

統合テストを実行する前に、MySQLデータベースを起動してください：

```bash
# データベースの起動（docker-compose使用）
cd ../../db
make up

# データベースの停止
make down

# データベースのクリーンアップ
make clean
```

### Lint実行

```bash
make lint
```

## 参考リンク

### 設計パターンとアーキテクチャ

- [CQRS Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs) - Microsoft Azure アーキテクチャセンター
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Robert C. Martin

### Go言語とベストプラクティス

- [Effective Go](https://golang.org/doc/effective_go) - Go言語公式ガイド
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Goコードレビューのベストプラクティス

### 依存性注入

- [Uber Fx Documentation](https://uber-go.github.io/fx/) - Uber Fx公式ドキュメント
- [Fx Modules Guide](https://uber-go.github.io/fx/modules.html) - モジュール設計パターン

### テスティング

- [testify](https://github.com/stretchr/testify) - アサーションライブラリ
- [gomock](https://github.com/golang/mock) - モック生成ツール

### プロジェクト内部リソース

- [プロジェクトスタイルガイド](../../../.gemini/styleguide.md) - コーディング規約と設計原則
- [Command Service README](../command/README.md) - Command Service（書き込み側）のドキュメント
