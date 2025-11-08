package presentation

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation/server"
	"go.uber.org/fx"
)

// Module はインフラストラクチャ層のFxモジュールです。
var Module = fx.Module(
	"presentation",
	infrastructure.Module,
	fx.Provide(
		server.NewCQRSServiceConfig,
		server.NewCQRSServiceHandler,
		server.NewCQRSServiceServer,
	),
	// ライフサイクルフックを登録
	fx.Invoke(server.RegisterLifecycleHooks),
)
