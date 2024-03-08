// Package osutil osutil方法
package osutil

import (
	"os"
)

// IsLocalEnv 是否是本地环境
func IsLocalEnv() bool {
	return os.Getenv("DOCKER_ENV") == ""
}
