//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 11:04 上午
package server

import (
	"github.com/Lenora-Z/low-code/service/application"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func (ds *defaultServer) CreateToken(uId, appId uint64, mobile, name string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uId,
		"exp":     time.Now().Add(utils.TokenExpire).Unix(),
		"mobile":  mobile,
		"account": name,
		"app_id":  appId,
	})

	token, err := at.SignedString([]byte(ds.conf.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (ds *defaultServer) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		g := strings.Split(header, " ")
		if g[0] != "Bearer" || g[1] == "" {
			ctx.JSON(http.StatusUnauthorized, ApiResponse{
				Code:    2002,
				Message: "invalid Authorization",
			})
			ctx.Abort()
			return
		}
		ret, err := jwt.Parse(g[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(ds.conf.JWTSecret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ApiResponse{
				Code:    2002,
				Message: "invalid Authorization",
			})
			ctx.Abort()
			return
		}

		cla := ret.Claims.(jwt.MapClaims)
		expire := cla["exp"].(float64)
		//重置密码校验
		us := user.NewUserService(ds.db)
		_, user := us.GetUser(uint64(cla["user_id"].(float64)))
		if user.ID == 0 {
			ctx.JSON(http.StatusOK, ApiResponse{
				Code:    2001,
				Message: "invalid token",
			})
			ctx.Abort()
			return
		}
		if *user.PwdStatus {
			ctx.JSON(http.StatusOK, ApiResponse{
				Code:    2001,
				Message: "user reset password",
			})
			ctx.Abort()
			return
		}
		if expire < float64(time.Now().Unix()) {
			ctx.JSON(http.StatusOK, ApiResponse{
				Code:    2001,
				Message: "user licence expired",
			})
			ctx.Abort()
			return
		}
		appId := uint64(cla["app_id"].(float64))
		appSrv := application.NewApplicationService(ds.db)
		status, app := appSrv.GetApplication(appId)
		if status {
			logrus.Error("app is not exist")
			ctx.JSON(http.StatusOK, ApiResponse{
				Code:    FAIL,
				Message: "app is not exist",
			})
			ctx.Abort()
			return
		}

		if app.Status != application.ONLINE {
			logrus.Error("app is not online")
			ctx.JSON(http.StatusOK, ApiResponse{
				Code:    FAIL,
				Message: "app is not online",
			})
			ctx.Abort()
			return
		}

		claims := UserClaims{
			AppId:   appId,
			UserId:  uint64(cla["user_id"].(float64)),
			Account: cla["account"].(string),
			GroupId: user.GroupID,
			AppName: app.Name,
		}

		ctx.Set("claims", &claims)
		ctx.Next()
	}
}
