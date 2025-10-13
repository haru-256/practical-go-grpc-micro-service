package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/models"
)

// categoryRepositoryImpl はカテゴリリポジトリのSQLBoilerを使用した実装です。
type categoryRepositoryImpl struct{}

// NewCategoryRepository は新しいCategoryRepositoryインスタンスを生成します。
// この関数は、カテゴリの挿入、更新、削除後に実行されるフックを登録します。
func NewCategoryRepository() categories.CategoryRepository {
	models.AddCategoryHook(boil.AfterInsertHook, CategoryAfterInsertHook)
	models.AddCategoryHook(boil.AfterUpdateHook, CategoryAfterUpdateHook)
	models.AddCategoryHook(boil.AfterDeleteHook, CategoryAfterDeleteHook)
	return &categoryRepositoryImpl{}
}

// ExistsByName は指定された名前のカテゴリが存在するかどうかをチェックします。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - name: チェック対象のカテゴリ名
//
// Returns:
//   - bool: カテゴリが存在する場合はtrue、存在しない場合はfalse
//   - error: データベースエラーが発生した場合
func (r *categoryRepositoryImpl) ExistsByName(ctx context.Context, tx *sql.Tx, name *categories.CategoryName) (bool, error) {
	condition := models.CategoryWhere.Name.EQ(name.Value())
	exists, err := models.Categories(condition).Exists(ctx, tx)
	if err != nil {
		return false, handler.DBErrHandler(err)
	}
	return exists, nil
}

// Create は新しいカテゴリをデータベースに追加します。
//
// NOTE: boil.Infer() を使用することで、auto-incrementのIDは無視され、
// DB側で自動採番された後、SQLBoiler側の構造体にセットされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - category: 作成するカテゴリ
//
// Returns:
//   - error: データベースエラーが発生した場合
func (r *categoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category *categories.Category) error {
	newCategory := models.Category{
		ObjID: category.Id().Value(),
		Name:  category.Name().Value(),
	}
	// NOTE: boil.Infer() でauto-incrementのIDは無視され、勝手にDB側で採番された後、sqlboiler側の構造体にセットされる
	if err := newCategory.Insert(ctx, tx, boil.Infer()); err != nil {
		return handler.DBErrHandler(err)
	}
	return nil
}

// UpdateById は指定されたIDのカテゴリ情報を更新します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - category: 更新するカテゴリ情報
//
// Returns:
//   - error: カテゴリが存在しない場合はNOT_FOUNDエラー、
//     データベースエラーが発生した場合はそのエラー
func (r *categoryRepositoryImpl) UpdateById(ctx context.Context, tx *sql.Tx, category *categories.Category) error {
	condition := models.CategoryWhere.ObjID.EQ(category.Id().Value())
	upModel, err := models.Categories(condition).One(ctx, tx)
	if upModel == nil {
		return errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("カテゴリ番号: %s は存在しないため、更新できませんでした。", category.Id().Value()))
	}
	if err != nil {
		return handler.DBErrHandler(err)
	}
	// Update the fields of upModel as needed
	upModel.ObjID = category.Id().Value()
	upModel.Name = category.Name().Value()
	if _, err := upModel.Update(ctx, tx, boil.Whitelist(models.CategoryColumns.ObjID, models.CategoryColumns.Name)); err != nil {
		return handler.DBErrHandler(err)
	}
	return nil
}

// DeleteById は指定されたIDのカテゴリをデータベースから削除します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - id: 削除するカテゴリのID
//
// Returns:
//   - error: カテゴリが存在しない場合はNOT_FOUNDエラー、
//     データベースエラーが発生した場合はそのエラー
func (r *categoryRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id *categories.CategoryId) error {
	condition := models.CategoryWhere.ObjID.EQ(id.Value())
	delModel, err := models.Categories(condition).One(ctx, tx)
	if delModel == nil {
		return errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("カテゴリ番号: %s は存在しないため、削除できませんでした。", id.Value()))
	}
	if err != nil {
		return handler.DBErrHandler(err)
	}
	if _, err := delModel.Delete(ctx, tx); err != nil {
		return handler.DBErrHandler(err)
	}
	return nil
}

// CategoryAfterInsertHook はカテゴリの挿入後に実行されるフックです。
// 新規作成されたカテゴリの情報をログに出力します。
//
// Parameters:
//   - ctx: コンテキスト
//   - exec: コンテキスト付きエグゼキューター
//   - category: 挿入されたカテゴリ
//
// Returns:
//   - error: 常にnilを返します
func CategoryAfterInsertHook(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
	log.Printf("カテゴリが新規作成されました。ID=%s, Name=%s", category.ObjID, category.Name)
	return nil
}

// CategoryAfterUpdateHook はカテゴリの更新後に実行されるフックです。
// 更新されたカテゴリの情報をログに出力します。
//
// Parameters:
//   - ctx: コンテキスト
//   - exec: コンテキスト付きエグゼキューター
//   - category: 更新されたカテゴリ
//
// Returns:
//   - error: 常にnilを返します
func CategoryAfterUpdateHook(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
	log.Printf("カテゴリが更新されました。ID=%s, Name=%s", category.ObjID, category.Name)
	return nil
}

// CategoryAfterDeleteHook はカテゴリの削除後に実行されるフックです。
// 削除されたカテゴリの情報をログに出力します。
//
// Parameters:
//   - ctx: コンテキスト
//   - exec: コンテキスト付きエグゼキューター
//   - category: 削除されたカテゴリ
//
// Returns:
//   - error: 常にnilを返します
func CategoryAfterDeleteHook(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
	log.Printf("カテゴリが削除されました。ID=%s, Name=%s", category.ObjID, category.Name)
	return nil
}
