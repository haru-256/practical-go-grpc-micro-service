# Command Service

CQRS（Command Query Responsibility Segregation）パターンにおけるCommand Serviceの実装です。
このサービスは、データの作成・更新・削除といった書き込み操作を担当します。

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
  - **products/**: 商品集約（Product エンティティ、ProductId、ProductName、ProductPrice、ProductRepository）
  - **categories/**: カテゴリ集約（Category エンティティ、CategoryId、CategoryName、CategoryRepository）

詳細は [Domain Layer README](./internal/domain/README.md) を参照してください。

#### ドメイン層の設計原則

##### ドメイン駆動設計（DDD）の適用

1. **エンティティ（Entity）**
   - 一意な識別子を持つドメインオブジェクト
   - ライフサイクルを通じて同一性が保たれる
   - 例: `Product`, `Category`

2. **値オブジェクト（Value Object）**
   - 属性によって識別されるオブジェクト
   - 不変（Immutable）である
   - 例: `ProductId`, `ProductName`, `ProductPrice`, `CategoryId`, `CategoryName`

3. **リポジトリ（Repository）**
   - 集約の永続化を抽象化するインターフェース
   - ドメイン層はインターフェースのみを定義し、実装は infrastructure 層に委譲

#### ドメインモデルの使用例

```go
// 新しい商品の作成
name, _ := products.NewProductName("サンプル商品")
price, _ := products.NewProductPrice(1000)
categoryName, _ := categories.NewCategoryName("カテゴリ1")
category, _ := categories.NewCategory(categoryName)

product, err := products.NewProduct(name, price, category)
if err != nil {
    // エラーハンドリング
}

// 商品情報の変更
newName, _ := products.NewProductName("新しい商品名")
product.ChangeName(newName)
```

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

#### エラーの種類

ドメイン層では、`errs.DomainError` を使用してドメインに関連するエラーを表現します。

**エラーコード:**

- **INVALID_ARGUMENT**: バリデーションエラー（不正な引数）
- **INTERNAL**: 内部エラー（UUID生成失敗など）

**使用例:**

```go
// バリデーションエラー
if count < MIN_LENGTH || count > MAX_LENGTH {
    return nil, errs.NewDomainError(
        "INVALID_ARGUMENT",
        fmt.Sprintf("商品名は%d文字以上%d文字以下で入力してください", MIN_LENGTH, MAX_LENGTH),
    )
}

// 原因付きエラー
if err != nil {
    return nil, errs.NewDomainErrorWithCause("INTERNAL", "商品IDの生成に失敗しました", err)
}
```

### 3. トランザクション管理

- アプリケーション層でトランザクション境界を管理
- ドメイン層は純粋なビジネスロジックのみに集中

### 4. テスト戦略

- 各層ごとにユニットテストを実装
- ドメイン層のテストは外部依存なしで実行可能
- インテグレーションテストでエンドツーエンドの動作を確認

#### ドメイン層のテスト

ドメイン層には、Ginkgo/Gomega を使用した包括的なテストが用意されています。

```bash
# すべてのドメイン層テストを実行
go test ./internal/domain/models/...

# カバレッジレポート付きで実行
go test -cover ./internal/domain/models/...
```

テスト構成:

- **値オブジェクトのテスト**: バリデーションルールの検証
- **エンティティのテスト**: 生成、再構築、変更、同一性検証
- **テーブル駆動テスト**: `DescribeTable` を使用した効率的なテストケース定義

#### バリデーションルール

##### 商品（Product）

| フィールド | 型 | 制約 |
|-----------|-----|------|
| ProductId | string | UUID形式、36文字 |
| ProductName | string | 1〜100文字 |
| ProductPrice | uint32 | 1〜1,000,000円 |
| Category | Category | 必須 |

##### カテゴリ（Category）

| フィールド | 型 | 制約 |
|-----------|-----|------|
| CategoryId | string | UUID形式、36文字 |
| CategoryName | string | 1〜100文字 |

## 開発時の注意点

1. **ドメイン駆動設計の原則を守る**
   - ビジネスルールはドメイン層に記述
   - 技術的な詳細はインフラストラクチャ層に分離
   - 値オブジェクトを積極的に使用（Primitive Obsession アンチパターンの回避）

2. **Command側の特徴を理解する**
   - 書き込み操作のみを担当
   - データの一貫性と整合性を重視
   - 必要に応じてイベントを発行

3. **パフォーマンスよりも整合性を優先**
   - トランザクション境界を適切に設定
   - 必要に応じて楽観的排他制御を実装

4. **不変性とカプセル化を維持する**
   - 値オブジェクトは不変（Immutable）
   - エンティティの変更は専用メソッド経由のみ
   - エンティティの同一性は ID で判断

5. **依存関係の向きを守る**
   - 依存関係は常にドメイン層に向かう（依存性逆転の原則）
   - リポジトリはインターフェースのみドメイン層で定義、実装はインフラ層

## ビルドとテスト

### ビルド

```bash
# プロジェクトルートから
make build

# または直接ビルド
go build -o bin/command-service ./cmd/server
```

### テスト実行

```bash
# すべてのテストを実行
make test

# カバレッジレポート付き
make test-coverage

# ドメイン層のみ
go test ./internal/domain/models/...

# 特定のパッケージ
go test ./internal/domain/models/products
go test ./internal/domain/models/categories
```

### Lint実行

```bash
make lint
```

## 参考リンク

- [CQRS Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)
