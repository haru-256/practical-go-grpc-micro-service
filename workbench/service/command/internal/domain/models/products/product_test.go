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

var _ = Describe("Productエンティティを構成する値オブジェクト", Ordered, Label("ProductId構造体の生成"), func() {
	// 文字数の検証
	Context("商品IDの文字数の検証", Label("文字数"), func() {
		_, err := NewProductId("")
		emptyStr, ok := err.(*errs.DomainError)
		Expect(ok).To(BeTrue())
		It("空文字列の場合、エラーになること", func() {
			Expect(emptyStr).To(Equal(errs.NewDomainError("INVALID_ARGUMENT", "商品IDの長さは36文字である必要があります")))
		})

		_, err = NewProductId("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		lengthOver, ok := err.(*errs.DomainError)
		Expect(ok).To(BeTrue())
		It("36文字以上の場合、エラーになること", func() {
			Expect(lengthOver).To(Equal(errs.NewDomainError("INVALID_ARGUMENT", "商品IDの長さは36文字である必要があります")))
		})
	})
	// UUID形式の検証
	Context("商品IDのUUID形式の検証", Label("UUID形式"), func() {
		_, err := NewProductId("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		notUUID, ok := err.(*errs.DomainError)
		Expect(ok).To(BeTrue())

		id, _ := uuid.NewRandom()
		uid := id.String()
		productId, _ := NewProductId(uid)

		It("UUID形式ではない場合、エラーになること", func() {
			Expect(notUUID).To(Equal(errs.NewDomainError("INVALID_ARGUMENT", "商品IDはUUIDの形式である必要があります")))
		})
		It("UUID形式の場合、エラーにならないこと", func() {
			id, _ := NewProductId(uid)
			Expect(productId).To(Equal(id))
		})
	})
})
