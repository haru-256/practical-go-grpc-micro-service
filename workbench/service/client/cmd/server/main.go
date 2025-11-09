package main

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation"
	"go.uber.org/fx"
)

// @title Client Service API
// @version 1.0
// @description CQRS Client Service API
// @BasePath /
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
