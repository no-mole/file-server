package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"smart.gitlab.biomind.com.cn/infrastructure/biogo/config"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/env"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/utils"
)

func PProf(_ context.Context) error {
	if config.GlobalConfig.Env.Mode == env.ModeDev || config.GlobalConfig.Env.Mode == env.ModeTest {
		runtime.SetBlockProfileRate(1) //对阻塞超过1纳秒的goroutine进行数据采集。
		go func() {
			port, _ := utils.GetAvailablePort()
			fmt.Println("pprof_port: ", port)
			http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		}()
	}
	return nil
}
