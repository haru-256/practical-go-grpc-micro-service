package products

import (
	"fmt"
	"unicode/utf8"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

type ProductName struct {
	value string // 商品名
}

// valueフィールドのゲッター
func (p *ProductName) Value() string {
	return p.value
}

// コンストラクタ
func NewProductName(value string) (*ProductName, error) {
	const MIN_LENGTH int = 1   // 最小文字数
	const MAX_LENGTH int = 100 // 最大文字数

	if count := utf8.RuneCountInString(value); count < MIN_LENGTH || count > MAX_LENGTH {
		return nil, errs.NewDomainError(fmt.Sprintf("商品名は%d文字以上%d文字以下で入力してください", MIN_LENGTH, MAX_LENGTH))
	}

	return &ProductName{value: value}, nil
}
