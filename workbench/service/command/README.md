# Command Service

CQRS（Command Query Responsibility Segregation）パターンにおけるCommand Serviceの実装です。
このサービスは、データの作成・更新・削除といった書き込み操作を担当します。

## 概要

Command Serviceは、CQRSアーキテクチャの一部として、以下の責務を持ちます：

- データの作成、更新、削除操作の処理
- ビジネスルールとドメインロジックの実行
- データの整合性を保証するトランザクション管理
- イベントの発行（Event Sourcing使用時）

## ディレクトリ構造

```text
command/
├── cmd/                    # アプリケーションのエントリポイント
│   └── server/            # サーバー起動処理
│       └── main.go        # メイン関数
├── internal/              # 内部パッケージ（外部から直接importできない）
│   ├── application/       # アプリケーション層
│   ├── domain/           # ドメイン層
│   │   └── models/       # ドメインモデル定義
│   ├── errs/             # エラー定義
│   ├── infrastructure/   # インフラストラクチャ層
│   └── presentation/     # プレゼンテーション層
└── README.md             # このファイル
```

## 各層の責務

### cmd/

アプリケーションのエントリポイントを含みます。

- **server/main.go**: アプリケーションの起動処理、依存関係の注入、サーバーの設定を行います

### internal/application/

アプリケーション層（Application Layer）を実装します。

- **責務**:
  - ユースケースの実装
  - ドメインオブジェクトの協調
  - トランザクション管理
  - 外部サービスとの連携

### internal/domain/

ドメイン層（Domain Layer）を実装します。CQRSのCoreとなる部分です。

- **責務**:
  - ビジネスルールとドメインロジック
  - エンティティとバリューオブジェクト
  - ドメインサービス
  - リポジトリインターフェース

- **models/**: ドメインモデル（エンティティ、バリューオブジェクト）の定義

### internal/errs/

アプリケーション全体で使用するエラー定義を管理します。

- ドメインエラー
- アプリケーションエラー
- インフラストラクチャエラー

### internal/infrastructure/

インフラストラクチャ層（Infrastructure Layer）を実装します。

- **責務**:
  - データベースへのアクセス（リポジトリ実装）
  - 外部APIとの通信
  - メッセージング（Event Bus、Message Queue）
  - ファイルシステムへのアクセス

### internal/presentation/

プレゼンテーション層（Presentation Layer）を実装します。

- **責務**:
  - gRPCハンドラーの実装
  - リクエスト/レスポンスの変換
  - 入力値の検証
  - 認証・認可の処理

## アーキテクチャパターン

このCommand Serviceは以下のアーキテクチャパターンを採用しています：

### クリーンアーキテクチャ

- 依存関係の向きがドメイン層に向かっている
- 外部の詳細（データベース、Webフレームワーク等）に依存しないドメイン層

### CQRS（Command Query Responsibility Segregation）

- コマンド（書き込み）とクエリ（読み込み）の責務を分離
- このサービスはコマンド側の実装を担当

### DDD（Domain-Driven Design）

- ドメインモデルを中心とした設計
- ビジネスロジックをドメイン層に集約

## 実装ガイドライン

### 1. 依存関係の管理

- 依存関係注入（DI）を使用して、各層の結合度を下げる
- インターフェースを活用して、実装の詳細を隠蔽する

### 2. エラーハンドリング

- ドメインエラーは`internal/errs`で定義
- 各層で適切なエラー変換を行う

### 3. トランザクション管理

- アプリケーション層でトランザクション境界を管理
- ドメイン層は純粋なビジネスロジックのみに集中

### 4. テスト戦略

- 各層ごとにユニットテストを実装
- ドメイン層のテストは外部依存なしで実行可能
- インテグレーションテストでエンドツーエンドの動作を確認

## 開発時の注意点

1. **ドメイン駆動設計の原則を守る**
   - ビジネスルールはドメイン層に記述
   - 技術的な詳細はインフラストラクチャ層に分離

2. **Command側の特徴を理解する**
   - 書き込み操作のみを担当
   - データの一貫性と整合性を重視
   - 必要に応じてイベントを発行

3. **パフォーマンスよりも整合性を優先**
   - トランザクション境界を適切に設定
   - 必要に応じて楽観的排他制御を実装

## 参考リンク

- [CQRS Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)
