package products

import (
	"fmt"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

type ProductPrice struct {
	value uint32 // 商品価格(単位: 円)
}

// valueフィールドのゲッター
func (p *ProductPrice) Value() uint32 {
	return p.value
}

// コンストラクタ
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
