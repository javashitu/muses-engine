package router

import (
	"log"
	"muses-engine/config"
	"muses-engine/internal/common"
	"muses-engine/internal/model/response"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var Jwtkey = []byte(config.Server.HttpConfig.JwtKey)

type UserClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func CreateToken(userId string) (string, error) {
	//过期时间
	expireTime := time.Now().Add(time.Duration(config.Server.HttpConfig.JwtExpireSeconds) * time.Second)
	log.Printf("expire time is %+v and config expire seconds %d \n", expireTime, config.Server.HttpConfig.JwtExpireSeconds)
	nowTime := time.Now()
	claims := UserClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString(Jwtkey)
}

func CheckToken(token string) (*UserClaims, bool) {
	log.Println("will beed check token is ", token)
	tokenObj, _ := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Jwtkey, nil
	})
	if key, _ := tokenObj.Claims.(*UserClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

func AuthByJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			log.Printf("no token in request %s \n ", c.Request.URL)
			c.JSON(http.StatusOK, response.FailureApi(common.UnauthenticatedError))
			c.Abort()
			return
		}
		tokenSlice := strings.SplitN(tokenStr, " ", 2)
		tokenStruck, ok := CheckToken(tokenSlice[1])
		if !ok {
			c.JSON(http.StatusOK, response.FailureApi(common.UnauthenticatedError))
			c.Abort()
			return
		}
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, response.FailureApi(common.TokenExpireError))
			c.Abort()
			return
		}
		c.Set("user_id", tokenStruck.UserId)

		c.Next()
	}
}
