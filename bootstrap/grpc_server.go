package bootstrap

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"math"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/app"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/middleware"
	fs "smart.gitlab.biomind.com.cn/intelligent-system/file-server/service/file_server"
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