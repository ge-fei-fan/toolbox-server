package cron

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path/filepath"
	"regexp"
	"strings"
	"toolbox-server/global"
	"toolbox-server/model/example"
	"toolbox-server/utils"
)

// b站用户视频采集
func AutoCollect(job gocron.Job) {
	var users []example.BilibiliCollect
	err := global.TOOL_DB.Where("is_collect =?", true).Find(&users).Error
	if err != nil {
		global.TOOL_LOG.Error("查询采集用户出错", zap.Error(err))
		return
	}
	if len(users) == 0 {
		global.TOOL_LOG.Info("没有需要采集的用户")
		return
	}
	ch := make(chan struct{}, 5)
	for _, user := range users {
		ch <- struct{}{}
		go func(u example.BilibiliCollect) {
			defer func() {
				<-ch
			}()
			//vals, _ := query.Values(&struct {
			//	Mid string
			//	Ps  string
			//	Pn  string
			//}{
			//	Mid: u.Mid,
			//	Ps:  "10",
			//	Pn:  "1",
			//})
			//params := vals.Encode()
			//if params == "" {
			//	global.TOOL_LOG.Error("params 为空")
			//	return
			//}
			urlStr := fmt.Sprintf("https://api.bilibili.com/x/space/wbi/arc/search?mid=%d&ps=%s&pn=%s", u.Mid, "10", "1")
			arcUrl := utils.SignURL(urlStr)

			var guvr example.GetUserVideosResult
			client := resty.New()
			_, err = client.R().SetHeaders(map[string]string{
				"user-agent": global.UserAgent,
			}).SetResult(&guvr).Get(arcUrl)
			if err != nil {
				return
			}
			if guvr.Message != "0" {
				global.TOOL_LOG.Error(u.Name + "GetUserVideo err:" + guvr.Message)
				return
			}
			//查询当前用户最新视频
			var lastVideo int64
			var collectV example.CollectVideo
			err = global.TOOL_DB.Where("bilibili_collect_id = ?", u.ID).Order("created desc").First(&collectV).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					lastVideo = 0
				} else {
					global.TOOL_LOG.Error("查询用户最新视频出错", zap.Error(err))
					return
				}
			}
			lastVideo = collectV.Created
			for _, v := range guvr.Data.List.Vlist {

				if v.Created > lastVideo {
					//采集到视频更记录到数据库
					err = global.TOOL_DB.Create(&example.CollectVideo{
						BilibiliCollectID: u.ID,
						SpaceVideo:        v,
					}).Error
					if err != nil {
						global.TOOL_LOG.Error("插入采集视频表数据错误", zap.Error(err))
						return
					}
					//下载视频
					err = downloadVideoByBvid(v.Bvid)
					if err != nil {
						global.TOOL_LOG.Error("下载视频失败", zap.Error(err))
						return
					}
					if lastVideo == 0 {
						return
					}
				} else {
					break
				}

			}
		}(user)
	}
}

func downloadVideoByBvid(bvid string) error {
	v := &example.VideoInfo{}
	client := resty.New()
	_, err := client.R().SetQueryParams(map[string]string{"bvid": bvid}).
		SetResult(v).Get("https://api.bilibili.com/x/web-interface/view")
	if err != nil {
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

// ===================================================

// Todo 未测试等下次账号过期后测试
func ReFreshToken(job gocron.Job) {
	/*
		while True:
		if 每日第一次访问接口:
			if 检查是否需要刷新(cookie):
				CorrespondPath = 生成CorrespondPath(当前毫秒时间戳)
				refresh_csrf = 获取refresh_csrf(CorrespondPath, cookie)
				refresh_token_old = refresh_token # 这一步必须保存旧的 refresh_token 备用
				cookie, refresh_token = 刷新Cookie(refresh_token, refresh_csrf, cookie)
				确认更新(refresh_token_old, cookie) # 这一步需要新的 Cookie 以及旧的 refresh_token
				SSO站点跨域登录(cookie)
		do_somethings(cookie) # 其他业务逻辑处理
	*/
	const infoUrl = "https://passport.bilibili.com/x/passport-login/web/cookie/info"
	const csrfUrl = "https://www.bilibili.com/correspond/1/%s"
	const freshUrl = "https://passport.bilibili.com/x/passport-login/web/cookie/refresh"
	cookie := fmt.Sprintf("SESSDATA=%s", global.TOOL_CONFIG.Bilibili.SessData)
	var ci example.CookieInfo
	client := resty.New()
	_, err := client.R().SetResult(&ci).SetHeaders(map[string]string{
		"cookie": cookie,
	}).Get(infoUrl)
	if err != nil {
		global.TOOL_LOG.Error("获取是否需要刷新接口失败：", zap.Error(err))
		return
	}
	if ci.Code == -101 {
		global.TOOL_LOG.Error(ci.Message)
		return
	}
	if ci.Code == 0 {
		if ci.Data.Refresh { //需要刷新token
			correspondPath, err := getCorrespondPath(ci.Data.Timestamp)
			if err != nil {
				global.TOOL_LOG.Error("获取correspondPath失败：", zap.Error(err))
				return
			}
			u := fmt.Sprintf(csrfUrl, correspondPath)
			resp, err := client.R().SetHeaders(map[string]string{
				"cookie": cookie,
			}).Get(u)
			var refreshCsrf string
			pattern := `<div id="1-name">(.*?)</div>` // 匹配 div 标签和其中的文本内容
			reg := regexp.MustCompile(pattern)        // 编译正则表达式
			matches := reg.FindStringSubmatch(resp.String())
			if len(matches) > 1 {
				refreshCsrf = matches[1] // 获取括号内的内容
			} else {
				// 未找到匹配结果
				global.TOOL_LOG.Error("获取refresh_csrf失败：")
				return
			}

			//刷新cookie
			var fci example.ReFreshCookieInfo
			resp, err = client.R().SetQueryParams(map[string]string{
				"csrf":          global.TOOL_CONFIG.Bilibili.BiliJct,
				"refresh_csrf":  refreshCsrf,
				"source":        "main_web",
				"refresh_token": global.TOOL_CONFIG.Bilibili.RefreshToken,
			}).SetHeaders(map[string]string{
				"cookie": cookie,
			}).SetResult(&fci).Post(freshUrl)

			switch fci.Code {
			case 0:
				//for _, ck := range resp.Cookies() {
				//	if ck.Name == "bili_jct" {
				//		err = checkFlash(ck.Value, global.TOOL_CONFIG.Bilibili.SessData)
				//		if err != nil {
				//			global.TOOL_LOG.Error("checkFlash err:", zap.Error(err))
				//			return
				//		}
				//		global.TOOL_CONFIG.Bilibili.BiliJct = ck.Value
				//		break
				//	}
				//}
				//for _, ck := range resp.Cookies() {
				//	if ck.Name == "SESSDATA" {
				//		global.TOOL_CONFIG.Bilibili.SessData = ck.Value
				//		break
				//	}
				//}
				//global.TOOL_CONFIG.Bilibili.RefreshToken = fci.Data.RefreshToken
				//err = utils.WriteConfig()
				//if err != nil {
				//	global.TOOL_LOG.Error("WriteConfig err:", zap.Error(err))
				//	return
				//}
				var sessdata, bilijct string
				for _, ck := range resp.Cookies() {
					switch ck.Name {
					case "SESSDATA":
						sessdata = ck.Value
						break
					case "bili_jct":
						bilijct = ck.Value
						break
					}
				}

				err = checkFlash(bilijct, global.TOOL_CONFIG.Bilibili.RefreshToken)
				if err != nil {
					global.TOOL_LOG.Error("checkFlash err:", zap.Error(err))
					return
				}
				global.TOOL_CONFIG.Bilibili.BiliJct = bilijct
				global.TOOL_CONFIG.Bilibili.SessData = sessdata
				global.TOOL_CONFIG.Bilibili.RefreshToken = fci.Data.RefreshToken
				err = global.WriteConfig()
				if err != nil {
					global.TOOL_LOG.Error("WriteConfig err:", zap.Error(err))
					return
				}
			case -101:
				global.TOOL_LOG.Error("账号未登录")
				return
			case -111:
				global.TOOL_LOG.Error("csrf 校验失败")
				return
			case 86095:
				global.TOOL_LOG.Error("refresh_csrf 错误或 refresh_token 与 cookie 不匹配")
				return

			}
		}
	}
}
func getCorrespondPath(ts int64) (string, error) {
	const publicKey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg\nUc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71\nnzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40\nJNrRuoEUXpabUzGB8QIDAQAB\n-----END PUBLIC KEY-----"
	//生成毫秒时间戳
	//ts := time.Now().UnixNano() / int64(time.Millisecond)
	// 解析公钥
	publicKeyBlock, _ := pem.Decode([]byte(publicKey))
	if publicKeyBlock == nil {
		return "", errors.New("failed to decode public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return "", err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("failed to parse RSA public key")
	}

	// 使用 RSA-OAEP 加密
	label := []byte("refresh_" + fmt.Sprintf("%d", ts))
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		rsaPubKey,
		label,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(ciphertext), nil
}
func checkFlash(csrf, refreshToken string) error {
	const checkFreshUrl = "https://passport.bilibili.com/x/passport-login/web/confirm/refresh"
	res := struct {
		Code    int16  `json:"code"`
		Message string `json:"message"`
		Ttl     int16  `json:"ttl"`
	}{}
	cookie := fmt.Sprintf("SESSDATA=%s", global.TOOL_CONFIG.Bilibili.SessData)
	client := resty.New()
	_, err := client.R().SetResult(&res).SetHeaders(map[string]string{
		"cookie": cookie,
	}).SetQueryParams(map[string]string{
		"csrf":          csrf,
		"refresh_token": refreshToken,
	}).Post(checkFreshUrl)
	if err != nil {
		return err
	}
	switch res.Code {
	case 0:
		return nil
	case -101:
		return errors.New("账号未登录")
	case -111:
		return errors.New("csrf 校验失败")
	case -400:
		return errors.New("请求错误")
	}
	return nil
}

//===================================================
