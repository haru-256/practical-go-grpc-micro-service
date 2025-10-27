//go:build integration || !ci

package presentation_test

import (
	"context"
	"testing"
	"time"

	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/presentation"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/presentation/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestPresentation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Presentation Module Suite")
}

var _ = Describe("Presentation Module", func() {
	var (
		app *fx.App
	)

	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)

	AfterEach(func() {
		if app != nil {
			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(app.Stop(stopCtx)).To(Succeed())
		}
	})

	Context("モジュールの初期化", func() {
		It("正常に初期化され、必要な依存関係が提供されること", func() {
			var categoryServiceHandler cmdconnect.CategoryServiceHandler
			var categoryService service.CategoryService

			app = fx.New(
				// configPathとconfigNameを提供
				configOption,
				presentation.Module,
				fx.Populate(&categoryServiceHandler, &categoryService),
				fx.NopLogger, // テスト時はログを抑制
			)

			Expect(app.Err()).ToNot(HaveOccurred(), "fx app should initialize without errors")
			Expect(categoryServiceHandler).ToNot(BeNil(), "category service handler should be provided")
			Expect(categoryServiceHandler).To(BeAssignableToTypeOf(&server.CategoryServiceHandlerImpl{}), "category service handler should be the correct type")
			// Application層のサービスも提供されていることを確認
			Expect(categoryService).ToNot(BeNil(), "category service should be provided")
		})
	})
})
