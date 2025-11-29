package main

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/presentation"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Supply(
			fx.Annotate("./", fx.ResultTags(`name:"configPath"`)),
			fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
		),
		presentation.Module,
	)
	app.Run()
}
