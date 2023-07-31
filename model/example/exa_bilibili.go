package example

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"toolbox-server/global"
	"toolbox-server/utils"
)

type BiliResponse struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

// bilibili账号信息
type Account struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Mid    int64  `json:"mid"`
		Uname  string `json:"uname"`
		UserId string `json:"userid"`
		Rank   string `json:"rank"`
	} `json:"data"`
}
type CookieInfo struct {
	BiliResponse
	Data struct {
		Refresh   bool  `json:"refresh"`
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
}

type ReFreshCookieInfo struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Status       int    `json:"status"`
		Message      string `json:"message"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

//-----------------------自动采集相关的-----------------------

// 用户基本信息
type UserDetail struct {
	Mid     int    `json:"mid"`
	Name    string `json:"name"`
	Sex     string `json:"sex"` //性别
	Face    string `json:"face"`
	Sign    string `json:"sign"`    //签名
	Silence int8   `json:"silence"` //封禁状态 0：正常 1：被封
}

// 用户空间详细信息响应
type UserSpaceResult struct {
	Code    int64      `json:"code"`
	Message string     `json:"message"`
	Data    UserDetail `json:"data"`
}

// 采集账号数据库表
type BilibiliCollect struct {
	global.TOOL_MODEL
	UserDetail
	IsCollect bool `json:"is_collect" gorm:"comment:是否采集"`
	Videos    []CollectVideo
	//SpaceVideo
}

// 采集视频表
type CollectVideo struct {
	global.TOOL_MODEL
	BilibiliCollect   BilibiliCollect //反向查询可以使用
	BilibiliCollectID uint
	SpaceVideo
}

// 投稿视频信息
type SpaceVideo struct {
	Aid          int    `json:"aid"`              // 稿件avid
	Author       string `json:"author"`           // 视频UP主，不一定为目标用户（合作视频）
	Bvid         string `json:"bvid"`             // 稿件bvid
	Comment      int    `json:"comment" `         // 视频评论数
	Copyright    string `json:"copyright" `       // 空，作用尚不明确
	Created      int64  `json:"created"`          // 投稿时间戳
	Description  string `json:"description"  `    // 视频简介
	HideClick    bool   `json:"hide_click"  `     // 固定值false，作用尚不明确
	IsPay        int    `json:"is_pay"  `         // 固定值0，作用尚不明确
	IsUnionVideo int    `json:"is_union_video"  ` // 是否为合作视频，0：否，1：是
	Length       string `json:"length"  `         // 视频长度，MM:SS
	Mid          int    `json:"mid"  `            // 视频UP主mid，不一定为目标用户（合作视频）
	Pic          string `json:"pic" `             // 视频封面
	Play         int    `json:"play" `            // 视频播放次数
	Review       int    `json:"review" `          // 固定值0，作用尚不明确
	Subtitle     string `json:"subtitle" `        // 固定值空，作用尚不明确
	Title        string `json:"title" `           // 视频标题
	Typeid       int    `json:"typeid"`           // 视频分区tid
	VideoReview  int    `json:"video_review" `    // 视频弹幕数
}

// 查询用户投稿视频响应
type GetUserVideosResult struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List struct { // 列表信息
			Tlist map[int]struct { // 投稿视频分区索引
				Count int    `json:"count"` // 投稿至该分区的视频数
				Name  string `json:"name"`  // 该分区名称
				Tid   int    `json:"tid"`   // 该分区tid
			} `json:"tlist"`
			Vlist []SpaceVideo `json:"vlist"` // 投稿视频列表
		} `json:"list"`
		Page struct { // 页面信息
			Count int `json:"count"` // 总计稿件数
			Pn    int `json:"pn"`    // 当前页码
			Ps    int `json:"ps"`    // 每页项数
		} `json:"page"`
		EpisodicButton struct { // “播放全部“按钮
			Text string `json:"text"` // 按钮文字
			Uri  string `json:"uri"`  // 全部播放页url
		} `json:"episodic_button"`
	} `json:"data"`
}

//-----------------------获取视频相关的-----------------------

type Page struct {
	Cid int64 `json:"cid"`
	//分批序号
	Page int16 `json:"page"`
	//分P标题
	Part string `json:"part"`
}

// VideoInfo 视频概要信息
type VideoInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Bvid   string `json:"bvid"`
		Aid    int64  `json:"aid"`
		Videos int64  `json:"videos"`
		Title  string `json:"title"`
		//子分区名称
		Tname string `json:"tname"`
		Owner struct {
			Mid  int64  `json:"mid"`
			Name string `json:"name"`
		} `json:"owner"`
		Pages []Page `json:"pages"`
	} `json:"data"`
}

// PlayInfo 视频文件详情信息
type PlayInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AcceptDescription []string `json:"accept_description"`
		AcceptQuality     []int16  `json:"accept_quality"`
		Dash              struct {
			Video []struct {
				Id        int16    `json:"id"`
				BaseUrl   string   `json:"baseUrl"`
				BackupUrl []string `json:"backupUrl"`
			} `json:"video"`
			Audio []struct {
				Id        int16    `json:"id"`
				BaseUrl   string   `json:"baseUrl"`
				BackupUrl []string `json:"backupUrl"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}

// 单个视频结数据库表

type VideoInstance struct {
	global.TOOL_MODEL
	Bvid     string `json:"bvid" form:"bvid"`
	Aid      string `json:"aid" form:"aid"`
	Cid      int64  `json:"cid" form:"cid"`
	Title    string `json:"title" form:"title" gorm:"comment:标题"`
	SavePath string `json:"savePath" gorm:"comment:保存路径"`
	Result   string `json:"result" gorm:"comment:下载结果"`
	Status   int8   `json:"status" form:"status" gorm:"comment:当前状态 0已完成 1下载视频 2下载音频 3合并音视频 -1下载失败"`
	Quality  string `json:"quality" form:"quality" gorm:"comment:视频质量"`
	//VideoTotalBytes     int64  `json:"videoTotalBytes"`
	//AudioTotalBytes     int64  `json:"audioTotalBytes"`
	//VideoCompletedBytes int64  `json:"videoCompletedBytes"`
	//AduioCompletedBytes int64  `json:"audioCompletedBytes"`
	CompletedBytes int64  `json:"completedBytes"`
	TotalBytes     int64  `json:"totalBytes"`
	Owner          string `json:"owner"`
	Progress       int    `json:"progress"`
}

const (
	//获取音视频下载链接
	playUrl = "https://api.bilibili.com/x/player/playurl"
)

var VideoDownloading = make([]*VideoInstance, 0, 50)
var VideoMutex sync.Mutex

func AddVideoDownloading(instance *VideoInstance) bool {
	VideoMutex.Lock()
	defer VideoMutex.Unlock()
	for _, item := range VideoDownloading {
		//切片里已经存在了
		if item.Cid == instance.Cid {
			return false
		}
	}
	VideoDownloading = append(VideoDownloading, instance)
	return true
}

// 根据bvid下载视频
// isRetry 是否重试请求
func (v *VideoInstance) Download(isRetry bool) {
	if isRetry {
		//获取视频的信息
		if err := global.TOOL_DB.Where("cid = ?", v.Cid).First(&VideoInstance{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 视频不存在
				global.TOOL_LOG.Warn("视频数据库信息不存在", zap.Error(err))
				return
				//err = global.TOOL_DB.Create(v).Error
				//if err != nil {
				//	global.TOOL_LOG.Warn("写入数据库失败", zap.Error(err))
				//}
			} else {
				// 查询出错
				global.TOOL_LOG.Warn("查询数据库出错", zap.Error(err))
				return
			}
		} else {
			v.Status = 1
			v.Save()
		}
	} else {
		err := global.TOOL_DB.Create(v).Error
		if err != nil {
			global.TOOL_LOG.Warn("写入数据库失败", zap.Error(err))
		}
	}
	client := resty.New()
	p := &PlayInfo{}
	cookie := fmt.Sprintf("buvid3=012A0511-C43C-4195-A0AA-697EECE05C21148831infoc;SESSDATA=%s", global.TOOL_CONFIG.Bilibili.SessData)
	_, err := client.R().SetQueryParams(map[string]string{
		"bvid":  v.Bvid,
		"cid":   strconv.Itoa(int(v.Cid)),
		"qn":    "0",
		"fnver": "0",
		"fnval": "208",
		"fourk": "1",
	}).SetHeaders(map[string]string{
		"cookie": cookie,
		//"User-Agent": "PostmanRuntime/7.32.2",
	}).SetResult(p).Get(playUrl)
	if err != nil {
		v.Result = err.Error()
		v.Save()
		return
	}
	if p.Code == 87007 {
		v.Result = "可能为付费视频"
		v.Status = -1
		v.Save()
		return
	}
	if p.Message != "0" {
		v.Result = p.Message
		v.Status = -1
		v.Save()
		return
	}

	//下载视频
	v.Quality = v.getVideoQuality(p.Data.Dash.Video[0].Id)
	v.Save()
	err = v.saveToFile(p.Data.Dash.Video[0].BaseUrl, ".video", v)
	if err != nil {
		v.Result = err.Error()
		v.Status = -1
		v.Save()
		return
	}
	//下载音频
	v.Status = 2
	v.Save()
	err = v.saveToFile(p.Data.Dash.Audio[0].BaseUrl, ".audio", v)
	if err != nil {
		v.Result = err.Error()
		v.Status = -1
		v.Save()
		return
	}
	//合并音视频
	v.Status = 3
	v.Save()
	err = v.mergeVideo()
	if err != nil {
		v.Result = err.Error()
		v.Status = -1
		v.Save()
		return
	}
	VideoMutex.Lock()
	for i, d := range VideoDownloading {
		if v == d {
			VideoDownloading = append(VideoDownloading[:i], VideoDownloading[i+1:]...)
			break
		}
	}
	VideoMutex.Unlock()

}
func (v *VideoInstance) saveToFile(url, suffix string, l utils.ProgressListener) error {
	err := os.MkdirAll(v.SavePath, 0666)
	if err != nil {
		return err
	}
	filePath := filepath.Join(v.SavePath, v.Title+suffix)
	client := resty.New().SetRetryCount(3)
	resp, err := client.R().SetDoNotParseResponse(true).SetHeader("referer", "https://www.bilibili.com").
		SetOutput(filePath).Get(url)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	contentLen, _ := strconv.ParseInt(resp.Header().Get("Content-Length"), 10, 64)
	body := utils.TeeReader(resp.RawBody(), contentLen, v)
	defer body.Close()
	_, err = io.Copy(f, body)
	if err != nil {
		return err
	}
	return nil
}
func (v *VideoInstance) mergeVideo() error {
	videoPath, _ := filepath.Abs(filepath.Join(v.SavePath, v.Title+".video"))
	audioPath, _ := filepath.Abs(filepath.Join(v.SavePath, v.Title+".audio"))
	outPath, _ := filepath.Abs(filepath.Join(v.SavePath, v.Title+".mp4"))
	cmdArguments := []string{"-i", videoPath, "-i", audioPath,
		"-c:v", "copy", "-c:a", "copy", "-f", "mp4", outPath}
	if global.TOOL_FFMPEG == "" {
		return errors.New("没有ffmpeg环境")
	}
	cmd := exec.Command(global.TOOL_FFMPEG, cmdArguments...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	v.Status = 0
	v.Save()

	//删除临时文件
	err = os.Remove(videoPath)
	if err != nil {
		return err
	}
	err = os.Remove(audioPath)
	if err != nil {
		return err
	}
	return nil
}
func (v *VideoInstance) getVideoQuality(q int16) string {
	switch q {
	case 6:
		return "240P 极速"
	case 16:
		return "240P 极速"
	case 32:
		return "480P 清晰"
	case 64:
		return "720P 高清"
	case 72:
		return "720P60 高帧率"
	case 80:
		return "1080P 高清"
	case 112:
		return "1080P+ 高码率"
	case 116:
		return "1080P60 高帧率"
	case 120:
		return "4K 超清"
	case 125:
		return "HDR 真彩色"
	case 126:
		return "杜比视界"
	case 127:
		return "8K 超高清"
	default:
		return "未知"
	}
}
func (v *VideoInstance) ProgressChanged(event *utils.ProgressEvent) {
	switch event.EventType {
	case utils.TransferStartedEvent:
		v.TotalBytes = event.TotalBytes
	case utils.TransferDataEvent:
		v.CompletedBytes = event.ConsumedBytes
		v.TotalBytes = event.TotalBytes
		v.Progress = int((float64(event.ConsumedBytes) / float64(event.TotalBytes)) * 100)
	case utils.TransferCompletedEvent:
		v.CompletedBytes = event.ConsumedBytes
		v.Progress = 100
	case utils.TransferFailedEvent:
		//v.Result = "TransferFail"
	default:
	}
}

// 更新数据库
func (v *VideoInstance) Save() {
	err := global.TOOL_DB.Save(v).Error
	if err != nil {
		global.TOOL_LOG.Warn("更新数据库失败", zap.Error(err))
	}
}
