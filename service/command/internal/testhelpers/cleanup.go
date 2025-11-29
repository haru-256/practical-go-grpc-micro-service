package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/onsi/ginkgo/v2"
)

// CleanupProduct はテスト終了後にテスト用商品をデータベースから削除します。
// この関数はAfterEachフックで呼び出され、テストの独立性を保証します。
// 商品がnilの場合や削除に失敗した場合もエラーを返さず、
// 次のテストに影響を与えないように設計されています。
//
// Parameters:
//   - tm: トランザクションマネージャー
//   - repo: 商品リポジトリ
//   - product: 削除対象の商品（nilの場合は何もしない）
func CleanupProduct(tm service.TransactionManager, repo products.ProductRepository, productDTO *dto.ProductDTO) (err error) {
	var (
		product *products.Product
		tx      *sql.Tx
	)

	if productDTO == nil {
		return nil
	}

	product, err = dto.ProductFromDTO(productDTO)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err = tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if completeErr := tm.Complete(ctx, tx, err); err == nil {
			err = completeErr
		}
	}()

	exists, err := repo.ExistsById(ctx, tx, product.Id())
	if err != nil || !exists {
		ginkgo.GinkgoWriter.Printf("⚠️  Cleanup: 商品存在確認失敗: %v\n", product.Id().Value())
		return nil
	}

	err = repo.DeleteById(ctx, tx, product.Id())
	if err != nil {
		return err
	}
	return nil
}

// CleanupCategory はテスト終了後にテスト用カテゴリをデータベースから削除します。
// この関数はAfterEachフックで呼び出され、テストの独立性を保証します。
// カテゴリがnilの場合や削除に失敗した場合もエラーを返さず、
// 次のテストに影響を与えないように設計されています。
//
// Parameters:
//   - tm: トランザクションマネージャー
//   - repo: カテゴリリポジトリ
//   - category: 削除対象のカテゴリ（nilの場合は何もしない）
func CleanupCategory(tm service.TransactionManager, repo categories.CategoryRepository, categoryDTO *dto.CategoryDTO) (err error) {
	var (
		category       *categories.Category
		categoryResult *categories.Category
		tx             *sql.Tx
	)

	if categoryDTO == nil {
		return nil
	}

	category, err = dto.CategoryFromDTO(categoryDTO)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err = tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if completeErr := tm.Complete(ctx, tx, err); err == nil {
			err = completeErr
		}
	}()

	categoryResult, err = repo.FindById(ctx, tx, category.Id())
	if err != nil || categoryResult == nil {
		ginkgo.GinkgoWriter.Printf("⚠️  Cleanup: カテゴリ存在確認失敗: %v\n", category.Id().Value())
		return nil
	}

	err = repo.DeleteById(ctx, tx, category.Id())
	if err != nil {
		return err
	}
	return nil
}

// CleanupProductCategory は商品とカテゴリを外部キー制約を考慮して順序良く削除します。
// 商品を先に削除してからカテゴリを削除することで、外部キー制約エラーを回避します。
// この関数を使用することで、個別にクリーンアップ関数を呼び出すよりも安全で確実な削除が可能です。
//
// Parameters:
//   - tm: トランザクションマネージャー
//   - productRepo: 商品リポジトリ
//   - categoryRepo: カテゴリリポジトリ
//   - productDTO: 削除対象の商品（nilの場合は商品削除をスキップ）
//   - categoryDTO: 削除対象のカテゴリ（nilの場合はカテゴリ削除をスキップ）
//
// Returns:
//   - error: 削除処理でエラーが発生した場合
func CleanupProductCategory(tm service.TransactionManager, productRepo products.ProductRepository, categoryRepo categories.CategoryRepository, productDTO *dto.ProductDTO, categoryDTO *dto.CategoryDTO) error {
	err := CleanupProduct(tm, productRepo, productDTO)
	if err != nil {
		return err
	}
	err = CleanupCategory(tm, categoryRepo, categoryDTO)
	if err != nil {
		return err
	}
	return nil
}

const (
	// testCategoryNameMaxDigits はユニークなカテゴリ名生成時の最大桁数
	// カテゴリ名の制限（20文字）に対して "TEST_" (5文字) + 8桁 = 13文字
	testCategoryNameMaxDigits = 100000000
	// testProductNameMaxDigits はユニークな商品名生成時の最大桁数
	// 商品名の制限（20文字）に対して "TEST_PRODUCT_" (14文字) + 8桁 = 22文字
	testProductNameMaxDigits = 100000000
)

// GenerateUniqueCategoryName は各テストケースで使用するユニークなカテゴリ名を生成します。
// カテゴリ名には100文字以下の制約があるため、ミリ秒単位のタイムスタンプを使用します。
//
// Returns:
//   - string: "TEST_" + ミリ秒タイムスタンプの形式のユニークなカテゴリ名
func GenerateUniqueCategoryName() string {
	return fmt.Sprintf("TEST_%d", time.Now().UnixMilli()%testCategoryNameMaxDigits)
}

// GenerateUniqueProductName は各テストケースで使用するユニークな商品名を生成します。
// 商品名には100文字以下の制約があるため、ミリ秒単位のタイムスタンプを使用します。
//
// Returns:
//   - string: "TEST_PRODUCT_" + ミリ秒タイムスタンプの形式のユニークな商品名
func GenerateUniqueProductName() string {
	return fmt.Sprintf("TEST_PRODUCT_%d", time.Now().UnixMilli()%testProductNameMaxDigits)
}
