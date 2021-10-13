package bootstrap

import (
	"context"
	"os"
	"path"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/logger"
)

func InitLogger(ctx context.Context) error {
	body, err := os.ReadFile(path.Join(config.GlobalConfig.ConfigPath, "logger.yml"))
	if err != nil {
		return err
	}
	return logger.Bootstrap(ctx, body)
}
