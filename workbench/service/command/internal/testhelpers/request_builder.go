// Package testhelpers provides common test helper functions for handler tests.
package testhelpers

import (
	"connectrpc.com/connect"
	cmd "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1"
	common "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/common/v1"
)

// CreateCategoryRequest はCreateCategoryRequestを生成するヘルパー関数です。
//
// Parameters:
//   - name: カテゴリ名
//
// Returns:
//   - *connect.Request[cmd.CreateCategoryRequest]: 生成されたリクエスト
func CreateCategoryRequest(name string) *connect.Request[cmd.CreateCategoryRequest] {
	categoryName := &common.CategoryName{}
	categoryName.SetValue(name)
	createCategoryReq := &cmd.CreateCategoryRequest{}
	createCategoryReq.SetName(categoryName)
	createCategoryReq.SetCrud(cmd.CRUD_CRUD_INSERT)
	return connect.NewRequest(createCategoryReq)
}

// CreateUpdateCategoryRequest はUpdateCategoryRequestを生成するヘルパー関数です。
//
// Parameters:
//   - id: カテゴリID
//   - name: カテゴリ名
//
// Returns:
//   - *connect.Request[cmd.UpdateCategoryRequest]: 生成されたリクエスト
func CreateUpdateCategoryRequest(id, name string) *connect.Request[cmd.UpdateCategoryRequest] {
	categoryId := &common.CategoryId{}
	categoryId.SetValue(id)
	categoryName := &common.CategoryName{}
	categoryName.SetValue(name)
	category := &cmd.UpdateCategoryRequest_Category{}
	category.SetId(categoryId)
	category.SetName(categoryName)
	updateCategoryReq := &cmd.UpdateCategoryRequest{}
	updateCategoryReq.SetCategory(category)
	updateCategoryReq.SetCrud(cmd.CRUD_CRUD_UPDATE)
	return connect.NewRequest(updateCategoryReq)
}

// CreateDeleteCategoryRequest はDeleteCategoryRequestを生成するヘルパー関数です。
//
// Parameters:
//   - id: カテゴリID
//
// Returns:
//   - *connect.Request[cmd.DeleteCategoryRequest]: 生成されたリクエスト
func CreateDeleteCategoryRequest(id string) *connect.Request[cmd.DeleteCategoryRequest] {
	categoryId := &common.CategoryId{}
	categoryId.SetValue(id)
	deleteCategoryReq := &cmd.DeleteCategoryRequest{}
	deleteCategoryReq.SetCategoryId(categoryId)
	deleteCategoryReq.SetCrud(cmd.CRUD_CRUD_DELETE)
	return connect.NewRequest(deleteCategoryReq)
}

// CreateProductRequest はCreateProductRequestを生成するヘルパー関数です。
//
// Parameters:
//   - name: 商品名
//   - price: 商品価格
//   - categoryID: カテゴリID
//   - categoryName: カテゴリ名
//
// Returns:
//   - *connect.Request[cmd.CreateProductRequest]: 生成されたリクエスト
func CreateProductRequest(name string, price uint32, categoryID, categoryName string) *connect.Request[cmd.CreateProductRequest] {
	productName := &common.ProductName{}
	productName.SetValue(name)
	productPrice := &common.ProductPrice{}
	productPrice.SetValue(int32(price))
	catID := &common.CategoryId{}
	catID.SetValue(categoryID)
	catName := &common.CategoryName{}
	catName.SetValue(categoryName)

	category := &cmd.CreateProductRequest_Product_Category{}
	category.SetId(catID)
	category.SetName(catName)

	product := &cmd.CreateProductRequest_Product{}
	product.SetName(productName)
	product.SetPrice(productPrice)
	product.SetCategory(category)

	createProductReq := &cmd.CreateProductRequest{}
	createProductReq.SetProduct(product)
	createProductReq.SetCrud(cmd.CRUD_CRUD_INSERT)
	return connect.NewRequest(createProductReq)
}

// UpdateProductRequest はUpdateProductRequestを生成するヘルパー関数です。
//
// Parameters:
//   - id: 商品ID
//   - name: 商品名
//   - price: 商品価格
//   - categoryID: カテゴリID
//
// Returns:
//   - *connect.Request[cmd.UpdateProductRequest]: 生成されたリクエスト
func UpdateProductRequest(id, name string, price uint32, categoryID string) *connect.Request[cmd.UpdateProductRequest] {
	productID := &common.ProductId{}
	productID.SetValue(id)
	productName := &common.ProductName{}
	productName.SetValue(name)
	productPrice := &common.ProductPrice{}
	productPrice.SetValue(int32(price))
	catID := &common.CategoryId{}
	catID.SetValue(categoryID)

	product := &cmd.UpdateProductRequest_Product{}
	product.SetId(productID)
	product.SetName(productName)
	product.SetPrice(productPrice)
	product.SetCategoryId(catID)

	updateProductReq := &cmd.UpdateProductRequest{}
	updateProductReq.SetProduct(product)
	updateProductReq.SetCrud(cmd.CRUD_CRUD_UPDATE)
	return connect.NewRequest(updateProductReq)
}

// DeleteProductRequest はDeleteProductRequestを生成するヘルパー関数です。
//
// Parameters:
//   - id: 商品ID
//
// Returns:
//   - *connect.Request[cmd.DeleteProductRequest]: 生成されたリクエスト
func DeleteProductRequest(id string) *connect.Request[cmd.DeleteProductRequest] {
	productID := &common.ProductId{}
	productID.SetValue(id)
	deleteProductReq := &cmd.DeleteProductRequest{}
	deleteProductReq.SetProductId(productID)
	return connect.NewRequest(deleteProductReq)
}
