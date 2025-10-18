//go:build integration || !ci

// Package impl_test provides integration tests for the application service layer.
// These tests verify the CategoryService implementation including Add, Update, and Delete operations
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
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// setupDatabase はテスト用のデータベース接続を初期化します。
// database.tomlファイルを読み込み、テストデータベースへの接続を確立します。
// この関数はBeforeAllフックで呼び出されることを想定しています。
func setupDatabase() {
	v := config.NewViper("../../../", "config")
	config, err := handler.NewDBConfig(v)
	Expect(err).NotTo(HaveOccurred(), "DBConfigの生成に失敗しました")

	_, err = handler.NewDatabase(config)
	Expect(err).NotTo(HaveOccurred(), "データベース接続が失敗したのでテストを中止します")
}

// generateUniqueCategoryName は各テストケースで使用するユニークなカテゴリ名を生成します。
// カテゴリ名には20文字以下の制約があるため、ミリ秒単位のタイムスタンプを使用し、
// 8桁に制限することで"TEST_12345678"のような形式（13文字）を生成します。
//
// Returns:
//   - string: ユニークなカテゴリ名（例: "TEST_12345678"）
func generateUniqueCategoryName() string {
	return fmt.Sprintf("TEST_%d", time.Now().UnixMilli()%100000000)
}

// cleanupCategory はテスト終了後にテスト用カテゴリをデータベースから削除します。
// この関数はAfterEachフックで呼び出され、テストの独立性を保証します。
// カテゴリがnilの場合や削除に失敗した場合もエラーを返さず、
// 次のテストに影響を与えないように設計されています。
//
// Parameters:
//   - tm: トランザクションマネージャー
//   - repo: カテゴリリポジトリ
//   - category: 削除対象のカテゴリ（nilの場合は何もしない）
func cleanupCategory(tm service.TransactionManager, repo categories.CategoryRepository, category *categories.Category) {
	if category == nil {
		return
	}

	ctx := context.Background()
	tx, err := tm.Begin(ctx)
	if err != nil {
		return
	}

	exists, err := repo.ExistsByName(ctx, tx, category.Name())
	if err != nil {
		_ = tm.Complete(ctx, tx, err)
		return
	}

	if exists {
		err = repo.DeleteByName(ctx, tx, category.Name())
	}
	_ = tm.Complete(ctx, tx, err)
}

var _ = Describe("CategoryService Integration Test", Ordered, func() {
	var (
		cs   service.CategoryService
		tm   service.TransactionManager
		repo categories.CategoryRepository
		ctx  context.Context
	)

	BeforeAll(func() {
		// データベース接続の初期化
		setupDatabase()

		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// サービスとリポジトリの初期化
		repo = repository.NewCategoryRepositoryImpl(logger)
		tm = repository.NewTransactionManagerImpl(logger)
		cs = NewCategoryServiceImpl(logger, repo, tm)
	})

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("Addメソッドの動作確認", func() {
		var testCategory *categories.Category

		BeforeEach(func() {
			// 各テストで新しいユニークなカテゴリを作成
			name, err := categories.NewCategoryName(generateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました")
			testCategory, err = categories.NewCategory(name)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました")
		})

		AfterEach(func() {
			// テスト終了後のクリーンアップ
			cleanupCategory(tm, repo, testCategory)
		})

		It("新しいカテゴリを追加できること", func() {
			err := cs.Add(ctx, testCategory)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの追加に失敗しました")

			// 追加されたカテゴリが存在することを確認
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = tm.Complete(ctx, tx, err)
			}()

			exists, err := repo.ExistsByName(ctx, tx, testCategory.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "追加したカテゴリが存在しません")
		})

		It("既存のカテゴリ名で追加しようとするとエラーになること", func() {
			// 最初の追加
			err := cs.Add(ctx, testCategory)
			Expect(err).NotTo(HaveOccurred())

			// 重複追加
			err = cs.Add(ctx, testCategory)
			Expect(err).To(HaveOccurred(), "重複するカテゴリ名で追加できてしまいました")

			// エラーの詳細を検証
			appErr, ok := err.(*errs.ApplicationError)
			Expect(ok).To(BeTrue(), "返されたエラーがApplicationErrorではありません")
			Expect(appErr.Code).To(Equal("CATEGORY_ALREADY_EXISTS"), "エラーコードが期待値と異なります")
		})
	})

	Context("Updateメソッドの動作確認", func() {
		var testCategory *categories.Category

		BeforeEach(func() {
			// テスト用カテゴリを作成して追加
			name, err := categories.NewCategoryName(generateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = categories.NewCategory(name)
			Expect(err).NotTo(HaveOccurred())

			err = cs.Add(ctx, testCategory)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの追加に失敗しました")
		})

		AfterEach(func() {
			// テスト終了後のクリーンアップ
			cleanupCategory(tm, repo, testCategory)
		})

		It("既存のカテゴリを更新できること", func() {
			// 更新されたカテゴリを作成
			updatedName, err := categories.NewCategoryName(generateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred(), "更新用カテゴリ名の生成に失敗しました")
			updatedCategory, err := categories.BuildCategory(testCategory.Id(), updatedName)
			Expect(err).NotTo(HaveOccurred(), "更新用カテゴリの生成に失敗しました")

			err = cs.Update(ctx, updatedCategory)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの更新に失敗しました")

			// 更新されたカテゴリが存在することを確認
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = tm.Complete(ctx, tx, err)
			}()

			exists, err := repo.ExistsByName(ctx, tx, updatedCategory.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "カテゴリ名が正しく更新されていません")

			// 古いカテゴリ名が存在しないことを確認
			oldExists, err := repo.ExistsByName(ctx, tx, testCategory.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(oldExists).To(BeFalse(), "古いカテゴリ名が残っています")

			// クリーンアップのために更新されたカテゴリを保存
			testCategory = updatedCategory
		})

		It("存在しないカテゴリを更新しようとするとエラーになること", func() {
			// 存在しないカテゴリを作成
			nonExistentName, err := categories.NewCategoryName(generateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentCategory, err := categories.NewCategory(nonExistentName)
			Expect(err).NotTo(HaveOccurred())

			err = cs.Update(ctx, nonExistentCategory)
			Expect(err).To(HaveOccurred(), "存在しないカテゴリを更新できてしまいました")

			// エラーの詳細を検証
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})

	Context("Deleteメソッドの動作確認", func() {
		var testCategory *categories.Category

		BeforeEach(func() {
			// テスト用カテゴリを作成して追加
			name, err := categories.NewCategoryName(generateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = categories.NewCategory(name)
			Expect(err).NotTo(HaveOccurred())

			err = cs.Add(ctx, testCategory)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの追加に失敗しました")
		})

		AfterEach(func() {
			// テスト終了後のクリーンアップ（削除失敗時のため）
			cleanupCategory(tm, repo, testCategory)
		})

		It("既存のカテゴリを削除できること", func() {
			err := cs.Delete(ctx, testCategory)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの削除に失敗しました")

			// 削除されたことを確認
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = tm.Complete(ctx, tx, err)
			}()

			exists, err := repo.ExistsByName(ctx, tx, testCategory.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "カテゴリが正しく削除されていません")
		})

		It("存在しないカテゴリを削除しようとするとエラーになること", func() {
			// 存在しないカテゴリを作成
			nonExistentName, err := categories.NewCategoryName(generateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentCategory, err := categories.NewCategory(nonExistentName)
			Expect(err).NotTo(HaveOccurred())

			err = cs.Delete(ctx, nonExistentCategory)
			Expect(err).To(HaveOccurred(), "存在しないカテゴリを削除できてしまいました")

			// エラーの詳細を検証
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})
})
