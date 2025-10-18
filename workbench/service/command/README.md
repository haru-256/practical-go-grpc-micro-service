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

- **server/main.go**: アプリケーションの起動処理、Uber Fxによる依存関係の注入、サーバーの設定を行います

### internal/application/

アプリケーション層（Application Layer）を実装します。

- **責務**:
  - ユースケースの実装
  - ドメインオブジェクトの協調
  - トランザクション管理
  - 外部サービスとの連携

- **service/**: サービスインターフェース定義
  - `ProductService`: 商品に関するビジネスロジック（Add/Update/Delete）
  - `CategoryService`: カテゴリに関するビジネスロジック（Add/Update/Delete）
  - `TransactionManager`: トランザクション管理

- **impl/**: サービス実装
  - `ProductServiceImpl`: ProductServiceの具象実装
  - `CategoryServiceImpl`: CategoryServiceの具象実装

- **module.go**: Uber Fxモジュール定義（アプリケーション層の依存関係を構成）

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

- **sqlboiler/**: SQLBoilerを使用したデータベース実装
  - **repository/**: リポジトリ実装
    - `ProductRepositoryImpl`: 商品リポジトリの具象実装
    - `CategoryRepositoryImpl`: カテゴリリポジトリの具象実装
    - `TransactionManagerImpl`: トランザクションマネージャーの具象実装
  - **handler/**: データベース接続管理
  - **module.go**: Uber Fxモジュール定義（インフラ層の依存関係を構成）

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

### 1. 依存関係の管理と設計パターン

#### 「インターフェースを受け入れ、具象（struct）を返す」原則

このプロジェクトでは、Goの重要な設計思想に従っています：

**コンストラクタの設計原則:**

```go
// ✅ 推奨: インターフェースを受け取り、具象型を返す
func NewProductServiceImpl(
    repo products.ProductRepository,  // インターフェースを受け入れる
    tm service.TransactionManager,    // インターフェースを受け入れる
) *ProductServiceImpl {               // 具象型を返す
    return &ProductServiceImpl{
        repo: repo,
        tm:   tm,
    }
}

// 使用例
svc := impl.NewProductServiceImpl(repo, tm)
var service service.ProductService = svc  // 必要に応じてインターフェースとして扱う
```

**この原則の利点:**

- **柔軟性**: 呼び出し側がインターフェースとして扱うか、具象型として扱うかを選択できる
- **明示性**: コンストラクタ名（`NewXxxImpl`）で具象型を返すことが明確
- **テスタビリティ**: モックを使った単体テストが容易
- **DIフレームワークとの親和性**: Uber Fxで簡単にインターフェースに変換可能

#### Uber Fxによる依存性注入

依存性注入には**Uber Fx**を使用しています：

```go
// internal/application/module.go
var Module = fx.Module(
    "application",
    sqlboiler.Module,  // インフラ層を含める
    fx.Provide(
        // 具象型を返すコンストラクタをインターフェースに変換
        fx.Annotate(
            impl.NewProductServiceImpl,
            fx.As(new(service.ProductService)),
        ),
        fx.Annotate(
            impl.NewCategoryServiceImpl,
            fx.As(new(service.CategoryService)),
        ),
    ),
)
```

**Fxの利点:**

- **自動的な依存解決**: コンストラクタの引数から依存関係を自動解決
- **ライフサイクル管理**: リソースの初期化と終了処理を自動管理
- **モジュラー設計**: 層ごとにモジュールを分離して可視化
- **型安全**: コンパイル時に依存関係の整合性を検証

詳細は [プロジェクトスタイルガイド](../../.gemini/styleguide.md#432-モジュール構造とuber-fxによる依存性注入) を参照してください。

### 2. エラーハンドリング

各層で適切なエラー型を使用し、エラー情報を保持しながら上位層に伝播させます。

#### エラーの種類

**ドメインエラー (`errs.DomainError`)**

ドメイン層で発生する、ビジネスルール違反を表すエラーです。

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

**アプリケーションエラー (`errs.ApplicationError`)**

アプリケーション層で発生する、ビジネスロジック上のエラーです。

**エラーコード例:**

- **PRODUCT_ALREADY_EXISTS**: 商品名の重複
- **CATEGORY_ALREADY_EXISTS**: カテゴリ名の重複

**使用例:**

```go
// 商品名の重複チェック
exists, err := s.repo.ExistsByName(ctx, tx, product.Name())
if err != nil {
    return err
}
if exists {
    return errs.NewApplicationError("PRODUCT_ALREADY_EXISTS", "商品名が既に存在します")
}
```

**CRUDエラー (`errs.CRUDError`)**

リポジトリ層で発生する、データアクセスに関するエラーです。

**エラーコード例:**

- **NOT_FOUND**: リソースが見つからない
- **DUPLICATE_KEY**: 主キーやユニークキーの重複

**使用例:**

```go
// リソースが見つからない場合
if errors.Is(err, sql.ErrNoRows) {
    return errs.NewCRUDError("NOT_FOUND", 
        fmt.Sprintf("商品番号: %s は存在しません", id.Value()))
}
```

#### エラーラッピング

エラーは`fmt.Errorf`の`%w`を使用してラッピングし、コンテキスト情報を追加します：

```go
if err := s.repo.Create(ctx, tx, product); err != nil {
    return fmt.Errorf("failed to create product %s: %w", product.Name().Value(), err)
}
```

### 3. トランザクション管理

- アプリケーション層でトランザクション境界を管理
- ドメイン層は純粋なビジネスロジックのみに集中

### 4. テスト戦略

各層ごとにユニットテストと統合テストを実装し、包括的なテストカバレッジを確保しています。

#### テストツール

- **Ginkgo v2**: BDDスタイルのテストフレームワーク
- **Gomega**: マッチャーライブラリ
- **gomock**: モック生成ツール

#### ドメイン層のテスト

ドメイン層には、外部依存なしで実行できる包括的なテストが用意されています。

```bash
# すべてのドメイン層テストを実行
go test ./internal/domain/models/...

# カバレッジレポート付きで実行
go test -cover ./internal/domain/models/...
```

**テスト構成:**

- **値オブジェクトのテスト**: バリデーションルールの検証
- **エンティティのテスト**: 生成、再構築、変更、同一性検証
- **テーブル駆動テスト**: `DescribeTable` を使用した効率的なテストケース定義

#### アプリケーション層のテスト

アプリケーション層には、ユニットテストと統合テストの両方があります。

**ユニットテスト（gomock使用）:**

```bash
# Product Service ユニットテスト
go test ./internal/application/impl/product_impl_test.go

# Category Service ユニットテスト  
go test ./internal/application/impl/category_impl_test.go
```

モックを使用して、依存関係を分離したテストを実行します。

**統合テスト（実際のDB使用）:**

```bash
# Product Service 統合テスト
go test ./internal/application/impl/product_impl_integration_test.go

# Category Service 統合テスト
go test ./internal/application/impl/category_impl_integration_test.go
```

実際のPostgreSQLデータベースを使用して、エンドツーエンドの動作を検証します。

#### すべてのテストを実行

Ginkgoを使用してプロジェクト全体のテストを実行：

```bash
# Makefileを使用（推奨）
make ginkgo

# または直接実行
go tool ginkgo run ./...
```

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

6. **「インターフェースを受け入れ、具象を返す」原則の遵守**
   - コンストラクタは必ず具象型（`*XxxImpl`）を返す
   - コンストラクタ名に`Impl`サフィックスを付ける（例：`NewProductServiceImpl`）
   - Fxモジュールで`fx.Annotate`と`fx.As`を使用してインターフェースに変換

7. **レイヤー間の責務分離**
   - **Service層**: ビジネスロジック（重複チェックなど）のみを実装
   - **Repository層**: データアクセスと存在確認を実装
   - Service層でのUpdate/Delete時は、存在チェックをRepository層に委譲

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
# すべてのテストを実行（Ginkgo使用）
make ginkgo

# または
go tool ginkgo run ./...

# 特定の層のテスト
go test ./internal/domain/models/...        # ドメイン層
go test ./internal/application/impl/...     # アプリケーション層
go test ./internal/infrastructure/...       # インフラ層

# 個別のパッケージ
go test ./internal/domain/models/products
go test ./internal/domain/models/categories
go test ./internal/application/impl -run TestProductService
go test ./internal/application/impl -run TestCategoryService
```

### データベースセットアップ

統合テストを実行する前に、PostgreSQLデータベースを起動してください：

```bash
# データベースの起動（docker-compose使用）
cd ../../db
make up

# データベースの停止
make down
```

### Lint実行

```bash
make lint
```

## 参考リンク

### 設計パターンとアーキテクチャ

- [CQRS Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs) - Microsoft Azure アーキテクチャセンター
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Robert C. Martin
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans

### Go言語とベストプラクティス

- [Effective Go](https://golang.org/doc/effective_go) - Go言語公式ガイド
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Goコードレビューのベストプラクティス

### 依存性注入

- [Uber Fx Documentation](https://uber-go.github.io/fx/) - Uber Fx公式ドキュメント
- [Fx Modules Guide](https://uber-go.github.io/fx/modules.html) - モジュール設計パターン

### テスティング

- [Ginkgo](https://onsi.github.io/ginkgo/) - BDDスタイルのテストフレームワーク
- [Gomega](https://onsi.github.io/gomega/) - マッチャーライブラリ
- [gomock](https://github.com/golang/mock) - モック生成ツール

### プロジェクト内部リソース

- [プロジェクトスタイルガイド](../../.gemini/styleguide.md) - コーディング規約と設計原則
- [Domain Layer README](./internal/domain/README.md) - ドメイン層の詳細
