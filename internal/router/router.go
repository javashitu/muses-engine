package router

import (
	"muses-engine/internal/api"
	"muses-engine/internal/middlewares"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("some"))
	// logging.HttpLogToFile("release")
	r.Use(sessions.Sessions("mysession", store))
	r.Use(middlewares.Cores())
	// v1 := r.Group("api/media").Use(AuthByJwt())
	v1 := r.Group("api/media").Use()
	{
		v1.GET("captchId", api.GetCaptchaId)
		v1.POST("upload", api.UploadFile)
		v1.POST("preSign", api.PreSign)
	}
	return r
}
