package server

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	cmd "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1"
	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
	common "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/common/v1"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// createCategoryFromDTO はDTOからProtobufのCategoryを作成します。
func createCategoryFromDTO(dto *dto.CategoryDTO) *common.Category {
	c := &common.Category{}
	c.SetId(dto.Id)
	c.SetName(dto.Name)
	return c
}

// createProductFromDTO はDTOからProtobufのProductを作成します。
func createProductFromDTO(dto *dto.ProductDTO) *common.Product {
	c := createCategoryFromDTO(dto.Category)
	p := &common.Product{}
	p.SetId(dto.Id)
	p.SetName(dto.Name)
	p.SetPrice(int32(dto.Price))
	p.SetCategory(c)
	return p
}

// CategoryServiceHandlerImpl はカテゴリサービスのgRPCハンドラー実装です。
// Connect RPCを使用してカテゴリの作成、更新、削除のエンドポイントを提供します。
type CategoryServiceHandlerImpl struct {
	logger *slog.Logger
	cs     service.CategoryService
	cmdconnect.UnimplementedCategoryServiceHandler
}

// NewCategoryServiceHandlerImpl は新しいCategoryServiceHandlerImplインスタンスを生成します。
//
// Parameters:
//   - logger: ロガーインスタンス
//   - cs: カテゴリサービスのインターフェース
//
// Returns:
//   - *CategoryServiceHandlerImpl: 生成されたハンドラーインスタンス
//   - error: エラーが発生した場合（現在は常にnil）
func NewCategoryServiceHandlerImpl(logger *slog.Logger, cs service.CategoryService) (*CategoryServiceHandlerImpl, error) {
	return &CategoryServiceHandlerImpl{
		logger: logger,
		cs:     cs,
	}, nil
}

// CreateCategory は新しいカテゴリを作成します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - req: カテゴリ作成リクエスト（カテゴリ名を含む）
//
// Returns:
//   - *connect.Response[cmd.CreateCategoryResponse]: 作成されたカテゴリ情報を含むレスポンス
//   - error: バリデーションエラーの場合はCodeInvalidArgument、サービス層エラーの場合はCodeInternal
func (s *CategoryServiceHandlerImpl) CreateCategory(ctx context.Context, req *connect.Request[cmd.CreateCategoryRequest]) (*connect.Response[cmd.CreateCategoryResponse], error) {
	createCategoryDTO := &dto.CreateCategoryDTO{
		Name: req.Msg.GetName().GetValue(),
	}

	categoryDTO, err := s.cs.Add(ctx, createCategoryDTO)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("create error: %w", err))
	}

	res := &cmd.CreateCategoryResponse{}
	res.SetCategory(createCategoryFromDTO(categoryDTO))
	res.SetTimestamp(timestamppb.Now())

	return connect.NewResponse(res), nil
}

// UpdateCategory は既存のカテゴリを更新します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - req: カテゴリ更新リクエスト（カテゴリIDと新しい名前を含む）
//
// Returns:
//   - *connect.Response[cmd.UpdateCategoryResponse]: 更新されたカテゴリ情報を含むレスポンス
//   - error: バリデーションエラーの場合はCodeInvalidArgument、サービス層エラーの場合はCodeInternal
func (s *CategoryServiceHandlerImpl) UpdateCategory(ctx context.Context, req *connect.Request[cmd.UpdateCategoryRequest]) (*connect.Response[cmd.UpdateCategoryResponse], error) {
	updateCategoryDTO := &dto.UpdateCategoryDTO{
		Id:   req.Msg.GetCategory().GetId().GetValue(),
		Name: req.Msg.GetCategory().GetName().GetValue(),
	}

	categoryDTO, err := s.cs.Update(ctx, updateCategoryDTO)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("update error: %w", err))
	}

	res := &cmd.UpdateCategoryResponse{}
	res.SetCategory(createCategoryFromDTO(categoryDTO))
	res.SetTimestamp(timestamppb.Now())

	return connect.NewResponse(res), nil
}

// DeleteCategory は既存のカテゴリを削除します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - req: カテゴリ削除リクエスト（カテゴリIDを含む）
//
// Returns:
//   - *connect.Response[cmd.DeleteCategoryResponse]: 削除されたカテゴリ情報を含むレスポンス
//   - error: バリデーションエラーの場合はCodeInvalidArgument、サービス層エラーの場合はCodeInternal
func (s *CategoryServiceHandlerImpl) DeleteCategory(ctx context.Context, req *connect.Request[cmd.DeleteCategoryRequest]) (*connect.Response[cmd.DeleteCategoryResponse], error) {
	deleteCategoryDTO := &dto.DeleteCategoryDTO{
		Id: req.Msg.GetCategoryId().GetValue(),
	}

	categoryDTO, err := s.cs.Delete(ctx, deleteCategoryDTO)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("delete error: %w", err))
	}

	res := &cmd.DeleteCategoryResponse{}
	res.SetCategory(createCategoryFromDTO(categoryDTO))
	res.SetTimestamp(timestamppb.Now())

	return connect.NewResponse(res), nil
}

// ProductServiceHandlerImpl は商品サービスのgRPCハンドラー実装です。
// Connect RPCを使用して商品の作成、更新、削除のエンドポイントを提供します。
type ProductServiceHandlerImpl struct {
	logger *slog.Logger
	ps     service.ProductService
	cmdconnect.UnimplementedProductServiceHandler
}

// NewProductServiceHandlerImpl は新しいProductServiceHandlerImplインスタンスを生成します。
//
// Parameters:
//   - logger: ロガーインスタンス
//   - ps: 商品サービスのインターフェース
//
// Returns:
//   - *ProductServiceHandlerImpl: 生成されたハンドラーインスタンス
//   - error: エラーが発生した場合（現在は常にnil）
func NewProductServiceHandlerImpl(logger *slog.Logger, ps service.ProductService) (*ProductServiceHandlerImpl, error) {
	return &ProductServiceHandlerImpl{
		logger: logger,
		ps:     ps,
	}, nil
}

// CreateProduct は新しい商品を作成します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - req: 商品作成リクエスト（商品名、価格、カテゴリ情報を含む）
//
// Returns:
//   - *connect.Response[cmd.CreateProductResponse]: 作成された商品情報を含むレスポンス
//   - error: バリデーションエラーの場合はCodeInvalidArgument、サービス層エラーの場合はCodeInternal
func (s *ProductServiceHandlerImpl) CreateProduct(ctx context.Context, req *connect.Request[cmd.CreateProductRequest]) (*connect.Response[cmd.CreateProductResponse], error) {
	createProductDTO := &dto.CreateProductDTO{
		Name:  req.Msg.GetProduct().GetName().GetValue(),
		Price: uint32(req.Msg.GetProduct().GetPrice().GetValue()),
		Category: &dto.CategoryDTO{
			Id:   req.Msg.GetProduct().GetCategory().GetId().GetValue(),
			Name: req.Msg.GetProduct().GetCategory().GetName().GetValue(),
		},
	}

	productDTO, err := s.ps.Add(ctx, createProductDTO)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("create error: %w", err))
	}

	res := &cmd.CreateProductResponse{}
	res.SetProduct(createProductFromDTO(productDTO))
	res.SetTimestamp(timestamppb.Now())

	return connect.NewResponse(res), nil
}

// UpdateProduct は既存の商品を更新します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - req: 商品更新リクエスト（商品ID、名前、価格、カテゴリ情報を含む）
//
// Returns:
//   - *connect.Response[cmd.UpdateProductResponse]: 更新された商品情報を含むレスポンス
//   - error: バリデーションエラーの場合はCodeInvalidArgument、サービス層エラーの場合はCodeInternal
func (s *ProductServiceHandlerImpl) UpdateProduct(ctx context.Context, req *connect.Request[cmd.UpdateProductRequest]) (*connect.Response[cmd.UpdateProductResponse], error) {
	updateProductDTO := &dto.UpdateProductDTO{
		Id:         req.Msg.GetProduct().GetId().GetValue(),
		Name:       req.Msg.GetProduct().GetName().GetValue(),
		Price:      uint32(req.Msg.GetProduct().GetPrice().GetValue()),
		CategoryId: req.Msg.GetProduct().GetCategoryId().GetValue(),
	}

	productDTO, err := s.ps.Update(ctx, updateProductDTO)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("update error: %w", err))
	}

	res := &cmd.UpdateProductResponse{}
	res.SetProduct(createProductFromDTO(productDTO))
	res.SetTimestamp(timestamppb.Now())

	return connect.NewResponse(res), nil
}

// DeleteProduct は既存の商品を削除します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - req: 商品削除リクエスト（商品IDを含む）
//
// Returns:
//   - *connect.Response[cmd.DeleteProductResponse]: 削除された商品情報を含むレスポンス
//   - error: バリデーションエラーの場合はCodeInvalidArgument、サービス層エラーの場合はCodeInternal
func (s *ProductServiceHandlerImpl) DeleteProduct(ctx context.Context, req *connect.Request[cmd.DeleteProductRequest]) (*connect.Response[cmd.DeleteProductResponse], error) {
	deleteProductDTO := &dto.DeleteProductDTO{
		Id: req.Msg.GetProductId().GetValue(),
	}

	productDTO, err := s.ps.Delete(ctx, deleteProductDTO)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("delete error: %w", err))
	}

	res := &cmd.DeleteProductResponse{}
	res.SetProduct(createProductFromDTO(productDTO))
	res.SetTimestamp(timestamppb.Now())

	return connect.NewResponse(res), nil
}
