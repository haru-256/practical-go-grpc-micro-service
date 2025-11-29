package categories

import (
	"fmt"
	"regexp"
	"sync"
	"unicode/utf8"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
)

// CategoryId はカテゴリIDを表す値オブジェクトです。
type CategoryId struct {
	value string
}

// Value はカテゴリIDの値を返します。
func (c *CategoryId) Value() string {
	return c.value
}

// Equals は2つのカテゴリIDの同一性を検証します。
func (c *CategoryId) Equals(other *CategoryId) bool {
	if other == nil {
		return false
	}
	if c == other { // アドレスが同じ?
		return true
	}
	// 値が同じ
	return c.value == other.Value()
}

var getUUIDRegexp = sync.OnceValue(func() *regexp.Regexp {
	const REGEXP string = "^([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})$"
	return regexp.MustCompile(REGEXP)
})

// NewCategoryId はカテゴリIDを生成します。
func NewCategoryId(value string) (*CategoryId, error) {
	// フィールドの長さ
	const LENGTH int = 36

	// 引数の文字数チェック
	// 引数の文字数チェック
	if utf8.RuneCountInString(value) != LENGTH {
		// TODO: 専用のエラーオブジェクトを作成する
		return nil, errs.NewDomainError(
			"INVALID_ARGUMENT", fmt.Sprintf("カテゴリIDの長さは%d文字である必要があります", LENGTH),
		)
	}
	// UUIDの形式チェック
	if !getUUIDRegexp().MatchString(value) {
		return nil, errs.NewDomainError("INVALID_ARGUMENT", "カテゴリIDはUUIDの形式である必要があります")
	}
	return &CategoryId{value: value}, nil
}
