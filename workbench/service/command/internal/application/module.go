package application

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/impl"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler"
	"go.uber.org/fx"
)

// Module はアプリケーション層のFxモジュールです。
// このモジュールは以下を提供します:
//   - カテゴリサービスの実装（NewCategoryServiceImpl → service.CategoryService）
//   - 商品サービスの実装（NewProductServiceImpl → service.ProductService）
//   - SQLBoilerインフラストラクチャ層（sqlboiler.Module）
//
// 各サービスは具象型（*Impl）として作成され、fx.Annotateによって
// 対応するインターフェース型に変換されます。
var Module = fx.Module(
	"application",
	sqlboiler.Module,
	fx.Provide(
		fx.Annotate(
			impl.NewCategoryServiceImpl,
			fx.As(new(service.CategoryService)),
		),
		fx.Annotate(
			impl.NewProductServiceImpl,
			fx.As(new(service.ProductService)),
		),
	),
)
