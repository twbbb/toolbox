package env

import (
	"os"

	"git.code.oa.com/trpc-go/trpc-go"
)

// IsLocalEnv 是否是本地环境
func IsLocalEnv() bool {
	return os.Getenv("DOCKER_ENV") == ""
}

// IsDevEnv 是否是本地环境
func IsDevEnv() bool {
	if trpc.GlobalConfig() == nil {
		return false
	}
	return trpc.GlobalConfig().Global.Namespace == "Development"
}

// IsTestEnv 是否是本地环境
func IsTestEnv() bool {
	if trpc.GlobalConfig() == nil {
		return false
	}
	return trpc.GlobalConfig().Global.Namespace == "Test"
}

// IsProEnv 是否是本地环境
func IsProEnv() bool {
	if trpc.GlobalConfig() == nil {
		return false
	}
	return trpc.GlobalConfig().Global.Namespace == "Production"
}
