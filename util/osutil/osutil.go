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
		defer signal.Stop(sigs)
		select {
		case <-sigs:
			callback() // 调用传入的回调函数
			cancel()   // 在接收到 SIGTERM 信号时取消上下文
		case <-ctx.Done():
			// 如果主上下文被取消，也停止此 goroutine
		}
	}()
}

// IsPdfFile 检测是否为PDF文件
func IsPdfFile(buffer []byte) (bool, error) {
	if len(buffer) < 8 {
		return false, fmt.Errorf("buffer too short")
	}

	// 检查是否是base64数据URI格式
	if bytes.HasPrefix(buffer, []byte("data:")) {
		// 查找base64数据开始位置
		base64Start := bytes.Index(buffer, []byte("base64,"))
		if base64Start != -1 {
			// 提取base64数据部分
			base64Data := buffer[base64Start+7:]
			if len(base64Data) < 8 {
				return false, fmt.Errorf("base64 data too short")
			}

			// 对base64数据部分进行文件类型检测
			return detectPdfFile(base64Data)
		}
	}
	// 如果不是base64格式，直接检测原始内容
	return detectPdfFile(buffer)
}

func detectPdfFile(buffer []byte) (bool, error) {
	if len(buffer) < 8 {
		return false, fmt.Errorf("buffer too short")
	}

	// 剥离各种BOM (Byte Order Mark)
	data := stripBOM(buffer)

	// 确保剥离BOM后数据仍然足够长
	if len(data) < 4 {
		return false, fmt.Errorf("buffer too short after BOM removal")
	}
	if isPDFFile(data) {
		return true, nil
	}
	return false, nil
}

// GuessFileType 根据文件头猜测文件类型，支持普通文件和base64数据URI格式
func GuessFileType(buffer []byte) (string, error) {
	if len(buffer) < 8 {
		return "", fmt.Errorf("buffer too short")
	}

	// 检查是否是base64数据URI格式
	if bytes.HasPrefix(buffer, []byte("data:")) {
		// 查找base64数据开始位置
		base64Start := bytes.Index(buffer, []byte("base64,"))
		if base64Start != -1 {
			// 提取base64数据部分
			base64Data := buffer[base64Start+7:]
			if len(base64Data) < 8 {
				return "", fmt.Errorf("base64 data too short")
			}

			// 对base64数据部分进行文件类型检测
			return detectFileTypeFromContent(base64Data)
		}
	}

	// 如果不是base64格式，直接检测原始内容
	return detectFileTypeFromContent(buffer)
}

// detectFileTypeFromContent 根据文件内容检测文件类型
func detectFileTypeFromContent(buffer []byte) (string, error) {
	if len(buffer) < 8 {
		return "", fmt.Errorf("buffer too short")
	}

	// 剥离各种BOM (Byte Order Mark)
	data := stripBOM(buffer)

	// 确保剥离BOM后数据仍然足够长
	if len(data) < 4 {
		return "", fmt.Errorf("buffer too short after BOM removal")
	}

	// 第一优先级：检测二进制文件格式（有明确的魔数）
	// 这些格式不会被误判为文本

	// JPEG
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "jpg", nil
	}

	// PNG
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return "png", nil
	}

	// GIF
	if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}) {
		return "gif", nil
	}

	// BMP
	if bytes.HasPrefix(data, []byte{0x42, 0x4D}) {
		return "bmp", nil
	}

	// TIFF (小端序)
	if bytes.HasPrefix(data, []byte{0x49, 0x49, 0x2A, 0x00}) {
		return "tiff", nil
	}

	// TIFF (大端序)
	if bytes.HasPrefix(data, []byte{0x4D, 0x4D, 0x00, 0x2A}) {
		return "tiff", nil
	}

	// ZIP压缩文件
	if bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}) {
		return "zip", nil
	}

	// RAR压缩文件
	if bytes.HasPrefix(data, []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}) {
		return "rar", nil
	}

	// 7Z压缩文件
	if bytes.HasPrefix(data, []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}) {
		return "7z", nil
	}

	// GZ压缩文件
	if bytes.HasPrefix(data, []byte{0x1F, 0x8B}) {
		return "gz", nil
	}

	// MP3 (ID3v2标签或MPEG帧)
	if bytes.HasPrefix(data, []byte("ID3")) ||
		(len(data) >= 2 && (data[0] == 0xFF && (data[1]&0xE0) == 0xE0)) {
		return "mp3", nil
	}

	// RIFF容器格式需要检查子类型
	if bytes.HasPrefix(data, []byte{0x52, 0x49, 0x46, 0x46}) && len(data) >= 12 {
		// WAV文件
		if bytes.Equal(data[8:12], []byte{0x57, 0x41, 0x56, 0x45}) {
			return "wav", nil
		}
		// AVI文件
		if bytes.Equal(data[8:12], []byte{0x41, 0x56, 0x49, 0x20}) {
			return "avi", nil
		}
		// WEBP
		if bytes.Equal(data[8:12], []byte{0x57, 0x45, 0x42, 0x50}) {
			return "webp", nil
		}
	}

	// MP4文件
	if (len(data) >= 12 && bytes.Equal(data[4:8], []byte("ftyp"))) ||
		bytes.HasPrefix(data, []byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}) {
		return "mp4", nil
	}

	// 第二优先级：检测有明确特征的文本格式
	// pdf (支持带BOM和不带BOM的PDF，以及前面有垃圾数据的PDF)
	// PDF标记可能不在文件开头，搜索前512字节
	if isPDFFile(data) {
		return "pdf", nil
	}

	// JSON文件
	if bytes.HasPrefix(data, []byte("{")) || bytes.HasPrefix(data, []byte("[")) {
		return "json", nil
	}

	// XML文件
	if bytes.HasPrefix(data, []byte("<?xml")) {
		return "xml", nil
	}

	// HTML文件
	if bytes.HasPrefix(data, []byte("<!DOCTYPE html")) || bytes.HasPrefix(data, []byte("<html")) {
		return "html", nil
	}

	// 第三优先级：CSV文件（需要更严格的判断）
	if isCSVFile(data) {
		return "csv", nil
	}

	// 第四优先级：普通文本文件（兜底）
	if isTextFile(data) {
		return "txt", nil
	}

	return "", fmt.Errorf("unknown file type")
}

// stripBOM 剥离各种BOM (Byte Order Mark)标记
// 支持的BOM类型：
// - UTF-8: EF BB BF
// - UTF-16 LE: FF FE
// - UTF-16 BE: FE FF
// - UTF-32 LE: FF FE 00 00
// - UTF-32 BE: 00 00 FE FF
func stripBOM(buffer []byte) []byte {
	if len(buffer) < 2 {
		return buffer
	}

	// 检查UTF-32 LE BOM (4字节): FF FE 00 00
	if len(buffer) >= 4 &&
		buffer[0] == 0xFF && buffer[1] == 0xFE &&
		buffer[2] == 0x00 && buffer[3] == 0x00 {
		return buffer[4:]
	}

	// 检查UTF-32 BE BOM (4字节): 00 00 FE FF
	if len(buffer) >= 4 &&
		buffer[0] == 0x00 && buffer[1] == 0x00 &&
		buffer[2] == 0xFE && buffer[3] == 0xFF {
		return buffer[4:]
	}

	// 检查UTF-8 BOM (3字节): EF BB BF
	if len(buffer) >= 3 &&
		buffer[0] == 0xEF && buffer[1] == 0xBB && buffer[2] == 0xBF {
		return buffer[3:]
	}

	// 检查UTF-16 LE BOM (2字节): FF FE
	if buffer[0] == 0xFF && buffer[1] == 0xFE {
		return buffer[2:]
	}

	// 检查UTF-16 BE BOM (2字节): FE FF
	if buffer[0] == 0xFE && buffer[1] == 0xFF {
		return buffer[2:]
	}

	// 没有BOM，返回原始数据
	return buffer
}

// isTextFile 检测是否为文本文件
func isTextFile(data []byte) bool {
	// 检查前100个字节是否都是可打印字符或常见控制字符
	checkLen := len(data)
	if checkLen > 100 {
		checkLen = 100
	}

	for i := 0; i < checkLen; i++ {
		b := data[i]
		// 允许的字符：可打印ASCII字符(32-126)、制表符(9)、换行符(10)、回车符(13)
		if b < 32 && b != 9 && b != 10 && b != 13 {
			return false
		}
		if b > 126 {
			return false
		}
	}
	return true
}

// isPDFFile 检测是否为PDF文件
// PDF文件以 %PDF- 开头，但可能前面有一些垃圾数据
// 根据PDF规范，PDF标记应该在文件的前1024字节内
func isPDFFile(data []byte) bool {
	pdfSignature := []byte{0x25, 0x50, 0x44, 0x46} // %PDF

	// 搜索范围：前512字节或整个数据（取较小值）
	searchLen := len(data)
	if searchLen > 512 {
		searchLen = 512
	}

	// 在搜索范围内查找PDF标记
	idx := bytes.Index(data[:searchLen], pdfSignature)
	return idx != -1
}

// isCSVFile 检测是否为CSV文件
// CSV文件的特征：
// 1. 包含分隔符（逗号、分号或制表符）
// 2. 有多行数据
// 3. 第一行通常是表头
// 4. 每行的字段数量相对一致
func isCSVFile(data []byte) bool {
	// 检查前500个字节或全部数据
	checkLen := len(data)
	if checkLen > 500 {
		checkLen = 500
	}

	// 必须是文本文件
	if !isTextFile(data) {
		return false
	}

	// 统计换行符和分隔符
	lineCount := 0
	separatorCount := 0

	for i := 0; i < checkLen; i++ {
		b := data[i]

		// 统计换行符
		if b == '\n' {
			lineCount++
		}

		// 统计分隔符
		if b == ',' || b == ';' || b == '\t' {
			separatorCount++
		}
	}

	// CSV文件必须有至少2行数据（表头+数据行）
	// 并且分隔符数量要足够多（至少每行平均有1个分隔符）
	return lineCount >= 2 && separatorCount >= lineCount
}
