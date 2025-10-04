package products

import (
	"fmt"
	"regexp"
	"sync"
	"unicode/utf8"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// ProductId は商品IDを表す値オブジェクトです。
// UUID形式の文字列を保持し、商品の一意な識別子として機能します。
type ProductId struct {
	value string // 商品番号(UUID)
}

// Value は商品IDの値を返します。
func (p *ProductId) Value() string {
	return p.value
}

// Equals は2つの商品IDの同一性を検証します。
// 値が一致する場合、またはアドレスが同じ場合にtrueを返します。
func (p *ProductId) Equals(other *ProductId) bool {
	if other == nil {
		return false
	}
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

// NewProductId は商品IDを生成します。
// 引数はUUID形式の文字列である必要があります（36文字）。
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
	// UUIDの形式チェック
	if !getUUIDRegexp().MatchString(value) {
		return nil, errs.NewDomainError("INVALID_ARGUMENT", "商品IDはUUIDの形式である必要があります")
	}

	return &ProductId{value: value}, nil
}
