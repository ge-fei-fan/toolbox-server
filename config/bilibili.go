package config

import (
	"os/user"
	"path/filepath"
)

type Bilibili struct {
	AutoDownload bool   `mapstructure:"auto-download" json:"auto-download" yaml:"auto-download"` //监控粘贴板下载
	Collect      bool   `mapstructure:"collect" json:"collect" yaml:"collect"`                   // 采集功能
	SessData     string `mapstructure:"sess-data" json:"sess-data" yaml:"sess-data"`
	RefreshToken string `mapstructure:"refresh-token" json:"refresh-token" yaml:"refresh-token"`
	DownloadPath string `mapstructure:"download-path" json:"download-path" yaml:"download-path"`
	BiliJct      string `mapstructure:"bili-jct" json:"bili-jct" yaml:"bili-jct"`
	//FfmpegPath   string `mapstructure:"ffmpeg-path" json:"ffmpeg-path" yaml:"ffmpeg-path"`
}

func initDefaultBilibili() Bilibili {
	usr, _ := user.Current()
	defalutPath := filepath.Join(usr.HomeDir, "Downloads")
	//currentPath, _ := os.Getwd()
	//defalutPath := filepath.Join(currentPath, "bilibili")
	//var ffmpegPath string
	//switch gin.Mode() {
	//case gin.DebugMode:
	//	ffmpegPath = currentPath
	//case gin.ReleaseMode:
	//	ffmpegPath = filepath.Join(currentPath, "resources")
	//}

	return Bilibili{
		AutoDownload: false,
		Collect:      false,
		SessData:     "",
		RefreshToken: "",
		DownloadPath: defalutPath,
		BiliJct:      "",
		//FfmpegPath:   ffmpegPath,
	}
}
