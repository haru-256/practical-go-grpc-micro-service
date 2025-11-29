package main

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/presentation"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Supply(
			// NOTE: バイナリを実行する位置からの相対パスで指定する
			fx.Annotate("./", fx.ResultTags(`name:"configPath"`)),
			fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
		),
		presentation.Module,
	)
	app.Run()
}
