//go:build integration || !ci

// Package impl provides integration tests for the application service layer.
// These tests verify the ProductService implementation including Add, Update, and Delete operations
// using a real database connection. Each test case is independent with automatic cleanup
// to ensure test isolation and prevent data pollution.
package impl

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// generateUniqueProductName は各テストケースで使用するユニークな商品名を生成します。
// 商品名には100文字以下の制約があるため、ミリ秒単位のタイムスタンプを使用します。
//
// Returns:
//   - string: ユニークな商品名（例: "TEST_PRODUCT_1234567890"）
func generateUniqueProductName() string {
	return fmt.Sprintf("TEST_PRODUCT_%d", time.Now().UnixMilli())
}

// cleanupProduct はテスト終了後にテスト用商品をデータベースから削除します。
// この関数はAfterEachフックで呼び出され、テストの独立性を保証します。
// 商品がnilの場合や削除に失敗した場合もエラーを返さず、
// 次のテストに影響を与えないように設計されています。
//
// Parameters:
//   - tm: トランザクションマネージャー
//   - repo: 商品リポジトリ
//   - product: 削除対象の商品（nilの場合は何もしない）
func cleanupProduct(tm service.TransactionManager, repo products.ProductRepository, product *products.Product) {
	if product == nil {
		return
	}

	ctx := context.Background()
	tx, err := tm.Begin(ctx)
	if err != nil {
		return
	}

	exists, err := repo.ExistsById(ctx, tx, product.Id())
	if err != nil {
		_ = tm.Complete(ctx, tx, err)
		return
	}

	if exists {
		err = repo.DeleteById(ctx, tx, product.Id())
	}
	_ = tm.Complete(ctx, tx, err)
}

var _ = Describe("ProductService Integration Test", Ordered, func() {
	var (
		ps           service.ProductService
		tm           service.TransactionManager
		productRepo  products.ProductRepository
		categoryRepo categories.CategoryRepository
		ctx          context.Context
		testCategory *categories.Category
	)

	BeforeAll(func() {
		// データベース接続の初期化（categoryのsetupDatabaseを再利用）
		setupDatabase()

		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// サービスとリポジトリの初期化
		productRepo = repository.NewProductRepositoryImpl(logger)
		categoryRepo = repository.NewCategoryRepositoryImpl(logger)
		tm = repository.NewTransactionManagerImpl(logger)
		ps = NewProductServiceImpl(logger, productRepo, tm)

		// テスト用カテゴリの作成
		ctx = context.Background()
		categoryName, err := categories.NewCategoryName(generateUniqueCategoryName())
		Expect(err).NotTo(HaveOccurred())
		testCategory, err = categories.NewCategory(categoryName)
		Expect(err).NotTo(HaveOccurred())

		// カテゴリをデータベースに追加
		tx, err := tm.Begin(ctx)
		Expect(err).NotTo(HaveOccurred())
		err = categoryRepo.Create(ctx, tx, testCategory)
		Expect(tm.Complete(ctx, tx, err)).To(Succeed())
	})

	AfterAll(func() {
		// テスト用カテゴリのクリーンアップ
		cleanupCategory(tm, categoryRepo, testCategory)
	})

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("Addメソッドの動作確認", func() {
		var testProduct *products.Product

		BeforeEach(func() {
			// 各テストで新しいユニークな商品を作成
			name, err := products.NewProductName(generateUniqueProductName())
			Expect(err).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました")
			price, err := products.NewProductPrice(1000)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました")
			testProduct, err = products.NewProduct(name, price, testCategory)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました")
		})

		AfterEach(func() {
			// テスト終了後のクリーンアップ
			cleanupProduct(tm, productRepo, testProduct)
		})

		It("新しい商品を追加できること", func() {
			err := ps.Add(ctx, testProduct)
			Expect(err).NotTo(HaveOccurred(), "商品の追加に失敗しました")

			// 追加された商品が存在することを確認
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = tm.Complete(ctx, tx, err)
			}()

			exists, err := productRepo.ExistsById(ctx, tx, testProduct.Id())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "追加した商品が存在しません")
		})

		It("既存の商品名で追加しようとするとエラーになること", func() {
			// 最初の追加
			err := ps.Add(ctx, testProduct)
			Expect(err).NotTo(HaveOccurred())

			// 重複追加（同じ名前の新しい商品を作成）
			duplicateName := testProduct.Name()
			duplicatePrice, err := products.NewProductPrice(2000)
			Expect(err).NotTo(HaveOccurred())
			duplicateProduct, err := products.NewProduct(duplicateName, duplicatePrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			err = ps.Add(ctx, duplicateProduct)
			Expect(err).To(HaveOccurred(), "重複する商品名で追加できてしまいました")

			// エラーの詳細を検証
			appErr, ok := err.(*errs.ApplicationError)
			Expect(ok).To(BeTrue(), "返されたエラーがApplicationErrorではありません")
			Expect(appErr.Code).To(Equal("PRODUCT_ALREADY_EXISTS"), "エラーコードが期待値と異なります")
		})
	})

	Context("Updateメソッドの動作確認", func() {
		var testProduct *products.Product

		BeforeEach(func() {
			// テスト用商品を作成して追加
			name, err := products.NewProductName(generateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			price, err := products.NewProductPrice(1000)
			Expect(err).NotTo(HaveOccurred())
			testProduct, err = products.NewProduct(name, price, testCategory)
			Expect(err).NotTo(HaveOccurred())

			err = ps.Add(ctx, testProduct)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品の追加に失敗しました")
		})

		AfterEach(func() {
			// テスト終了後のクリーンアップ
			cleanupProduct(tm, productRepo, testProduct)
		})

		It("既存の商品を更新できること", func() {
			// 更新された商品を作成
			updatedName, err := products.NewProductName(generateUniqueProductName())
			Expect(err).NotTo(HaveOccurred(), "更新用商品名の生成に失敗しました")
			updatedPrice, err := products.NewProductPrice(2000)
			Expect(err).NotTo(HaveOccurred(), "更新用商品価格の生成に失敗しました")
			updatedProduct, err := products.BuildProduct(testProduct.Id(), updatedName, updatedPrice, testCategory)
			Expect(err).NotTo(HaveOccurred(), "更新用商品の生成に失敗しました")

			err = ps.Update(ctx, updatedProduct)
			Expect(err).NotTo(HaveOccurred(), "商品の更新に失敗しました")

			// 更新された商品が存在することを確認
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = tm.Complete(ctx, tx, err)
			}()

			exists, err := productRepo.ExistsByName(ctx, tx, updatedProduct.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "商品名が正しく更新されていません")

			// 古い商品名が存在しないことを確認
			oldExists, err := productRepo.ExistsByName(ctx, tx, testProduct.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(oldExists).To(BeFalse(), "古い商品名が残っています")

			// クリーンアップのために更新された商品を保存
			testProduct = updatedProduct
		})

		It("存在しない商品を更新しようとするとエラーになること", func() {
			// 存在しない商品を作成
			nonExistentName, err := products.NewProductName(generateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentPrice, err := products.NewProductPrice(3000)
			Expect(err).NotTo(HaveOccurred())
			nonExistentProduct, err := products.NewProduct(nonExistentName, nonExistentPrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			err = ps.Update(ctx, nonExistentProduct)
			Expect(err).To(HaveOccurred(), "存在しない商品を更新できてしまいました")

			// エラーの詳細を検証（リポジトリが返すエラー）
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})

	Context("Deleteメソッドの動作確認", func() {
		var testProduct *products.Product

		BeforeEach(func() {
			// テスト用商品を作成して追加
			name, err := products.NewProductName(generateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			price, err := products.NewProductPrice(1000)
			Expect(err).NotTo(HaveOccurred())
			testProduct, err = products.NewProduct(name, price, testCategory)
			Expect(err).NotTo(HaveOccurred())

			err = ps.Add(ctx, testProduct)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品の追加に失敗しました")
		})

		AfterEach(func() {
			// テスト終了後のクリーンアップ（削除失敗時のため）
			cleanupProduct(tm, productRepo, testProduct)
		})

		It("既存の商品を削除できること", func() {
			err := ps.Delete(ctx, testProduct)
			Expect(err).NotTo(HaveOccurred(), "商品の削除に失敗しました")

			// 削除されたことを確認
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = tm.Complete(ctx, tx, err)
			}()

			exists, err := productRepo.ExistsById(ctx, tx, testProduct.Id())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "商品が正しく削除されていません")
		})

		It("存在しない商品を削除しようとするとエラーになること", func() {
			// 存在しない商品を作成
			nonExistentName, err := products.NewProductName(generateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentPrice, err := products.NewProductPrice(3000)
			Expect(err).NotTo(HaveOccurred())
			nonExistentProduct, err := products.NewProduct(nonExistentName, nonExistentPrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			err = ps.Delete(ctx, nonExistentProduct)
			Expect(err).To(HaveOccurred(), "存在しない商品を削除できてしまいました")

			// エラーの詳細を検証（リポジトリが返すエラー）
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})
})
