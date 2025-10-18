package config

import "go.uber.org/fx"

var Module = fx.Module(
	"config",
	fx.Provide(
		fx.Annotate(
			NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
	),
)
