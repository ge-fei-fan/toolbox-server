package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestF(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	r, err := CheckFFmpeg()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(r)
}
