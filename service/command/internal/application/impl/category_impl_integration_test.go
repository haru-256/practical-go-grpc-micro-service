//go:build integration || !ci

// Package impl_test provides integration tests for the application service layer.
// These tests verify the CategoryService implementation including Add, Update, and Delete operations
// using a real database connection. Each test case is independent with automatic cleanup
// to ensure test isolation and prevent data pollution.
package impl

import (
	"context"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/testhelpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CategoryService Integration Test", Label("IntegrationTests"), Label("CategoryService"), Ordered, func() {
	var (
		cs   service.CategoryService
		tm   service.TransactionManager
		repo categories.CategoryRepository
		ctx  context.Context
	)

	BeforeAll(func() {
		// データベース接続の初期化
		err := testhelpers.SetupDatabase("../../../", "config")
		Expect(err).NotTo(HaveOccurred())

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
		var testCategoryName *categories.CategoryName

		BeforeEach(func() {
			// 各テストで新しいユニークなカテゴリを作成
			var err error
			testCategoryName, err = categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました")
		})

		It("新しいカテゴリを追加できること", func() {
			createDTO := &dto.CreateCategoryDTO{
				Name: testCategoryName.Value(),
			}
			result, err := cs.Add(ctx, createDTO)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの追加に失敗しました")
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(testCategoryName.Value()))
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, result)
			testCategory, err := dto.CategoryFromDTO(result)
			Expect(err).NotTo(HaveOccurred(), "カテゴリDTOからドメインモデルへの変換に失敗しました")

			// 追加されたカテゴリが存在することを確認
			exists, err := testhelpers.VerifyCategoryById(ctx, tm, repo, testCategory.Id())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "作成したカテゴリがDBに存在しません")
		})

		It("既存のカテゴリ名で追加しようとするとエラーになること", func() {
			createDTO := &dto.CreateCategoryDTO{
				Name: testCategoryName.Value(),
			}
			// 最初の追加
			result, err := cs.Add(ctx, createDTO)
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, result)

			// 重複追加
			duplicateResult, err := cs.Add(ctx, createDTO)
			Expect(err).To(HaveOccurred(), "重複するカテゴリ名で追加できてしまいました")
			Expect(duplicateResult).To(BeNil())

			// エラーの詳細を検証
			appErr, ok := err.(*errs.ApplicationError)
			Expect(ok).To(BeTrue(), "返されたエラーがApplicationErrorではありません")
			Expect(appErr.Code).To(Equal("CATEGORY_ALREADY_EXISTS"), "エラーコードが期待値と異なります")
		})
	})

	Context("Updateメソッドの動作確認", func() {
		var testCategory *categories.Category
		var testCategoryDTO *dto.CategoryDTO

		BeforeEach(func() {
			var err error
			// テスト用カテゴリを作成して追加
			createDTO := &dto.CreateCategoryDTO{
				Name: testhelpers.GenerateUniqueCategoryName(),
			}
			testCategoryDTO, err = cs.Add(ctx, createDTO)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの追加に失敗しました")
			testCategory, err = dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリDTOからドメインモデルへの変換に失敗しました")
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, testCategoryDTO)
		})

		It("既存のカテゴリを更新できること", func() {
			// 更新されたカテゴリを作成
			updatedName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred(), "更新用カテゴリ名の生成に失敗しました")

			updateDTO := &dto.UpdateCategoryDTO{
				Id:   testCategory.Id().Value(),
				Name: updatedName.Value(),
			}
			result, err := cs.Update(ctx, updateDTO)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの更新に失敗しました")
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(updatedName.Value()))

			// 更新されたカテゴリが存在することを確認
			exists, err := testhelpers.VerifyCategoryById(ctx, tm, repo, testCategory.Id())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新したカテゴリIDがDBに存在しません")
			exists, err = testhelpers.VerifyCategoryByName(ctx, tm, repo, updatedName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新したカテゴリ名がDBに存在しません")
			exists, err = testhelpers.VerifyCategoryByName(ctx, tm, repo, testCategory.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "古いカテゴリ名がDBに存在しています")
		})

		It("存在しないカテゴリを更新しようとするとエラーになること", func() {
			// 存在しないカテゴリを作成
			nonExistentName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentCategory, err := categories.NewCategory(nonExistentName)
			Expect(err).NotTo(HaveOccurred())

			updateDTO := &dto.UpdateCategoryDTO{
				Id:   nonExistentCategory.Id().Value(),
				Name: nonExistentCategory.Name().Value(),
			}
			result, err := cs.Update(ctx, updateDTO)
			Expect(err).To(HaveOccurred(), "存在しないカテゴリを更新できてしまいました")
			Expect(result).To(BeNil())

			// エラーの詳細を検証
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})

	Context("Deleteメソッドの動作確認", func() {
		var testCategory *categories.Category
		var testCategoryDTO *dto.CategoryDTO

		BeforeEach(func() {
			// テスト用カテゴリを作成して追加
			name, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = categories.NewCategory(name)
			Expect(err).NotTo(HaveOccurred())

			createDTO := &dto.CreateCategoryDTO{
				Name: testCategory.Name().Value(),
			}
			testCategoryDTO, err = cs.Add(ctx, createDTO)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの追加に失敗しました")
			testCategory, err = dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリDTOからドメインモデルへの変換に失敗しました")
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, testCategoryDTO)
		})

		It("既存のカテゴリを削除できること", func() {
			deleteDTO := &dto.DeleteCategoryDTO{
				Id: testCategory.Id().Value(),
			}
			result, err := cs.Delete(ctx, deleteDTO)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの削除に失敗しました")
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(testCategory.Name().Value()))

			// 削除されたことを確認
			exists, _ := testhelpers.VerifyCategoryById(ctx, tm, repo, testCategory.Id())
			Expect(exists).To(BeFalse(), "削除したカテゴリがDBに存在します")
		})

		It("存在しないカテゴリを削除しようとするとエラーになること", func() {
			// 存在しないカテゴリを作成
			nonExistentName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentCategory, err := categories.NewCategory(nonExistentName)
			Expect(err).NotTo(HaveOccurred())

			deleteDTO := &dto.DeleteCategoryDTO{
				Id: nonExistentCategory.Id().Value(),
			}
			result, err := cs.Delete(ctx, deleteDTO)
			Expect(err).To(HaveOccurred(), "存在しないカテゴリを削除できてしまいました")
			Expect(result).To(BeNil())

			// エラーの詳細を検証
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})
})
