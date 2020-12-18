package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"mygin_websrv/conf"
	"mygin_websrv/models"
	"mygin_websrv/modules/cache"
	"mygin_websrv/modules/response"
	"mygin_websrv/public/common"
	"strconv"
	"time"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	var u User
	err := ctx.BindJSON(&u)
	if err != nil {
		response.ShowError(ctx, "fail")
		return
	}
	if u.Username == "" || u.Password == "" {
		response.ShowError(ctx, "fail")
		return
	}
	user := models.SystemUser{Name: u.Username}
	if exist := user.UserExist(); !exist {
		response.ShowError(ctx, "fail")
		return
	}
	if common.Sha1En(u.Password+user.Salt) != user.Password {
		response.ShowError(ctx, "err_password")
		return
	}
	session := sessions.Default(ctx)
	var data = make(map[string]interface{}, 0)
	token := session.Get(conf.Cfg.Token)
	if token == nil {
		curTime := time.Now()
		timestamps := curTime.UnixNano()
		times := strconv.FormatInt(timestamps, 10)
		token = common.Md5En(common.GetRandomString(16) + times)
		session.Set(conf.Cfg.Token, token)
		session.Set(token, user.Id)
		err = session.Save()
		log.Println("session config success")
	}
	data[conf.Cfg.Token] = token
	response.ShowData(ctx, data)
	return
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	token := session.Get(conf.Cfg.Token)
	id := session.Get(token)
	strId := strconv.Itoa(id.(int))
	menuKey := conf.Cfg.RedisPre + "menu." + strId
	rc := cache.RedisClient.Get()
	defer rc.Close()
	if _, err := rc.Do("del", menuKey); err != nil { // redis key not del ???
		log.Println("redis delete failed", err)
	}
	session.Clear()
	session.Save()
	response.ShowSuccess(ctx, "success")
	return
}
