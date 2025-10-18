package service

import (
	"context"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
)

// ProductService は商品に関するアプリケーションサービスのインターフェースです。
// 商品の追加、更新、削除といったビジネスロジックを提供し、
// トランザクション管理とドメインルールの適用を担います。
type ProductService interface {
	// Add は新しい商品を追加します。
	// 商品名の重複をチェックし、既に存在する場合はApplicationErrorを返します。
	// トランザクション内で実行され、エラー時は自動的にロールバックされます。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - product: 追加する商品エンティティ
	//
	// Returns:
	//   - error: 追加に失敗した場合のエラー
	//     - ApplicationError (PRODUCT_ALREADY_EXISTS): 同じ名前の商品が既に存在する場合
	//     - CRUDError: データベース操作に失敗した場合
	//     - InternalError: トランザクション開始に失敗した場合
	Add(ctx context.Context, product *products.Product) error

	// Update は既存の商品情報を更新します。
	// リポジトリ層で商品の存在確認を行い、存在しない場合はCRUDErrorを返します。
	// トランザクション内で実行され、エラー時は自動的にロールバックされます。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - product: 更新する商品エンティティ（IDで識別）
	//
	// Returns:
	//   - error: 更新に失敗した場合のエラー
	//     - CRUDError (NOT_FOUND): 指定されたIDの商品が存在しない場合
	//     - CRUDError: データベース操作に失敗した場合
	//     - InternalError: トランザクション開始に失敗した場合
	Update(ctx context.Context, product *products.Product) error

	// Delete は指定された商品を削除します。
	// リポジトリ層で商品の存在確認を行い、存在しない場合はCRUDErrorを返します。
	// トランザクション内で実行され、エラー時は自動的にロールバックされます。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - product: 削除する商品エンティティ（IDで識別）
	//
	// Returns:
	//   - error: 削除に失敗した場合のエラー
	//     - CRUDError (NOT_FOUND): 指定されたIDの商品が存在しない場合
	//     - CRUDError: データベース操作に失敗した場合
	//     - InternalError: トランザクション開始に失敗した場合
	Delete(ctx context.Context, product *products.Product) error
}
