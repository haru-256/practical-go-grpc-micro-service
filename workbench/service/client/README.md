# Client Service

CQRS（Command Query Responsibility Segregation）パターンにおけるClient Serviceの実装です。
このサービスは、Command ServiceとQuery Serviceの両方と連携し、REST APIを提供するフロントエンドサービスとして機能します。

## 目次

- [概要](#概要)
- [アーキテクチャ](#アーキテクチャ)
- [技術スタック](#技術スタック)
- [ディレクトリ構造](#ディレクトリ構造)
- [各層の責務](#各層の責務)
- [エンドポイント](#エンドポイント)
- [設定](#設定)
- [起動方法](#起動方法)
- [ビルドとテスト](#ビルドとテスト)
- [開発時の注意点](#開発時の注意点)

## 概要

Client Serviceは、CQRSアーキテクチャの一部として、以下の責務を持ちます：

- REST APIの提供（Echo Web Framework使用）
- Command ServiceとQuery Serviceの橋渡し
- HTTPクライアントによるgRPC-over-HTTP通信（Connect RPC）
- リクエストのバリデーション
- Swagger/OpenAPIドキュメントの提供

## アーキテクチャ

```
┌──────────────────────┐
│   Client Service     │
│   (REST API)         │
│   Echo + Connect     │
└──────────┬───────────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
┌─────────┐  ┌─────────┐
│ Command │  │  Query  │
│ Service │  │ Service │
│ (Write) │  │ (Read)  │
└────┬────┘  └────┬────┘
     │            │
     ▼            ▼
┌─────────┐  ┌─────────┐
│Command  │  │ Query   │
│   DB    │←─│   DB    │
│(MySQL)  │  │(MySQL)  │
└─────────┘  └─────────┘
```

### 通信フロー

1. クライアント（ブラウザ/CLI等）がREST APIを呼び出し
2. Client ServiceがHTTP経由でCommand/Query Serviceを呼び出し
3. Command ServiceまたはQuery Serviceがデータベースにアクセス
4. レスポンスをクライアントに返却

## 技術スタック

### Webフレームワーク

- **Echo v4**: 高速で使いやすいGo Web Framework
    - ミドルウェア: リクエストログ、ボディダンプ
    - バリデーション: go-playground/validator/v10

### RPC通信

- **Connect RPC**: gRPC-over-HTTP通信
    - Protocol Buffers
    - HTTP/1.1およびHTTP/2サポート

### 依存性注入

- **Uber Fx**: 依存性注入とライフサイクル管理

### 設定管理

- **Viper**: 設定ファイル（TOML）と環境変数の読み込み

### ドキュメント

- **Swagger/OpenAPI**: API仕様の自動生成
    - swaggo/swag: Goのコメントから生成

### ログ

- **slog**: 構造化ログ（Go標準ライブラリ）

### テスト

- **testify**: アサーションとモック
- **gomock**: モック生成

## ディレクトリ構造

```text
client/
├── cmd/                        # アプリケーションのエントリポイント
│   └── server/
│       └── main.go            # メイン関数（Swagger設定含む）
├── docs/                      # Swagger生成ファイル
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/                  # 内部パッケージ
│   ├── domain/               # ドメイン層
│   │   ├── models/          # ドメインモデル
│   │   │   ├── categories.go
│   │   │   └── products.go
│   │   └── repository/      # リポジトリインターフェース
│   │       └── repository.go
│   ├── infrastructure/       # インフラストラクチャ層
│   │   ├── config/          # 設定管理
│   │   │   ├── config.go
│   │   │   └── config_test.go
│   │   ├── cqrs/            # CQRS実装
│   │   │   ├── client.go        # HTTPクライアント生成
│   │   │   ├── command.go       # Command Serviceクライアント
│   │   │   ├── query.go         # Query Serviceクライアント
│   │   │   ├── repository.go    # CQRSリポジトリ実装
│   │   │   └── repository_test.go
│   │   ├── module.go        # Fxモジュール定義
│   │   └── module_test.go
│   ├── presentation/         # プレゼンテーション層
│   │   ├── dto/             # Data Transfer Objects
│   │   │   └── dto.go
│   │   ├── server/          # HTTPサーバー
│   │   │   ├── handler.go       # HTTPハンドラ
│   │   │   ├── handler_test.go
│   │   │   └── server.go        # Echoサーバー設定
│   │   ├── module.go        # Fxモジュール定義
│   │   └── module_test.go
│   └── mock/                # モック（gomockで生成）
│       └── repository/
│           └── repository_mock.go
├── config.toml               # 設定ファイル
└── README.md                # このファイル
```

## 各層の責務

### cmd/server/

アプリケーションのエントリポイント。

- **main.go**:
    - Uber Fxによる依存性注入の設定
    - アプリケーションの起動
    - Swaggerアノテーション（@title, @version, @BasePath）

### internal/domain/

ドメイン層（Domain Layer）。

- **models/**: ドメインモデル定義
    - `Category`: カテゴリエンティティ（ID、名前）
    - `Product`: 商品エンティティ（ID、名前、価格、カテゴリ）
    - 読み取り専用のシンプルなモデル（Getterのみ）

- **repository/**: リポジトリインターフェース
    - `CQRSRepository`: Command/Query Serviceへの操作を抽象化
        - カテゴリ操作: Create/Update/Delete/List/FindById
        - 商品操作: Create/Update/Delete/List/FindById/FindByKeyword

### internal/infrastructure/

インフラストラクチャ層（Infrastructure Layer）。

- **config/**: 設定管理
    - `NewViper()`: Viperによる設定ファイルと環境変数の読み込み

- **cqrs/**: CQRS通信の実装
    - `CQRSServiceConfig`: Command/Query ServiceのURL、タイムアウト設定
    - `NewClient()`: 接続プーリング設定を持つHTTPクライアント生成
    - `CommandServiceClient`: Command Serviceへのクライアント
    - `QueryServiceClient`: Query Serviceへのクライアント
    - `CQRSRepositoryImpl`: リポジトリインターフェースの実装
    - `RegisterLifecycleHooks()`: HTTPクライアントのクリーンアップ

- **module.go**: インフラ層のFxモジュール定義

### internal/presentation/

プレゼンテーション層（Presentation Layer）。

- **dto/**: データ転送オブジェクト
    - リクエスト/レスポンスの構造体
    - バリデーションタグ（required, min, max, uuid4）

- **server/**: HTTPサーバー実装
    - `CQRSServiceHandler`: REST APIハンドラ
    - `CQRSServiceServer`: Echoサーバー設定
    - `CustomValidator`: リクエストバリデーション
    - `RegisterLifecycleHooks()`: サーバーのグレースフルシャットダウン

- **module.go**: プレゼンテーション層のFxモジュール定義

## エンドポイント

### Swagger UI

- `GET /swagger/*`: Swagger UIによるAPI仕様閲覧

### カテゴリ操作

- `POST /categories`: カテゴリ作成
- `GET /categories`: カテゴリ一覧取得
- `GET /categories/:id`: カテゴリ取得
- `PUT /categories/:id`: カテゴリ更新
- `DELETE /categories/:id`: カテゴリ削除

### 商品操作

- `POST /products`: 商品作成
- `GET /products`: 商品一覧取得
- `GET /products?keyword=xxx`: 商品検索（キーワード指定）
- `GET /products/:id`: 商品取得
- `PUT /products/:id`: 商品更新
- `DELETE /products/:id`: 商品削除

## 設定

### config.toml

```toml
[server]
port = "8080"

[cqrs]
command_service_url = "http://localhost:50051"
query_service_url = "http://localhost:50052"
request_timeout = "30s"
tcp_timeout = "10s"
tcp_keep_alive = "30s"
tls_handshake_timeout = "10s"
response_header_timeout = "10s"
idle_conn_timeout = "90s"
max_idle_conns = 100
max_idle_conns_per_host = 10
```

### 環境変数

環境変数は設定ファイルの値を上書きできます（`.` を `_` に置換）：

```bash
export SERVER_PORT=8080
export CQRS_COMMAND_SERVICE_URL=http://localhost:50051
export CQRS_QUERY_SERVICE_URL=http://localhost:50052
```

## 起動方法

### 前提条件

- Go 1.25.1以上
- Command ServiceとQuery Serviceが起動していること

### ローカル起動

```bash
# 依存関係のインストール
go mod download

# Swaggerドキュメントの生成
swag init -g cmd/server/main.go -o docs

# アプリケーションの起動
go run cmd/server/main.go
```

### Docker起動（docker-compose使用）

```bash
# すべてのサービスを起動
docker-compose up -d

# ログの確認
docker-compose logs -f client
```

### アクセス確認

```bash
# ヘルスチェック（カテゴリ一覧）
curl http://localhost:8080/categories

# Swagger UI
open http://localhost:8080/swagger/index.html
```

## ビルドとテスト

### ビルド

```bash
# バイナリをビルド
go build -o bin/client cmd/server/main.go

# 実行
./bin/client
```

### テスト

```bash
# 全テストを実行
go test ./...

# カバレッジ付きで実行
go test -cover ./...

# 詳細なテスト結果
go test -v ./...

# 特定のパッケージのみ
go test -v ./internal/presentation/server/
```

### モック生成

```bash
# リポジトリのモックを生成
go generate ./internal/domain/repository/...
```

## 開発時の注意点

### HTTPクライアントの設定

- **接続プーリング**: 効率的な接続再利用のため、`MaxIdleConns`と`MaxIdleConnsPerHost`を適切に設定
- **タイムアウト**: `RequestTimeout`, `TCPTimeout`, `ResponseHeaderTimeout`を用途に応じて調整
- **グレースフルシャットダウン**: `RegisterLifecycleHooks`でアイドル接続をクリーンアップ

### バリデーション

- リクエストDTOには適切なバリデーションタグを設定
- `CustomValidator`でvalidator/v10を使用
- エラーメッセージは400 Bad Requestとして返却

### ルーティング

- `/products/search`ではなく`/products?keyword=xxx`を使用（RESTful設計）
- パスパラメータ（`:id`）より具体的なパスを先に定義

### Swagger生成

- `swag init`でドキュメントを生成
- ハンドラに`@tags`, `@Summary`, `@Description`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`を記述
- DTOは明示的なフィールドを持つ構造体（埋め込みは避ける）

### テスト

- ハンドラのテスト: gomockでリポジトリをモック化
- 統合テスト: `setupTestApp()`ヘルパーでFxアプリを構築
- `testify/require`と`testify/assert`を使用

### ログ

- 構造化ログ（slog）を使用
- エラーレベルを適切に設定（Error, Warn, Info）
- リクエスト/レスポンスのログはミドルウェアで自動出力
