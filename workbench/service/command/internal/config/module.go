package config

import "go.uber.org/fx"

// Module は設定管理のFxモジュールです。
var Module = fx.Module(
	"config",
	fx.Provide(
		fx.Annotate(
			NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
	),
)
