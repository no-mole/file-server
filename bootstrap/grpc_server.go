package bootstrap

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	fsPb "github.com/no-mole/file-server/protos/file_server"
	fs "github.com/no-mole/file-server/service/file_server"
	"github.com/no-mole/neptune/app"
	"google.golang.org/grpc"
	"math"
	middleware "smart.gitlab.biomind.com.cn/infrastructure/middlewares"
)

func InitGrpcServer(_ context.Context) error {
	s := app.NewGrpcServer(
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc_middleware.WithUnaryServerChain(
			middleware.TracingServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			middleware.TracingServerStreamInterceptor(),
		),
	)
	s.RegisterService(&fsPb.Metadata().ServiceDesc, fs.New())

	return nil
}
