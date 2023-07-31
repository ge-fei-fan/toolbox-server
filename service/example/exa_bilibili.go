package example

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"path/filepath"
	"strconv"
	"toolbox-server/global"
	"toolbox-server/model/example/request"
	"toolbox-server/utils"

	"toolbox-server/model/example"
)

type Bilibili struct{}

const (
	accountInfoUrl = "https://api.bilibili.com/x/member/web/account"
)

func (b *Bilibili) GetAccountInfo() (*example.Account, error) {
	client := resty.New()
	cookie := fmt.Sprintf("SESSDATA=%s", global.TOOL_CONFIG.Bilibili.SessData)
	client.SetHeader("cookie", cookie)
	a := &example.Account{}
	_, err := client.R().SetResult(a).Get(accountInfoUrl)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (b *Bilibili) DownloadVideo(s string) error {
	return DownloadVideo(s)
}
func (b *Bilibili) ReDownloadVideo(cid string) error {
	c, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return err
	}
	//var instance example.VideoInstance
	//err = global.TOOL_DB.Where("cid = ?",c).First(&instance).Error
	//if err != nil {
	//	return err
	//}
	//if instance.Status == -1{
	//
	//}
	example.VideoMutex.Lock()
	defer example.VideoMutex.Unlock()
	for _, v := range example.VideoDownloading {
		if v.Cid == c {
			if v.Status == -1 {
				go v.Download(true)
				return nil
			}
			return errors.New("重新下载失败，文件正在下载")
		}
	}
	return errors.New("重新下载失败，cid错误")
}
func (b *Bilibili) GetVideoList(info request.ExaVideoSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize                    //5  5 5
	offset := info.PageSize * (info.Page - 1) //0 5 10
	// 创建db
	db := global.TOOL_DB.Model(&example.VideoInstance{})
	videoInstances := make([]*example.VideoInstance, 0)
	if info.Status == 0 { //下载完成的
		db = db.Where("`status` = ?", 0)
	} else {
		//db = db.Not("`status` = ?", 0)
		total = int64(len(example.VideoDownloading))
		if len(example.VideoDownloading) > 0 {
			if len(example.VideoDownloading) > offset+limit {
				videoInstances = example.VideoDownloading[offset : offset+limit]
			} else if len(example.VideoDownloading) < offset+limit && len(example.VideoDownloading) > offset {
				videoInstances = example.VideoDownloading[offset:]
			}
		}

		return videoInstances, total, nil
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&videoInstances).Error
	return videoInstances, total, err
}
func (b *Bilibili) DeleteVideo(instance example.VideoInstance) error {
	//根据id拿到数据库里所有信息
	err := global.TOOL_DB.First(&instance).Error
	if err != nil {
		return err
	}
	//在下载列表 并且下载失败的视频  可以删除
	if len(example.VideoDownloading) != 0 {
		example.VideoMutex.Lock()
		defer example.VideoMutex.Unlock()
		for i, v := range example.VideoDownloading {
			if v.ID == instance.ID {
				if v.Status == -1 {
					example.VideoDownloading = append(example.VideoDownloading[:i], example.VideoDownloading[i+1:]...)
					break
				} else {
					return errors.New("视频正在下载，无法删除")
				}
			}
		}

	}
	//删除视频文件
	videoPath, _ := filepath.Abs(filepath.Join(instance.SavePath, instance.Title+".video"))
	audioPath, _ := filepath.Abs(filepath.Join(instance.SavePath, instance.Title+".audio"))
	outPath, _ := filepath.Abs(filepath.Join(instance.SavePath, instance.Title+".mp4"))
	has, _ := utils.FileExists(videoPath)
	if has {
		err = utils.DeleteFile(videoPath)
		if err != nil {
			return err
		}
	}
	has, _ = utils.FileExists(audioPath)
	if has {
		err = utils.DeleteFile(audioPath)
		if err != nil {
			return err
		}
	}
	has, _ = utils.FileExists(outPath)
	if has {
		err = utils.DeleteFile(outPath)
		if err != nil {
			return err
		}
	}

	//数据库删除
	err = global.TOOL_DB.Delete(&instance).Error
	if err != nil {
		return err
	}

	return err
}
func (b *Bilibili) AddCollectUser(user example.BilibiliCollect) error {
	err := global.TOOL_DB.Where("mid = ?", user.Mid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//获取用户信息
			urlStr := fmt.Sprintf("https://api.bilibili.com/x/space/wbi/acc/info?mid=%d", user.Mid)
			AccUrl := utils.SignURL(urlStr)
			client := resty.New()
			cookie := fmt.Sprintf("SESSDATA=%s", global.TOOL_CONFIG.Bilibili.SessData)
			var userSpaceRes example.UserSpaceResult
			_, err = client.R().SetHeaders(map[string]string{
				"cookie":     cookie,
				"User-Agent": global.UserAgent,
			}).SetResult(&userSpaceRes).Get(AccUrl)
			if err != nil {
				return err
			}
			if userSpaceRes.Code != 0 {
				return errors.New(userSpaceRes.Message)
			}
			user.UserDetail = userSpaceRes.Data
			user.IsCollect = true
			err = global.TOOL_DB.Create(&user).Error
			if err != nil {
				return err
			}
		} else {
			// 查询出错
			return err
		}
	} else {
		return errors.New("用户已存在")
	}
	return nil
}

func (b *Bilibili) CollectUserList(info request.ExaCollectUserSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize                    //5  5 5
	offset := info.PageSize * (info.Page - 1) //0 5 10
	// 创建db
	db := global.TOOL_DB.Model(&example.BilibiliCollect{})

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var users []example.BilibiliCollect
	err = db.Limit(limit).Offset(offset).Order("bilibili_collects.created_at desc").Find(&users).Error
	if err != nil {
		return
	}
	return users, total, err

}
func (b *Bilibili) CollectUserStatus(user example.BilibiliCollect) error {
	err := global.TOOL_DB.Where("ID = ?", user.ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		} else {
			// 查询出错
			return err
		}
	}
	if user.IsCollect {
		user.IsCollect = false
		err = global.TOOL_DB.Save(user).Error
		if err != nil {
			return err
		}
	} else {
		user.IsCollect = true
		err = global.TOOL_DB.Save(user).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bilibili) CollectVideoList(info request.ExaCollectVideoSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.TOOL_DB.Model(&example.CollectVideo{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var videos []example.CollectVideo
	err = db.Limit(limit).Offset(offset).Order("collect_videos.created_at desc").Joins("BilibiliCollect").Find(&videos).Error
	if err != nil {
		return
	}
	return videos, total, nil
}
func (b *Bilibili) VideoDetail(v example.VideoInstance) (video *example.VideoInstance, err error) {
	err = global.TOOL_DB.Where("ID = ?", v.ID).First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("视频不存在")
		} else {
			// 查询出错
			return nil, err
		}
	}
	return &v, nil
}
