package example

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"toolbox-server/global"
	"toolbox-server/model/common/response"
	"toolbox-server/model/example"
	"toolbox-server/model/example/request"
)

type BilibiliApi struct{}

func (b *BilibiliApi) Login(c *gin.Context) {

}

func (b *BilibiliApi) Generate(c *gin.Context) {
	qr, err := biliQrcodeService.GetQrcode()
	if err != nil {
		global.TOOL_LOG.Error("获取登录二维码失败", zap.Error(err))
		response.FailWithMessage("获取登录二维码失败", c)
		return
	}
	if qr.Message != "0" {
		response.FailWithMessage(qr.Message, c)
		return
	}
	if qr.Data.Url == "" {
		response.FailWithMessage("获取登录二维码连接为空", c)
		return
	}
	response.OkWithDetailed(qr.Data, "获取二维码连接成功", c)
}

func (b *BilibiliApi) Poll(c *gin.Context) {
	key, ok := c.GetQuery("qrcodekey")
	if !ok {
		response.FailWithMessage("qrcodekey为空", c)
		return
	}
	pr, err := biliQrcodeService.PullQrcode(key)
	if err != nil {
		global.TOOL_LOG.Error("获取扫码状态失败：", zap.Error(err))
		response.FailWithMessage("获取扫码状态失败", c)
	}
	response.OkWithDetailed(pr.Data, "获取扫码状态成功", c)
}

func (b *BilibiliApi) Account(c *gin.Context) {
	a, err := biliService.GetAccountInfo()
	if err != nil {
		global.TOOL_LOG.Error("获取账号信息失败", zap.Error(err))
		response.FailWithMessage("获取账号信息失败", c)
	}
	response.OkWithDetailed(a, "获取账号信息成功", c)

}

// 网页链接直接下载
func (b *BilibiliApi) Download(c *gin.Context) {
	key, ok := c.GetQuery("url")
	if !ok || key == "" {
		response.FailWithMessage("链接为空", c)
		return
	}
	err := biliService.DownloadVideo(key)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("请求成功", c)
}

func (b *BilibiliApi) GetVideoList(c *gin.Context) {
	var pageInfo request.ExaVideoSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//todo 判断校验
	list, total, err := biliService.GetVideoList(pageInfo)
	if err != nil {
		global.TOOL_LOG.Error("获取视频数据失败!", zap.Error(err))
		response.FailWithMessage("获取视频数据失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
func (b *BilibiliApi) DeleteVideo(c *gin.Context) {
	var video example.VideoInstance
	err := c.ShouldBindJSON(&video)
	if err != nil {
		global.TOOL_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	if video.ID == 0 {
		response.FailWithMessage("视频id为空", c)
		return
	}
	err = biliService.DeleteVideo(video)
	if err != nil {
		global.TOOL_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// 重新下载
func (b *BilibiliApi) ReDownload(c *gin.Context) {
	key, ok := c.GetQuery("cid")
	if !ok || key == "" {
		response.FailWithMessage("cid为空", c)
		return
	}
	err := biliService.ReDownloadVideo(key)
	if err != nil {
		global.TOOL_LOG.Error("重新下载失败!", zap.Error(err))
		response.FailWithMessage("重新下载失败", c)
		return
	}
	response.OkWithMessage("请求成功", c)
}

// 新增采集用户
func (b *BilibiliApi) AddCollectUser(c *gin.Context) {
	var user example.BilibiliCollect
	err := c.ShouldBind(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if user.UserDetail.Mid == 0 {
		response.FailWithMessage("mid为空", c)
		return
	}
	err = biliService.AddCollectUser(user)
	if err != nil {
		global.TOOL_LOG.Error("添加采集用户失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("请求成功", c)
}

// 采集用户列表
func (b *BilibiliApi) CollectUserList(c *gin.Context) {
	var pageInfo request.ExaCollectUserSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := biliService.CollectUserList(pageInfo)
	if err != nil {
		global.TOOL_LOG.Error("查询采集用户信息失败!", zap.Error(err))
		response.FailWithMessage("查询采集用户信息失败!", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 采集用户启用禁用
func (b *BilibiliApi) CollectUserStatus(c *gin.Context) {
	var user example.BilibiliCollect
	err := c.ShouldBind(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if user.ID == 0 {
		response.FailWithMessage("id为空", c)
		return
	}
	err = biliService.CollectUserStatus(user)
	if err != nil {
		global.TOOL_LOG.Error("修改采集用户状态失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("请求成功", c)
}

// 采集视频列表
func (b *BilibiliApi) CollectVideoList(c *gin.Context) {
	var pageInfo request.ExaCollectVideoSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := biliService.CollectVideoList(pageInfo)
	if err != nil {
		global.TOOL_LOG.Error("查询采集用视频列表失败!", zap.Error(err))
		response.FailWithMessage("查询采集用视频列表失败!", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (b *BilibiliApi) VideoDetail(c *gin.Context) {
	var video example.VideoInstance
	err := c.ShouldBindQuery(&video)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if video.ID == 0 {
		response.FailWithMessage("视频id为空", c)
		return
	}
	result, err := biliService.VideoDetail(video)
	if err != nil {
		global.TOOL_LOG.Error("查询采集用视频信息失败!", zap.Error(err))
		response.FailWithMessage("查询采集用视频信息失败!", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
