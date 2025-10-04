package categories

import (
	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// Category はカテゴリエンティティを表すドメインオブジェクトです。
// カテゴリは、IDと名前を持ち、商品を分類するために使用されます。
type Category struct {
	id   *CategoryId   // カテゴリID
	name *CategoryName // カテゴリ名
}

// Id はカテゴリIDを返します。
func (c *Category) Id() *CategoryId {
	return c.id
}

// Name はカテゴリ名を返します。
func (c *Category) Name() *CategoryName {
	return c.name
}

// ChangeName はカテゴリ名を変更します。
func (c *Category) ChangeName(name *CategoryName) {
	c.name = name
}

// Equals は2つのカテゴリエンティティの同一性を検証します。
// カテゴリIDが一致する場合、同一のカテゴリとして扱います。
// otherがnilの場合はエラーを返します。
func (c *Category) Equals(other *Category) (bool, error) {
	if other == nil {
		return false, errs.NewDomainError("INVALID_ARGUMENT", "比較対象のCategoryがnilです")
	}
	return c.id.Equals(other.Id()), nil
}

// NewCategory は新しいカテゴリエンティティを生成します。
// カテゴリIDは自動的にUUIDとして生成されます。
// name: カテゴリ名
func NewCategory(name *CategoryName) (*Category, error) {
	if uid, err := uuid.NewRandom(); err != nil {
		return nil, errs.NewDomainErrorWithCause("INTERNAL", "カテゴリIDの生成に失敗しました", err)
	} else {
		if id, err := NewCategoryId(uid.String()); err != nil {
			return nil, errs.NewDomainErrorWithCause("INTERNAL", "カテゴリIDの生成に失敗しました", err)
		} else {
			// 生成に成功した場合は何もしない
			return &Category{id: id, name: name}, nil
		}
	}
}

// BuildCategory は既存のカテゴリIDを使用してカテゴリエンティティを再構築します。
// データベースから取得したカテゴリデータを復元する際などに使用します。
// id: カテゴリID
// name: カテゴリ名
func BuildCategory(id *CategoryId, name *CategoryName) (*Category, error) {
	category := Category{
		id:   id,
		name: name,
	}
	return &category, nil
}
