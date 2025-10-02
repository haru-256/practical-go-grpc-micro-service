package products

import (
	"testing"

	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEntityPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "domain/models/products packageのテスト")
}

var _ = Describe("Productエンティティを構成する値オブジェクト", Label("ProductId構造体の生成"), func() {
	DescribeTable("商品IDのバリデーション",
		func(input string, expectError bool, expectedErrorCode string, expectedErrorMsg string) {
			productId, err := NewProductId(input)

			if expectError {
				Expect(err).To(HaveOccurred())
				domainErr, ok := err.(*errs.DomainError)
				Expect(ok).To(BeTrue())
				Expect(domainErr).To(Equal(errs.NewDomainError(expectedErrorCode, expectedErrorMsg)))
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
			"商品IDの長さは36文字である必要があります",
		),
		Entry(
			"36文字以上の場合、エラーになること",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			true,
			"INVALID_ARGUMENT",
			"商品IDの長さは36文字である必要があります",
		),
		Entry(
			"UUID形式ではない場合、エラーになること",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			true,
			"INVALID_ARGUMENT",
			"商品IDはUUIDの形式である必要があります",
		),
		Entry(
			"正常なUUID形式の場合、エラーにならないこと",
			uuid.New().String(),
			false,
			"",
			"",
		),
	)
})
