// Package rainbowconfig rainbowconfig
package rainbowconfig

import (
	"context"
	"os"
	"sync/atomic"

	_ "git.code.oa.com/trpc-go/trpc-config-rainbow" // 使用rainbow配置时引入注册
	"git.code.oa.com/trpc-go/trpc-go/config"
	rainbowLogger "git.code.oa.com/trpc-go/trpc-go/log"
	"gopkg.in/yaml.v2"
)

// DBInfo 配置
type DBInfo struct {
	DSN            string `yaml:"dsn"`
	MaxIdleConnNum int    `yaml:"maxIdleConnNum"`
	MaxConnNum     int    `yaml:"maxConnNum"`
}

// TrpcInfo 配置
type TrpcInfo struct {
	Namespace string `yaml:"namespace"`
	Env       string `yaml:"env"`
	Target    string `yaml:"target"`
	Timeout   int32  `yaml:"timeout"`
}

// TafInfo 配置
type TafInfo struct {
	Namespace string `yaml:"namespace"`
	Env       string `yaml:"env"`
	Target    string `yaml:"target"`
	Timeout   int32  `yaml:"timeout"`
	TafObj    string `yaml:"tafObj"`
}

// RedisInfo redis配置
type RedisInfo struct {
	Address string `yaml:"address"`
	Passwd  string `yaml:"passwd"`
	DB      int    `yaml:"db"`
}

// 并发安全的Value
var atomicValue atomic.Value

// InitCfg 初始化七彩石配置
func InitCfg(filename string, appConf interface{}) (err error) {
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
	atomicValue.Store(appConf)
	// rainbowLogger.Infof("conf %v", appConf)
	if !runInLocal {
		configChan, watchErr := config.Get("rainbow").Watch(context.TODO(), "app.yaml")
		if watchErr == nil {
			// 监听后续rainbowConfig的变化
			go func() {
				for r := range configChan {
					if watchErr = yaml.Unmarshal([]byte(r.Value()), appConf); watchErr == nil {
						atomicValue.Store(appConf)
					}
					rainbowLogger.Infof("event: %d, reload config ", r.Event())
				}
			}()
		}
	}
	rainbowLogger.Infof("rainbow init end")

	return err
}
