package bootstrap

import (
	"context"
	fsPb "github.com/no-mole/file-server/protos/file_server"
	fs "github.com/no-mole/file-server/service/file_server"
	"github.com/no-mole/neptune/app"
	"google.golang.org/grpc"
	"math"
)

func InitGrpcServer(_ context.Context) error {
	s := app.NewGrpcServer(
		grpc.MaxRecvMsgSize(math.MaxInt32),
	)
	s.RegisterService(&fsPb.Metadata().ServiceDesc, fs.New())

	return nil
}
