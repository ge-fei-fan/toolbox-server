package example

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"toolbox-server/global"
	"toolbox-server/model/example"
	"toolbox-server/utils"
)

const (
	//获取视频cid
	viewUrl = "https://api.bilibili.com/x/web-interface/view"
)

// 根据分享链接获取bvid链接，下载视频
//func ParseUrl(url string) error {
//	u := utils.FilterString("https://.*", url)
//	if u == "" {
//		return errors.New("过滤分享链接为空")
//	}
//	client := resty.New()
//	resp, err := client.R().Get(u)
//	if err != nil {
//		return err
//	}
//	err = DownloadVideo(resp.Request.URL)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func GetBvid(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		global.TOOL_LOG.Error("解析bvid失败:", zap.Error(err))
		return "", errors.New("解析bvid失败")
	}
	bvid := path.Base(u.Path)
	if bvid == "" {
		return "", errors.New("获取bvid为空")
	}
	return bvid, nil
}
func DownloadVideo(s string) error {
	bvid, err := GetBvid(s)
	if err != nil {
		return err
	}
	v := &example.VideoInfo{}
	client := resty.New()
	_, err = client.R().SetQueryParams(map[string]string{"bvid": bvid}).
		SetResult(v).Get(viewUrl)
	if err != nil {
		global.TOOL_LOG.Error("获取视频cid失败:", zap.Error(err))
		return errors.New("获取视频cid失败")
	}
	if v.Message != "0" {
		errStr := fmt.Sprintf("%s 获取cid失败: %s", bvid, v.Message)
		return errors.New(errStr)
	}
	if v.Data.Videos == 1 {
		//todo 查数据库检查cid是否存在
		result := global.TOOL_DB.Where("cid = ?", v.Data.Pages[0].Cid).Limit(1).Find(&example.VideoInstance{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return errors.New(v.Data.Title + ":文件已存在")
		}
		instance := &example.VideoInstance{
			Owner:    v.Data.Owner.Name,
			Bvid:     bvid,
			Cid:      v.Data.Pages[0].Cid,
			Title:    utils.FilterFilename(v.Data.Title),
			SavePath: global.TOOL_CONFIG.Bilibili.DownloadPath,
			Status:   1,
		}
		ok := example.AddVideoDownloading(instance)
		if ok {
			go instance.Download(false)
		} else {
			return errors.New(instance.Title + ":文件已存在")
		}
	} else {
		var parts []string
		for _, p := range v.Data.Pages {
			//todo 查数据库检查cid是否存在
			result := global.TOOL_DB.Where("cid = ?", v.Data.Pages[0].Cid).Limit(1).Find(&example.VideoInstance{})
			if result.Error != nil {
				global.TOOL_LOG.Error(utils.FilterFilename(p.Part)+"查询数据库失败", zap.Error(result.Error))
				continue
			}
			if result.RowsAffected > 0 {
				parts = append(parts, utils.FilterFilename(p.Part))
				continue
			}
			instance := &example.VideoInstance{
				Owner:    v.Data.Owner.Name,
				Bvid:     bvid,
				Cid:      p.Cid,
				Title:    utils.FilterFilename(p.Part),
				SavePath: filepath.Join(global.TOOL_CONFIG.Bilibili.DownloadPath, utils.FilterFilename(v.Data.Title)),
				Status:   1,
			}
			ok := example.AddVideoDownloading(instance)
			if ok {
				go instance.Download(false)
			} else {
				parts = append(parts, instance.Title)
			}
		}
		if len(parts) > 0 {
			return errors.New(strings.Join(parts, ",") + ":文件已存在")
		}
	}

	return nil
}
