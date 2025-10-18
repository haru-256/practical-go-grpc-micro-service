package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/models"
)

// CategoryRepositoryImpl はカテゴリリポジトリのSQLBoilerを使用した実装です。
type CategoryRepositoryImpl struct {
	logger *slog.Logger
}

// NewCategoryRepositoryImpl は新しいCategoryRepositoryImplインスタンスを生成します。
// この関数は、カテゴリの挿入、更新、削除後に実行されるフックを登録します。
// フック関数はファクトリー関数により生成され、渡されたloggerをクロージャーに保持します。
// 具象型を返すことで、呼び出し側が必要に応じてインターフェースとして扱えるようにします。
//
// 使用例:
//
//	repo := repository.NewCategoryRepositoryImpl(logger)
//	var categoryRepo categories.CategoryRepository = repo  // インターフェースとして使用
func NewCategoryRepositoryImpl(logger *slog.Logger) *CategoryRepositoryImpl {
	models.AddCategoryHook(boil.AfterInsertHook, categoryHookFactory(boil.AfterInsertHook, logger))
	models.AddCategoryHook(boil.AfterUpdateHook, categoryHookFactory(boil.AfterUpdateHook, logger))
	models.AddCategoryHook(boil.AfterDeleteHook, categoryHookFactory(boil.AfterDeleteHook, logger))
	return &CategoryRepositoryImpl{logger: logger}
}

// ExistsByName は指定されたカテゴリ名が既に存在するかをチェックします。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - name: チェックするカテゴリ名
//
// Returns:
//   - bool: カテゴリが存在する場合はtrue、存在しない場合はfalse
//   - error: データベースエラーが発生した場合
func (r *CategoryRepositoryImpl) ExistsByName(ctx context.Context, tx *sql.Tx, name *categories.CategoryName) (bool, error) {
	condition := models.CategoryWhere.Name.EQ(name.Value())
	exists, err := models.Categories(condition).Exists(ctx, tx)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if category exists", slog.Any("error", err))
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
func (r *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category *categories.Category) error {
	newCategory := models.Category{
		ObjID: category.Id().Value(),
		Name:  category.Name().Value(),
	}
	// NOTE: boil.Infer() でauto-incrementのIDは無視され、勝手にDB側で採番された後、sqlboiler側の構造体にセットされる
	if err := newCategory.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.ErrorContext(ctx, "Failed to create category", slog.Any("error", err))
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
func (r *CategoryRepositoryImpl) UpdateById(ctx context.Context, tx *sql.Tx, category *categories.Category) error {
	condition := models.CategoryWhere.ObjID.EQ(category.Id().Value())
	upModel, err := models.Categories(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDErrorWithCause("NOT_FOUND", fmt.Sprintf("カテゴリ番号: %s は存在しないため、更新できませんでした。", category.Id().Value()), err)
		}
		return handler.DBErrHandler(err)
	}
	// Update the fields of upModel as needed
	upModel.ObjID = category.Id().Value()
	upModel.Name = category.Name().Value()
	if _, updateErr := upModel.Update(ctx, tx, boil.Whitelist(models.CategoryColumns.ObjID, models.CategoryColumns.Name)); updateErr != nil {
		return handler.DBErrHandler(updateErr)
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
func (r *CategoryRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id *categories.CategoryId) error {
	condition := models.CategoryWhere.ObjID.EQ(id.Value())
	delModel, err := models.Categories(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("カテゴリ番号: %s は存在しないため、削除できませんでした。", id.Value()))
		}
		return handler.DBErrHandler(err)
	}
	if _, deleteErr := delModel.Delete(ctx, tx); deleteErr != nil {
		return handler.DBErrHandler(deleteErr)
	}
	return nil
}

// DeleteByName は指定されたカテゴリ名のカテゴリをデータベースから削除します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - name: 削除するカテゴリ名
//
// Returns:
//   - error: カテゴリが存在しない場合はNOT_FOUNDエラー、
//     データベースエラーが発生した場合はそのエラー
func (r *CategoryRepositoryImpl) DeleteByName(ctx context.Context, tx *sql.Tx, name *categories.CategoryName) error {
	condition := models.CategoryWhere.Name.EQ(name.Value())
	delModel, err := models.Categories(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("カテゴリ名: %s は存在しないため、削除できませんでした。", name.Value()))
		}
		return handler.DBErrHandler(err)
	}
	if _, deleteErr := delModel.Delete(ctx, tx); deleteErr != nil {
		return handler.DBErrHandler(deleteErr)
	}
	return nil
}

// categoryHookFactory は指定されたフックタイプに応じたカテゴリ用フック関数を生成します。
// loggerをクロージャーに保持し、各フック実行時に構造化ログを出力します。
//
// Parameters:
//   - hookType: フックのタイプ（AfterInsertHook, AfterUpdateHook, AfterDeleteHook）
//   - logger: 構造化ログ出力に使用するlogger
//
// Returns:
//   - func: SQLBoilerのフック関数
//
// Panics:
//   - hookTypeが想定外の値の場合
func categoryHookFactory(hookType boil.HookPoint, logger *slog.Logger) func(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
	switch hookType {
	case boil.AfterInsertHook:
		return func(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
			logger.InfoContext(ctx, "カテゴリが新規作成されました",
				slog.String("obj_id", category.ObjID),
				slog.String("name", category.Name))
			return nil
		}
	case boil.AfterUpdateHook:
		return func(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
			logger.InfoContext(ctx, "カテゴリが更新されました",
				slog.String("obj_id", category.ObjID),
				slog.String("name", category.Name))
			return nil
		}
	case boil.AfterDeleteHook:
		return func(ctx context.Context, exec boil.ContextExecutor, category *models.Category) error {
			logger.InfoContext(ctx, "カテゴリが削除されました",
				slog.String("obj_id", category.ObjID),
				slog.String("name", category.Name))
			return nil
		}
	default:
		panic("Invalid hookType")
	}
}
