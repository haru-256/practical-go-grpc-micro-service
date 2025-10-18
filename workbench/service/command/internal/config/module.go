package config

import "go.uber.org/fx"

// Module は設定管理のFxモジュールです。
// このモジュールは以下を提供します:
//   - Viperインスタンス（*viper.Viper）
//
// 依存関係:
//   - configPath (string): 設定ファイルの検索パス（name:"configPath"）
//   - configName (string): 設定ファイルの名前（name:"configName"）
//
// これらの名前付きパラメータは、このモジュールを使用する際に
// 外部から提供される必要があります。
var Module = fx.Module(
	"config",
	fx.Provide(
		fx.Annotate(
			NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
	),
)
