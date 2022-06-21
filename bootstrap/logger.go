package bootstrap

import (
	"context"
	"os"
	"path"

	_ "github.com/no-mole/log-collector/service/logger_center"
	"github.com/no-mole/neptune/config"
	"github.com/no-mole/neptune/logger"
)

func InitLogger(ctx context.Context) error {
	body, err := os.ReadFile(path.Join(config.GlobalConfig.ConfigPath, "logger.yml"))
	if err != nil {
		return err
	}
	return logger.Bootstrap(ctx, body)
}
