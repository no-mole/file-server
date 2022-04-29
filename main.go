package main

import (
	"context"
	"smart.gitlab.biomind.com.cn/infrastructure/file-server/bootstrap"

	_ "go.uber.org/automaxprocs"

	biogo "smart.gitlab.biomind.com.cn/infrastructure/biogo/app"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/config"
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
