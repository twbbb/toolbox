// Package osutil osutil方法
package osutil

import (
	"bytes"
	"context"
	"fmt"
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

// ListenForSIGTERM 启动一个 goroutine 来监听 SIGTERM 信号，并在收到信号时执行提供的回调函数。
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

// GuessFileType 根据文件头猜测文件类型
func GuessFileType(buffer []byte) (string, error) {
	if len(buffer) < 8 {
		return "", fmt.Errorf("buffer too short")
	}

	// pdf
	if bytes.HasPrefix(buffer, []byte{0x25, 0x50, 0x44, 0x46}) {
		return "pdf", nil
	}

	// JPEG
	if bytes.HasPrefix(buffer, []byte{0xFF, 0xD8, 0xFF}) {
		return "jpg", nil
	}

	// PNG
	if bytes.HasPrefix(buffer, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return "png", nil
	}

	// GIF
	if bytes.HasPrefix(buffer, []byte{0x47, 0x49, 0x46, 0x38}) {
		return "gif", nil
	}

	// BMP
	if bytes.HasPrefix(buffer, []byte{0x42, 0x4D}) {
		return "bmp", nil
	}

	// WEBP
	if bytes.HasPrefix(buffer, []byte{0x52, 0x49, 0x46, 0x46}) && bytes.HasSuffix(buffer, []byte{0x57, 0x45, 0x42, 0x50}) {
		return "webp", nil
	}

	return "", fmt.Errorf("unknown image type")
}
