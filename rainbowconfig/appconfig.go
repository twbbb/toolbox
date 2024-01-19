// Package rainbowconfig rainbowconfig
package rainbowconfig

import (
	"context"
	"os"

	_ "git.code.oa.com/trpc-go/trpc-config-rainbow" // 使用rainbow配置时引入注册
	"git.code.oa.com/trpc-go/trpc-go/config"
	rainbowLogger "git.code.oa.com/trpc-go/trpc-go/log"
	"gopkg.in/yaml.v2"
)

// InitCfg 初始化七彩石配置
func InitCfg(
	filename string, appConf interface{}, callback func(appConf interface{}, err error),
) (err error) {
	rainbowLogger.Infof("rainbow init start")
	// 使用trpc-go/config中Watch接口监听七彩石远程配置变化
	var c config.Config
	runInLocal := os.Getenv("DOCKER_ENV") == ""
	if runInLocal {
		c, err = config.Load(filename, config.WithCodec("yaml"), config.WithProvider("file"))
	} else {
		c, err = config.Load(filename, config.WithCodec("yaml"), config.WithProvider("rainbow"))
	}

	if err != nil {
		rainbowLogger.Errorf("Load rainbow config: %v ", err)
		return err
	}

	// 读取配置, appname为配置中的：TConf
	if err = c.Unmarshal(appConf); err != nil {
		rainbowLogger.Errorf("Load rainbow conifg fail:%v", err)
		return err
	}
	if runInLocal {
		rainbowLogger.Infof("rainbow local init end")

		return nil
	}

	configChan, watchErr := config.Get("rainbow").Watch(context.TODO(), filename)
	if watchErr == nil {
		// 监听后续rainbowConfig的变化
		go func() {
			for r := range configChan {
				watchErr := yaml.Unmarshal([]byte(r.Value()), appConf)
				if watchErr != nil {
					rainbowLogger.Errorf("watchErr event: %d, reload fail,%+v ", r.Event(), watchErr)
				} else {
					rainbowLogger.Infof("event: %d, reload config ", r.Event())
				}
				if callback != nil {
					callback(appConf, watchErr)
				}
			}
		}()
	}
	rainbowLogger.Infof("rainbow init end")

	return err
}
