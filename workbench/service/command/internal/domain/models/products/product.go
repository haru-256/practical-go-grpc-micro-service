package products

import (
	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// Product は商品エンティティを表すドメインオブジェクトです。
// 商品は、ID、名前、価格、カテゴリを持ち、ビジネスロジックをカプセル化します。
type Product struct {
	id       *ProductId           // 商品ID
	name     *ProductName         // 商品名
	price    *ProductPrice        // 商品価格(単位: 円)
	category *categories.Category // カテゴリ
}

// Id は商品IDを返します。
func (p *Product) Id() *ProductId {
	return p.id
}

// Name は商品名を返します。
func (p *Product) Name() *ProductName {
	return p.name
}

// Price は商品価格を返します。
func (p *Product) Price() *ProductPrice {
	return p.price
}

// Category はカテゴリを返します。
func (p *Product) Category() *categories.Category {
	return p.category
}

// ChangeName は商品名を変更します。
func (p *Product) ChangeName(name *ProductName) {
	p.name = name
}

// ChangePrice は商品価格を変更します。
func (p *Product) ChangePrice(price *ProductPrice) {
	p.price = price
}

// ChangeCategory はカテゴリを変更します。
func (p *Product) ChangeCategory(category *categories.Category) {
	p.category = category
}

// Equals は2つの商品エンティティの同一性を検証します。
// 商品IDが一致する場合、同一の商品として扱います。
// otherがnilの場合はエラーを返します。
func (p *Product) Equals(other *Product) (bool, error) {
	if other == nil {
		return false, errs.NewDomainError("INVALID_ARGUMENT", "比較対象のProductがnilです")
	}
	return p.id.Equals(other.Id()), nil
}

// NewProduct は新しい商品エンティティを生成します。
// 商品IDは自動的にUUIDとして生成されます。
// name: 商品名
// price: 商品価格
// category: カテゴリ
func NewProduct(name *ProductName, price *ProductPrice, category *categories.Category) (*Product, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, errs.NewDomainErrorWithCause("INTERNAL", "商品IDの生成に失敗しました", err)
	}

	id, err := NewProductId(uid.String())
	if err != nil {
		return nil, errs.NewDomainErrorWithCause("INTERNAL", "商品IDの生成に失敗しました", err)
	}

	return &Product{id: id, name: name, price: price, category: category}, nil
}

// BuildProduct は既存の商品IDを使用して商品エンティティを再構築します。
// データベースから取得した商品データを復元する際などに使用します。
// id: 商品ID
// name: 商品名
// price: 商品価格
// category: カテゴリ
func BuildProduct(id *ProductId, name *ProductName, price *ProductPrice, category *categories.Category) (*Product, error) {
	product := Product{
		id:       id,
		name:     name,
		price:    price,
		category: category,
	}
	return &product, nil
}
