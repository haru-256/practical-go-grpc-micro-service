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
	os.Setenv("DATABASE_TOML_PATH", absPath)
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
			catName, err := categories.NewCategoryName(name)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, err := categories.NewCategory(catName)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")
			result, err := rep.ExistsByName(ctx, tx, category.Name())
			Expect(err).NotTo(HaveOccurred(), "ExistsByIdの実行に失敗しました。")
			Expect(result).To(Equal(expected), "存在するカテゴリIDに対してExistsByIdがfalseを返しました。")
		},
		Entry("文房具", "文房具", true),
		Entry("食品", "食品", false),
	)

	Context("Createの動作確認", func() {
		It("新しいカテゴリを作成できること", func() {
			name, err := categories.NewCategoryName("食品")
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, err := categories.NewCategory(name)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")

			err = rep.Create(ctx, tx, category)
			Expect(err).NotTo(HaveOccurred(), "カテゴリの作成に失敗しました。")
		})
		It("重複するobj_idでエラーになること", func() {
			// 文房具のobj_idを指定してカテゴリを作成する
			id, err := categories.NewCategoryId("b1524011-b6af-417e-8bf2-f449dd58b5c0")
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリIDの生成に失敗しました。")
			name, err := categories.NewCategoryName("新しい文房具")
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました。")
			category, err := categories.BuildCategory(id, name)
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリの生成に失敗しました。")

			err = rep.Create(ctx, tx, category)
			Expect(err).To(HaveOccurred())
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue())
			Expect(crudErr.Code).To(Equal("DB_UNIQUE_CONSTRAINT_VIOLATION"))
			Expect(crudErr.Message).To(ContainSubstring("一意制約違反です。"))
		})
	})
})
