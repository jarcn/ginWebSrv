package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"mygin_websrv/conf"
	"mygin_websrv/models"
	"mygin_websrv/modules/request"
	"mygin_websrv/modules/response"
	"mygin_websrv/public/common"
	"strconv"
	"time"
)

type Userinfo struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}

type UserDetail struct {
	models.SystemUser
	CheckedRoles []string `json:"checkedRoles"`
}

func Reg(c *gin.Context) {
	name := c.PostForm("name")
	passwd := c.PostForm("passwd")
	if name == "" || passwd == "" {
		response.ShowError(c, "fail")
		return
	}
	salt := common.GetRandomBoth(4)
	passwd = common.Sha1En(passwd + salt)
}

func Info(ctx *gin.Context) {
	session := sessions.Default(ctx)
	token := session.Get(conf.Cfg.Token)
	if token == nil {
		response.ShowError(ctx, "not_login")
		return
	}
	uid := session.Get(token)
	user := models.SystemUser{Id: uid.(int)}
	if exist := user.UserExist(); !exist {
		response.ShowError(ctx, "user_error")
		return
	}
	roles := models.SystemUserRole{SystemUserId: uid.(int)}
	role, _ := roles.GetRowByUid()
	var info Userinfo
	info.Roles = role
	info.Name = user.Name
	info.Avatar = user.Avatar
	info.Introduction = user.Introduction
	response.ShowData(ctx, info)
	return
}

func Search(ctx *gin.Context) {
	name, exist := ctx.GetQuery("name")
	if !exist {
		response.ShowErrorParams(ctx, "name")
		return
	}
	user := models.SystemUser{}
	userBeans, _ := user.SelectByName(name)
	nameList := make(map[string][]models.SearchUser, 0)
	nameList["items"] = userBeans
	response.ShowData(ctx, nameList)
	return
}

func Detail(ctx *gin.Context) {
	id, exist := ctx.GetQuery("id")
	if !exist {
		response.ShowErrorParams(ctx, "id")
		return
	}
	user := models.SystemUser{}
	user.Id, _ = strconv.Atoi(id)
	exist = user.UserExist()
	if !exist {
		response.ShowError(ctx, "user_error")
		return
	}
	userRole := models.SystemUserRole{SystemUserId: user.Id}
	role, _ := userRole.GetRowByUid()
	detail := UserDetail{}
	detail.CheckedRoles = role
	detail.Id = user.Id
	detail.Name = user.Name
	detail.Nickname = user.Nickname
	detail.Phone = user.Phone
	detail.Status = user.Status
	response.ShowData(ctx, detail)
	return
}

func Index(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)

	paging := &common.Paging{Page: page, PageSize: limit}
	userModel := models.SystemUser{}
	userArr, err := userModel.SelectByPage(paging)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	data := make(map[string]interface{})
	data["items"] = userArr
	data["total"] = paging.Total
	response.ShowData(c, data)
	return
}

func Create(c *gin.Context) {
	data, err := request.GetJson(c)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["name"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["nickname"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["password"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["repassword"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["status"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	userModel := models.SystemUser{}
	userModel.Name = data["name"].(string)
	has := userModel.UserExist()
	if has {
		response.ShowError(c, "name_exists")
		return
	}
	userModel.Password = data["password"].(string)
	if userModel.Password != data["repassword"].(string) {
		response.ShowError(c, "fail")
		return
	}
	userModel.Salt = common.GetRandomBoth(4)
	userModel.Password = common.Sha1En(userModel.Password + userModel.Salt)
	userModel.Name = data["name"].(string)
	userModel.Nickname = data["nickname"].(string)
	if _, ok := data["phone"]; ok {
		userModel.Phone = data["phone"].(string)
	}
	if _, ok := data["status"]; ok && data["status"].(bool) {
		userModel.Status = 1
	}
	userModel.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	userModel.Ctime = time.Now()
	if _, ok := data["checkedRoles"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	roles := data["checkedRoles"].([]interface{})
	_, err = userModel.Insert(roles)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowData(c, userModel)
	return
}

func Delete(c *gin.Context) {
	id, has := c.GetQuery("id")
	if !has {
		response.ShowErrorParams(c, "id")
		return
	}
	user := models.SystemUser{}
	user.Id, _ = strconv.Atoi(id)
	err := user.Delete()
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowData(c, "success")
	return
}

func Edit(c *gin.Context) {
	data, err := request.GetJson(c)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["id"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	userModel := models.SystemUser{}
	userModel.Id = int(data["id"].(float64))
	has := userModel.UserExist()
	if !has {
		response.ShowError(c, "user_error")
		return
	}
	if _, ok := data["nickname"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["status"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["status"]; ok && data["status"].(bool) {
		userModel.Status = 1
	} else {
		userModel.Status = 0
	}
	userModel.Nickname = data["nickname"].(string)
	if _, ok := data["phone"]; ok {
		userModel.Phone = data["phone"].(string)
	}
	if _, ok := data["checkedRoles"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	roles := data["checkedRoles"].([]interface{})
	err = userModel.Update(roles)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowData(c, userModel)
	return
}
func Repasswd(c *gin.Context) {
	data, err := request.GetJson(c)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["id"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	userModel := models.SystemUser{}
	userModel.Id = int(data["id"].(float64))
	has := userModel.UserExist()
	if !has {
		response.ShowError(c, "user_error")
		return
	}
	if userModel.Name == "admin" {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["password"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	if _, ok := data["repassword"]; !ok {
		response.ShowError(c, "fail")
		return
	}
	userModel.Password = data["password"].(string)
	if userModel.Password != data["repassword"].(string) {
		response.ShowError(c, "fail")
		return
	}
	userModel.Salt = common.GetRandomBoth(4)
	userModel.Password = common.Sha1En(userModel.Password + userModel.Salt)
	err = userModel.UpdatePasswd()
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowData(c, userModel)
	return
}
