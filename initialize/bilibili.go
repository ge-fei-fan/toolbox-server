package initialize

import (
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"toolbox-server/global"
	"toolbox-server/model/example"
)

func InitDownloading() error {
	db := global.TOOL_DB.Model(&example.VideoInstance{})
	db = db.Where("`status` = ?", -1)
	err := db.Find(&example.VideoDownloading).Error
	if err != nil {
		return err
	}
	return nil
}

// 把所有下载中的视频都设置为下载错误
func DealDownloading() {
	for _, v := range example.VideoDownloading {
		v.Status = -1
		v.Save()
	}
}

func CheckFFmpeg() error {
	var ffmpegPath string
	currentPath, _ := os.Getwd()
	switch gin.Mode() {
	case gin.DebugMode:
		ffmpegPath = filepath.Join(currentPath, "ffmpeg")
	case gin.ReleaseMode:
		ffmpegPath = filepath.Join(currentPath, "resources", "ffmpeg")
	}
	global.TOOL_LOG.Info(ffmpegPath)
	_, err := getFFmpegVersion(ffmpegPath)
	if err != nil {
		_, err = getFFmpegVersion("ffmpeg")
		if err != nil {
			global.TOOL_LOG.Info("没有ffmpeg环境")
			global.TOOL_FFMPEG = ""
			return errors.New("没有ffmpeg环境")
		}
		global.TOOL_LOG.Info("使用全局ffmpeg")
		global.TOOL_FFMPEG = "ffmpeg"
		return nil
	}
	global.TOOL_LOG.Info("使用程序自带ffmpeg")
	global.TOOL_FFMPEG = ffmpegPath
	return nil

}

func getFFmpegVersion(path string) (string, error) {
	cmd := exec.Command(path, "-version")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	res, err := cmd.Output()
	if err != nil {
		return string(res), err
	}
	return string(res), nil
}
