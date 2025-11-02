package categories

import (
	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
)

// Category はカテゴリエンティティを表すドメインオブジェクトです。
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
func (c *Category) Equals(other *Category) (bool, error) {
	if other == nil {
		return false, errs.NewDomainError("INVALID_ARGUMENT", "比較対象のCategoryがnilです")
	}
	return c.id.Equals(other.Id()), nil
}

// NewCategory は新しいカテゴリエンティティを生成します。
func NewCategory(name *CategoryName) (*Category, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, errs.NewDomainErrorWithCause("INTERNAL", "カテゴリIDの生成に失敗しました", err)
	}

	id, err := NewCategoryId(uid.String())
	if err != nil {
		return nil, errs.NewDomainErrorWithCause("INTERNAL", "カテゴリIDの生成に失敗しました", err)
	}

	return &Category{id: id, name: name}, nil
}

// BuildCategory は既存のカテゴリIDを使用してカテゴリエンティティを再構築します。
func BuildCategory(id *CategoryId, name *CategoryName) (*Category, error) {
	category := Category{
		id:   id,
		name: name,
	}
	return &category, nil
}
