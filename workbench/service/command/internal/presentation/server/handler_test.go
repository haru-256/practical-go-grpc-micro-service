package server_test

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/presentation/server"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/testhelpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CategoryServiceHandler Unit Test", Label("UnitTests"), func() {
	var (
		ctrl                *gomock.Controller
		mockCategoryService *service.MockCategoryService
		csh                 *server.CategoryServiceHandlerImpl
		ctx                 context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockCategoryService = service.NewMockCategoryService(ctrl)
		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		var err error
		csh, err = server.NewCategoryServiceHandlerImpl(logger, mockCategoryService)
		Expect(err).NotTo(HaveOccurred())
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("CreateCategory", func() {
		Context("正常系: カテゴリが正常に作成される場合", func() {
			It("サービス層から返されたDTOをレスポンスとして返すこと", func() {
				// Arrange
				categoryName := "TestCategory"
				req := testhelpers.CreateCategoryRequest(categoryName)

				expectedDTO := &dto.CategoryDTO{
					Id:   "test-category-id",
					Name: categoryName,
				}

				mockCategoryService.EXPECT().
					Add(ctx, &dto.CreateCategoryDTO{Name: categoryName}).
					Return(expectedDTO, nil)

				// Act
				resp, err := csh.CreateCategory(ctx, req)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Msg.GetCategory().GetId()).To(Equal(expectedDTO.Id))
				Expect(resp.Msg.GetCategory().GetName()).To(Equal(expectedDTO.Name))
				Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())
			})
		})

		Context("異常系: バリデーションエラーが発生する場合", func() {
			It("空の名前で InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateCategoryRequest("")

				// Act
				resp, err := csh.CreateCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})
		})

		Context("異常系: サービス層でエラーが発生する場合", func() {
			It("Internal エラーを返すこと", func() {
				// Arrange
				categoryName := "TestCategory"
				req := testhelpers.CreateCategoryRequest(categoryName)

				expectedErr := errors.New("service error")
				mockCategoryService.EXPECT().
					Add(ctx, &dto.CreateCategoryDTO{Name: categoryName}).
					Return(nil, expectedErr)

				// Act
				resp, err := csh.CreateCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
			})
		})
	})

	Describe("UpdateCategory", func() {
		Context("正常系: カテゴリが正常に更新される場合", func() {
			It("サービス層から返されたDTOをレスポンスとして返すこと", func() {
				// Arrange
				categoryId := "test-category-id"
				newName := "UpdatedCategory"
				req := testhelpers.CreateUpdateCategoryRequest(categoryId, newName)

				expectedDTO := &dto.CategoryDTO{
					Id:   categoryId,
					Name: newName,
				}

				mockCategoryService.EXPECT().
					Update(ctx, &dto.UpdateCategoryDTO{Id: categoryId, Name: newName}).
					Return(expectedDTO, nil)

				// Act
				resp, err := csh.UpdateCategory(ctx, req)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Msg.GetCategory().GetId()).To(Equal(expectedDTO.Id))
				Expect(resp.Msg.GetCategory().GetName()).To(Equal(expectedDTO.Name))
				Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())
			})
		})

		Context("異常系: バリデーションエラーが発生する場合", func() {
			It("空のIDで InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateUpdateCategoryRequest("", "NewName")

				// Act
				resp, err := csh.UpdateCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})

			It("空の名前で InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateUpdateCategoryRequest("test-id", "")

				// Act
				resp, err := csh.UpdateCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})
		})

		Context("異常系: サービス層でエラーが発生する場合", func() {
			It("Internal エラーを返すこと", func() {
				// Arrange
				categoryId := "test-category-id"
				newName := "UpdatedCategory"
				req := testhelpers.CreateUpdateCategoryRequest(categoryId, newName)

				expectedErr := errors.New("service error")
				mockCategoryService.EXPECT().
					Update(ctx, &dto.UpdateCategoryDTO{Id: categoryId, Name: newName}).
					Return(nil, expectedErr)

				// Act
				resp, err := csh.UpdateCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
			})
		})
	})

	Describe("DeleteCategory", func() {
		Context("正常系: カテゴリが正常に削除される場合", func() {
			It("サービス層から返されたDTOをレスポンスとして返すこと", func() {
				// Arrange
				categoryId := "test-category-id"
				req := testhelpers.CreateDeleteCategoryRequest(categoryId)

				expectedDTO := &dto.CategoryDTO{
					Id:   categoryId,
					Name: "DeletedCategory",
				}

				mockCategoryService.EXPECT().
					Delete(ctx, &dto.DeleteCategoryDTO{Id: categoryId}).
					Return(expectedDTO, nil)

				// Act
				resp, err := csh.DeleteCategory(ctx, req)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Msg.GetCategory().GetId()).To(Equal(expectedDTO.Id))
				Expect(resp.Msg.GetCategory().GetName()).To(Equal(expectedDTO.Name))
				Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())
			})
		})

		Context("異常系: バリデーションエラーが発生する場合", func() {
			It("空のIDで InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateDeleteCategoryRequest("")

				// Act
				resp, err := csh.DeleteCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})
		})

		Context("異常系: サービス層でエラーが発生する場合", func() {
			It("Internal エラーを返すこと", func() {
				// Arrange
				categoryId := "test-category-id"
				req := testhelpers.CreateDeleteCategoryRequest(categoryId)

				expectedErr := errors.New("service error")
				mockCategoryService.EXPECT().
					Delete(ctx, &dto.DeleteCategoryDTO{Id: categoryId}).
					Return(nil, expectedErr)

				// Act
				resp, err := csh.DeleteCategory(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
			})
		})
	})
})

var _ = Describe("ProductServiceHandler Unit Test", Label("UnitTests"), func() {
	var (
		ctrl               *gomock.Controller
		mockProductService *service.MockProductService
		psh                *server.ProductServiceHandlerImpl
		ctx                context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockProductService = service.NewMockProductService(ctrl)
		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		var err error
		psh, err = server.NewProductServiceHandlerImpl(logger, mockProductService)
		Expect(err).NotTo(HaveOccurred())
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("CreateProduct", func() {
		Context("正常系: 商品が正常に作成される場合", func() {
			It("サービス層から返されたDTOをレスポンスとして返すこと", func() {
				// Arrange
				productName := "TestProduct"
				productPrice := uint32(1000)
				categoryId := "test-category-id"
				categoryName := "TestCategory"
				req := testhelpers.CreateProductRequest(productName, productPrice, categoryId, categoryName)

				expectedDTO := &dto.ProductDTO{
					Id:    "test-product-id",
					Name:  productName,
					Price: productPrice,
					Category: &dto.CategoryDTO{
						Id:   categoryId,
						Name: categoryName,
					},
				}

				mockProductService.EXPECT().
					Add(ctx, &dto.CreateProductDTO{
						Name:  productName,
						Price: productPrice,
						Category: &dto.CategoryDTO{
							Id:   categoryId,
							Name: categoryName,
						},
					}).
					Return(expectedDTO, nil)

				// Act
				resp, err := psh.CreateProduct(ctx, req)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Msg.GetProduct().GetId()).To(Equal(expectedDTO.Id))
				Expect(resp.Msg.GetProduct().GetName()).To(Equal(expectedDTO.Name))
				Expect(resp.Msg.GetProduct().GetPrice()).To(Equal(int32(expectedDTO.Price)))
				Expect(resp.Msg.GetProduct().GetCategory().GetId()).To(Equal(expectedDTO.Category.Id))
				Expect(resp.Msg.GetProduct().GetCategory().GetName()).To(Equal(expectedDTO.Category.Name))
				Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())
			})
		})

		Context("異常系: バリデーションエラーが発生する場合", func() {
			It("空の商品名で InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateProductRequest("", 1000, "cat-id", "Category")

				// Act
				resp, err := psh.CreateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})

			It("価格が0で InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateProductRequest("Product", 0, "cat-id", "Category")

				// Act
				resp, err := psh.CreateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})

			It("空のカテゴリIDで InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.CreateProductRequest("Product", 1000, "", "Category")

				// Act
				resp, err := psh.CreateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})
		})

		Context("異常系: サービス層でエラーが発生する場合", func() {
			It("Internal エラーを返すこと", func() {
				// Arrange
				productName := "TestProduct"
				productPrice := uint32(1000)
				categoryId := "test-category-id"
				categoryName := "TestCategory"
				req := testhelpers.CreateProductRequest(productName, productPrice, categoryId, categoryName)

				expectedErr := errors.New("service error")
				mockProductService.EXPECT().
					Add(ctx, &dto.CreateProductDTO{
						Name:  productName,
						Price: productPrice,
						Category: &dto.CategoryDTO{
							Id:   categoryId,
							Name: categoryName,
						},
					}).
					Return(nil, expectedErr)

				// Act
				resp, err := psh.CreateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
			})
		})
	})

	Describe("UpdateProduct", func() {
		Context("正常系: 商品が正常に更新される場合", func() {
			It("サービス層から返されたDTOをレスポンスとして返すこと", func() {
				// Arrange
				productId := "test-product-id"
				productName := "UpdatedProduct"
				productPrice := uint32(2000)
				categoryId := "test-category-id"
				req := testhelpers.UpdateProductRequest(productId, productName, productPrice, categoryId)

				expectedDTO := &dto.ProductDTO{
					Id:    productId,
					Name:  productName,
					Price: productPrice,
					Category: &dto.CategoryDTO{
						Id:   categoryId,
						Name: "TestCategory",
					},
				}

				mockProductService.EXPECT().
					Update(ctx, &dto.UpdateProductDTO{
						Id:         productId,
						Name:       productName,
						Price:      productPrice,
						CategoryId: categoryId,
					}).
					Return(expectedDTO, nil)

				// Act
				resp, err := psh.UpdateProduct(ctx, req)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Msg.GetProduct().GetId()).To(Equal(expectedDTO.Id))
				Expect(resp.Msg.GetProduct().GetName()).To(Equal(expectedDTO.Name))
				Expect(resp.Msg.GetProduct().GetPrice()).To(Equal(int32(expectedDTO.Price)))
				Expect(resp.Msg.GetProduct().GetCategory().GetId()).To(Equal(expectedDTO.Category.Id))
				Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())
			})
		})

		Context("異常系: バリデーションエラーが発生する場合", func() {
			It("空の商品IDで InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.UpdateProductRequest("", "Product", 1000, "cat-id")

				// Act
				resp, err := psh.UpdateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})

			It("空の商品名で InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.UpdateProductRequest("prod-id", "", 1000, "cat-id")

				// Act
				resp, err := psh.UpdateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})

			It("価格が0で InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.UpdateProductRequest("prod-id", "Product", 0, "cat-id")

				// Act
				resp, err := psh.UpdateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})

			It("空のカテゴリIDで InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.UpdateProductRequest("prod-id", "Product", 1000, "")

				// Act
				resp, err := psh.UpdateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})
		})

		Context("異常系: サービス層でエラーが発生する場合", func() {
			It("Internal エラーを返すこと", func() {
				// Arrange
				productId := "test-product-id"
				productName := "UpdatedProduct"
				productPrice := uint32(2000)
				categoryId := "test-category-id"
				req := testhelpers.UpdateProductRequest(productId, productName, productPrice, categoryId)

				expectedErr := errors.New("service error")
				mockProductService.EXPECT().
					Update(ctx, &dto.UpdateProductDTO{
						Id:         productId,
						Name:       productName,
						Price:      productPrice,
						CategoryId: categoryId,
					}).
					Return(nil, expectedErr)

				// Act
				resp, err := psh.UpdateProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
			})
		})
	})

	Describe("DeleteProduct", func() {
		Context("正常系: 商品が正常に削除される場合", func() {
			It("サービス層から返されたDTOをレスポンスとして返すこと", func() {
				// Arrange
				productId := "test-product-id"
				req := testhelpers.DeleteProductRequest(productId)

				expectedDTO := &dto.ProductDTO{
					Id:    productId,
					Name:  "DeletedProduct",
					Price: 1000,
					Category: &dto.CategoryDTO{
						Id:   "cat-id",
						Name: "Category",
					},
				}

				mockProductService.EXPECT().
					Delete(ctx, &dto.DeleteProductDTO{Id: productId}).
					Return(expectedDTO, nil)

				// Act
				resp, err := psh.DeleteProduct(ctx, req)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Msg.GetProduct().GetId()).To(Equal(expectedDTO.Id))
				Expect(resp.Msg.GetProduct().GetName()).To(Equal(expectedDTO.Name))
				Expect(resp.Msg.GetProduct().GetPrice()).To(Equal(int32(expectedDTO.Price)))
				Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())
			})
		})

		Context("異常系: バリデーションエラーが発生する場合", func() {
			It("空の商品IDで InvalidArgument エラーを返すこと", func() {
				// Arrange
				req := testhelpers.DeleteProductRequest("")

				// Act
				resp, err := psh.DeleteProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
			})
		})

		Context("異常系: サービス層でエラーが発生する場合", func() {
			It("Internal エラーを返すこと", func() {
				// Arrange
				productId := "test-product-id"
				req := testhelpers.DeleteProductRequest(productId)

				expectedErr := errors.New("service error")
				mockProductService.EXPECT().
					Delete(ctx, &dto.DeleteProductDTO{Id: productId}).
					Return(nil, expectedErr)

				// Act
				resp, err := psh.DeleteProduct(ctx, req)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				var connectErr *connect.Error
				Expect(errors.As(err, &connectErr)).To(BeTrue())
				Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
			})
		})
	})
})
