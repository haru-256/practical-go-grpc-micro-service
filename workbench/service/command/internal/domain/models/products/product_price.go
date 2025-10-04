package products

import (
	"fmt"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// ProductPrice は商品価格を表す値オブジェクトです。
// 1円以上1,000,000円以下の価格を保持します。
type ProductPrice struct {
	value uint32 // 商品価格(単位: 円)
}

// Value は商品価格の値を返します（単位: 円）。
func (p *ProductPrice) Value() uint32 {
	return p.value
}

// NewProductPrice は商品価格を生成します。
// 1円以上1,000,000円以下の価格である必要があります。
func NewProductPrice(value uint32) (*ProductPrice, error) {
	const MIN_VALUE uint32 = 1       // 最小値(1円)
	const MAX_VALUE uint32 = 1000000 // 最大値(1,000,000円)

	if value < MIN_VALUE || value > MAX_VALUE {
		return nil, errs.NewDomainError(
			"INVALID_ARGUMENT",
			fmt.Sprintf("商品価格は%d円以上%d円以下で入力してください", MIN_VALUE, MAX_VALUE),
		)
	}

	return &ProductPrice{value: value}, nil
}
