package infrastructure

import (
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/log"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure/cqrs"
	"go.uber.org/fx"
)

// Module はインフラストラクチャ層のFxモジュールです。
var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		fx.Annotate(
			config.NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
		log.NewLogger,
		cqrs.NewCQRSServiceConfig,
		cqrs.NewClient,
		cqrs.NewCommandServiceClient,
		cqrs.NewQueryServiceClient,
		fx.Annotate(
			cqrs.NewCQRSRepositoryImpl,
			fx.As(new(repository.CQRSRepository)),
		),
	),
	// ライフサイクルフックを登録
	fx.Invoke(cqrs.RegisterLifecycleHooks),
)
