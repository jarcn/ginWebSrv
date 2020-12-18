package fileopt

import (
	"github.com/gin-gonic/gin"
	"mygin_websrv/conf"
	"mygin_websrv/modules/response"
	"mygin_websrv/public/common"
	"os"
	"path"
	"strings"
	"time"
)

func ImgUpload(ctx *gin.Context) {
	tmpFile, err := ctx.FormFile("file")
	if err != nil {
		response.ShowError(ctx, "fail")
		return
	}
	timeNow := time.Now().Unix()
	dayTime := time.Unix(timeNow, 0)

	uploadDir := "upload/"
	relativeDir := dayTime.Format("20060102") + "/"
	os.MkdirAll(uploadDir+relativeDir, os.ModePerm)

	ext := path.Ext(tmpFile.Filename)
	relativeDir = relativeDir + common.GetRandomBoth(32) + ext

	ctx.SaveUploadedFile(tmpFile, uploadDir+relativeDir)
	data := conf.Cfg.Host + "/showimage?imgname=upload/" + relativeDir
	response.ShowData(ctx, data)
	return
}

func DelImage(ctx *gin.Context) {
	url, has := ctx.GetQuery("url")
	if !has {
		response.ShowErrorParams(ctx, "url")
		return
	}
	url = common.SubstrContains(url, "upload/")
	//过虑危险字符
	url = strings.Replace(url, "../", "", -1)
	if common.IsFile(url) {
		err := os.Remove(url)
		if err != nil {
			response.ShowError(ctx, "fail")
			return
		}
	}
	response.ShowData(ctx, "success")
	return
}
