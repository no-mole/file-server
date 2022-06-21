package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/no-mole/file-server/model"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/no-mole/neptune/config"
	"github.com/no-mole/neptune/logger"
	"github.com/no-mole/neptune/utils"
)

const (
	_ttl = 5
)

type ServerNode struct {
	NodeName string `json:"node_name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DirSize  int64  `json:"dir_size"`
}

func DirSizeB(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func InitFileServer(ctx context.Context) error {
	size, err := DirSizeB(path.Join(utils.GetCurrentAbPath(), model.RootDir))
	if err != nil {
		return err
	}
	node := &ServerNode{
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

func ReLoadStoreRate(ctx context.Context) error {
	size, err := DirSizeB(path.Join(utils.GetCurrentAbPath(), model.RootDir))
	if err != nil {
		return err
	}
	node := &ServerNode{
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
