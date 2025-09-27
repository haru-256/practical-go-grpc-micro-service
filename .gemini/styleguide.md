# プロジェクト開発スタイルガイド

## 1. 基本方針

### 1.1. 目的と価値観

- **品質重視**: コードの品質、保守性、一貫性を向上させ、長期的な開発効率を追求する。
- **チーム成長**: 開発者の技術的成長を支援し、知識の共有を促進する
- **ユーザー価値**: 最終的にはユーザーに価値を提供することを最優先とする

### 1.2. コミュニケーション指針

- **言語**: レビュー、コメント、ドキュメントは日本語で統一
- **建設的姿勢**: 問題指摘時は必ず改善案とその理由を併記
- **学習支援**: 関連するベストプラクティスや学習リソースを積極的に共有
- **心理的安全性**: 質問や議論を歓迎し、失敗から学ぶ文化を醸成

## 2. コードレビューガイドライン

> 📋 **このセクションのサマリー**  
> 効果的なコードレビューの手法を定義。優先度付きの観点、構造化されたコメント手法、ポジティブフィードバックの重要性を説明します。

コードレビューは単なる品質チェックではなく、知識共有とチーム学習の貴重な機会です。レビュワーは教師として、レビュイーは学習者として、互いに成長できる場にします。効果的なレビューは、コードの品質向上だけでなく、チーム全体の技術力底上げにも寄与します。

### 2.1. レビューの優先順位

レビューでは無数の改善点が見つかりますが、すべてを同じ重要度で扱うと、本質的な問題が見落とされがちです。以下の優先順位に従って、重要な問題から順に取り組みます：

以下の観点を**優先度順**で評価し、レビューコメントを作成する。

| 優先度 | 観点 | チェック項目例 | 影響範囲 |
| :---: | :--- | :--- | :--- |
| **🔴 最高** | **Correctness (正確性)** | 仕様を満たしているか？バグやエッジケースは考慮されているか？ | ユーザー体験直結 |
| **🟠 高** | **Security (セキュリティ)** | SQLインジェクション等の脆弱性はないか？入力値のバリデーションは適切か？ | システム全体 |
| **🟡 中高** | **Performance (パフォーマンス)** | 非効率なアルゴリズムはないか？N+1問題やメモリリークはないか？ | システム性能 |
| **🔵 中** | **Maintainability (保守性)** | 責務は適切に分離されているか？将来の変更は容易か？複雑度は適切か？ | 開発効率 |
| **🟢 低** | **Readability (可読性)** | 命名は適切か？ロジックは追いやすいか？コメントは必要十分か？ | 開発者体験 |
| **⚪ 最低** | **Consistency (一貫性)** | プロジェクトのコーディング規約や設計パターンに準拠しているか？ | コード統一性 |

### 2.2. レビュー時の着眼点

#### 🎯 特に注目すべきポイント

- **エラーハンドリング**: 例外的なケースでの動作は安全か？
- **リソース管理**: DB接続、ファイルハンドル、goroutineの適切なクリーンアップ
- **並行処理**: データ競合やデッドロックの可能性はないか？
- **テスト容易性**: モックやスタブが使いやすい設計になっているか？

## 3. 効果的なレビューコメントの書き方

> 📋 **このセクションのサマリー**  
> 建設的で実用的なレビューコメントの作成方法。STARメソッドによる構造化、重要度の明示、学習リソースの提供を通じて、真に価値あるフィードバックを実現します。

レビューコメントは、相手に行動を促す重要なコミュニケーションツールです。単なる指摘ではなく、相手の理解を深め、具体的な改善行動につながるコメントを心がけます。効果的なコメントは、即座に問題を解決するだけでなく、将来的に同様の問題を避ける学習効果も生み出します。

### 3.1. 重要度の明示とアクションの明確化

コメントを受ける側が適切な優先順位で対応できるよう、重要度とアクションを明確に示します。これにより、限られた時間の中で最も効果的な改善が可能になります：

コメントには重要度と必要なアクションを明記し、開発者が優先順位を判断できるようにする。

| 重要度 | アクション | 対象例 | タイムライン |
| :---: | :--- | :--- | :--- |
| **🔴 CRITICAL** | **マージ前に必須修正** | バグ、セキュリティ脆弱性、データ破損リスク | 即座 |
| **🟠 HIGH** | **マージ前の修正を強く推奨** | パフォーマンス問題、保守性の問題 | 24時間以内 |
| **🟡 MEDIUM** | **次回イテレーションで修正** | リファクタリング、テスト追加 | 1週間以内 |
| **🟢 LOW** | **時間のあるときに修正** | 命名改善、コメント追加 | バックログ |
| **🔵 INFO** | **情報提供・質問** | ベストプラクティス紹介、代替手法の提案 | - |

### 3.2. 効果的なコメントの構造

コメントは以下の**STARメソッド**で構成する。

```

🎨 **[SEVERITY] 状況 (Situation)**
🔍 **原因 (Task/Problem)**
📝 **提案 (Action)**
🎆 **期待される結果 (Result)**

```

#### 悪い例：指摘のみ

```

この変数名は分かりにくい。

```

#### 良い例：STARメソッド適用

```

🎨 **[LOW] 変数名の可読性問題**
変数 `d` が何を表しているのかが一目でわかりません。

🔍 **原因**
短い略記名はコードの可読性を下げ、メンテナンスコストを増加させます。

📝 **提案**

```go
// Before
d := time.Since(startTime)

// After
elapsedTimeInSeconds := time.Since(startTime)
```

🎆 **期待される結果**
コードの意図が明確になり、新しいメンバーでも理解しやすくなります。

```

### 3.3. コード例とリソースの提供

#### コード例のベストプラクティス

1. **Before/After形式**: 現在のコードと改善後のコードを並べて表示
2. **実行可能なコード**: コピー＆ペーストですぐに試せるように
3. **コメント付き**: 重要な部分には説明コメントを追加

#### 学習リソースの提供

コメントには関連する学習リソースを積極的に含める。

- **公式ドキュメント**: Goの公式ドキュメントやEffective Goへのリンク
- **ベストプラクティス記事**: 信頼できる技術ブログや書籍への参照
- **内部リソース**: プロジェクト内の関連コードやドキュメントへのリンク

### 3.4. ポジティブフィードバックの実践

#### 認めるべき優れた点

- **✨ 優雅な設計**: シンプルで理解しやすいアーキテクチャ
- **🎨 美しいコード**: 読みやすく、一貫性のあるコード
- **🔧 巧妙な解法**: パフォーマンスやメモリ効率を考慮した実装
- **🛡️ 強固なエラーハンドリング**: 例外的ケースを適切に処理
- **🧪 網羅的なテスト**: エッジケースを含む十分なテストカバレッジ

#### ポジティブコメントの例

```

🎆 **素晴らしい実装です！**
このコンテキストパッケージの設計は、タイムアウトとキャンセルを適切に処理しており、Goのベストプラクティスに完全に準拠しています。

特に、deferを使ったリソースのクリーンアップが美しく、メモリリークやゴルーチンリークの心配がありません。

📚 **参考**: [Effective Go - コンテキスト](https://golang.org/doc/effective_go#concurrency)

```

## 4. プロジェクト固有の技術ガイドライン

> 📋 **このセクションのサマリー**  
> gRPC/Protocol Buffers、Go言語、マイクロサービスアーキテクチャに特化した実践的ガイドライン。具体的なコード例と実装パターンを通じて、プロジェクトの技術スタックを最大限に活用する方法を説明します。

本プロジェクトでは、モダンなマイクロサービスアーキテクチャを採用しています。gRPCによる高性能な内部通信、型安全なProtocol Buffers、Go言語の並行性機能を組み合わせて、スケーラブルで保守性の高いシステムを構築します。

以下のガイドラインは、これらの技術を効果的に活用し、チーム全体で一貫した実装を行うための指針です。単なるルールではなく、なぜそのパターンが推奨されるのかの理由と、具体的な実装方法を含めて説明します。

### 4.1. 🚀 gRPC & Protocol Buffers (buf + connect-go)

gRPCとProtocol Buffersは、マイクロサービス間の効率的で型安全な通信を実現します。しかし、その強力さを最大限に活用するには、適切な設計パターンと運用ルールが不可欠です。

**なぜこれらの技術を選択したのか：**
- **パフォーマンス**: バイナリプロトコルによる高速通信
- **型安全性**: コンパイル時の型チェックによるバグの早期発見
- **言語間互換性**: 複数の言語で同じスキーマを共有可能
- **後方互換性**: 適切に設計すれば、APIの進化が容易

#### 4.1.1. アーキテクチャ設計原則

適切なパッケージングと命名は、APIの可読性と保守性を大きく左右します。サービスが成長しても継続的に管理できるよう、一貫した戦略を立てます。

##### 📦 パッケージング戦略

バージョニングを含む明確なパッケージ名は、APIの進化と互換性の管理を容易にします：

```protobuf
// ✅ 推奨: 明確な階層構造
package myservice.command.v1;
package myservice.query.v1;
package myservice.common.v1;

// ❌ 非推奨: バージョンなし、不明確な構造
package myservice;
package commands;
```

**このパターンが重要な理由：**

推奨パターンでは、サービス名、機能カテゴリ、バージョンが明確に分かれます。これにより、開発者はどのAPIがどのサービスのどのバージョンに属するかを一目で理解でき、コードの可読性と保守性が大幅に向上します。また、バージョン情報が含まれることで、以降の破壊的変更の管理が容易になります。

##### 🎯 命名規約の統一

| 要素 | 規約 | 良い例 | 悪い例 |
| :--- | :--- | :--- | :--- |
| **RPCメソッド** | `VerbNoun` | `CreateProduct`, `ListOrders` | `ProductCreate`, `GetAllOrders` |
| **メッセージ** | `MethodNameRequest/Response` | `CreateProductRequest` | `ProductCreateReq`, `CreateReq` |
| **サービス** | `NounService` | `ProductService`, `OrderService` | `Products`, `OrderManager` |
| **フィールド** | `snake_case` | `user_id`, `created_at` | `userId`, `CreatedAt` |

#### 4.1.2. 🎨 メッセージ設計パターン

ユーザーフレンドリーでスケーラブルなAPIを設計するための実証済みパターンです。これらのパターンは、Google API Design Guideや一般的なベストプラクティスに基づいています。

##### 📄 ページネーション (必須実装)

大量のデータを扱うAPIでは、ページネーションは必須です。トークンベースのページネーションを採用することで、データの一貫性を保ち、パフォーマンスを最適化します：

```protobuf
// ✅ 推奨: トークンベースページネーション
message ListProductsRequest {
  int32 page_size = 1;   // 最大100、デフォルト20
  string page_token = 2; // Base64エンコードされたカーソル
  
  // フィルター条件（オプショナル）
  string category = 3;
  google.protobuf.Timestamp created_after = 4;
}

message ListProductsResponse {
  repeated Product products = 1;
  string next_page_token = 2;  // 次ページなしの場合は空文字列
  int32 total_count = 3;       // 可能な場合のみ提供
}
```

**ページネーションの本質と利点：**

このパターンは、大量のデータを扱うAPIでの標準的なアプローチです。`page_token`は不透明な文字列であり、サーバー側でデータの位置や状態を管理します。これにより、ページング中にデータが変更されても一貫性を保つことができ、ユーザーは重複や欠損なしにすべてのデータを取得できます。また、フィルター条件を含めることで、クライアントは必要なデータのみを効率的に取得でき、ネットワーク帯域と処理時間を節約できます。

##### 🎯 部分更新パターン

```protobuf
import "google/protobuf/field_mask.proto";

// ✅ 推奨: FieldMaskを使用した部分更新
message UpdateProductRequest {
  Product product = 1;
  google.protobuf.FieldMask update_mask = 2;
  
  // バージョニング（楽観的ロック）
  string etag = 3;
}

// 使用例（クライアント側）
// update_mask: "name,price,description"
```

**FieldMaskの革新性と実用性：**

`FieldMask`は部分更新の標準的な手法で、REST APIのPATCHメソッドに相当します。クライアントは更新したいフィールドのみを指定でき、サーバー側ではそのフィールドのみを更新します。これにより、不意なデータの上書きを防げると同時に、ネットワーク帯域を節約できます。`etag`を組み合わせることで、楽観的ロッキングを実現し、同時更新によるデータの破損を防ぐことができます。

##### 🔄 冪等性とリトライ対応

```protobuf
message CreateOrderRequest {
  Order order = 1;
  
  // 冪等性キー（UUID v4推奨）
  string idempotency_key = 2 [
    (validate.rules).string = {
      pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
      ignore_empty: true
    }
  ];
}
```

**冪等性キーの重要性と実装詳細：**

冪等性キーは、マイクロサービス環境での信頼性の高い通信のために不可欠です。ネットワーク障害やタイムアウトにより、クライアントが同じリクエストを繰り返し送信することがあります。冪等性キーを使用することで、サーバー側では同じキーでの重複リクエストを検知し、安全に無視または同じ結果を返すことができます。UUID v4の使用と正規表現による験証は、キーの一意性とフォーマットの正当性を保証します。

##### ✅ バリデーション戦略

```protobuf
import "validate/validate.proto";

message CreateUserRequest {
  string email = 1 [(validate.rules).string.email = true];
  string name = 2 [(validate.rules).string = {
    min_len: 1
    max_len: 100
    pattern: "^[a-zA-Z0-9\\s\\-_]+$"
  }];
  int32 age = 3 [(validate.rules).int32 = {
    gte: 0
    lte: 150
  }];
}
```

#### 4.1.3. 🔄 API進化戦略

APIの進化は、マイクロサービスアーキテクチャで最も重要な課題の一つです。不適切な変更は、システム全体の停止やデータの不整合を引き起こす可能性があります。

**API進化の金則：**

- **互換性第一**: 既存クライアントを壊さない
- **漸進的移行**: 段階的な変更でリスクを最小化
- **明確なコミュニケーション**: 変更の理由と影響を明確に伝える

##### 📋 変更管理ルール

以下のルールに従って、安全なAPIの進化を実現します：

| 変更タイプ | 許可 | 実装方法 | 注意点 |
| :--- | :---: | :--- | :--- |
| **フィールド追加** | ✅ | 新しいフィールド番号で追加 | 常にオプショナル |
| **オプショナルフィールド削除** | ⚠️ | `reserved` で予約 | フィールド番号は再利用禁止 |
| **フィールド名変更** | ❌ | 新規追加→古いものをreserved | 段階的な移行が必要 |
| **フィールド型変更** | ❌ | 新しいフィールドとして追加 | 互換性なし |
| **必須フィールド化** | ❌ | バリデーション層で対応 | protoではオプショナル維持 |

##### 🛡️ 破壊的変更の防止

```yaml
# .github/workflows/buf.yml
- name: Breaking Change Detection
  uses: bufbuild/buf-action@v1
  with:
    breaking_against: 'https://github.com/${{ github.repository }}.git#branch=main,subdir=api'
    lint: true
    breaking: true
```

##### 📝 フィールド廃止のベストプラクティス

```protobuf
message Product {
  string id = 1;
  string name = 2;
  
  // 廃止されたフィールドの予約
  reserved 3;           // old_price フィールド（削除済み）
  reserved "old_price"; // フィールド名も予約
  
  // 新しい価格フィールド
  Price price = 4;
}
```

### 4.2. 🛡️ エラーハンドリング戦略

マイクロサービスアーキテクチャでは、エラーは避けられない現実です。重要なのは、エラーが発生した際にシステム全体の信頼性を保ち、適切な復旧措置を取れることです。

**効果的なエラーハンドリングの原則：**

- **透明性**: エラーの原因と発生箇所を明確に特定できる
- **回復可能性**: 一時的なエラーは自動的に回復を試行する
- **ユーザビリティ**: エンドユーザーに分かりやすいエラーメッセージを提供
- **運用性**: 運用チームが迅速に問題を特定・解決できる情報を提供

#### 4.2.1. エラーラッピングのベストプラクティス

Goのエラーハンドリングは、適切にラッピングすることでスタックトレースのような情報を保持できます。これにより、デバッグ時に問題の根本原因を特定しやすくなります：

```go
// ✅ 推奨: コンテキスト付きエラーラッピング
func (s *ProductService) GetProduct(ctx context.Context, id string) (*Product, error) {
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get product %s: %w", id, err)
    }
    return product, nil
}

// ❌ 非推奨: コンテキスト情報の欠如
func (s *ProductService) GetProduct(ctx context.Context, id string) (*Product, error) {
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err // コンテキスト情報が失われる
    }
    return product, nil
}
```

**エラーラッピングの哲学と実践的利益：**

推奨パターンでは、`fmt.Errorf`の`%w`ベーブを使って元のエラーをラッピングし、同時にコンテキスト情報（この例では商品ID）を追加しています。これにより、エラーが発生した際に、開発者は具体的にどの商品の取得に失敗したのかを理解でき、デバッグ作業が大幅に効率化されます。また、元のエラーが保持されるため、エラーチェーンを追跡して根本原因を特定できます。

#### 4.2.2. エラータイプ判定とハンドリング

```go
// カスタムエラー型の定義
type DomainError struct {
    Code    string
    Message string
    Cause   error
}

func (e *DomainError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
    return e.Cause
}

// エラータイプ別のハンドリング
func (h *Handler) handleError(err error) *connect.Error {
    var domainErr *DomainError
    if errors.As(err, &domainErr) {
        switch domainErr.Code {
        case "NOT_FOUND":
            return connect.NewError(connect.CodeNotFound, domainErr)
        case "INVALID_ARGUMENT":
            return connect.NewError(connect.CodeInvalidArgument, domainErr)
        case "PERMISSION_DENIED":
            return connect.NewError(connect.CodePermissionDenied, domainErr)
        default:
            return connect.NewError(connect.CodeInternal, domainErr)
        }
    }
    
    // システムエラーの場合
    return connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
}
```

**カスタムエラー型の設計意図と効果：**

この`DomainError`構造体は、ビジネスロジック層で発生するエラーを構造化し、標準化するためのパターンです。`Code`フィールドはエラーの種類を機械的に判定できるようにし、`Message`フィールドは人間が読める説明を提供します。`Cause`フィールドにより元のエラーを保持し、エラーチェーンを構築できます。`handleError`関数では、ドメインエラーを適切なgRPCステータスコードにマッピングし、クライアントが適切な対応を取れるようにしています。

#### 4.2.3. ログ出力とメトリクス

```go
// エラーレベル別のログ出力
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) error {
    if err := s.validator.Validate(req); err != nil {
        // クライアントエラーはINFOレベル
        slog.InfoContext(ctx, "validation failed", 
            "error", err, 
            "request", req)
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if err := s.repo.Create(ctx, req.Product); err != nil {
        // システムエラーはERRORレベル
        slog.ErrorContext(ctx, "failed to create product", 
            "error", err,
            "product_id", req.Product.Id)
        return fmt.Errorf("failed to create product: %w", err)
    }
    
    return nil
}
```

### 4.3. 🔗 依存関係管理

適切な依存関係管理は、テスタブルで保守性の高いコードの基盤です。依存性注入（DI）パターンを活用することで、各コンポーネントの責務を明確にし、結合度を下げます。

**依存関係管理の利点：**

- **テスタビリティ**: モックを使った単体テストが容易
- **柔軟性**: 実装を簡単に切り替え可能
- **保守性**: 変更の影響範囲を限定
- **再利用性**: コンポーネントの独立性が高まる

#### 4.3.1. 依存性注入 (DI) パターン

Goでは明示的な依存性注入フレームワークは必須ではありませんが、Googleが開発したWireを使用することで、コンパイル時に依存関係を解決し、ランタイムエラーを防げます：

```go
// インターフェースの定義（ドメイン層）
type ProductRepository interface {
    FindByID(ctx context.Context, id string) (*Product, error)
    Create(ctx context.Context, product *Product) error
    Update(ctx context.Context, product *Product) error
    Delete(ctx context.Context, id string) error
}

type ProductService interface {
    GetProduct(ctx context.Context, id string) (*Product, error)
    CreateProduct(ctx context.Context, product *Product) error
}

// サービスの実装（アプリケーション層）
type productService struct {
    repo      ProductRepository // インターフェースに依存
    validator Validator
    logger    *slog.Logger
}

// コンストラクタインジェクション
func NewProductService(
    repo ProductRepository,
    validator Validator,
    logger *slog.Logger,
) ProductService {
    return &productService{
        repo:      repo,
        validator: validator,
        logger:    logger,
    }
}
```

**依存性注入パターンの設計哲学：**

このパターンの本質は、具体的な実装ではなく抽象的なインターフェースに依存することです。`productService`は`ProductRepository`インターフェースに依存し、具体的なデータベース実装（MySQL、PostgreSQLなど）を知りません。これにより、サービスロジックを変更することなく、データストレージを切り替えたり、テスト時にモックを使用したりすることができます。コンストラクタでの依存性注入は、依存関係を明示的にし、初期化時に必要なすべてのコンポーネントが提供されることを保証します。

#### 4.3.2. モジュール構造とWireによるDI

```go
// cmd/server/wire.go
//+build wireinject

package main

import (
    "github.com/google/wire"
    "yourproject/internal/domain"
    "yourproject/internal/infrastructure"
    "yourproject/internal/application"
)

// Provider Sets
var InfrastructureSet = wire.NewSet(
    infrastructure.NewMySQLProductRepository,
    wire.Bind(new(domain.ProductRepository), new(*infrastructure.MySQLProductRepository)),
)

var ApplicationSet = wire.NewSet(
    application.NewProductService,
    wire.Bind(new(domain.ProductService), new(*application.ProductService)),
)

var PresentationSet = wire.NewSet(
    presentation.NewProductHandler,
)

// Wireによる依存関係の解決
func InitializeServer() (*Server, error) {
    wire.Build(
        InfrastructureSet,
        ApplicationSet,
        PresentationSet,
        NewServer,
    )
    return nil, nil
}
```

**Wireの強力さと設計の优雅さ：**

Google Wireはコンパイル時に依存関係を解決するため、ランタイムエラーやパフォーマンスのオーバーヘッドがありません。`wire.NewSet`でグループ化されたプロバイダーセットは、アーキテクチャの層を明確に表現し、依存関係の正当性を保証します。`wire.Bind`はインターフェースと実装を結びつけ、具体的な型ではなく抽象的なインターフェースを通じてコンポーネントが連携できるようにします。このアプローチにより、コードの可読性、テスタビリティ、保守性が大幅に向上します。

#### 4.3.3. バージョン管理戦略

```go
// go.mod
module github.com/yourorg/yourproject

go 1.21

require (
    connectrpc.com/connect v1.11.1
    github.com/bufbuild/protovalidate-go v0.4.2
    google.golang.org/protobuf v1.31.0
)

// バージョンアップデートの手順
// 1. Dependabotの設定
// 2. 定期的な`go mod tidy && go mod verify`
// 3. セキュリティアップデートの優先的適用
```

### 4.4. 🌐 Goマイクロサービス実装

Goの言語特性（並行性、シンプルさ、パフォーマンス）を活かしたマイクロサービスの実装パターンを説明します。本番環境での運用を考慮した、実用的で堅牢な実装方法を重視します。

**Goがマイクロサービスに適している理由：**

- **軽量**: 小さなメモリフットプリントで多数のサービスを稼働可能
- **並行性**: Goroutineによる効率的な並行処理
- **デプロイ**: 単一バイナリでの簡単なデプロイメント
- **パフォーマンス**: 低レイテンシ、高スループットを実現

#### 4.4.1. 📊 設定管理パターン

マイクロサービスでは、環境ごとに異なる設定を柔軟に管理する必要があります。Twelve-Factor Appの原則に従い、環境変数を中心とした設定管理を行います：

```go
// config/config.go
type Config struct {
    Server   ServerConfig   `env:",prefix=SERVER_"`
    Database DatabaseConfig `env:",prefix=DB_"`
    Redis    RedisConfig    `env:",prefix=REDIS_"`
    Log      LogConfig      `env:",prefix=LOG_"`
}

type ServerConfig struct {
    Port         int           `env:"PORT,default=8080"`
    ReadTimeout  time.Duration `env:"READ_TIMEOUT,default=30s"`
    WriteTimeout time.Duration `env:"WRITE_TIMEOUT,default=30s"`
    GrpcPort     int           `env:"GRPC_PORT,default=9090"`
}

type DatabaseConfig struct {
    Host         string        `env:"HOST,default=localhost"`
    Port         int           `env:"PORT,default=3306"`
    User         string        `env:"USER,required"`
    Password     string        `env:"PASSWORD,required"`
    Database     string        `env:"NAME,required"`
    MaxOpenConns int           `env:"MAX_OPEN_CONNS,default=25"`
    MaxIdleConns int           `env:"MAX_IDLE_CONNS,default=5"`
    MaxLifetime  time.Duration `env:"MAX_LIFETIME,default=5m"`
}

// 設定の読み込みとバリデーション
func Load() (*Config, error) {
    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    return &cfg, nil
}

func (c *Config) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return errors.New("invalid server port")
    }
    // その他のバリデーション...
    return nil
}
```

**構造化設定管理の哲学と実装の妙味：**

このパターンは、Twelve-Factor Appの「設定は環境変数に格納」原則をGoで実現するエレガントな方法です。構造体タグでプレフィックス、デフォルト値、必須フィールドを宣言的に記述でき、設定の意図がコード自体で明確になります。環境変数名のプレフィックス化により、複数のサービスが同じコンテナで動作しても設定の競合を避けられ、`Validate`メソッドにより起動時に設定の整合性を確認できます。これにより、訳のわからないランタイムエラーを防げることができます。

#### 4.4.2. ロギング

マイクロサービス環境では、ログはデバッグ、モニタリング、トラブルシューティングの主要な手段です。構造化されたログは、自動解析ツールやモニタリングシステムとの連携を容易にします。

**原則**:

- **構造化**: 標準ライブラリの `log/slog` を使用し、本番環境ではJSON形式で出力する。
- **コンテキスト**: リクエストIDやトレースIDなどのコンテキスト情報をすべてのログに含める。
- **レベル**: クライアント起因のエラーは`INFO`、サーバー内部の問題は`ERROR`レベルで記録する。

```go
// pkg/logging/logger.go

// NewLogger は環境に応じたslog.Loggerを生成します。
func NewLogger(level string, isDevelopment bool) *slog.Logger {
 var logLevel slog.Level
 switch strings.ToLower(level) {
 case "debug":
  logLevel = slog.LevelDebug
 case "warn":
  logLevel = slog.LevelWarn
 case "error":
  logLevel = slog.LevelError
 default:
  logLevel = slog.LevelInfo
 }

 opts := &slog.HandlerOptions{
  Level:     logLevel,
  AddSource: isDevelopment, // 開発時のみソース位置を出力
 }

 var handler slog.Handler
 if isDevelopment {
  handler = slog.NewTextHandler(os.Stdout, opts)
 } else {
  handler = slog.NewJSONHandler(os.Stdout, opts)
 }

 logger := slog.New(handler)
 slog.SetDefault(logger) // グローバルロガーとして設定
 return logger
}
```

#### 4.4.3. 🔄 コンテキスト管理パターン

Goの`context.Context`は、リクエストスコープの情報、タイムアウト、キャンセルシグナルを管理するための標準的な手段です。適切に使用することで、ユーザーのリクエストキャンセルやタイムアウトに適切に対応できます。

**原則**:

- **伝播**: すべての関数呼び出しで第一引数として渡す。
- **キー管理**: キーの衝突を避けるため、独自型を定義する。
- **値のスコープ**: リクエストスコープの情報（リクエストID、認証情報など）のみを格納する。

```go
// context/context.go
type contextKey string

const (
 RequestIDKey contextKey = "request_id"
 LoggerKey    contextKey = "logger"
)

// WithLogger はコンテキストにロガーを格納します。
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
 return context.WithValue(ctx, LoggerKey, logger)
}

// FromContext はコンテキストからロガーを取得します。
func FromContext(ctx context.Context) *slog.Logger {
 if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
  return logger
 }
 return slog.Default() // 見つからない場合はデフォルトロガーを返す
}

// 使用例 (Interceptor)
func NewLoggingInterceptor() connect.UnaryInterceptor {
 return func(next connect.UnaryFunc) connect.UnaryFunc {
  return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
   requestID := uuid.NewString() // or from header
   logger := slog.With("request_id", requestID, "procedure", req.Spec().Procedure)
   ctx = WithLogger(ctx, logger)

   logger.InfoContext(ctx, "request started")
   res, err := next(ctx, req)
   if err != nil {
    logger.ErrorContext(ctx, "request failed", "error", err)
   } else {
    logger.InfoContext(ctx, "request finished")
   }
   return res, err
  }
 }
}

// コンテキストから値を取得
func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        return requestID
    }
    return ""
}

// タイムアウトとキャンセルのハンドリング
func (s *ProductService) GetProduct(ctx context.Context, id string) (*Product, error) {
    // タイムアウト付きコンテキストの作成
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // キャンセルチェック
    select {
    case <-ctx.Done():
        return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
    default:
    }
    
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get product: %w", err)
    }
    
    return product, nil
}
```

#### 4.4.4. Graceful Shutdown

サービスの停止時に、処理中のリクエストを適切に完了させ、リソースを正常にクリーンアップすることで、データの整合性を保ち、ユーザー体験を向上させます。

**原則**:

- **シグナルハンドリング**: `os.Signal` を捕捉し、`SIGINT` や `SIGTERM` を受け取った際にシャットダウン処理を開始する。
- **並行クリーンアップ**: `errgroup` を使用して、HTTPサーバー、DBコネクション、その他のリソースのクローズ処理を並行して行い、シャットダウン時間を短縮する。
- **タイムアウト**: シャットダウン処理全体にタイムアウトを設定し、無期限にブロックされるのを防ぐ。

```go
// cmd/server/main.go
func main() {
 ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
 defer stop()

 // ... （設定読み込み、ロガー・サーバー初期化）
 // server は *http.Server を持つ構造体
 // cleanup はDBなどのリソースを解放する関数

 // サーバーをgoroutineで起動
 go func() {
  slog.Info("starting server", "addr", server.Addr)
  if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
   slog.Error("failed to start server", "error", err)
   stop() // サーバー起動に失敗したらシャットダウン
  }
 }()

 // シャットダウンシグナルを待機
 <-ctx.Done()
 slog.Info("shutdown signal received, starting graceful shutdown")

 // Graceful shutdownのタイムアウト付きコンテキスト
 shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
 defer cancel()

 // HTTPサーバーのシャットダウン
 if err := server.Shutdown(shutdownCtx); err != nil {
  slog.Error("HTTP server shutdown error", "error", err)
 }

 // その他のリソースのクリーンアップ
 if err := cleanup(); err != nil {
  slog.Error("cleanup error", "error", err)
 }

 slog.Info("server stopped gracefully")
}
```

#### 4.4.5. 🧪 テスト戦略

**テストピラミッドの構成：**

- **結合テスト** (20%): コンポーネント間の連携テスト
- **E2Eテスト** (10%): システム全体の動作確認

```go
// テーブル駆動テストのベストプラクティス
func TestProductService_CreateProduct(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateProductRequest
        setup   func(*testing.T) ProductRepository
        want    *Product
        wantErr bool
        errType error
    }
    {
        {
            name: "正常ケース",
            input: &CreateProductRequest{
                Product: &Product{
                    Name:     "Test Product",
                    Price:    1000,
                    Category: "Electronics",
                },
            },
            setup: func(t *testing.T) ProductRepository {
                repo := &MockProductRepository{}
                repo.On("Create", mock.Anything, mock.Anything).Return(nil)
                return repo
            },
            wantErr: false,
        },
        {
            name: "バリデーションエラー",
            input: &CreateProductRequest{
                Product: &Product{
                    Name:  "", // 空文字列
                    Price: -100, // 負の値
                },
            },
            setup: func(t *testing.T) ProductRepository {
                return &MockProductRepository{}
            },
            wantErr: true,
            errType: &ValidationError{},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := tt.setup(t)
            service := NewProductService(repo, NewValidator(), slog.Default())
            
            ctx := context.Background()
            got, err := service.CreateProduct(ctx, tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.ErrorAs(t, err, &tt.errType)
                }
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}

// インテグレーションテストの例
func TestProductHandler_Integration(t *testing.T) {
    // テスト用DBのセットアップ
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // テストサーバーの起動
    server := setupTestServer(t, db)
    defer server.Close()
    
    client := connectclient.NewProductServiceClient(
        http.DefaultClient,
        server.URL,
    )
    
    // テスト実行
    res, err := client.CreateProduct(context.Background(), connect.NewRequest(&CreateProductRequest{
        Product: &Product{
            Name:     "Integration Test Product",
            Price:    2000,
            Category: "Test",
        },
    }))
    
    assert.NoError(t, err)
    assert.NotEmpty(t, res.Msg.Product.Id)
    assert.Equal(t, "Integration Test Product", res.Msg.Product.Name)
}
```

**テーブル駆動テストの哲学と実装美学：**

テーブル駆動テストは、同じロジックを異なる入力で網羅的にテストするためのエレガントな手法です。各テストケースは構造体で定義され、`name`フィールドで明確な説明、`setup`関数でテスト固有のモック設定、`wantErr`と`errType`でエラーケースの験証を行います。このパターンの強力さは、新しいテストケースの追加が簡単で、テストコードの重複を最小限に抑えられることです。インテグレーションテストでは、実際のデータベースとHTTPサーバーを使用して、システム全体の結合部を検証し、単体テストでは捕捉できない問題を発見できます。

---

## 5. 🚀 まとめとベストプラクティス

> 📋 **このセクションのサマリー**  
> プロジェクト全体を通じて重要な原則とチェックリストをまとめ。日常的な開発からリリースまで、品質を保つためのガイダンスを提供します。

このスタイルガイドで説明した各要素は、独立したテクニックではなく、相互に関連し合うシステムです。コードレビューの文化、技術的な実装パターン、チームコミュニケーションが組み合わさることで、初めて持続可能で価値ある開発が実現されます。

**スタイルガイドの真の目的：**

1. **効率性の向上**: 一貫したパターンにより、認知負荷を減らし開発速度を向上
2. **品質の保証**: 体系的なアプローチで、バグや設計問題を早期に発見・修正
3. **知識の共有**: チーム全体で技術的な知見を共有し、集合知を活用
4. **持続可能性**: 長期的にメンテナンスしやすいコードベースの構築

### 5.1. 重要な原則

以下の4つの原則は、すべての技術的判断の基盤となります：

1. **シンプルさを優先** - 複雑さよりもシンプルさを選ぶ
2. **一貫性を維持** - プロジェクト全体で統一されたスタイルを保つ
3. **フィードバックを大切に** - 建設的な議論で品質を向上
4. **継続的改善** - 小さな改善を積み重ねる

### 5.2. チェックリスト

#### 📄 コードレビュー時

- [ ] 正確性: 仕様を満たし、バグがないか？
- [ ] セキュリティ: 脆弱性や入力バリデーション漏れはないか？
- [ ] パフォーマンス: N+1問題やメモリリークはないか？
- [ ] テスト: 十分なテストカバレッジがあるか？
- [ ] ドキュメント: コメントやREADMEが更新されているか？

#### 🔧 リリース前

- [ ] CI/CD: すべてのテストが成功しているか？
- [ ] API: 破壊的変更がないかチェック済みか？
- [ ] ログ: 総合テストでログ出力を確認したか？
- [ ] メトリクス: パフォーマンスメトリクスを確認したか？

### 5.3. 参考リソース

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Connect-Go Documentation](https://connectrpc.com/docs/go/getting-started)
- [Buf Documentation](https://buf.build/docs)
- [Google API Design Guide](https://cloud.google.com/apis/design)
