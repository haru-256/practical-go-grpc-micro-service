package categories

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEntityPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "domain/models/categories packageのテスト")
}

var _ = Describe("Categoryエンティティを構成する値オブジェクト", Label("CategoryId構造体の生成"), func() {
	DescribeTable("CategoryIDのバリデーション",
		func(input string, expectError bool, expectedErrorCode string, expectedErrorMsg string) {
			categoryId, err := NewCategoryId(input)
			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr).To(Equal(errs.NewDomainError(expectedErrorCode, expectedErrorMsg)))
				Expect(categoryId).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(categoryId).NotTo(BeNil())
			}
		},
		Entry(
			"空文字列の場合、エラーになること",
			"",
			true,
			"INVALID_ARGUMENT",
			"カテゴリIDの長さは36文字である必要があります",
		),
		Entry(
			"36文字以上の場合、エラーになること",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			true,
			"INVALID_ARGUMENT",
			"カテゴリIDの長さは36文字である必要があります",
		),
		Entry(
			"UUID形式ではない場合、エラーになること",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			true,
			"INVALID_ARGUMENT",
			"カテゴリIDはUUIDの形式である必要があります",
		),
		Entry(
			"正常なUUID形式の場合、エラーにならないこと",
			uuid.New().String(),
			false,
			"",
			"",
		),
	)

	DescribeTable("商品名のバリデーション",
		func(name string, expectError bool, expectedErrorCode string) {
			categoryName, err := NewCategoryName(name)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(categoryName).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(categoryName).NotTo(BeNil())
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
})

var _ = Describe("Categoryエンティティオブジェクト", Label("Categoryエンティティ"), func() {
	var (
		validCategoryName *CategoryName
	)

	BeforeEach(func() {
		var err error
		validCategoryName, err = NewCategoryName("テストカテゴリ")
		Expect(err).NotTo(HaveOccurred())
	})

	DescribeTable("NewCategory",
		func(setupFunc func() *CategoryName, expectError bool, expectedErrorCode string) {
			name := setupFunc()
			category, err := NewCategory(name)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr.Code).To(Equal(expectedErrorCode))
				Expect(category).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(category).NotTo(BeNil())
				Expect(category.Id()).NotTo(BeNil())
				Expect(category.Name()).To(Equal(name))
			}
		},
		Entry(
			"正常な引数でカテゴリを生成できること",
			func() *CategoryName {
				name, _ := NewCategoryName("テストカテゴリ")
				return name
			},
			false,
			"",
		),
		Entry(
			"別の名前でカテゴリを生成できること",
			func() *CategoryName {
				name, _ := NewCategoryName("別のカテゴリ")
				return name
			},
			false,
			"",
		),
	)

	DescribeTable("BuildCategory",
		func(setupFunc func() (*CategoryId, *CategoryName), expectError bool) {
			id, name := setupFunc()
			category, err := BuildCategory(id, name)

			if expectError {
				Expect(err).To(HaveOccurred())
				Expect(category).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(category).NotTo(BeNil())
				Expect(category.Id()).To(Equal(id))
				Expect(category.Name()).To(Equal(name))
			}
		},
		Entry(
			"正常な引数でカテゴリを再構築できること",
			func() (*CategoryId, *CategoryName) {
				id, _ := NewCategoryId(uuid.New().String())
				name, _ := NewCategoryName("テストカテゴリ")
				return id, name
			},
			false,
		),
		Entry(
			"別のIDと名前でカテゴリを再構築できること",
			func() (*CategoryId, *CategoryName) {
				id, _ := NewCategoryId(uuid.New().String())
				name, _ := NewCategoryName("別のカテゴリ")
				return id, name
			},
			false,
		),
	)

	DescribeTable("ChangeName",
		func(newNameStr string, expectValidNewName bool) {
			category, err := NewCategory(validCategoryName)
			Expect(err).NotTo(HaveOccurred())

			newName, err := NewCategoryName(newNameStr)
			if !expectValidNewName {
				Expect(err).To(HaveOccurred())
				return
			}
			Expect(err).NotTo(HaveOccurred())

			category.ChangeName(newName)
			Expect(category.Name()).To(Equal(newName))
		},
		Entry(
			"カテゴリ名を変更できること",
			"新しいカテゴリ名",
			true,
		),
		Entry(
			"別のカテゴリ名に変更できること",
			"さらに新しいカテゴリ名",
			true,
		),
	)

	DescribeTable("Equals",
		func(setupFunc func() (*Category, *Category), expectEqual bool, expectError bool, expectedErrorCode string) {
			category1, category2 := setupFunc()
			equal, err := category1.Equals(category2)

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
			"同じIDを持つカテゴリは等しいこと",
			func() (*Category, *Category) {
				categoryId, _ := NewCategoryId(uuid.New().String())
				name, _ := NewCategoryName("テストカテゴリ")

				category1, _ := BuildCategory(categoryId, name)
				category2, _ := BuildCategory(categoryId, name)
				return category1, category2
			},
			true,
			false,
			"",
		),
		Entry(
			"異なるIDを持つカテゴリは等しくないこと",
			func() (*Category, *Category) {
				name, _ := NewCategoryName("テストカテゴリ")

				category1, _ := NewCategory(name)
				category2, _ := NewCategory(name)
				return category1, category2
			},
			false,
			false,
			"",
		),
		Entry(
			"nilとの比較でエラーになること",
			func() (*Category, *Category) {
				name, _ := NewCategoryName("テストカテゴリ")
				category, _ := NewCategory(name)
				return category, nil
			},
			false,
			true,
			"INVALID_ARGUMENT",
		),
	)
})
