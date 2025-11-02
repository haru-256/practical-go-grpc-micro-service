package impl

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	mock_service "github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/mock/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("handleTransactionComplete", Label("UnitTests"), func() {
	var (
		ctrl   *gomock.Controller
		mockTm *mock_service.MockTransactionManager
		ctx    context.Context
		mockTx *sql.Tx
		logger *slog.Logger
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockTm = mock_service.NewMockTransactionManager(ctrl)
		ctx = context.Background()
		mockTx = &sql.Tx{}
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when business logic succeeds (err == nil)", func() {
		It("should commit successfully and not modify err or result", func() {
			// Arrange
			var err error = nil
			result := &dto.CategoryDTO{Id: "test-id", Name: "test-name"}
			mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil)

			// Act
			handleTransactionComplete(ctx, mockTm, mockTx, &err, &result, logger)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Id).To(Equal("test-id"))
		})

		It("should set err and clear result when commit fails", func() {
			// Arrange
			var err error = nil
			result := &dto.CategoryDTO{Id: "test-id", Name: "test-name"}
			commitErr := fmt.Errorf("commit failed")
			mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(commitErr)

			// Act
			handleTransactionComplete(ctx, mockTm, mockTx, &err, &result, logger)

			// Assert
			Expect(err).To(Equal(commitErr))
			Expect(result).To(BeNil())
		})
	})

	Context("when business logic fails (err != nil)", func() {
		It("should rollback successfully, keep original error, and clear result", func() {
			// Arrange
			originalErr := fmt.Errorf("business logic error")
			err := originalErr
			result := &dto.CategoryDTO{Id: "test-id", Name: "test-name"}
			mockTm.EXPECT().Complete(ctx, mockTx, originalErr).Return(nil)

			// Act
			handleTransactionComplete(ctx, mockTm, mockTx, &err, &result, logger)

			// Assert
			Expect(err).To(Equal(originalErr))
			Expect(result).To(BeNil())
		})

		It("should keep original error and clear result when rollback also fails", func() {
			// Arrange
			originalErr := fmt.Errorf("business logic error")
			rollbackErr := fmt.Errorf("rollback failed")
			err := originalErr
			result := &dto.CategoryDTO{Id: "test-id", Name: "test-name"}
			mockTm.EXPECT().Complete(ctx, mockTx, originalErr).Return(rollbackErr)

			// Act
			handleTransactionComplete(ctx, mockTm, mockTx, &err, &result, logger)

			// Assert
			// 元のエラーを保持する（ロールバックエラーはログに記録されるのみ）
			Expect(err).To(Equal(originalErr))
			Expect(result).To(BeNil())
		})
	})
})
