package application

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/impl"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure"
	"go.uber.org/fx"
)

// Module はアプリケーション層のFxモジュールです。
var Module = fx.Module(
	"application",
	infrastructure.Module,
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
