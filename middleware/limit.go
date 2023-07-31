package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
	"toolbox-server/global"
	"toolbox-server/model/common/response"
)

const (
	LimitDefault int = 0
)

var LimitMap = make(map[string][]time.Duration)
var mx sync.Mutex

type LimitConfig struct {
	// GenerationKey 根据业务生成key 下面CheckOrMark查询生成
	GenerationKey func(c *gin.Context) string
	// 检查函数,用户可修改具体逻辑,更加灵活
	CheckOrMark func(key string, expire int, limit int) error
	// Expire key 过期时间
	Expire int
	// Limit 周期时间
	Limit int
}

func (l LimitConfig) LimitWithTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := l.CheckOrMark(l.GenerationKey(c), l.Expire, l.Limit); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": response.LIMITERROR, "msg": err.Error()})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}

// DefaultGenerationKey 默认生成key
func DefaultGenerationKey(c *gin.Context) string {
	//return "GVA_Limit" + c.ClientIP()
	return "TOOL_Limit" + c.FullPath()
}
func DefaultCheckOrMark(key string, expire int, limit int) (err error) {
	//使用非redis方法实现限流
	if LimitMap == nil {
		return err
	}
	mx.Lock()
	defer mx.Unlock()
	v, has := LimitMap[key]
	//访问过
	if has {
		//访问次数大于限制次数
		if len(v) >= limit {
			//判断是否在限制期限内
			now := time.Duration(time.Now().Unix())
			r := (now - v[0]) * time.Second
			if r >= time.Duration(expire)*time.Second {
				//应该情况列表，然后设置第一个为本次时间
				v = v[:0]
				LimitMap[key] = append(v, now)
				return nil
			}
			return errors.New("访问限制")
		} else { //小于访问次数限制，把这次访问加入切片中
			LimitMap[key] = append(v, time.Duration(time.Now().Unix()))
			return nil
		}
	} else { //没有访问过
		var timeSlice []time.Duration
		LimitMap[key] = append(timeSlice, time.Duration(time.Now().Unix()))
		return nil
	}
}

func DefaultLimit(count, expire int) gin.HandlerFunc {
	if count == LimitDefault {
		count = global.TOOL_CONFIG.System.LimitCountIP
	}
	if expire == LimitDefault {
		expire = global.TOOL_CONFIG.System.LimitCountIP
	}
	return LimitConfig{
		GenerationKey: DefaultGenerationKey,
		CheckOrMark:   DefaultCheckOrMark,
		Expire:        expire,
		Limit:         count,
	}.LimitWithTime()
}
