package main

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
	"mygin_websrv/conf"
	"mygin_websrv/ctrl/article"
	"mygin_websrv/ctrl/fileopt"
	"mygin_websrv/ctrl/menu"
	"mygin_websrv/ctrl/role"
	"mygin_websrv/ctrl/user"
	"mygin_websrv/models"
	"mygin_websrv/modules/cache"
	"mygin_websrv/modules/response"
	"mygin_websrv/public/common"
	"net/http"
	"net/url"
)

func init() {
	c := conf.Config{}
	c.Routers = []string{"/ping", "/login", "/role/index", "/info", "/dashboard", "/logout"}
	conf.Set(c)
	switch {
	case c.Env == "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func sessionStoreInit(engine *gin.Engine) {
	if sessionStore, err := redis.NewStoreWithPool(cache.RedisClient, []byte("secret")); err != nil {
		log.Println(err)
		panic(errors.New("redis client start error"))
	} else {
		engine.Use(sessions.Sessions("gsession", sessionStore))
	}
}

//api url config
func routerInit(engine *gin.Engine) {
	//心跳接口
	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})
	engine.POST("/upload/image", fileopt.ImgUpload)
	engine.GET("/del/image", fileopt.DelImage)

	engine.POST("/login", user.Login)
	engine.POST("/logout", user.Logout)

	engine.GET("/info", user.Info)
	engine.GET("/routes", menu.List)
	engine.GET("/dashboard", menu.Dashboard)
	engine.GET("/role/list", menu.Roles)
	engine.GET("/menu", menu.Index)
	engine.POST("/menu", menu.Create)
	engine.PUT("/menu", menu.Edit)
	engine.DELETE("/menu", menu.Delete)
	engine.GET("/user", user.Index)
	engine.GET("/user/detail", user.Detail)
	engine.GET("/user/search", user.Search)
	engine.POST("/user/create", user.Create)
	engine.POST("/user/edit", user.Edit)
	engine.POST("/user/repasswd", user.Repasswd)
	engine.GET("/user/delete", user.Delete)
	engine.POST("/role/delete/:name", role.DeleteRole)
	engine.POST("/role/update", role.UpdateRole)
	engine.POST("/role/add", role.AddRole)
	engine.GET("/role/index", role.Index)
	engine.POST("/reg", user.Reg)
	engine.POST("/articles/create", article.Create)
	engine.POST("/articles/edit", article.Edit)
	engine.GET("/articles/list", article.Index)
	engine.GET("/articles/detail", article.Detail)
	engine.GET("/showimage", article.ShowImage)

}

//跨域处理
func corsConfig(engine *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://chenjia.joyveb.com", "http://localhost:9529", "http://localhost:9528", "http://localhost:9527", "http://localhost"}
	config.AllowMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"x-requested-with", "Content-Type", "AccessToken", "X-CSRF-Token", "X-Token", "Authorization", "token"}
	engine.Use(cors.New(config))
}

func authConfig(engine *gin.Engine) {
	auth := func(context *gin.Context) {
		u, err := url.Parse(context.Request.RequestURI)
		if err != nil {
			log.Println("auth config err", err)
			panic(err)
		}
		if common.InArrayString(u.Path, &conf.Cfg.Routers) {
			context.Next()
			return
		}
		session := sessions.Default(context)
		token := session.Get(conf.Cfg.Token)
		if token == nil {
			context.Abort()
			response.ShowError(context, "not_login")
			return
		}
		uid := session.Get(token)
		user := models.SystemUser{Id: uid.(int), Status: 1}
		exist := user.UserExist()
		if !exist {
			context.Abort()
			response.ShowError(context, "user_error")
			return
		}
		//特殊账号
		if user.Name == conf.Cfg.Super {
			return
		}
		menuModel := models.SystemMenu{}
		menuMap, err := menuModel.GetRouteByUid(uid)
		if err != nil {
			context.Abort()
			response.ShowError(context, "unauthorized")
			return
		}
		if _, ok := menuMap[u.Path]; !ok {
			context.Abort()
			response.ShowError(context, "unauthorized")
			return
		}
		// access the status we are sending
		context.Writer.Status()
		context.Next()
	}
	engine.Use(auth)
}

func main() {
	engine := gin.Default()  //初始化gin
	sessionStoreInit(engine) //redis初始化
	corsConfig(engine)       //跨域问题
	authConfig(engine)       //权限认证
	routerInit(engine)       //路由初始化
	engine.Run()
}
