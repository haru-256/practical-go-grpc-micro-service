package models

// カテゴリエンティティ
type Category struct {
	id   string
	name string
}

// NewCategory はCategoryを生成します。
//
// Parameters:
//   - id: カテゴリID
//   - name: カテゴリ名
//
// Returns:
//   - *Category: Categoryポインタ
func NewCategory(id string, name string) *Category {
	return &Category{id: id, name: name}
}

// Id はカテゴリIDを返します。
//
// Returns:
//   - string: カテゴリID
func (c *Category) Id() string {
	return c.id
}

// Name はカテゴリ名を返します。
//
// Returns:
//   - string: カテゴリ名
func (c *Category) Name() string {
	return c.name
}
