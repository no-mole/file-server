package bootstrap

import (
	"context"
	"github.com/no-mole/file-server/model"
	"github.com/no-mole/neptune/config"
	"github.com/no-mole/neptune/config/center"
	"github.com/no-mole/neptune/redis"
)

var redisNames = []string{
	model.RedisEngine,
}

func InitRedis(ctx context.Context) error {

	configCenterClient := config.GetClient()
	for _, redisName := range redisNames {
		conf, err := configCenterClient.Get(ctx, redisName)
		if err != nil {
			return err
		}
		redis.Init(redisName, conf.GetValue())
		// 监听修改
		configCenterClient.Watch(ctx, conf, func(item *center.Item) {
			redis.Init(item.Key, item.GetValue())
		})
	}

	return nil
}
