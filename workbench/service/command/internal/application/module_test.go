//go:build integration || !ci

package application

import (
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Module Suite")
}

var _ = Describe("Application Module", func() {
	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)

	Context("when creating an fx app with the application module", func() {
		It("should successfully initialize without errors", func() {
			var categoryService service.CategoryService
			var productService service.ProductService

			app := fx.New(
				// configPathとconfigNameを提供
				configOption,
				Module,
				fx.Populate(&categoryService, &productService),
				fx.NopLogger, // テスト時はログを抑制
			)

			Expect(app.Err()).ToNot(HaveOccurred(), "fx app should initialize without errors")
		})

		It("should provide all required services", func() {
			var categoryService service.CategoryService
			var productService service.ProductService

			app := fx.New(
				configOption,
				Module,
				fx.Populate(&categoryService, &productService),
				fx.NopLogger,
			)

			Expect(app.Err()).ToNot(HaveOccurred())
			Expect(categoryService).ToNot(BeNil(), "category service should be provided")
			Expect(productService).ToNot(BeNil(), "product service should be provided")
		})

		It("should properly wire dependencies from infrastructure layer", func() {
			var categoryService service.CategoryService
			var productService service.ProductService
			var tm service.TransactionManager

			app := fx.New(
				configOption,
				Module,
				fx.Populate(&categoryService, &productService, &tm),
				fx.NopLogger,
			)

			Expect(app.Err()).ToNot(HaveOccurred())

			// インフラ層から提供されるTransactionManagerも取得できることを確認
			Expect(tm).ToNot(BeNil(), "transaction manager from infrastructure layer should be accessible")
		})
	})
})
