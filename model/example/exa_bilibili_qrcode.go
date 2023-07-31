package example

// BiliQrcode bilibili登录二维码连接响应
type BiliQrcode struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    qrcode `json:"data"`
}

type qrcode struct {
	Url       string `json:"url"`        // 二维码内容url
	QrcodeKey string `json:"qrcode_key"` // 扫码登录秘钥
}

// PullQrResp 扫描二维码状态响应
type PullQrcode struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RefreshToken string `json:"refresh_token"`
		Code         int64  `json:"code"`
		Message      string `json:"message"`
	} `json:"data"`
}
