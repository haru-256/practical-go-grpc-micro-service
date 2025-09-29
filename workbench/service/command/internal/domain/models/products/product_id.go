package products

import (
	"fmt"
	"regexp"
	"sync"
	"unicode/utf8"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// 商品番号を保持する値オブジェクト(UUIDを保持する)
type ProductId struct {
	value string // 商品番号(UUID)
}

// valueフィールドのゲッター
func (p *ProductId) Value() string {
	return p.value
}

func (p *ProductId) Equals(other *ProductId) bool {
	if p == other { // アドレスが同じ?
		return true
	}
	// 値の比較
	return p.value == other.Value()
}

var getUUIDRegexp = sync.OnceValue(func() *regexp.Regexp {
	const REGEXP string = "^([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})$"
	return regexp.MustCompile(REGEXP)
})

// コンストラクタ
func NewProductId(value string) (*ProductId, error) {
	// フィールドの長さ
	const LENGTH int = 36

	// 引数の文字数チェック
	if utf8.RuneCountInString(value) != LENGTH {
		// TODO: 専用のエラーオブジェクトを作成する
		return nil, errs.NewDomainError(
			"INVALID_ARGUMENT", fmt.Sprintf("商品IDの長さは%d文字である必要があります", LENGTH),
		)
	}
	// 引数の文字数チェック
	if !getUUIDRegexp().MatchString(value) {
		return nil, errs.NewDomainError("INVALID_ARGUMENT", "商品IDはUUIDの形式である必要があります")
	}

	return &ProductId{value: value}, nil
}
