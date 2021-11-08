package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"smart.gitlab.biomind.com.cn/intelligent-system/file-server/model"
	"time"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
	fs "smart.gitlab.biomind.com.cn/intelligent-system/biogo/file_server"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/logger"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"
)

const (
	_ttl = 5
)

func InitFileServer(ctx context.Context) error {
	size, err := fs.DirSizeB(path.Join(utils.GetCurrentAbPath(), model.RootDir))
	if err != nil {
		return err
	}
	node := &fs.ServerNode{
		NodeName: config.GlobalConfig.Name,
		Host:     config.GlobalConfig.IP,
		Port:     config.GlobalConfig.GrpcPort,
		DirSize:  size,
	}

	values, err := json.Marshal(node)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s",
		model.FileServerNodePrefix,
		config.GlobalConfig.Name,
	)
	Register(ctx, key, string(values))
	go loopStoreRate()
	return nil
}

func Register(ctx context.Context, key, value string) {
	err := config.GetClient().SetExKeepAlive(ctx, key, value, _ttl)
	if err != nil {
		logger.Error(ctx, "register", err, logger.WithField("msg", fmt.Sprintf("保持连接失败：%s", err.Error())))
	}
}

func loopStoreRate() {
	ctx := context.Background()
	C := time.Tick(20 * time.Second)
	for {
		select {
		case <-C:
			err := ReLoadStoreRate(ctx)
			if err != nil {
				logger.Error(ctx, "loopStoreRate.ReLoadStoreRate", err)
			}
		}
	}
}

func  ReLoadStoreRate(ctx context.Context) error {
	size, err := fs.DirSizeB(path.Join(utils.GetCurrentAbPath(), model.RootDir))
	if err != nil {
		return err
	}
	node := &fs.ServerNode{
		NodeName: config.GlobalConfig.Name,
		Host:     config.GlobalConfig.IP,
		Port:     config.GlobalConfig.GrpcPort,
		DirSize:  size,
	}

	values, err := json.Marshal(node)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s",
		model.FileServerNodePrefix,
		config.GlobalConfig.Name,
	)

	return config.GetClient().Set(ctx, key, string(values))
}