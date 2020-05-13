package auth

import (
	"github.com/astaxie/beego/context"
	"github.com/wonderivan/logger"
	"encoding/json"
	"strings"
	"ant-api/cache"
)


// 定义返回结果
type ResResult struct {
	Status string      // Success Failure NoLogin  用户没有登录则跳转到登录页面
  Data string
}

// 将ResResult转换成字符串
func (res ResResult) ConvertToString() string {
	data,_ := json.Marshal(res)
	return string(data)
}

// 用户token验证
var FilterAuth = func(ctx *context.Context) {
    var res ResResult
		// 除/login外的请求必须进行用户验证
    if !strings.HasPrefix(ctx.Input.URL(), "/user/login/") {
			username := ctx.GetCookie("username")
			usertoken := ctx.GetCookie("usertoken")
			// 从redis中获取token
			sessionToekn, _ := cache.Get(username)
			// 如果token验证不通过
			if usertoken == "" || usertoken != sessionToekn {
				logger.Info(username + " 用户token过期")
				res.Status = "NoLogin"
				res.Data = "用户token过期"
				result := res.ConvertToString()
				ctx.WriteString(result)
				return
			}
    }
}
