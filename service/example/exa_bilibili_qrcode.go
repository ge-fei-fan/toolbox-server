package example

import (
	"github.com/go-resty/resty/v2"
	"toolbox-server/global"
	"toolbox-server/model/example"
)

type BilibiliQrcode struct{}

const (
	generateQrcodeUrl = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	pullQrcodeUrl     = "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
)

func (b *BilibiliQrcode) GetQrcode() (*example.BiliQrcode, error) {
	client := resty.New()
	//var qr example.BiliQrcode
	qr := &example.BiliQrcode{}
	_, err := client.R().SetResult(qr).Get(generateQrcodeUrl)
	if err != nil {
		return nil, err
	}
	return qr, err
}

func (b *BilibiliQrcode) PullQrcode(key string) (*example.PullQrcode, error) {
	client := resty.New()
	client.SetHeader("user-agent", global.UserAgent)
	client.SetQueryParam("qrcode_key", key)
	pr := &example.PullQrcode{}
	resp, err := client.R().SetResult(pr).Get(pullQrcodeUrl)
	if err != nil {
		return nil, err
	}
	if pr.Data.Code == 0 {
		global.TOOL_CONFIG.Bilibili.RefreshToken = pr.Data.RefreshToken
		for _, ck := range resp.Cookies() {
			switch ck.Name {
			case "SESSDATA":
				global.TOOL_CONFIG.Bilibili.SessData = ck.Value
				break
			case "bili_jct":
				global.TOOL_CONFIG.Bilibili.BiliJct = ck.Value
				break
			}
		}
		err = global.WriteConfig()
		if err != nil {
			return nil, err
		}
	}
	return pr, err

}
