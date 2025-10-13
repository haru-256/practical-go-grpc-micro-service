//go:build integration || !ci

package repository

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepImplPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Implementation Suite")
}

var _ = BeforeSuite(func() {
	absPath, _ := filepath.Abs("../config/database.toml")
	Expect(os.Setenv("DATABASE_TOML_PATH", absPath)).To(Succeed())
	config, err := handler.NewDBConfig()
	Expect(err).NotTo(HaveOccurred(), "DBConfigの生成に失敗しました。")
	_, err = handler.NewDatabase(config)
	Expect(err).NotTo(HaveOccurred(), "データベース接続が失敗したのでテストを中止します。")
})

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
})
