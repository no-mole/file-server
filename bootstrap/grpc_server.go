package bootstrap

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"math"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/app"
	fs "smart.gitlab.biomind.com.cn/infrastructure/file-server/service/file_server"
	middleware "smart.gitlab.biomind.com.cn/infrastructure/middlewares"
	fsPb "smart.gitlab.biomind.com.cn/intelligent-system/protos/file_server"
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
