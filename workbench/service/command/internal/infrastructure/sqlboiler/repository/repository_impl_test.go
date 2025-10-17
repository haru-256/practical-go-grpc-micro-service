//go:build integration || !ci

// Package repository_test provides integration tests for repository implementations.
// These tests verify the CategoryRepository and ProductRepository implementations
// using a real database connection with transaction rollback for test isolation.
package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("categoryRepositoryImpl構造体", Ordered, Label("CategoryRepositoryインターフェースメソッドのテスト"), func() {
	var rep categories.CategoryRepository
	var ctx context.Context
	var tx *sql.Tx
	var err error

	BeforeAll(func() {
		rep = NewCategoryRepository()
	})

	BeforeEach(func() {
		ctx = context.Background()
		tx, err = boil.BeginTx(ctx, nil)
		Expect(err).NotTo(HaveOccurred(), "トランザクションの開始に失敗しました。")
	})

	AfterEach(func() {
		err = tx.Rollback()
		Expect(err).NotTo(HaveOccurred(), "トランザクションのロールバックに失敗しました。")
	})

	DescribeTable(
		"ExistsByNameの動作確認",
		func(name string, expected bool) {
			catName, cateNameErr := categories.NewCategoryName(name)
			Expect(cateNameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, cateErr := categories.NewCategory(catName)
			Expect(cateErr).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")
			result, existsErr := rep.ExistsByName(ctx, tx, category.Name())
			Expect(existsErr).NotTo(HaveOccurred(), "ExistsByNameの実行に失敗しました。")
			Expect(result).To(Equal(expected), "存在するカテゴリ名に対してExistsByNameがfalseを返しました。")
		},
		Entry("文房具", "文房具", true),
		Entry("食品", "食品", false),
	)

	Context("Createの動作確認", func() {
		It("新しいカテゴリを作成できること", func() {
			name, nameErr := categories.NewCategoryName("食品")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, categoryErr := categories.NewCategory(name)
			Expect(categoryErr).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")

			createErr := rep.Create(ctx, tx, category)
			Expect(createErr).NotTo(HaveOccurred(), "カテゴリの作成に失敗しました。")
		})
		It("重複するobj_idでエラーになること", func() {
			// 文房具のobj_idを指定してカテゴリを作成する
			id, idErr := categories.NewCategoryId("b1524011-b6af-417e-8bf2-f449dd58b5c0")
			Expect(idErr).NotTo(HaveOccurred(), "テスト用カテゴリIDの生成に失敗しました。")
			name, nameErr := categories.NewCategoryName("新しい文房具")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, categoryErr := categories.BuildCategory(id, name)
			Expect(categoryErr).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")

			createErr := rep.Create(ctx, tx, category)
			Expect(createErr).To(HaveOccurred())
			crudErr, ok := createErr.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("DB_UNIQUE_CONSTRAINT_VIOLATION"))
			Expect(crudErr.Message).To(ContainSubstring("一意制約違反です。"))
		})
	})

	Context("DeleteByIdの動作確認", func() {
		It("既存のカテゴリを削除できること", func() {
			// まず新しいカテゴリを作成
			name, nameErr := categories.NewCategoryName("ID削除テストカテゴリ")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, categoryErr := categories.NewCategory(name)
			Expect(categoryErr).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")

			createErr := rep.Create(ctx, tx, category)
			Expect(createErr).NotTo(HaveOccurred(), "カテゴリの作成に失敗しました。")

			// 作成したカテゴリをIDで削除
			deleteErr := rep.DeleteById(ctx, tx, category.Id())
			Expect(deleteErr).NotTo(HaveOccurred(), "カテゴリの削除に失敗しました。")

			// 削除されたカテゴリが存在しないことを確認
			exists, existsErr := rep.ExistsByName(ctx, tx, category.Name())
			Expect(existsErr).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "削除したカテゴリがまだ存在しています。")
		})

		It("存在しないカテゴリIDで削除しようとするとエラーになること", func() {
			id, idErr := categories.NewCategoryId("00000000-0000-0000-0000-000000000000")
			Expect(idErr).NotTo(HaveOccurred(), "テスト用カテゴリIDの生成に失敗しました。")

			deleteErr := rep.DeleteById(ctx, tx, id)
			Expect(deleteErr).To(HaveOccurred())
			crudErr, ok := deleteErr.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("NOT_FOUND"))
			Expect(crudErr.Message).To(ContainSubstring("存在しないため、削除できませんでした。"))
		})
	})

	Context("DeleteByNameの動作確認", func() {
		It("既存のカテゴリを削除できること", func() {
			// まず新しいカテゴリを作成
			name, nameErr := categories.NewCategoryName("削除テストカテゴリ")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, categoryErr := categories.NewCategory(name)
			Expect(categoryErr).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")

			createErr := rep.Create(ctx, tx, category)
			Expect(createErr).NotTo(HaveOccurred(), "カテゴリの作成に失敗しました。")

			// 作成したカテゴリを削除
			deleteErr := rep.DeleteByName(ctx, tx, category.Name())
			Expect(deleteErr).NotTo(HaveOccurred(), "カテゴリの削除に失敗しました。")

			// 削除されたカテゴリが存在しないことを確認
			exists, existsErr := rep.ExistsByName(ctx, tx, category.Name())
			Expect(existsErr).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "削除したカテゴリがまだ存在しています。")
		})

		It("存在しないカテゴリ名で削除しようとするとエラーになること", func() {
			name, nameErr := categories.NewCategoryName("存在しないカテゴリ")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")

			deleteErr := rep.DeleteByName(ctx, tx, name)
			Expect(deleteErr).To(HaveOccurred())
			crudErr, ok := deleteErr.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("NOT_FOUND"))
			Expect(crudErr.Message).To(ContainSubstring("存在しないため、削除できませんでした。"))
		})
	})
})

var _ = Describe("productRepositoryImpl構造体", Ordered, Label("ProductRepositoryインターフェースメソッドのテスト"), func() {
	var rep products.ProductRepository
	var ctx context.Context
	var tx *sql.Tx
	var err error
	var testCategory *categories.Category

	BeforeAll(func() {
		rep = NewProductRepository()
		// テスト用のカテゴリを作成
		catName, catNameErr := categories.NewCategoryName("文房具")
		Expect(catNameErr).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
		catId, catIdErr := categories.NewCategoryId("b1524011-b6af-417e-8bf2-f449dd58b5c0")
		Expect(catIdErr).NotTo(HaveOccurred(), "テスト用カテゴリIDの生成に失敗しました。")
		testCategory, err = categories.BuildCategory(catId, catName)
		Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")
	})

	BeforeEach(func() {
		ctx = context.Background()
		tx, err = boil.BeginTx(ctx, nil)
		Expect(err).NotTo(HaveOccurred(), "トランザクションの開始に失敗しました。")
	})

	AfterEach(func() {
		err = tx.Rollback()
		Expect(err).NotTo(HaveOccurred(), "トランザクションのロールバックに失敗しました。")
	})

	DescribeTable(
		"ExistsByIdの動作確認",
		func(id string, expected bool) {
			productId, idErr := products.NewProductId(id)
			Expect(idErr).NotTo(HaveOccurred(), "テスト用商品IDの生成に失敗しました。")
			result, existsErr := rep.ExistsById(ctx, tx, productId)
			Expect(existsErr).NotTo(HaveOccurred(), "ExistsByIdの実行に失敗しました。")
			Expect(result).To(Equal(expected))
		},
		Entry("存在する商品ID(水性ボールペン(黒))", "ac413f22-0cf1-490a-9635-7e9ca810e544", true),
		Entry("存在しない商品ID", "00000000-0000-0000-0000-000000000000", false),
	)

	DescribeTable(
		"ExistsByNameの動作確認",
		func(name string, expected bool) {
			productName, nameErr := products.NewProductName(name)
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました。")
			result, existsErr := rep.ExistsByName(ctx, tx, productName)
			Expect(existsErr).NotTo(HaveOccurred(), "ExistsByNameの実行に失敗しました。")
			Expect(result).To(Equal(expected))
		},
		Entry("存在する商品名", "水性ボールペン(黒)", true),
		Entry("存在しない商品名", "存在しない商品", false),
	)

	Context("Createの動作確認", func() {
		It("新しい商品を作成できること", func() {
			name, nameErr := products.NewProductName("新しい商品")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました。")
			price, priceErr := products.NewProductPrice(1000)
			Expect(priceErr).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました。")
			product, productErr := products.NewProduct(name, price, testCategory)
			Expect(productErr).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました。")

			createErr := rep.Create(ctx, tx, product)
			Expect(createErr).NotTo(HaveOccurred(), "商品の作成に失敗しました。")

			// 作成された商品が存在することを確認
			exists, existsErr := rep.ExistsById(ctx, tx, product.Id())
			Expect(existsErr).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "作成した商品が存在しません。")
		})

		It("重複するobj_idでエラーになること", func() {
			// 水性ボールペン(黒)のobj_idを指定して商品を作成する
			id, idErr := products.NewProductId("ac413f22-0cf1-490a-9635-7e9ca810e544")
			Expect(idErr).NotTo(HaveOccurred(), "テスト用商品IDの生成に失敗しました。")
			name, nameErr := products.NewProductName("重複商品")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました。")
			price, priceErr := products.NewProductPrice(500)
			Expect(priceErr).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました。")
			product, productErr := products.BuildProduct(id, name, price, testCategory)
			Expect(productErr).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました。")

			createErr := rep.Create(ctx, tx, product)
			Expect(createErr).To(HaveOccurred())
			crudErr, ok := createErr.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("DB_UNIQUE_CONSTRAINT_VIOLATION"))
			Expect(crudErr.Message).To(ContainSubstring("一意制約違反です。"))
		})
	})

	Context("UpdateByIdの動作確認", func() {
		It("既存の商品を更新できること", func() {
			// 既存の商品(水性ボールペン(黒))を取得して更新
			id, idErr := products.NewProductId("ac413f22-0cf1-490a-9635-7e9ca810e544")
			Expect(idErr).NotTo(HaveOccurred(), "テスト用商品IDの生成に失敗しました。")
			newName, nameErr := products.NewProductName("更新された商品名")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました。")
			newPrice, priceErr := products.NewProductPrice(150)
			Expect(priceErr).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました。")
			product, productErr := products.BuildProduct(id, newName, newPrice, testCategory)
			Expect(productErr).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました。")

			updateErr := rep.UpdateById(ctx, tx, product)
			Expect(updateErr).NotTo(HaveOccurred(), "商品の更新に失敗しました。")
		})

		It("存在しない商品IDで更新しようとするとエラーになること", func() {
			id, idErr := products.NewProductId("00000000-0000-0000-0000-000000000000")
			Expect(idErr).NotTo(HaveOccurred(), "テスト用商品IDの生成に失敗しました。")
			name, nameErr := products.NewProductName("存在しない商品")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました。")
			price, priceErr := products.NewProductPrice(100)
			Expect(priceErr).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました。")
			product, productErr := products.BuildProduct(id, name, price, testCategory)
			Expect(productErr).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました。")

			updateErr := rep.UpdateById(ctx, tx, product)
			Expect(updateErr).To(HaveOccurred())
			crudErr, ok := updateErr.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("NOT_FOUND"))
			Expect(crudErr.Message).To(ContainSubstring("存在しないため、更新できませんでした。"))
		})
	})

	Context("DeleteByIdの動作確認", func() {
		It("既存の商品を削除できること", func() {
			// まず新しい商品を作成
			name, nameErr := products.NewProductName("削除テスト商品")
			Expect(nameErr).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました。")
			price, priceErr := products.NewProductPrice(200)
			Expect(priceErr).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました。")
			product, productErr := products.NewProduct(name, price, testCategory)
			Expect(productErr).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました。")

			createErr := rep.Create(ctx, tx, product)
			Expect(createErr).NotTo(HaveOccurred(), "商品の作成に失敗しました。")

			// 作成した商品を削除
			deleteErr := rep.DeleteById(ctx, tx, product.Id())
			Expect(deleteErr).NotTo(HaveOccurred(), "商品の削除に失敗しました。")

			// 削除された商品が存在しないことを確認
			exists, existsErr := rep.ExistsById(ctx, tx, product.Id())
			Expect(existsErr).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "削除した商品がまだ存在しています。")
		})

		It("存在しない商品IDで削除しようとするとエラーになること", func() {
			id, idErr := products.NewProductId("00000000-0000-0000-0000-000000000000")
			Expect(idErr).NotTo(HaveOccurred(), "テスト用商品IDの生成に失敗しました。")

			deleteErr := rep.DeleteById(ctx, tx, id)
			Expect(deleteErr).To(HaveOccurred())
			crudErr, ok := deleteErr.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("NOT_FOUND"))
			Expect(crudErr.Message).To(ContainSubstring("存在しないため、削除できませんでした。"))
		})
	})

})
