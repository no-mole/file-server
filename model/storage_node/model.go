package storage_node

import (
	"context"
	"file-server/model"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/database"
)

type Model struct {
	ctx context.Context
	db  *database.BaseDb
}

func New(ctx context.Context) *Model {
	db := &database.BaseDb{}
	db.SetEngine(ctx, model.MysqlEngineBar)
	return &Model{
		ctx: ctx,
		db:  db,
	}
}
