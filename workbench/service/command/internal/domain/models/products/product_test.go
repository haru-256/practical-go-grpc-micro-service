package products

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEntityPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "domain/models/products packageのテスト")
}

var _ = Describe("Productエンティティを構成する値オブジェクト", Label("Productの値オブジェクト"), func() {
	DescribeTable("商品IDのバリデーション",
		func(id string, expectError bool, expectedErrorCode string) {
			productId, err := NewProductId(id)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(productId).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(productId).NotTo(BeNil())
			}
		},
		Entry(
			"空文字列の場合、エラーになること",
			"",
			true,
			"INVALID_ARGUMENT",
		),
		Entry(
			"36文字以上の場合、エラーになること",
			strings.Repeat("a", 36),
			true,
			"INVALID_ARGUMENT",
		),
		Entry(
			"UUID形式ではない場合、エラーになること",
			strings.Repeat("a", 30),
			true,
			"INVALID_ARGUMENT",
		),
		Entry(
			"正常なUUID形式の場合、エラーにならないこと",
			uuid.New().String(),
			false,
			"",
		),
	)

	DescribeTable("商品名のバリデーション",
		func(name string, expectError bool, expectedErrorCode string) {
			productName, err := NewProductName(name)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(productName).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(productName).NotTo(BeNil())
			}
		},
		Entry(
			"空文字列の場合、エラーになること",
			"",
			true,
			"INVALID_ARGUMENT",
		),
		Entry(
			"100文字以上の場合、エラーになること",
			strings.Repeat("a", 101),
			true,
			"INVALID_ARGUMENT",
		),
	)

	DescribeTable("商品価格のバリデーション",
		func(price uint32, expectError bool, expectedErrorCode string) {
			productPrice, err := NewProductPrice(price)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(productPrice).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(productPrice).NotTo(BeNil())
			}
		},
		Entry(
			"0円の場合、エラーになること",
			uint32(0),
			true,
			"INVALID_ARGUMENT",
		),
		Entry(
			"1,000,001円以上の場合、エラーになること",
			uint32(1000001),
			true,
			"INVALID_ARGUMENT",
		),
		Entry(
			"1円以上1,000,000円以下の場合、エラーにならないこと",
			uint32(1000000),
			false,
			"",
		),
	)
})

var _ = Describe("Productエンティティオブジェクト", Label("Productエンティティ"), func() {
	var (
		validProductName  *ProductName
		validProductPrice *ProductPrice
		validCategory     *categories.Category
	)

	BeforeEach(func() {
		var err error
		validProductName, err = NewProductName("テスト商品")
		Expect(err).NotTo(HaveOccurred())
		validProductPrice, err = NewProductPrice(uint32(1000))
		Expect(err).NotTo(HaveOccurred())
		categoryName, err := categories.NewCategoryName("カテゴリ1")
		Expect(err).NotTo(HaveOccurred())
		validCategory, err = categories.NewCategory(categoryName)
		Expect(err).NotTo(HaveOccurred())
	})

	DescribeTable("NewProduct",
		func(setupFunc func() (*ProductName, *ProductPrice, *categories.Category), expectError bool, expectedErrorCode string) {
			name, price, category := setupFunc()
			product, err := NewProduct(name, price, category)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(product).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(product).NotTo(BeNil())
				Expect(product.Id()).NotTo(BeNil())
				Expect(product.Name()).To(Equal(name))
				Expect(product.Price()).To(Equal(price))
				Expect(product.Category()).To(Equal(category))
			}
		},
		Entry(
			"正常な引数で商品を生成できること",
			func() (*ProductName, *ProductPrice, *categories.Category) {
				name, _ := NewProductName("テスト商品")
				price, _ := NewProductPrice(uint32(1000))
				categoryName, _ := categories.NewCategoryName("カテゴリ1")
				category, _ := categories.NewCategory(categoryName)
				return name, price, category
			},
			false,
			"",
		),
	)

	DescribeTable("BuildProduct",
		func(setupFunc func() (*ProductId, *ProductName, *ProductPrice, *categories.Category), expectError bool) {
			id, name, price, category := setupFunc()
			product, err := BuildProduct(id, name, price, category)

			if expectError {
				Expect(err).To(HaveOccurred())
				Expect(product).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(product).NotTo(BeNil())
				Expect(product.Id()).To(Equal(id))
				Expect(product.Name()).To(Equal(name))
				Expect(product.Price()).To(Equal(price))
				Expect(product.Category()).To(Equal(category))
			}
		},
		Entry(
			"正常な引数で商品を再構築できること",
			func() (*ProductId, *ProductName, *ProductPrice, *categories.Category) {
				id, _ := NewProductId(uuid.New().String())
				name, _ := NewProductName("テスト商品")
				price, _ := NewProductPrice(uint32(1000))
				categoryName, _ := categories.NewCategoryName("カテゴリ1")
				category, _ := categories.NewCategory(categoryName)
				return id, name, price, category
			},
			false,
		),
	)

	DescribeTable("ChangeName",
		func(newNameStr string, expectValidNewName bool) {
			product, err := NewProduct(validProductName, validProductPrice, validCategory)
			Expect(err).NotTo(HaveOccurred())

			newName, err := NewProductName(newNameStr)
			if !expectValidNewName {
				Expect(err).To(HaveOccurred())
				return
			}
			Expect(err).NotTo(HaveOccurred())

			product.ChangeName(newName)
			Expect(product.Name()).To(Equal(newName))
		},
		Entry(
			"商品名を変更できること",
			"新しい商品名",
			true,
		),
		Entry(
			"別の商品名に変更できること",
			"さらに新しい商品名",
			true,
		),
	)

	DescribeTable("ChangePrice",
		func(newPriceValue uint32, expectValidNewPrice bool) {
			product, err := NewProduct(validProductName, validProductPrice, validCategory)
			Expect(err).NotTo(HaveOccurred())

			newPrice, err := NewProductPrice(newPriceValue)
			if !expectValidNewPrice {
				Expect(err).To(HaveOccurred())
				return
			}
			Expect(err).NotTo(HaveOccurred())

			product.ChangePrice(newPrice)
			Expect(product.Price()).To(Equal(newPrice))
		},
		Entry(
			"商品価格を変更できること",
			uint32(2000),
			true,
		),
		Entry(
			"別の商品価格に変更できること",
			uint32(500000),
			true,
		),
	)

	DescribeTable("ChangeCategory",
		func(newCategoryNameStr string, expectValidNewCategory bool) {
			product, err := NewProduct(validProductName, validProductPrice, validCategory)
			Expect(err).NotTo(HaveOccurred())

			newCategoryName, err := categories.NewCategoryName(newCategoryNameStr)
			if !expectValidNewCategory {
				Expect(err).To(HaveOccurred())
				return
			}
			Expect(err).NotTo(HaveOccurred())
			newCategory, err := categories.NewCategory(newCategoryName)
			Expect(err).NotTo(HaveOccurred())

			product.ChangeCategory(newCategory)
			Expect(product.Category()).To(Equal(newCategory))
		},
		Entry(
			"カテゴリを変更できること",
			"カテゴリ2",
			true,
		),
		Entry(
			"別のカテゴリに変更できること",
			"カテゴリ3",
			true,
		),
	)

	DescribeTable("Equals",
		func(setupFunc func() (*Product, *Product), expectEqual bool, expectError bool, expectedErrorCode string) {
			product1, product2 := setupFunc()
			equal, err := product1.Equals(product2)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(equal).To(BeFalse())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(equal).To(Equal(expectEqual))
			}
		},
		Entry(
			"同じIDを持つ商品は等しいこと",
			func() (*Product, *Product) {
				productId, _ := NewProductId(uuid.New().String())
				name, _ := NewProductName("テスト商品")
				price, _ := NewProductPrice(uint32(1000))
				categoryName, _ := categories.NewCategoryName("カテゴリ1")
				category, _ := categories.NewCategory(categoryName)

				product1, _ := BuildProduct(productId, name, price, category)
				product2, _ := BuildProduct(productId, name, price, category)
				return product1, product2
			},
			true,
			false,
			"",
		),
		Entry(
			"異なるIDを持つ商品は等しくないこと",
			func() (*Product, *Product) {
				name, _ := NewProductName("テスト商品")
				price, _ := NewProductPrice(uint32(1000))
				categoryName, _ := categories.NewCategoryName("カテゴリ1")
				category, _ := categories.NewCategory(categoryName)

				product1, _ := NewProduct(name, price, category)
				product2, _ := NewProduct(name, price, category)
				return product1, product2
			},
			false,
			false,
			"",
		),
		Entry(
			"nilとの比較でエラーになること",
			func() (*Product, *Product) {
				name, _ := NewProductName("テスト商品")
				price, _ := NewProductPrice(uint32(1000))
				categoryName, _ := categories.NewCategoryName("カテゴリ1")
				category, _ := categories.NewCategory(categoryName)

				product, _ := NewProduct(name, price, category)
				return product, nil
			},
			false,
			true,
			"INVALID_ARGUMENT",
		),
	)
})
