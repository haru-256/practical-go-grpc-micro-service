package models

// Product は商品エンティティ
type Product struct {
	id       string    // 商品ID
	name     string    // 商品名
	price    uint32    // 価格
	category *Category // カテゴリ
}

// NewProduct はProductを生成します。
//
// Parameters:
//   - id: 商品ID
//   - name: 商品名
//   - price: 価格
//   - category: カテゴリ
//
// Returns:
//   - *Product: Productポインタ
func NewProduct(id string, name string, price uint32, category *Category) *Product {
	return &Product{id: id, name: name, price: price, category: category}
}

// Id は商品IDを返します。
//
// Returns:
//   - string: 商品ID
func (p *Product) Id() string {
	return p.id
}

// Name は商品名を返します。
//
// Returns:
//   - string: 商品名
func (p *Product) Name() string {
	return p.name
}

// Price は価格を返します。
//
// Returns:
//   - uint32: 価格
func (p *Product) Price() uint32 {
	return p.price
}

// Category はカテゴリを返します。
//
// Returns:
//   - *Category: Categoryポインタ
func (p *Product) Category() *Category {
	return p.category
}
