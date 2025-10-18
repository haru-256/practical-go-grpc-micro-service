//go:build integration || !ci

package sqlboiler

import (
	"database/sql"
	"log/slog"
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestSqlboiler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqlboiler Module Suite")
}

var _ = Describe("Sqlboiler Module", func() {
	configOption := fx.Supply(
		fx.Annotate("../../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)

	Context("when creating an fx app with the sqlboiler module", func() {
		It("should successfully initialize without errors", func() {
			var db *sql.DB
			var logger *slog.Logger
			var categoryRepo categories.CategoryRepository
			var productRepo products.ProductRepository
			var tm service.TransactionManager

			app := fx.New(
				configOption,
				Module,
				fx.Populate(&db, &logger, &categoryRepo, &productRepo, &tm),
				fx.NopLogger, // テスト時はログを抑制
			)

			Expect(app.Err()).ToNot(HaveOccurred(), "fx app should initialize without errors")
		})

		It("should provide all required dependencies", func() {
			var db *sql.DB
			var logger *slog.Logger
			var categoryRepo categories.CategoryRepository
			var productRepo products.ProductRepository
			var tm service.TransactionManager

			app := fx.New(
				configOption,
				Module,
				fx.Populate(&db, &logger, &categoryRepo, &productRepo, &tm),
				fx.NopLogger,
			)

			Expect(app.Err()).ToNot(HaveOccurred())
			Expect(db).ToNot(BeNil(), "database connection should be provided")
			Expect(logger).ToNot(BeNil(), "logger should be provided")
			Expect(categoryRepo).ToNot(BeNil(), "category repository should be provided")
			Expect(productRepo).ToNot(BeNil(), "product repository should be provided")
			Expect(tm).ToNot(BeNil(), "transaction manager should be provided")
		})

		It("should verify database connection is functional", func() {
			var db *sql.DB

			app := fx.New(
				configOption,
				Module,
				fx.Populate(&db),
				fx.NopLogger,
			)

			Expect(app.Err()).ToNot(HaveOccurred())
			Expect(db).ToNot(BeNil())

			// データベース接続が有効であることを確認
			err := db.Ping()
			Expect(err).ToNot(HaveOccurred(), "database connection should be alive")
		})
	})
})
