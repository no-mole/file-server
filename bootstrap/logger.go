package bootstrap

import (
	"context"
	"os"
	"path"

	"smart.gitlab.biomind.com.cn/infrastructure/biogo/config"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/logger"
	_ "smart.gitlab.biomind.com.cn/infrastructure/logger_center/service/logger_center"
)

func InitLogger(ctx context.Context) error {
	body, err := os.ReadFile(path.Join(config.GlobalConfig.ConfigPath, "logger.yml"))
	if err != nil {
		return err
	}
	return logger.Bootstrap(ctx, body)
}
