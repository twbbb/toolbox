// Package initdata initdata初始化
package initdata

import (
	"time"

	"github.com/twbbb/toolbox/trpcutil/entity"
	"trpc.group/trpc-go/trpc-go/client"
)

// GetTrpcOpts 获取trpc的opts
func GetTrpcOpts(trpConfig entity.TrpcInfo) []client.Option {
	options := []client.Option{
		client.WithTarget(trpConfig.Target),
		client.WithTimeout(time.Duration(trpConfig.Timeout) * time.Millisecond),
		client.WithNamespace(trpConfig.Namespace),
		client.WithEnvKey(trpConfig.Env),
		client.WithCalleeMetadata("env", trpConfig.Env),
	}
	return options
}

// GetTafOpts 获取taf的opts
func GetTafOpts(tafConfig entity.TafInfo) (string, []client.Option) {
	options := []client.Option{
		client.WithTarget(tafConfig.Target),
		client.WithTimeout(time.Duration(tafConfig.Timeout) * time.Millisecond),
		client.WithNamespace(tafConfig.Namespace),
		client.WithEnvKey(tafConfig.Env),
		client.WithCalleeMetadata("env", tafConfig.Env),
	}
	return tafConfig.TafObj, options
}
