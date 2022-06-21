package main

import (
	"context"
	"github.com/no-mole/file-server/bootstrap"

	_ "go.uber.org/automaxprocs"

	biogo "github.com/no-mole/neptune/app"
	"github.com/no-mole/neptune/config"
)

func main() {
	ctx := context.Background()

	biogo.NewApp(ctx)

	biogo.AddHook(
		config.Init, //初始化配置
		bootstrap.InitRedis,
		bootstrap.InitLogger, //初始化日志 bootstrap.InitFileServer,
		bootstrap.InitFileServer,
		bootstrap.InitGrpcServer, //初始化grpc server
		bootstrap.PProf,
	)

	if err := biogo.Start(); err != nil {
		panic(err)
	}

	err := <-biogo.ErrorCh()
	biogo.Stop()
	panic(err)
}
