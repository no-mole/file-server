package bootstrap

import (
	"context"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config/center"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/redis"
	"smart.gitlab.biomind.com.cn/intelligent-system/file-server/model"
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
