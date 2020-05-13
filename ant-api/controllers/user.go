package controllers

import (
	"ant-api/models"
	"ant-api/cache"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/google/uuid"
	// "fmt"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Login() {
	var user models.User
	var res ResResult
	// 将请求的body数据存入到user中
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)

	// 进行登录验证
	login := models.Login(user)
	Uuid := uuid.New()
	userToken := Uuid.String()
	// 如果验证成功，则在redis中保存token信息
	if login == "Success" {
		cache.Set(user.Username, userToken)
	}
	// 设置返回结果
	res.Status = login
	res.Data = userToken
	u.Ctx.WriteString(res.ConvertToString())
}
