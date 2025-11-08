//go:build integration || !ci

package presentation_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"buf.build/gen/go/grpc/grpc/connectrpc/go/grpc/health/v1/healthv1connect"
	healthv1 "buf.build/gen/go/grpc/grpc/protocolbuffers/go/grpc/health/v1"
	"connectrpc.com/connect"
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

var _ = Describe("Presentation Module", Label("Module"), func() {
	var (
		app *fx.App
	)

	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)

	appCleanup := func() {
		if app != nil {
			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(app.Stop(stopCtx)).To(Succeed())
			app = nil
		}
	}

	Context("モジュールの初期化", func() {
		It("正常に初期化され、必要な依存関係が提供されること", func() {
			var (
				categoryServiceHandler cmdconnect.CategoryServiceHandler
				categoryService        service.CategoryService
				productServiceHandler  cmdconnect.ProductServiceHandler
				productService         service.ProductService
				commandServer          *server.CommandServer
			)
			// Act
			app = fx.New(
				// configPathとconfigNameを提供
				configOption,
				presentation.Module,
				fx.Populate(&categoryServiceHandler, &categoryService, &productServiceHandler, &productService, &commandServer),
				fx.NopLogger, // テスト時はログを抑制
			)
			DeferCleanup(appCleanup)

			// Assert
			Expect(app.Err()).ToNot(HaveOccurred(), "fx app should initialize without errors")
			Expect(categoryServiceHandler).ToNot(BeNil(), "category service handler should be provided")
			Expect(categoryServiceHandler).To(BeAssignableToTypeOf(&server.CategoryServiceHandlerImpl{}), "category service handler should be the correct type")
			Expect(categoryService).ToNot(BeNil(), "category service should be provided")
			Expect(productServiceHandler).ToNot(BeNil(), "product service handler should be provided")
			Expect(productServiceHandler).To(BeAssignableToTypeOf(&server.ProductServiceHandlerImpl{}), "product service handler should be the correct type")
			Expect(productService).ToNot(BeNil(), "product service should be provided")
			Expect(commandServer).ToNot(BeNil(), "command server should be provided")
		})
	})

	Context("ヘルスチェックエンドポイント", func() {
		var commandServer *server.CommandServer

		BeforeEach(func() {
			app = fx.New(
				// configPathとconfigNameを提供
				configOption,
				presentation.Module,
				fx.Populate(&commandServer),
				fx.NopLogger, // テスト時はログを抑制
			)
			DeferCleanup(appCleanup)
			Expect(app.Err()).ToNot(HaveOccurred(), "fx app should initialize without errors")
			Expect(commandServer).ToNot(BeNil(), "command server should be provided")

			// アプリケーションの起動
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			DeferCleanup(cancel)
			Expect(app.Start(ctx)).ToNot(HaveOccurred(), "fx app should start without errors")
			// サーバーが初期化されていることを確認
			Expect(commandServer.Server.Addr).ToNot(BeEmpty(), "server address should be set")
		})

		It("SERVINGステータスを返すこと", func() {
			// Arrange
			client := healthv1connect.NewHealthClient(
				http.DefaultClient,
				"http://"+commandServer.Server.Addr,
				connect.WithGRPC(),
			)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			DeferCleanup(cancel)

			Eventually(func(g Gomega) {
				// Act
				res, err := client.Check(ctx, connect.NewRequest(&healthv1.HealthCheckRequest{}))

				// Assert
				g.Expect(err).ToNot(HaveOccurred(), "health check should succeed")
				g.Expect(res.Msg).ToNot(BeNil(), "response message should not be nil")
				g.Expect(res.Msg.Status).To(Equal(healthv1.HealthCheckResponse_SERVING), "response status should be SERVING")
			}).WithTimeout(3 * time.Second).WithPolling(100 * time.Millisecond).Should(Succeed())
		})
	})
})
