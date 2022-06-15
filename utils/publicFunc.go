package utils

import "strings"

// GetFileType 根据文件后缀名获取文件类型
func GetFileType(subffix string) uint8 {
	subffix = strings.ToLower(subffix)
	switch subffix {
	case "txt":
	case "log":
		return 1 // 文本
	case "zip":
	case "rar":
	case "7z":
	case "tar":
	case "gz":
	case "tar.gz":
		return 2 // 压缩包
	case "png":
	case "jpg":
	case "jpeg":
	case "gif":
		return 3 // 图片
	case "avi":
	case "mp4":
	case "mkv":
	case "flv":
		return 4 // 视频
	case "mp3":
	case "wma":
	case "wav":
		return 5 // 音频
	case "exe":
	case "msi":
	case "ico":
		return 6 // 可执行文件
	}
	return 0 // 未知文件
}
