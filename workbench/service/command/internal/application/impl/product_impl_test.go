package impl

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("ProductService", Label("UnitTests"), func() {
	var (
		ctrl        *gomock.Controller
		mockRepo    *products.MockProductRepository
		mockTm      *service.MockTransactionManager
		ps          service.ProductService
		ctx         context.Context
		mockTx      *sql.Tx
		testProduct *products.Product
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = products.NewMockProductRepository(ctrl)
		mockTm = service.NewMockTransactionManager(ctrl)
		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		ps = NewProductServiceImpl(logger, mockRepo, mockTm)
		ctx = context.Background()

		// モックトランザクション（実際のオブジェクトをシミュレート）
		mockTx = &sql.Tx{}

		// テスト用商品の作成（ユニットテストではユニークである必要はない）
		name, err := products.NewProductName("TestProduct")
		Expect(err).NotTo(HaveOccurred())
		price, err := products.NewProductPrice(1000)
		Expect(err).NotTo(HaveOccurred())

		// テスト用カテゴリの作成
		categoryName, err := categories.NewCategoryName("TestCategory")
		Expect(err).NotTo(HaveOccurred())
		category, err := categories.NewCategory(categoryName)
		Expect(err).NotTo(HaveOccurred())

		testProduct, err = products.NewProduct(name, price, category)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Add", func() {
		Context("when product name does not exist", func() {
			It("should successfully add a new product", func() {
				// Arrange: モックの期待値を順序付きで設定
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, testProduct).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				err := ps.Add(ctx, testProduct)

				// Assert
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when product name already exists", func() {
			It("should return ApplicationError with PRODUCT_ALREADY_EXISTS code", func() {
				// Arrange
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(true, nil),
					mockTm.EXPECT().Complete(ctx, mockTx, gomock.Any()).
						Do(func(ctx context.Context, tx *sql.Tx, err error) {
							// Completeに渡されるエラーがApplicationErrorであることを検証
							Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
							appErr := err.(*errs.ApplicationError)
							Expect(appErr.Code).To(Equal("PRODUCT_ALREADY_EXISTS"))
						}).
						Return(nil),
				)

				// Act
				err := ps.Add(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
				appErr := err.(*errs.ApplicationError)
				Expect(appErr.Code).To(Equal("PRODUCT_ALREADY_EXISTS"))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				err := ps.Add(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(beginErr))
			})
		})

		Context("when ExistsByName fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				existsErr := fmt.Errorf("database error")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(false, existsErr),
					mockTm.EXPECT().Complete(ctx, mockTx, existsErr).Return(nil),
				)

				// Act
				err := ps.Add(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(existsErr))
			})
		})

		Context("when Create fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				createErr := errs.NewCRUDError("DB_ERROR", "failed to create product")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, testProduct).Return(createErr),
					mockTm.EXPECT().Complete(ctx, mockTx, createErr).Return(nil),
				)

				// Act
				err := ps.Add(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(createErr))
			})
		})
	})

	Describe("Update", func() {
		Context("when update is successful", func() {
			It("should successfully update the product", func() {
				// Arrange
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().UpdateById(ctx, mockTx, testProduct).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				err := ps.Update(ctx, testProduct)

				// Assert
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when UpdateById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				updateErr := errs.NewCRUDError("NOT_FOUND", "product not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().UpdateById(ctx, mockTx, testProduct).Return(updateErr),
					mockTm.EXPECT().Complete(ctx, mockTx, updateErr).Return(nil),
				)

				// Act
				err := ps.Update(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(updateErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				err := ps.Update(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(beginErr))
			})
		})
	})

	Describe("Delete", func() {
		Context("when delete is successful", func() {
			It("should successfully delete the product", func() {
				// Arrange
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().DeleteById(ctx, mockTx, testProduct.Id()).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				err := ps.Delete(ctx, testProduct)

				// Assert
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when DeleteById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				deleteErr := errs.NewCRUDError("NOT_FOUND", "product not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().DeleteById(ctx, mockTx, testProduct.Id()).Return(deleteErr),
					mockTm.EXPECT().Complete(ctx, mockTx, deleteErr).Return(nil),
				)

				// Act
				err := ps.Delete(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(deleteErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				err := ps.Delete(ctx, testProduct)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(beginErr))
			})
		})
	})
})
