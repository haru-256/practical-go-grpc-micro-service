package categories

import (
	"fmt"
	"unicode/utf8"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
)

// CategoryName はカテゴリ名を表す値オブジェクトです。
type CategoryName struct {
	value string // カテゴリ名
}

// Value はカテゴリ名の値を返します。
func (c *CategoryName) Value() string {
	return c.value
}

// NewCategoryName はカテゴリ名を生成します。
func NewCategoryName(value string) (*CategoryName, error) {
	const MIN_LENGTH int = 1  // 最小文字数
	const MAX_LENGTH int = 20 // 最大文字数

	if count := utf8.RuneCountInString(value); count < MIN_LENGTH || count > MAX_LENGTH {
		return nil, errs.NewDomainError(
			"INVALID_ARGUMENT",
			fmt.Sprintf("カテゴリ名は%d文字以上%d文字以下で入力してください。Got: len(%s) = %d", MIN_LENGTH, MAX_LENGTH, value, count),
		)
	}

	return &CategoryName{value: value}, nil
}
