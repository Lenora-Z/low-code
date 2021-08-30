// Package server Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 11:03 上午
package server

import (
	"github.com/Lenora-Z/low-code/service/application"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (ds *defaultServer) appAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		app := context.Request.Header.Get("app-id")
		if app == "" {
			context.JSON(http.StatusBadRequest, ApiResponse{
				Code:    ACCESS_DENY,
				Message: "invalid app-id",
			})
			context.Abort()
			return
		}
		srv := application.NewApplicationService(ds.db)
		appId := utils.NewStr(app)
		status, item := srv.GetApplication(appId.Uint64())
		if status {
			context.JSON(http.StatusBadRequest, ApiResponse{
				Code:    ACCESS_DENY,
				Message: "app not exist!",
			})
			context.Abort()
			return
		}
		claims := context.MustGet(CLAIMS).(*CustomClaims)
		if claims == nil {
			claims = &CustomClaims{
				AppId:     item.ID,
				UUid:      item.AppHash,
				AppStatus: item.Status,
			}
		} else {
			claims.AppId = item.ID
			claims.UUid = item.AppHash
			claims.AppStatus = item.Status
		}
		context.Set(CLAIMS, claims)
		context.Next()
	}
}

func (ds *defaultServer) appHashAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		hash := context.Request.Header.Get("app-hash")
		if hash == "" {
			context.JSON(http.StatusBadRequest, ApiResponse{
				Code: ACCESS_DENY,
			})
			context.Abort()
			return
		}
		srv := application.NewApplicationService(ds.db)
		status, item := srv.GetApplicationByHash(hash)
		if status {
			context.JSON(http.StatusBadRequest, ApiResponse{
				Code:    ACCESS_DENY,
				Message: "app not exist!",
			})
			context.Abort()
			return
		}
		claims := CustomClaims{
			AppId: item.ID,
			UUid:  item.AppHash,
		}
		context.Set(CLAIMS, &claims)
		context.Next()
	}
}

func (ds *defaultServer) tritiumAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		users := context.Request.Header.Get("x-dataqin-userid")
		logrus.Info(users)
		if users == "" || utils.NewStr(users).Uint64() == 0 {
			context.JSON(http.StatusNonAuthoritativeInfo, ApiResponse{
				Code:    ACCESS_DENY,
				Message: "invalid x-dataqin-userid",
			})
			context.Abort()
			return
		}
		context.Set(CLAIMS, &CustomClaims{UserId: utils.NewStr(users).Uint64()})
		context.Next()
	}
}

func (ds *defaultServer) checkUserApiPermission() gin.HandlerFunc {
	return func(context *gin.Context) {
		routeId := context.Request.Header.Get("route-id")
		if routeId == "" {
			logrus.Error("route empty")
			context.JSON(http.StatusBadRequest, ApiResponse{
				Code:    ACCESS_DENY,
				Message: "invalid route-id",
			})
			context.Abort()
			return
		}
		claims := context.MustGet(CLAIMS).(*UserClaims)
		userSrv := user.NewUserService(ds.db)
		err, _, roles := userSrv.GetRelationByOrgId(claims.GroupId)
		if err != nil {
			logrus.Error("get organization relation error:", err)
			ds.ResponseError(context, ACCESS_DENY)
			context.Abort()
			return
		}

		route := utils.NewStr(routeId)
		status, _ := userSrv.SearchPermission(roles.RoleIds(), 4, route.Uint64())

		if status {
			logrus.Error("access deny")
			ds.ResponseError(context, ACCESS_DENY)
			context.Abort()
			return
		}
		context.Next()
	}
}
