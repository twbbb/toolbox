// Package osutil osutil方法
package osutil

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// IsLocalEnv 是否是本地环境
func IsLocalEnv() bool {
	return os.Getenv("DOCKER_ENV") == ""
}

// GetPodIP 获取podip
func GetPodIP() string {
	return os.Getenv("POD_IP")
}

// listenForSIGTERM 启动一个 goroutine 来监听 SIGTERM 信号，并在收到信号时执行提供的回调函数。
func ListenForSIGTERM(ctx context.Context, cancel context.CancelFunc, callback func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		select {
		case <-sigs:
			callback() // 调用传入的回调函数
			cancel()   // 在接收到 SIGTERM 信号时取消上下文
			return
		case <-ctx.Done():
			return // 如果主上下文被取消，也停止此 goroutine
		}
	}()
}
