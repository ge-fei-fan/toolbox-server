//go:build release

package main

import (
	"github.com/gin-gonic/gin"
)

func init() {
	//fmt.Println("Release Mode")
	gin.SetMode(gin.ReleaseMode)
}

// 此处可以添加其它你需要在发布模式下执行的代码
