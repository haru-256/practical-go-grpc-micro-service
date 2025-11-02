package db

type Product struct {
	Id         int    `gorm:"column:id;primaryKey"`
	ObjId      string `gorm:"column:obj_id;primaryKey"`
	Name       string `gorm:"column:name"`
	Price      uint32 `gorm:"column:price"`
	CategoryId string `gorm:"column:category_id"`

	Category Category `gorm:"foreignKey:CategoryId;references:ObjId"`
}

func (Product) TableName() string {
	return "product"
}

type Category struct {
	Id    int    `gorm:"column:id;primaryKey"`
	ObjId string `gorm:"column:obj_id;primaryKey"`
	Name  string `gorm:"column:name"`
}

func (Category) TableName() string {
	return "category"
}
