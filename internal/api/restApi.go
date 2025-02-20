package api

import (
	"log"
	"muses-engine/config"
	"muses-engine/internal/app/logic/fileSystem"
	"muses-engine/internal/common"
	"muses-engine/internal/model/request"
	"muses-engine/internal/model/response"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetCaptchaId(c *gin.Context) {
	captchId := captcha.NewLen(4)
	c.JSON(200, response.SuccessApi(captchId))
}

func UploadFile(c *gin.Context) {
	log.Println("server receive uplaod file request")

	form, _ := c.MultipartForm()
	file := form.File["file"]

	uploadFileReq := &request.UploadFileReq{}
	if err := c.ShouldBind(&uploadFileReq); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailureApi(common.BindParamError))
		log.Println("bind param error ", err)
		return
	}
	result := fileSystem.SaveFile(file, uploadFileReq)

	c.JSON(http.StatusOK, result)
}

func PreSign(c *gin.Context) {
	visitFileUrlReq := request.VisitFileUrlReq{}
	if err := c.ShouldBind(&visitFileUrlReq); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailureApi(common.BindParamError))
		log.Println("bind param error ", err)
		return
	}
	if visitFileUrlReq.ExpireSeconds < 0 || visitFileUrlReq.ExpireSeconds > config.MinioConfig.AccessExpires {
		c.JSON(http.StatusInternalServerError, response.FailureApi(common.ParamError))
		log.Println("param wrong ")
		return
	}
	apiResult := fileSystem.GenFileVisitUrl(visitFileUrlReq)

	c.JSON(http.StatusOK, apiResult)
}
