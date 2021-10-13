package bootstrap

import (
	"context"
	"encoding/json"

	"file-server/model"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/app"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config/center"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/database"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/env"
)

var dbNames = []string{
	model.MysqlEngineBar,
}

func InitDatabase(ctx context.Context) error {
	configCenterClient := config.GetClient()
	for _, dbName := range dbNames {
		conf, err := configCenterClient.Get(ctx, dbName)
		if err != nil {
			return err
		}
		err = initDatabaseDrive(dbName, conf.GetValue())
		if err != nil {
			return err
		}
		// 监听修改
		configCenterClient.Watch(ctx, conf, func(item *center.Item) {
			err := initDatabaseDrive(item.Key, item.GetValue())
			if err != nil {
				app.Error(err)
				return
			}
		})
	}
	return nil
}

func initDatabaseDrive(dbName, confStr string) error {
	conf := &database.Config{
		Driver:       "mysql",
		Host:         "localhost",
		Port:         3306,
		WriteTimeout: 1000,
		ReadTimeout:  2000,
	}
	err := json.Unmarshal([]byte(confStr), conf)
	if err != nil {
		return err
	}
	return database.Init(dbName, conf, env.GetEnvDebug())
}
