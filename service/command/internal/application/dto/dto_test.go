package dto_test

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCategoryDTOFromEntity", func() {
	It("Categoryエンティティから正しくDTOを生成できる", func() {
		// Arrange
		categoryName, err := categories.NewCategoryName("Electronics")
		Expect(err).NotTo(HaveOccurred())
		category, err := categories.NewCategory(categoryName)
		Expect(err).NotTo(HaveOccurred())

		// Act
		result := dto.NewCategoryDTOFromEntity(category)

		// Assert
		Expect(result).NotTo(BeNil())
		Expect(result.Id).To(Equal(category.Id().Value()))
		Expect(result.Name).To(Equal("Electronics"))
	})
})

var _ = Describe("CategoryFromCreateDTO", func() {
	Context("正常系", func() {
		It("CreateCategoryDTOから正しくCategoryエンティティを生成できる", func() {
			// Arrange
			createDTO := &dto.CreateCategoryDTO{
				Name: "Books",
			}

			// Act
			result, err := dto.CategoryFromCreateDTO(createDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Name().Value()).To(Equal("Books"))
		})
	})

	Context("異常系", func() {
		It("空のNameの場合エラーを返す", func() {
			// Arrange
			createDTO := &dto.CreateCategoryDTO{
				Name: "",
			}

			// Act
			_, err := dto.CategoryFromCreateDTO(createDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("CategoryFromUpdateDTO", func() {
	Context("正常系", func() {
		It("UpdateCategoryDTOから正しくCategoryエンティティを再構築できる", func() {
			// Arrange
			updateDTO := &dto.UpdateCategoryDTO{
				Id:   "550e8400-e29b-41d4-a716-446655440000",
				Name: "Updated Category",
			}

			// Act
			result, err := dto.CategoryFromUpdateDTO(updateDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Id().Value()).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
			Expect(result.Name().Value()).To(Equal("Updated Category"))
		})
	})

	Context("異常系", func() {
		It("不正なIdの場合エラーを返す", func() {
			// Arrange
			updateDTO := &dto.UpdateCategoryDTO{
				Id:   "",
				Name: "Valid Name",
			}

			// Act
			_, err := dto.CategoryFromUpdateDTO(updateDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})

		It("不正なNameの場合エラーを返す", func() {
			// Arrange
			updateDTO := &dto.UpdateCategoryDTO{
				Id:   "cat-123",
				Name: "",
			}

			// Act
			_, err := dto.CategoryFromUpdateDTO(updateDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("NewProductDTOFromEntity", func() {
	It("Productエンティティから正しくDTOを生成できる", func() {
		// Arrange
		categoryName, err := categories.NewCategoryName("Electronics")
		Expect(err).NotTo(HaveOccurred())
		category, err := categories.NewCategory(categoryName)
		Expect(err).NotTo(HaveOccurred())

		productName, err := products.NewProductName("Laptop")
		Expect(err).NotTo(HaveOccurred())
		productPrice, err := products.NewProductPrice(100000)
		Expect(err).NotTo(HaveOccurred())
		product, err := products.NewProduct(productName, productPrice, category)
		Expect(err).NotTo(HaveOccurred())

		// Act
		result := dto.NewProductDTOFromEntity(product)

		// Assert
		Expect(result).NotTo(BeNil())
		Expect(result.Id).To(Equal(product.Id().Value()))
		Expect(result.Name).To(Equal("Laptop"))
		Expect(result.Price).To(Equal(uint32(100000)))
		Expect(result.Category.Id).To(Equal(category.Id().Value()))
		Expect(result.Category.Name).To(Equal("Electronics"))
	})
})

var _ = Describe("ProductFromCreateDTO", func() {
	Context("正常系", func() {
		It("CreateProductDTOから正しくProductエンティティを生成できる", func() {
			// Arrange
			createDTO := &dto.CreateProductDTO{
				Name:  "Smartphone",
				Price: 50000,
				Category: &dto.CategoryDTO{
					Id:   "550e8400-e29b-41d4-a716-446655440000",
					Name: "Electronics",
				},
			}

			// Act
			result, err := dto.ProductFromCreateDTO(createDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Name().Value()).To(Equal("Smartphone"))
			Expect(result.Price().Value()).To(Equal(uint32(50000)))
			Expect(result.Category().Name().Value()).To(Equal("Electronics"))
			Expect(result.Category().Id().Value()).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
		})
	})

	Context("異常系", func() {
		It("不正なProductNameの場合エラーを返す", func() {
			// Arrange
			createDTO := &dto.CreateProductDTO{
				Name:  "",
				Price: 50000,
				Category: &dto.CategoryDTO{
					Id:   "550e8400-e29b-41d4-a716-446655440000",
					Name: "Electronics",
				},
			}

			// Act
			_, err := dto.ProductFromCreateDTO(createDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})

		It("不正なCategoryの場合エラーを返す", func() {
			// Arrange
			createDTO := &dto.CreateProductDTO{
				Name:  "Smartphone",
				Price: 50000,
				Category: &dto.CategoryDTO{
					Id:   "",
					Name: "Electronics",
				},
			}

			// Act
			_, err := dto.ProductFromCreateDTO(createDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("ProductFromUpdateDTO", func() {
	Context("正常系", func() {
		It("UpdateProductDTOから正しくProductエンティティを再構築できる", func() {
			// Arrange
			updateDTO := &dto.UpdateProductDTO{
				Id:         "650e8400-e29b-41d4-a716-446655440000",
				Name:       "Updated Product",
				Price:      75000,
				CategoryId: "550e8400-e29b-41d4-a716-446655440000",
			}
			categoryName, err := categories.NewCategoryName("Updated Category")
			Expect(err).NotTo(HaveOccurred())

			// Act
			result, err := dto.ProductFromUpdateDTO(updateDTO, categoryName)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Id().Value()).To(Equal("650e8400-e29b-41d4-a716-446655440000"))
			Expect(result.Name().Value()).To(Equal("Updated Product"))
			Expect(result.Price().Value()).To(Equal(uint32(75000)))
			Expect(result.Category().Id().Value()).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
			Expect(result.Category().Name().Value()).To(Equal("Updated Category"))
		})
	})

	Context("異常系", func() {
		It("不正なProductIdの場合エラーを返す", func() {
			// Arrange
			updateDTO := &dto.UpdateProductDTO{
				Id:         "",
				Name:       "Valid Name",
				Price:      75000,
				CategoryId: "550e8400-e29b-41d4-a716-446655440000",
			}
			categoryName, err := categories.NewCategoryName("Updated Category")
			Expect(err).NotTo(HaveOccurred())

			// Act
			_, err = dto.ProductFromUpdateDTO(updateDTO, categoryName)

			// Assert
			Expect(err).To(HaveOccurred())
		})

		It("不正なProductNameの場合エラーを返す", func() {
			// Arrange
			updateDTO := &dto.UpdateProductDTO{
				Id:         "650e8400-e29b-41d4-a716-446655440000",
				Name:       "",
				Price:      75000,
				CategoryId: "550e8400-e29b-41d4-a716-446655440000",
			}
			categoryName, err := categories.NewCategoryName("Updated Category")
			Expect(err).NotTo(HaveOccurred())

			// Act
			_, err = dto.ProductFromUpdateDTO(updateDTO, categoryName)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("CategoryIdFromDeleteDTO", func() {
	Context("正常系", func() {
		It("DeleteCategoryDTOから正しくCategoryIdを取得できる", func() {
			// Arrange
			deleteDTO := &dto.DeleteCategoryDTO{
				Id: "550e8400-e29b-41d4-a716-446655440000",
			}

			// Act
			result, err := dto.CategoryIdFromDeleteDTO(deleteDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Value()).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
		})
	})

	Context("異常系", func() {
		It("不正なIdの場合エラーを返す", func() {
			// Arrange
			deleteDTO := &dto.DeleteCategoryDTO{
				Id: "",
			}

			// Act
			_, err := dto.CategoryIdFromDeleteDTO(deleteDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})

		It("UUID形式でないIdの場合エラーを返す", func() {
			// Arrange
			deleteDTO := &dto.DeleteCategoryDTO{
				Id: "invalid-id",
			}

			// Act
			_, err := dto.CategoryIdFromDeleteDTO(deleteDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("ProductIdFromDeleteDTO", func() {
	Context("正常系", func() {
		It("DeleteProductDTOから正しくProductIdを取得できる", func() {
			// Arrange
			deleteDTO := &dto.DeleteProductDTO{
				Id: "650e8400-e29b-41d4-a716-446655440000",
			}

			// Act
			result, err := dto.ProductIdFromDeleteDTO(deleteDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Value()).To(Equal("650e8400-e29b-41d4-a716-446655440000"))
		})
	})

	Context("異常系", func() {
		It("不正なIdの場合エラーを返す", func() {
			// Arrange
			deleteDTO := &dto.DeleteProductDTO{
				Id: "",
			}

			// Act
			_, err := dto.ProductIdFromDeleteDTO(deleteDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})

		It("UUID形式でないIdの場合エラーを返す", func() {
			// Arrange
			deleteDTO := &dto.DeleteProductDTO{
				Id: "invalid-id",
			}

			// Act
			_, err := dto.ProductIdFromDeleteDTO(deleteDTO)

			// Assert
			Expect(err).To(HaveOccurred())
		})
	})
})
