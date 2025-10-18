package config

import "go.uber.org/fx"

var Module = fx.Module(
	"config",
	fx.Provide(
		fx.Annotate(
			func() (string, string) {
				return "../../", "config"
			},
			fx.ResultTags(`name:"configPath"`, `name:"configName"`),
		),
		fx.Annotate(
			NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
	),
)
