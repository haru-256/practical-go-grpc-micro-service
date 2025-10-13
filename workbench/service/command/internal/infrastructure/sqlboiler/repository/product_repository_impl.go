package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/models"
)

// productRepositoryImpl は商品リポジトリのSQLBoilerを使用した実装です。
type productRepositoryImpl struct{}

// NewProductRepository は新しいProductRepositoryインスタンスを生成します。
// この関数は、商品の挿入、更新、削除後に実行されるフックを登録します。
func NewProductRepository() products.ProductRepository {
	models.AddProductHook(boil.AfterInsertHook, ProductAfterInsertHook)
	models.AddProductHook(boil.AfterUpdateHook, ProductAfterUpdateHook)
	models.AddProductHook(boil.AfterDeleteHook, ProductAfterDeleteHook)
	return &productRepositoryImpl{}
}

// ExistsById は指定された商品IDが存在するかどうかをチェックします。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - id: チェック対象の商品ID
//
// Returns:
//   - bool: 商品が存在する場合はtrue、存在しない場合はfalse
//   - error: データベースエラーが発生した場合
func (r *productRepositoryImpl) ExistsById(ctx context.Context, tx *sql.Tx, id *products.ProductId) (bool, error) {
	condition := models.ProductWhere.ObjID.EQ(id.Value())
	exists, err := models.Products(condition).Exists(ctx, tx)
	if err != nil {
		return false, handler.DBErrHandler(err)
	}
	return exists, nil
}

// ExistsByName は指定された名前の商品が存在するかどうかをチェックします。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - name: チェック対象の商品名
//
// Returns:
//   - bool: 商品が存在する場合はtrue、存在しない場合はfalse
//   - error: データベースエラーが発生した場合
func (r *productRepositoryImpl) ExistsByName(ctx context.Context, tx *sql.Tx, name *products.ProductName) (bool, error) {
	condition := models.ProductWhere.Name.EQ(name.Value())
	exists, err := models.Products(condition).Exists(ctx, tx)
	if err != nil {
		return false, handler.DBErrHandler(err)
	}
	return exists, nil
}

// Create は新しい商品をデータベースに追加します。
//
// NOTE: boil.Infer() を使用することで、auto-incrementのIDは無視され、
// DB側で自動採番された後、SQLBoiler側の構造体にセットされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - Product: 作成する商品
//
// Returns:
//   - error: データベースエラーが発生した場合
func (r *productRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, product *products.Product) error {
	newProduct := models.Product{
		ObjID:      Product.Id().Value(),
		Name:       Product.Name().Value(),
		Price:      int(Product.Price().Value()),
		CategoryID: Product.Category().Id().Value(),
	}
	// NOTE: boil.Infer() でauto-incrementのIDは無視され、勝手にDB側で採番された後、sqlboiler側の構造体にセットされる
	if err := newProduct.Insert(ctx, tx, boil.Infer()); err != nil {
		return handler.DBErrHandler(err)
	}
	return nil
}

// UpdateById は指定されたIDの商品情報を更新します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - Product: 更新する商品情報
//
// Returns:
//   - error: 商品が存在しない場合はNOT_FOUNDエラー、
//     データベースエラーが発生した場合はそのエラー
func (r *productRepositoryImpl) UpdateById(ctx context.Context, tx *sql.Tx, Product *products.Product) error {
	condition := models.ProductWhere.ObjID.EQ(Product.Id().Value())
	upModel, err := models.Products(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDErrorWithCause("NOT_FOUND", fmt.Sprintf("商品番号: %s は存在しないため、更新できませんでした。", Product.Id().Value()), err)
		}
		return handler.DBErrHandler(err)
	}
	// Update the fields of upModel as needed
	upModel.ObjID = Product.Id().Value()
	upModel.Name = Product.Name().Value()
	upModel.Price = int(Product.Price().Value())
	upModel.CategoryID = Product.Category().Id().Value()
	if _, updateErr := upModel.Update(ctx, tx, boil.Whitelist(
		models.ProductColumns.ObjID,
		models.ProductColumns.Name,
		models.ProductColumns.Price,
		models.ProductColumns.CategoryID,
	)); updateErr != nil {
		return handler.DBErrHandler(updateErr)
	}
	return nil
}

// DeleteById は指定されたIDの商品をデータベースから削除します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - id: 削除する商品のID
//
// Returns:
//   - error: 商品が存在しない場合はNOT_FOUNDエラー、
//     データベースエラーが発生した場合はそのエラー
func (r *productRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id *products.ProductId) error {
	condition := models.ProductWhere.ObjID.EQ(id.Value())
	delModel, err := models.Products(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("商品番号: %s は存在しないため、削除できませんでした。", id.Value()))
		}
		return handler.DBErrHandler(err)
	}
	if _, deleteErr := delModel.Delete(ctx, tx); deleteErr != nil {
		return handler.DBErrHandler(deleteErr)
	}
	return nil
}

// ProductAfterInsertHook は商品の挿入後に実行されるフックです。
// 新規作成された商品の情報をログに出力します。
//
// Parameters:
//   - ctx: コンテキスト
//   - exec: コンテキスト付きエグゼキューター
//   - product: 挿入された商品
//
// Returns:
//   - error: 常にnilを返します
func ProductAfterInsertHook(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
	log.Printf("商品ID:%s 商品名:%s 単価:%d カテゴリ番号: %s を登録しました。\n",
		product.ObjID, product.Name, product.Price, product.CategoryID)
	return nil
}

// ProductAfterUpdateHook は商品の更新後に実行されるフックです。
// 更新された商品の情報をログに出力します。
//
// Parameters:
//   - ctx: コンテキスト
//   - exec: コンテキスト付きエグゼキューター
//   - product: 更新された商品
//
// Returns:
//   - error: 常にnilを返します
func ProductAfterUpdateHook(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
	log.Printf("商品ID:%s 商品名:%s 単価:%d カテゴリ番号: %s を変更しました。\n",
		product.ObjID, product.Name, product.Price, product.CategoryID)
	return nil
}

// ProductAfterDeleteHook は商品の削除後に実行されるフックです。
// 削除された商品の情報をログに出力します。
//
// Parameters:
//   - ctx: コンテキスト
//   - exec: コンテキスト付きエグゼキューター
//   - product: 削除された商品
//
// Returns:
//   - error: 常にnilを返します
func ProductAfterDeleteHook(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
	log.Printf("商品ID:%s 商品名:%s 単価:%d カテゴリ番号: %s を削除しました。\n",
		product.ObjID, product.Name, product.Price, product.CategoryID)
	return nil
}
