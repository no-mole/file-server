package bar

import (
	"github.com/gin-gonic/gin"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/grpc_pool"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/logger"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/output"
	"smart.gitlab.biomind.com.cn/intelligent-system/enum"

	barPb "smart.gitlab.biomind.com.cn/intelligent-system/protos/bar"
)

type SayHelloParams struct {
	Say string `json:"say" form:"say" binding:"required,min=1,max=10"`
}

func SayHello(ctx *gin.Context) {
	p := &SayHelloParams{}
	err := ctx.ShouldBindQuery(p)
	if err != nil {
		output.Json(ctx, enum.IllegalParam, nil)
		return
	}

	conn, err := grpc_pool.GetConnection(barPb.Metadata())
	if err != nil {
		output.Json(ctx, enum.ErrorGrpcConnect, nil)
		logger.Error(ctx, "SayHello", err)
		return
	}
	defer conn.Close()

	cli := barPb.NewServiceClient(conn.Value())
	resp, err := cli.SayHelly(ctx, &barPb.SayHelloRequest{Say: p.Say})
	if err != nil {
		output.Json(ctx, enum.ErrorGrpcConnect, nil)
		logger.Error(ctx, "SayHello", err)
		return
	}
	output.Json(ctx, enum.Success, resp.Reply)
}
