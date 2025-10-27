package impl

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
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
		ctrl             *gomock.Controller
		mockProductRepo  *products.MockProductRepository
		mockCategoryRepo *categories.MockCategoryRepository
		mockTm           *service.MockTransactionManager
		ps               service.ProductService
		ctx              context.Context
		mockTx           *sql.Tx
		testProduct      *products.Product
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockProductRepo = products.NewMockProductRepository(ctrl)
		mockCategoryRepo = categories.NewMockCategoryRepository(ctrl)
		mockTm = service.NewMockTransactionManager(ctrl)
		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		ps = NewProductServiceImpl(logger, mockProductRepo, mockCategoryRepo, mockTm)
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
				// Arrange
				createDTO := &dto.CreateProductDTO{
					Name: "TestProduct",
					Category: &dto.CategoryDTO{
						Id:   testProduct.Category().Id().Value(),
						Name: testProduct.Category().Name().Value(),
					},
					Price: 1000,
				}

				// Arrange: モックの期待値を順序付きで設定
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockProductRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(false, nil),
					mockProductRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				result, err := ps.Add(ctx, createDTO)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("TestProduct"))
				Expect(result.Price).To(Equal(uint32(1000)))
			})
		})

		Context("when product name already exists", func() {
			It("should return ApplicationError with PRODUCT_ALREADY_EXISTS code", func() {
				// Arrange
				createDTO := &dto.CreateProductDTO{
					Name: "TestProduct",
					Category: &dto.CategoryDTO{
						Id:   testProduct.Category().Id().Value(),
						Name: testProduct.Category().Name().Value(),
					},
					Price: 1000,
				}

				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockProductRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(true, nil),
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
				result, err := ps.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
				appErr := err.(*errs.ApplicationError)
				Expect(appErr.Code).To(Equal("PRODUCT_ALREADY_EXISTS"))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				createDTO := &dto.CreateProductDTO{
					Name: "TestProduct",
					Category: &dto.CategoryDTO{
						Id:   testProduct.Category().Id().Value(),
						Name: testProduct.Category().Name().Value(),
					},
					Price: 1000,
				}
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				result, err := ps.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(beginErr))
			})
		})

		Context("when ExistsByName fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				createDTO := &dto.CreateProductDTO{
					Name: "TestProduct",
					Category: &dto.CategoryDTO{
						Id:   testProduct.Category().Id().Value(),
						Name: testProduct.Category().Name().Value(),
					},
					Price: 1000,
				}
				existsErr := fmt.Errorf("database error")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockProductRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(false, existsErr),
					mockTm.EXPECT().Complete(ctx, mockTx, existsErr).Return(nil),
				)

				// Act
				result, err := ps.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(existsErr))
			})
		})

		Context("when Create fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				createDTO := &dto.CreateProductDTO{
					Name: "TestProduct",
					Category: &dto.CategoryDTO{
						Id:   testProduct.Category().Id().Value(),
						Name: testProduct.Category().Name().Value(),
					},
					Price: 1000,
				}
				createErr := errs.NewCRUDError("DB_ERROR", "failed to create product")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockProductRepo.EXPECT().ExistsByName(ctx, mockTx, testProduct.Name()).Return(false, nil),
					mockProductRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(createErr),
					mockTm.EXPECT().Complete(ctx, mockTx, createErr).Return(nil),
				)

				// Act
				result, err := ps.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(createErr))
			})
		})
	})

	Describe("Update", func() {
		Context("when update is successful", func() {
			It("should successfully update the product", func() {
				// Arrange
				updateDTO := &dto.UpdateProductDTO{
					Id:         testProduct.Id().Value(),
					Name:       "UpdatedProduct",
					CategoryId: testProduct.Category().Id().Value(),
					Price:      2000,
				}

				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockCategoryRepo.EXPECT().FindById(ctx, mockTx, gomock.Any()).Do(
						func(ctx context.Context, tx *sql.Tx, id *categories.CategoryId) {
							Expect(id.Value()).To(Equal(updateDTO.CategoryId))
						}).Return(testProduct.Category(), nil),
					mockProductRepo.EXPECT().UpdateById(ctx, mockTx, gomock.Any()).Do(
						func(ctx context.Context, tx *sql.Tx, product *products.Product) {
							Expect(product.Name().Value()).To(Equal(updateDTO.Name))
							Expect(product.Price().Value()).To(Equal(updateDTO.Price))
							Expect(product.Category().Id().Value()).To(Equal(updateDTO.CategoryId))
							Expect(product.Id().Value()).To(Equal(testProduct.Id().Value()))
						}).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				result, err := ps.Update(ctx, updateDTO)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("UpdatedProduct"))
				Expect(result.Price).To(Equal(uint32(2000)))
			})
		})

		Context("when UpdateById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				updateDTO := &dto.UpdateProductDTO{
					Id:         testProduct.Id().Value(),
					Name:       "UpdatedProduct",
					CategoryId: testProduct.Category().Id().Value(),

					Price: 2000,
				}
				updateErr := errs.NewCRUDError("NOT_FOUND", "product not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockCategoryRepo.EXPECT().FindById(ctx, mockTx, gomock.Any()).Do(
						func(ctx context.Context, tx *sql.Tx, id *categories.CategoryId) {
							Expect(id.Value()).To(Equal(updateDTO.CategoryId))
						}).Return(testProduct.Category(), nil),
					mockProductRepo.EXPECT().UpdateById(ctx, mockTx, gomock.Any()).Return(updateErr),
					mockTm.EXPECT().Complete(ctx, mockTx, updateErr).Return(nil),
				)

				// Act
				result, err := ps.Update(ctx, updateDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(updateErr))
			})
		})

		Context("when FindById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				updateDTO := &dto.UpdateProductDTO{
					Id:         testProduct.Id().Value(),
					Name:       "UpdatedProduct",
					CategoryId: testProduct.Category().Id().Value(),
					Price:      2000,
				}
				findErr := errs.NewCRUDError("NOT_FOUND", "category not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockCategoryRepo.EXPECT().FindById(ctx, mockTx, gomock.Any()).Do(
						func(ctx context.Context, tx *sql.Tx, id *categories.CategoryId) {
							Expect(id.Value()).To(Equal(updateDTO.CategoryId))
						}).Return(nil, findErr),
					mockTm.EXPECT().Complete(ctx, mockTx, findErr).Return(nil),
				)

				// Act
				result, err := ps.Update(ctx, updateDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(findErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				updateDTO := &dto.UpdateProductDTO{
					Id:         testProduct.Id().Value(),
					Name:       "UpdatedProduct",
					CategoryId: testProduct.Category().Id().Value(),

					Price: 2000,
				}
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				result, err := ps.Update(ctx, updateDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(beginErr))
			})
		})
	})

	Describe("Delete", func() {
		Context("when delete is successful", func() {
			It("should successfully delete the product", func() {
				// Arrange
				deleteDTO := &dto.DeleteProductDTO{
					Id: testProduct.Id().Value(),
				}

				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockProductRepo.EXPECT().FindById(ctx, mockTx, testProduct.Id()).Return(testProduct, nil),
					mockProductRepo.EXPECT().DeleteById(ctx, mockTx, testProduct.Id()).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				result, err := ps.Delete(ctx, deleteDTO)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("TestProduct"))
				Expect(result.Price).To(Equal(uint32(1000)))
			})
		})

		Context("when DeleteById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				deleteDTO := &dto.DeleteProductDTO{
					Id: testProduct.Id().Value(),
				}
				deleteErr := errs.NewCRUDError("NOT_FOUND", "product not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockProductRepo.EXPECT().FindById(ctx, mockTx, testProduct.Id()).Return(nil, deleteErr),
					mockTm.EXPECT().Complete(ctx, mockTx, deleteErr).Return(nil),
				)

				// Act
				result, err := ps.Delete(ctx, deleteDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(deleteErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				deleteDTO := &dto.DeleteProductDTO{
					Id: testProduct.Id().Value(),
				}
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				result, err := ps.Delete(ctx, deleteDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(beginErr))
			})
		})
	})
})
