package testhelpers

import (
	"context"
	"database/sql"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
)

// SetupDatabase はテスト用のデータベース接続を初期化します。
// 指定された設定ファイルを読み込み、テストデータベースへの接続を確立します。
// この関数は統合テストのBeforeAllフックで呼び出されることを想定しています。
//
// Parameters:
//   - configPath: 設定ファイルが配置されているディレクトリパス
//   - configName: 設定ファイル名（拡張子なし）
//
// Returns:
//   - error: データベース接続の初期化でエラーが発生した場合
//
// Usage:
//
//	err := testhelpers.SetupDatabase("../../../config", "database")
func SetupDatabase(configPath string, configName string) error {
	v := config.NewViper(configPath, configName)
	config, err := handler.NewDBConfig(v)
	if err != nil {
		return err
	}
	_, err = handler.NewDatabase(config)
	if err != nil {
		return err
	}
	return nil
}

// VerifyCategoryById はカテゴリがデータベースに存在するかを検証するヘルパー関数です。
// テスト内でカテゴリの存在確認を行う際に使用します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tm: トランザクションマネージャー
//   - repo: カテゴリリポジトリ
//   - id: 検証対象のカテゴリID
//
// Returns:
//   - exists: カテゴリが存在するかどうか
//   - err: データベースアクセスでエラーが発生した場合
func VerifyCategoryById(ctx context.Context, tm service.TransactionManager, repo categories.CategoryRepository, id *categories.CategoryId) (exists bool, err error) {
	var (
		tx       *sql.Tx
		category *categories.Category
	)
	tx, err = tm.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if compErr := tm.Complete(ctx, tx, err); err == nil {
			err = compErr
		}
	}()

	category, err = repo.FindById(ctx, tx, id)
	if err != nil {
		return false, err
	}
	if category != nil {
		return true, nil
	}
	return false, nil
}

// VerifyCategoryByName はカテゴリがデータベースに存在するかを名前で検証するヘルパー関数です。
// テスト内でカテゴリの存在確認を行う際に使用します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tm: トランザクションマネージャー
//   - repo: カテゴリリポジトリ
//   - name: 検証対象のカテゴリ名
//
// Returns:
//   - exists: カテゴリが存在するかどうか
//   - err: データベースアクセスでエラーが発生した場合
func VerifyCategoryByName(ctx context.Context, tm service.TransactionManager, repo categories.CategoryRepository, name *categories.CategoryName) (exists bool, err error) {
	var tx *sql.Tx
	tx, err = tm.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if compErr := tm.Complete(ctx, tx, err); err == nil {
			err = compErr
		}
	}()

	exists, err = repo.ExistsByName(ctx, tx, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// VerifyProductById は商品がDBに存在することを検証するヘルパー関数です。
//
// Parameters:
//   - ctx: コンテキスト
//   - tm: トランザクションマネージャー
//   - repo: 商品リポジトリ
//   - id: 検証対象の商品ID
//
// Returns:
//   - exists: 商品が存在するかどうか
//   - err: エラー（存在しない場合はerrはnilでexistsがfalse）
func VerifyProductById(ctx context.Context, tm service.TransactionManager, repo products.ProductRepository, id *products.ProductId) (exists bool, err error) {
	var tx *sql.Tx
	tx, err = tm.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if compErr := tm.Complete(ctx, tx, err); err == nil {
			err = compErr
		}
	}()

	exists, err = repo.ExistsById(ctx, tx, id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// VerifyProductByName は商品がDBに存在することを検証するヘルパー関数です。
//
// Parameters:
//   - ctx: コンテキスト
//   - tm: トランザクションマネージャー
//   - repo: 商品リポジトリ
//   - name: 検証対象の商品名
//
// Returns:
//   - exists: 商品が存在するかどうか
//   - err: エラー（存在しない場合はerrはnilでexistsがfalse）
func VerifyProductByName(ctx context.Context, tm service.TransactionManager, repo products.ProductRepository, name *products.ProductName) (exists bool, err error) {
	var tx *sql.Tx
	tx, err = tm.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if compErr := tm.Complete(ctx, tx, err); err == nil {
			err = compErr
		}
	}()

	exists, err = repo.ExistsByName(ctx, tx, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}
