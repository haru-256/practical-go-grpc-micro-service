# Practical Go gRPC Micro Service - Workbench

このリポジトリは、Go言語でgRPCマイクロサービスを実装する実践的なプロジェクトです。CQRS（Command Query Responsibility Segregation）パターンを採用し、商品管理システムを題材としたマイクロサービスアーキテクチャを学習できます。

## 🎯 プロジェクト概要

### 目的

- gRPCを使用したマイクロサービスの実装方法を学習
- CQRSパターンによる読み取り・書き込み責務の分離
- Protocol Buffersによるスキーマファーストな開発
- 実用的なAPIバリデーションとエラーハンドリング

### 機能

- **商品管理**: 商品のCRUD操作とカテゴリ管理
- **Command Service**: 書き込み専用サービス（作成、更新、削除）
- **Query Service**: 読み取り専用サービス（一覧取得、詳細取得、キーワード検索）
- **Client Service**: REST APIを提供するフロントエンドサービス（Swagger UI付き）
- **スキーマバリデーション**: protovalidateによるフィールドレベル検証

## 📁 ディレクトリ構成

```text
.
├── api/                          # API定義・生成コード
│   ├── proto/                    # Protocol Buffers定義
│   └── gen/                      # 自動生成コード
│
├── k8s/                          # Kubernetesマニフェスト
│   ├── base/                     # 共通マニフェスト
│   └── overlays/                 # 環境別オーバーレイ
│
├── service/                      # マイクロサービス実装
│   ├── client/                   # REST APIクライアント（Swagger UI付き）
│   ├── command/                  # コマンドサービス（書き込み専用）
│   └── query/                    # クエリサービス（読み取り専用）
│
├── db/                           # データベース関連
│   ├── command/                  # コマンド用データベース
│   ├── query/                    # クエリ用データベース
│   └── logs/                     # データベースログ
│
└── pkg/                          # 共通ライブラリ
    └── connect/interceptor/      # Connect RPC向けのロギング/バリデーション共通インターセプター
```

## 🏗️ アーキテクチャ

このプロジェクトは**CQRS（Command Query Responsibility Segregation）**パターンを採用しています：

```mermaid
graph TB
    subgraph "Client Applications"
        Browser[Browser]
        CLI[CLI Tool]
    end
    
    subgraph "REST API Layer"
        Client[Client Service<br/>REST API + Swagger<br/>:8080]
    end
    
    subgraph "gRPC Services"
        subgraph "Command Side (書き込み)"
            CS[Command Service<br/>gRPC<br/>:50051]
            CDB[(Command DB<br/>Write Optimized)]
        end
        
        subgraph "Query Side (読み取り)"
            QS[Query Service<br/>gRPC<br/>:50052]
            QDB[(Query DB<br/>Read Optimized)]
        end
    end
    
    Browser -->|HTTP/REST| Client
    CLI -->|HTTP/REST| Client
    
    Client -->|Connect RPC| CS
    Client -->|Connect RPC| QS
    
    CS --> CDB
    QS --> QDB
    CDB -.->|Replication| QDB
    
    classDef clientService fill:#f9f,stroke:#333,stroke-width:2px
    classDef commandService fill:#fff3e0
    classDef queryService fill:#e1f5fe
    classDef database fill:#e8f5e8
    classDef client fill:#f3e5f5
    
    class Client clientService
    class CS commandService
    class QS queryService
    class CDB,QDB database
    class Browser,CLI client
```

### 特徴

- **3層アーキテクチャ**: Client Service（REST） → Command/Query Service（gRPC） → Database
- **責務の分離**: 読み取りと書き込みを独立したサービスに分離
- **REST + gRPC**: クライアント向けREST API、サービス間通信はgRPC
- **スケーラビリティ**: 各サービスを独立してスケール可能
- **データベース最適化**: 用途に応じたデータベース設計
- **型安全性**: Protocol Buffersによる厳密な型定義
- **構造化ログ**: slogによるコンテキスト対応の構造化ログ
- **依存性注入**: Uber Fxによる型安全な依存関係管理
- **共通Connectインターセプター**: slogロギングとProtovalidate検証を共通パッケージで提供

## 🚀 クイックスタート

### 前提条件

- Go 1.25.1+
- [mise](https://mise.jdx.dev/) (開発環境管理)
- Docker & Docker Compose
- [buf](https://buf.build/) CLI

### 設定ファイル

各サービスは `config.toml` で設定を管理します：

```toml
[log]
level = "info"     # ログレベル: debug, info, warn, error
format = "text"    # ログフォーマット: text, json

[mysql]
dbname = "command_db"
host = "localhost"
port = 3306
user = "root"
pass = "password"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = "1h"
conn_max_idle_time = "10m"
```

環境変数で設定を上書き可能：

- `LOG_LEVEL`, `LOG_FORMAT`: ログ設定
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`: データベース設定

### セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/haru-256/practical-go-grpc-micro-service.git
cd practical-go-grpc-micro-service

# 開発環境のセットアップ
mise install

# APIコードの生成
cd api
make generate

# 依存関係の解決
cd ..
go mod tidy
```

### サービスの起動

#### Docker Composeで起動（推奨）

```bash
# すべてのサービス（データベース＋アプリケーション）を起動
docker compose up -d --build

# サービスの起動確認
docker compose ps

# Client Service（API Gateway）にアクセス
# Swagger UIでAPIを確認
open http://localhost:8090/swagger/index.html

# ログの確認
docker compose logs -f client_service

# サービスの停止
docker compose down

# データベースも含めて完全削除
docker compose down -v
```

**ポートマッピング:**

- Client Service (REST API): `http://localhost:8090`
- Command Service (gRPC): `http://localhost:8083`
- Query Service (gRPC): `http://localhost:8085`
- Command DB (MySQL): `localhost:3306`
- Query DB (MySQL): `localhost:3307`
- phpMyAdmin: `http://localhost:3100`

**初回起動時の注意:**

Docker Composeで起動すると、データベースは自動的に作成されますが、CQRSパターンのレプリケーション設定は手動で行う必要があります。詳細は [Database README](./db/README.md) を参照してください。

```bash
# データベースのレプリケーション設定
cd db
make create-data      # テストデータ作成
make dump             # Command DBをダンプ
make restore          # Query DBにリストア
make start-replication # レプリケーション開始
```

#### ローカル開発モード

開発時はローカルでサービスを起動することもできます：

```bash
# データベースのみ起動
docker compose up -d command_db query_db db_admin

# コマンドサービス（ポート8083）
cd service/command
go run cmd/server/main.go

# クエリサービス（ポート8085）
cd service/query
go run cmd/server/main.go

# クライアントサービス（ポート8090）- Command/Queryサービスが必要
cd service/client
go run cmd/server/main.go
```

### テストの実行

```bash
# すべてのテストを実行
make test

# Command Serviceのテスト
cd service/command
make test

# Query Serviceのテスト
cd service/query
make test

# Client Serviceのテスト
cd service/client
go test ./...

# 統合テストを含む（データベースが必要）
cd service/command
go test -tags=integration ./...

cd service/query
go test -tags=integration ./...
```

## ☸️ Kubernetes (k8s)

このプロジェクトでは、[Kustomize](https://kustomize.io/) を使用してKubernetesマニフェストを管理しています。

### ディレクトリ構成

- `k8s/base/`: 全ての環境で共通のベースとなるマニフェスト
    - `namespace.yaml`: プロジェクト用の名前空間
    - `db/`: データベース関連のマニフェスト
    - `services/`: 各マイクロサービスのマニフェスト
- `k8s/overlays/dev/`: `dev`環境用の差分マニフェスト（例: リソース割り当て、レプリカ数など）

### デプロイ方法

Kustomizeを使用して、特定の環境（例: `dev`）にデプロイするには、以下のコマンドを実行します。

```bash
# dev環境のマニフェストを適用
kubectl apply -k k8s/overlays/dev
```

これにより、`base`のマニフェストと`dev`オーバーレイが結合されたマニフェストがクラスターに適用されます。

## 🛠️ 開発ワークフロー

### API仕様の変更

1. `api/proto/`でProtocol Buffersファイルを編集
2. `cd api && make generate`でコード生成
3. `go mod tidy`で依存関係更新
4. サービス実装を更新

### データベーススキーマの変更

1. `db/command/ddl/`または`db/query/ddl/`でDDLを編集
2. `cd db && make reset`でデータベースリセット
3. 新しいスキーマでサービスを再起動

## 📚 学習リソース

### サービスドキュメント

- **[Client Service](./service/client/README.md)** - REST APIサービスの実装詳細（Swagger付き）
- **[Command Service](./service/command/README.md)** - 書き込み専用サービスの実装詳細
- **[Query Service](./service/query/README.md)** - 読み取り専用サービスの実装詳細

### その他のドキュメント

- **[API仕様書](./api/README.md)** - 詳細なAPI仕様とサンプル
- **[データベース設計](./db/README.md)** - DB設計とCQRS実装
- **[プロジェクトスタイルガイド](../.gemini/styleguide.md)** - コーディング規約と設計原則

### 公式ドキュメント

- [gRPC Go](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://protobuf.dev/)
- [buf](https://buf.build/docs/)
- [Connect RPC](https://connectrpc.com/)
- [Uber Fx](https://uber-go.github.io/fx/)

## 🤝 コントリビューション

1. Issueで問題を報告または新機能を提案
2. フィーチャーブランチを作成
3. 変更をコミット（コミットメッセージは[Conventional Commits](https://www.conventionalcommits.org/)に従う）
4. プルリクエストを作成

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。
